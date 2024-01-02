package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	startUpdateJob()

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
