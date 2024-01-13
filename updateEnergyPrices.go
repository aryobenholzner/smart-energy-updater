package main

import (
	"context"
	"encoding/json"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func updateEnergyPrices() {
	log.Println("Price update started")
	var err error

	retryCount := 5
	for i := 0; i < retryCount; i++ {
		responseData, err := fetchEnergyPrices()
		if err == nil {
			writePriceDataToDb(responseData)
			return
		}
		time.Sleep(time.Minute)

	}
	if err != nil {
		log.Println("Fetching energy prices failed")
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

func writePriceDataToDb(priceData *Response) {
	influxHost := os.Getenv("INFLUX_HOST")
	token := os.Getenv("INFLUX_TOKEN")
	bucket := os.Getenv("INFLUX_BUCKET")
	org := os.Getenv("INFLUX_ORG")
	targetMeasurement := os.Getenv("INFLUX_PRICE_MEASUREMENT")

	influxClient := influxdb2.NewClient(influxHost, token)
	writeApi := influxClient.WriteAPIBlocking(org, bucket)

	flatFeeEnv, flatFeeEnvExists := os.LookupEnv("FLAT_FEE")
	var flatFee float64 = 0
	var err error
	if flatFeeEnvExists {
		flatFee, err = strconv.ParseFloat(flatFeeEnv, 32)
		if err != nil {
			log.Println("flat fee env could not be parsed to float: ", err)
		}
	}

	for _, data := range priceData.Data {
		point := influxdb2.NewPoint(
			targetMeasurement,
			map[string]string{"unit": priceData.Unit},
			map[string]interface{}{"value": data.Value + flatFee},
			data.Date.Time,
		)

		err := writeApi.WritePoint(context.Background(), point)
		if err != nil {
			log.Println("error while writing to db: ", err)
		}
	}
	influxClient.Close()
	log.Printf("Wrote %d price points to db\n", len(priceData.Data))
}
