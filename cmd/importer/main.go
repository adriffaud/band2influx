package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	b2i "driffaud.fr/adrien/band2influx"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var hostname string
var username string
var password string

func getDataPoints(dbFile string) ([]b2i.Datapoint, error) {
	db, err := sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	datapoints := []b2i.Datapoint{}
	sqlStmt := `
		SELECT TIMESTAMP, RAW_INTENSITY, STEPS, RAW_KIND, HEART_RATE
		FROM MI_BAND_ACTIVITY_SAMPLE
		ORDER BY TIMESTAMP DESC;
	`
	err = db.Select(&datapoints, sqlStmt)
	if err != nil {
		return nil, err
	}

	return datapoints, nil
}

func sendDataPoints(dpts []b2i.Datapoint) error {
	jsonBody, _ := json.Marshal(dpts)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:8080/api/datapoints", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	return nil
}

func main() {
	fmt.Println("MiBand Gadgetbridge database importer")
	fmt.Println("=====================================")
	fmt.Println()

	flag.StringVar(&hostname, "h", "http://localhost:8080", "Endpoint of the server")
	flag.StringVar(&username, "u", "", "Username used for authentication")
	flag.StringVar(&password, "p", "", "Password used for authentication")
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Database file must be passed as an argument")
		return
	}

	databaseFile := os.Args[len(os.Args)-1]
	fmt.Printf("Opening file \"%s\"...\n", databaseFile)

	dpts, err := getDataPoints(databaseFile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d datapoints\n", len(dpts))

	err = sendDataPoints(dpts)
	if err != nil {
		fmt.Println("Could not write datapoints to InfluxDB:", err)
	}
}
