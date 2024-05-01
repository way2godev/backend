package parsers

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
	"way2go/domain/entities"
	"way2go/infraestructure/database"
	"way2go/microservices/gtfs-parser/constants"
	"way2go/microservices/gtfs-parser/csv"
)

type gtfsStop struct {
	StopID             string
	StopCode           string
	StopName           string
	StopDesc           string
	StopLat            string
	StopLon            string
	ZoneID             string
	StopUrl            string
	LocationType       string
	ParentStation      string
	StopTimezone       string
	WheelchairBoarding string
}

func Stops() {
	stops, err := csv.Read(fmt.Sprintf("%s/%s", constants.WORKDIR, constants.STOPS_FILE))
	if err != nil {
		log.Fatalf("Error reading CSV: %v", err)
		return
	}

	log.Printf("Found in total %d stops\n", len(stops))
	var parsedStops []gtfsStop
	for _, stop := range stops {
		parsedStop := gtfsStop{
			StopID:             stop["stop_id"],
			StopCode:           stop["stop_code"],
			StopName:           stop["stop_name"],
			StopDesc:           stop["stop_desc"],
			StopLat:            stop["stop_lat"],
			StopLon:            stop["stop_lon"],
			ZoneID:             stop["zone_id"],
			StopUrl:            stop["stop_url"],
			LocationType:       stop["location_type"],
			ParentStation:      stop["parent_station"],
			StopTimezone:       stop["stop_timezone"],
			WheelchairBoarding: stop["wheelchair_boarding"],
		}
		parsedStops = append(parsedStops, parsedStop)
	}
	log.Println("Stops parsed successfully")

	startTime := time.Now()
	var wg sync.WaitGroup
	chunkSize := len(parsedStops) / constants.PROCESSING_CHUNKS
	for chunk := 0; chunk < constants.PROCESSING_CHUNKS; chunk++ {
		wg.Add(1)
		go func(chunkIndex int) {
			defer wg.Done()
			startIndex := chunkIndex * chunkSize
			endIndex := (chunkIndex + 1) * chunkSize
			for i := startIndex; i < endIndex; i++ {
				parsedStops[i].saveToDatabase()
			}
		}(chunk)
	}
	wg.Wait()
	log.Printf("Stops saved to database in %s\n", time.Since(startTime))
}

func (s *gtfsStop) saveToDatabase() {
	lat, err := strconv.ParseFloat(s.StopLat, 64)
	if err != nil {
		log.Printf("Error parsing latitude on stop %s: %v (value: %s)\n", s.StopName, err, s.StopLat)
		return
	}
	lon, err := strconv.ParseFloat(s.StopLon, 64)
	if err != nil {
		log.Printf("Error parsing longitude on stop %s: %v (value: %s)\n", s.StopName, err, s.StopLon)
		return
	}

	var wheelchairBoarding *bool // If "1" then true, if "2" then false
	if s.WheelchairBoarding == "1" {
		wheelchairBoarding = new(bool)
		*wheelchairBoarding = true
	} else if s.WheelchairBoarding == "2" {
		wheelchairBoarding = new(bool)
		*wheelchairBoarding = false
	}

	var locType *int
	if s.LocationType != "" {
		locTypeInt, err := strconv.Atoi(s.LocationType)
		if err != nil {
			log.Printf("Error parsing location type on stop %s: %v (value: %s)\n", s.StopName, err, s.LocationType)
			return
		}
		locType = &locTypeInt
	}

	var timeZone *string // This will be nil if the timezone is empty
	if s.StopTimezone != "" {
		timeZone = &s.StopTimezone
	}

	stop := entities.Stop{
		Name:               s.StopName,
		Description:        &s.StopDesc,
		Latitude:           lat,
		Longitude:          lon,
		WheelchairBoarding: wheelchairBoarding,
		GtfsStopId:         s.StopID,
		GtfsStopCode:       &s.StopCode,
		GtfsLocationType:   locType,
		GtfsStopTimezone:   timeZone,
	}

	// Check if the stop already exists
	db := database.GetDB()
	var existingStop entities.Stop
	db.Where("gtfs_stop_id = ?", s.StopID).First(&existingStop)

	if existingStop.ID != 0 { // Stop already exists
		log.Printf("Stop %s already exists\n", s.StopName)
		return
	} else {
		db.Create(&stop)
		log.Printf("Stop %s created\n", s.StopName)
	}
}
