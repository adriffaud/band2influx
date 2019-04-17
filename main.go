package main

import (
	"fmt"
	"log"
	"os"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const influxURL = "http://localhost:8086"
const database = "miband"

// Datapoint is sensor data at a given point in time
type Datapoint struct {
	Timestamp    int64 `db:"TIMESTAMP"`
	RawIntensity int   `db:"RAW_INTENSITY"`
	Steps        int   `db:"STEPS"`
	RawKind      int   `db:"RAW_KIND"`
	HeartRate    int   `db:"HEART_RATE"`
}

func getDataPoints(dbFile string) ([]Datapoint, error) {
	db, err := sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	datapoints := []Datapoint{}
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

func main() {
	log.Println("MiBand Gadgetbridge database importer")

	if len(os.Args) < 2 {
		log.Fatal("Database file must be passed as an argument")
	}

	databaseFile := os.Args[1]
	log.Printf("Opening %s...\n", databaseFile)

	dpts, err := getDataPoints(databaseFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%d datapoints\n", len(dpts))

	c, err := client.NewHTTPClient(client.HTTPConfig{Addr: "http://localhost:8086"})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
	}
	defer c.Close()

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  database,
		Precision: "s",
	})

	log.Println("Preparing data for addition into InfluxDB")
	for _, dp := range dpts {
		fields := map[string]interface{}{
			"raw-intensity": dp.RawIntensity,
			"raw-kind":      dp.RawKind,
			"steps":         dp.Steps,
		}

		if dp.HeartRate < 250 && dp.HeartRate > 0 {
			fields["heart-rate"] = dp.HeartRate
		}

		pt, err := client.NewPoint("activity", nil, fields, time.Unix(dp.Timestamp, 0))
		if err != nil {
			log.Println("Error:", err.Error())
		}
		bp.AddPoint(pt)
	}

	log.Println("Writing into InfluxDB")
	err = c.Write(bp)
	if err != nil {
		log.Fatal(err.Error())
	}
}
