package main

import (
	"log"
	"sync"
	"time"
	"way2go/api"
	"way2go/bootstrap"
)

func main() {
	bootstrap.Init()

	var wg sync.WaitGroup
	wg.Add(1)

	// Start the server
	startTime := time.Now()
	serverStartedSignal := make(chan bool)
	go api.Server.Run(serverStartedSignal)
	<-serverStartedSignal
	log.Printf("Server started in %v\n", time.Since(startTime))

	wg.Wait()
}