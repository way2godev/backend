package main

import (
	"log"
	"sync"
	"time"
	"way2go/api"
	"way2go/bootstrap"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	bootstrap.Init()

	var wg sync.WaitGroup
	wg.Add(1)

	startTime := time.Now()
	serverStartedSignal := make(chan bool)
	go api.Server.Run(serverStartedSignal)
	<-serverStartedSignal
	log.Printf("Server started in %v\n", time.Since(startTime))

	wg.Wait()
}
