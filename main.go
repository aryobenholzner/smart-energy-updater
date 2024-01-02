package main

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("hello world")

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		fmt.Println("could not create scheduler", err)
		return
	}

	_, err = scheduler.NewJob(
		gocron.DurationJob(2*time.Second),
		gocron.NewTask(func() {
			fmt.Println("hello from the job")
		}),
	)
	if err != nil {
		fmt.Println("could not create job", err)
		return
	}

	scheduler.Start()

	go forever()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	//time for cleanup before exit
	fmt.Println("Adios!")
}

func forever() {
	for {
		time.Sleep(time.Second)
	}
}
