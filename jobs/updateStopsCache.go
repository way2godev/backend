package jobs

import (
	"log"
	"time"
	"way2go/domain/services"
)

var ForceUpdateStopsCache chan bool = make(chan bool)

func UpdateStopsCache() {
	log.Printf("UpdateStopsCache job initialized")

	for {
		services.StopSearchServiceInstance.Update()
		log.Printf("UpdateStopsCache job: update finished")

		select {
		case <-ForceUpdateStopsCache:
			log.Printf("UpdateStopsCache job: forced update")
		case <-time.After(24 * time.Hour):
			log.Printf("UpdateStopsCache job: regular update")
		}
	}
}