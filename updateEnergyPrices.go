package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"io"
	"log"
	"net/http"
	"os"
)

func startUpdateJob() {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal("could not create scheduler", err)
	}

	_, err = scheduler.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartImmediately()),
		gocron.NewTask(updateEnergyPrices),
	)
	if err != nil {
		log.Fatal("could not create job", err)
	}

	scheduler.Start()
}

func updateEnergyPrices() {
	//todo reschedule retry
	responseData, err := fetchEnergyPrices()
	if err != nil {
		log.Fatal("error while fetching API: ", err)
		return
	}

	fmt.Println(responseData.Tariff)
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
