package parsers

import (
	"fmt"
	"log"
	"sync"
	"time"
	"way2go/domain/entities"
	"way2go/infraestructure/database"
	"way2go/microservices/gtfs-parser/constants"
	"way2go/microservices/gtfs-parser/csv"
)

type gtfsTrip struct {
	RouteID             string
	ServiceID           string
	TripID              string
	TripHeadsign        string
	TripShortName       string
	DirectionID         string
	BlockID             string
	ShapeID             string
	WheelchairAccesible string
	BikesAllowed        string
}

func Trips() {
	trips, err := csv.Read(fmt.Sprintf("%s/%s", constants.WORKDIR, constants.TRIPS_FILE))
	if err != nil {
		log.Fatalf("Error reading CSV: %v", err)
		return
	}

	log.Printf("Found in total %d trips\n", len(trips))
	var parsedTrips []gtfsTrip
	for _, trip := range trips {
		parsedTrip := gtfsTrip{
			RouteID:             trip["route_id"],
			ServiceID:           trip["service_id"],
			TripID:              trip["trip_id"],
			TripHeadsign:        trip["trip_headsign"],
			TripShortName:       trip["trip_short_name"],
			DirectionID:         trip["direction_id"],
			BlockID:             trip["block_id"],
			ShapeID:             trip["shape_id"],
			WheelchairAccesible: trip["wheelchair_accessible"],
			BikesAllowed:        trip["bikes_allowed"],
		}
		parsedTrips = append(parsedTrips, parsedTrip)
	}
	log.Println("Trips parsed successfully")

	startTime := time.Now()
	var wg sync.WaitGroup
	chunkSize := len(parsedTrips) / constants.PROCESSING_CHUNKS
	for i := 0; i < constants.PROCESSING_CHUNKS; i++ {
		wg.Add(1)
		go func(chunkIndex int) {
			defer wg.Done()
			startIndex := chunkIndex * chunkSize
			endIndex := (chunkIndex + 1) * chunkSize
			for i := startIndex; i < endIndex; i++ {
				parsedTrips[i].saveToDatabase()
			}
		}(i)
	}
	wg.Wait()
	log.Printf("Trips saved to database in %v\n", time.Since(startTime))
}

func (t *gtfsTrip) saveToDatabase() {
	db := database.GetDB()
	var line entities.Line
	db.Where("gtfs_route_id = ?", t.RouteID).First(&line)
	if line.ID == 0 {
		log.Printf("Line with GTFS ID %s not found\n", t.RouteID)
		return
	}
	var name string
	if t.TripHeadsign != "" {
		name = t.TripHeadsign
	} else {
		name = t.TripShortName
	}

	var bikesAllowed *bool // If "1" true, if "2" false
	if t.BikesAllowed == "1" {
		bikesAllowed = new(bool)
		*bikesAllowed = true
	} else if t.BikesAllowed == "2" {
		bikesAllowed = new(bool)
		*bikesAllowed = false
	}

	trip := entities.Schedule{
		LineID:            line.ID,
		Name:              name,
		GtfsServiceId:     t.ServiceID,
		GtfsTripId:        t.TripID,
		GtfsTripShortName: &t.TripShortName,
		GtfsHeadsign:      &t.TripHeadsign,
		GtfsBikesAllowed:  bikesAllowed,
	}
	var existing entities.Schedule
	db.Where("gtfs_trip_id = ?", t.TripID).First(&existing)
	if existing.ID != 0 {
		trip.ID = existing.ID
	}
	db.Save(&trip)
}
