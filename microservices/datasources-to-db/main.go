package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"way2go/bootstrap"
	"way2go/constants"
	"way2go/domain/entities"
	"way2go/infraestructure/database"
)

func main() {
	bootstrap.Init()

	// Open the constants.SOURCE_FILE and read the csv
	file, err := os.Open(constants.GTFS_SOURCES_FILE)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	log.Printf("Found in total %d sources\n", len(records))

	db := database.GetDB()
	for i, record := range records {
		if i == 0 {
			// Skip the header
			continue
		}

		// Check if the record already exists in the database
		var existingDatasource entities.Datasource
		db.Where("url = ?", record[3]).First(&existingDatasource)
		if existingDatasource.ID != 0 {
			log.Printf("Datasource %s already exists in the database\n", record[0])
			continue
		} else {
			// Create a new datasource
			// Feed name, Provider, Location, Feed Url, expired, download_url
			expired, err := strconv.ParseBool(record[4])
			if err != nil {
				log.Fatalf("Error parsing expired field: %v, value: %s", err, record[4])
				continue
			}
			db.Save(&entities.Datasource{
				Name:     record[0],
				Provider: record[1],
				Location: record[2],
				Url:      record[5],
				Expired:  expired,
				Comments: nil,
			})
		}	
	}

	log.Println("Job done!")
}


