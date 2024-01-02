package main

import (
	"context"
	"encoding/json"
	"github.com/go-co-op/gocron/v2"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func startUpdateJob() {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal("could not create scheduler", err)
	}

	_, err = scheduler.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartImmediately()),
		gocron.NewTask(func() { updateEnergyPrices() }),
	)
	if err != nil {
		log.Fatal("could not create job", err)
	}

	scheduler.Start()
}

func updateEnergyPrices() {
	var err error

	retryCount := 60 * 24
	for i := 0; i < retryCount; i++ {
		responseData, err := fetchEnergyPrices()
		if err == nil {
			writeToDb(responseData)
			return
		}
		time.Sleep(time.Minute)

	}
	if err != nil {
		log.Println("Fetch request failed")
	}
}

func fetchEnergyPrices() (*Response, error) {
	response, err := http.Get(os.Getenv("PRICE_ENDPOINT"))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var parsedResponse Response

	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		return nil, err
	}

	return &parsedResponse, nil
}

func writeToDb(priceData *Response) {
	influxHost := os.Getenv("INFLUX_HOST")
	token := os.Getenv("INFLUX_TOKEN")
	bucket := os.Getenv("INFLUX_BUCKET")
	org := os.Getenv("INFLUX_ORG")

	influxClient := influxdb2.NewClient(influxHost, token)
	writeApi := influxClient.WriteAPIBlocking(org, bucket)

	for _, data := range priceData.Data {
		point := influxdb2.NewPoint(
			"energy-price",
			map[string]string{"unit": priceData.Unit},
			//todo add flat tax variable
			map[string]interface{}{"value": data.Value + 1.44},
			data.Date.Time,
		)

		err := writeApi.WritePoint(context.Background(), point)
		if err != nil {
			log.Println("error while writing to db: ", err)
		}
	}
	influxClient.Close()
	log.Println("Wrote %d points to db", len(priceData.Data))
}
