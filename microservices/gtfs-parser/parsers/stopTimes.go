package parsers

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
	"way2go/constants"
	"way2go/domain/entities"
	"way2go/infraestructure/database"
	"way2go/microservices/gtfs-parser/csv"
)

type gtfsStopTime struct {
	TripID            string
	ArrivalTime       string
	DepartureTime     string
	StopID            string
	StopSequence      string
	StopHeadsign      string
	ShapeDistTraveled string
}

func StopTimes() {
	stopTimes, err := csv.Read(fmt.Sprintf("%s/%s", constants.GTFS_PARSER_WORKDIR, constants.GTFS_STOP_TIMES_FILE))
	if err != nil {
		log.Fatalf("Error reading CSV: %v", err)
		return
	}

	log.Printf("Found in total %d stop times\n", len(stopTimes))
	var parsedStopTimes []gtfsStopTime
	for _, stopTime := range stopTimes {
		parsedStopTime := gtfsStopTime{
			TripID:            stopTime["trip_id"],
			ArrivalTime:       stopTime["arrival_time"],
			DepartureTime:     stopTime["departure_time"],
			StopID:            stopTime["stop_id"],
			StopSequence:      stopTime["stop_sequence"],
			StopHeadsign:      stopTime["stop_headsign"],
			ShapeDistTraveled: stopTime["shape_dist_traveled"],
		}
		parsedStopTimes = append(parsedStopTimes, parsedStopTime)
	}

	log.Println("Stop times parsed successfully")

	startTime := time.Now()
	var wg sync.WaitGroup
	chunkSize := len(parsedStopTimes) / constants.GTFS_PROCESSING_CHUNKS
	for chunk := 0; chunk < constants.GTFS_PROCESSING_CHUNKS; chunk++ {
		wg.Add(1)
		go func(chunk int) {
			defer wg.Done()
			startIndex := chunk * chunkSize
			endIndex := (chunk + 1) * chunkSize
			for i := startIndex; i < endIndex; i++ {
				parsedStopTimes[i].saveToDatabase()
			}
		}(chunk)
	}
	wg.Wait()
	log.Printf("Stop times saved to database in %v\n", time.Since(startTime))
}

func (st *gtfsStopTime) saveToDatabase() {
	db := database.GetDB()

	var schedule entities.Schedule
	db.Where("gtfs_trip_id = ?", st.TripID).First(&schedule)
	if schedule.ID == 0 {
		log.Printf("Schedule not found for trip_id %s\n", st.TripID)
		return
	}

	var stop entities.Stop
	db.Where("gtfs_stop_id = ?", st.StopID).First(&stop)
	if stop.ID == 0 {
		log.Printf("Stop not found for stop_id %s\n", st.StopID)
		return
	}

	stopSequence, err := strconv.Atoi(st.StopSequence)
	if err != nil {
		log.Printf("Error converting stop sequence to int: %v\n", err)
		return
	}

	var shapeDistTraveled *float64
	if st.ShapeDistTraveled == "" {
		shapeDistTraveled = nil
	} else {
		shapeDistTraveledValue, err := strconv.ParseFloat(st.ShapeDistTraveled, 64)
		if err != nil {
			log.Printf("Error converting shape_dist_traveled to float: %v\n", err)
			return
		}
		shapeDistTraveled = &shapeDistTraveledValue
	}

	scheduleStop := entities.ScheduleStop{
		ScheduleID:            schedule.ID,
		StopID:                stop.ID,
		StopSequence:          uint(stopSequence),
		ArrivalTime:           st.ArrivalTime,
		DepartureTime:         st.DepartureTime,
		ShapeDistanceTraveled: shapeDistTraveled,
	}
	var existingScheduleStop entities.ScheduleStop
	db.Where("schedule_id = ? AND stop_id = ? AND stop_sequence = ?", scheduleStop.ScheduleID, scheduleStop.StopID, scheduleStop.StopSequence).First(&existingScheduleStop)
	if existingScheduleStop.ID != 0 {
		scheduleStop.ID = existingScheduleStop.ID
	}
	db.Save(&scheduleStop)
}
