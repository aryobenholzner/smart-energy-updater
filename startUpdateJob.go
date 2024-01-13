package main

import (
	"github.com/go-co-op/gocron/v2"
	"log"
	"os"
)

func startUpdateJob() {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal("could not create scheduler", err)
	}

	_, err = scheduler.NewJob(
		gocron.CronJob(os.Getenv("CRON_SCHEDULE_PRICE"), false),
		gocron.NewTask(updateEnergyPrices),
	)
	if err != nil {
		log.Fatal("could not create job", err)
	}
	log.Println("price update job created")

	_, err = scheduler.NewJob(
		gocron.CronJob(os.Getenv("CRON_SCHEDULE_CONSUMPTION"), false),
		gocron.NewTask(updateEnergyConsumption),
	)
	if err != nil {
		log.Fatal("could not create job", err)
	}
	log.Println("consumption update job created")

	scheduler.Start()
}
