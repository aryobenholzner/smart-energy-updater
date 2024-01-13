package main

import (
	"github.com/go-co-op/gocron/v2"
	"log"
)

func startUpdateJob() {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal("could not create scheduler", err)
	}

	/*
		_, err = scheduler.NewJob(
			gocron.CronJob(os.Getenv("CRON_SCHEDULE"), false),
			gocron.NewTask(updateEnergyPrices),
		)
		if err != nil {
			log.Fatal("could not create job", err)
		}
	*/

	_, err = scheduler.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartImmediately()),
		gocron.NewTask(updateEnergyConsumption),
	)
	if err != nil {
		log.Fatal("could not create job", err)
	}

	scheduler.Start()
}
