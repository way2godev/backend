package main

import (
	"log"
	"os"
	"sync"
	"time"
	"way2go/api"
	"way2go/infraestructure/database"
	"way2go/jobs"

	"way2go/domain/services"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// Initialize the database
	if !database.InitDB() {
		log.Fatalf("Failed to initialize the database")
	}

	// Migrate the models if the --migrate flag is set on argv or the app is in production
	if len(os.Args) > 1 && os.Args[1] == "--migrate" || os.Getenv("APP_ENV") == "production" {
		services.MigrateModelsServiceInstance.MigrateModels()
	} else if len(os.Args) > 1 && os.Args[1] == "--drop" {
		services.MigrateModelsServiceInstance.DropModels()
		services.MigrateModelsServiceInstance.MigrateModels()
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// Start the server
	startTime := time.Now()
	serverStartedSignal := make(chan bool)
	go api.Server.Run(serverStartedSignal)
	<-serverStartedSignal
	log.Printf("Server started in %v\n", time.Since(startTime))

	// Start the jobs
	startTime = time.Now()
	go jobs.RunJobs(serverStartedSignal)
	<-serverStartedSignal
	log.Printf("Jobs started in %v\n", time.Since(startTime))

	wg.Wait()
}