package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	b2i "driffaud.fr/adrien/band2influx"
	"github.com/gorilla/mux"
)

var port int

var influxEndpoint string
var database string

func datapointHandler(w http.ResponseWriter, req *http.Request) {
	var datapoints []b2i.Datapoint
	err := json.NewDecoder(req.Body).Decode(&datapoints)
	if err != nil {
		fmt.Println("Error decoding json body", err)
	}

	fmt.Printf("%d datapoints received\n", len(datapoints))

	jsonBody, _ := json.Marshal(datapoints)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBody)
}

func main() {
	fmt.Println("Band2Influx server")
	fmt.Println("==================")
	fmt.Println()

	flag.IntVar(&port, "p", 8080, "the server port")
	flag.StringVar(&influxEndpoint, "influxEndpoint", "http://localhost:8086", "the influxdb host")
	flag.StringVar(&database, "db", "", "the influx database")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/api/datapoints", b2i.BasicAuth(datapointHandler, "adrien", "adrien", "toto")).Methods("POST")
	http.Handle("/", r)

	// Server configuration
	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		fmt.Printf("Server listening on http://localhost:%d\n", port)
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
	fmt.Println("shutting down")
	os.Exit(0)

	// c, err := client.NewHTTPClient(client.HTTPConfig{Addr: "http://localhost:8086"})
	// if err != nil {
	// 	fmt.Println("Error creating InfluxDB Client: ", err.Error())
	// }
	// defer c.Close()

	// bp, _ := client.NewBatchPoints(client.BatchPointsConfig{Database: database, Precision: "s"})

	// log.Println("Preparing data for addition into InfluxDB")
	// for _, dp := range dpts {
	// 	fields := map[string]interface{}{
	// 		"raw-intensity": dp.RawIntensity,
	// 		"raw-kind":      dp.RawKind,
	// 		"steps":         dp.Steps,
	// 	}

	// 	if dp.HeartRate < 250 && dp.HeartRate > 0 {
	// 		fields["heart-rate"] = dp.HeartRate
	// 	}

	// 	pt, err := client.NewPoint("activity", nil, fields, time.Unix(dp.Timestamp, 0))
	// 	if err != nil {
	// 		log.Println("Error:", err.Error())
	// 	}
	// 	bp.AddPoint(pt)
	// }

	// log.Println("Writing into InfluxDB")
	// err = c.Write(bp)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
}
