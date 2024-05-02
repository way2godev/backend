package parsers

import (
	"fmt"
	"log"
	"way2go/constants"
	"way2go/domain/entities"
	"way2go/infraestructure/database"
	"way2go/microservices/gtfs-parser/csv"
)

type gtfsAgency struct {
	AgencyID       string
	AgencyName     string
	AgencyUrl      string
	AgencyTimezone string
	AgencyLang     string
	AgencyPhone    string
	AgencyFareUrl  string
	AgencyEmail    string
}

func Agencies() {
	agencies, err := csv.Read(fmt.Sprintf("%s/%s", constants.GTFS_PARSER_WORKDIR, constants.GTFS_AGENCIES_FILE))
	if err != nil {
		log.Fatalf("Error reading CSV: %v", err)
		return
	}

	log.Printf("Found in total %d agencies\n", len(agencies))
	var parsedAgencies []gtfsAgency
	for _, agency := range agencies {
		parsedAgency := gtfsAgency{
			AgencyID:       agency["agency_id"],
			AgencyName:     agency["agency_name"],
			AgencyUrl:      agency["agency_url"],
			AgencyTimezone: agency["agency_timezone"],
			AgencyLang:     agency["agency_lang"],
			AgencyPhone:    agency["agency_phone"],
		}
		parsedAgencies = append(parsedAgencies, parsedAgency)
	}
	log.Println("Agencies parsed successfully")

	for _, a := range parsedAgencies {
		a.saveToDatabase()
	}
}

func (a *gtfsAgency) saveToDatabase() {
	agency := entities.Agency{
		Name:               a.AgencyName,
		GtfsAgencyId:       a.AgencyID,
		GtfsAgencyUrl:      a.AgencyUrl,
		GtfsAgencyTimezone: a.AgencyTimezone,
		GtfsAgencyLang:     &a.AgencyLang,
		GtfsAgencyPhone:    &a.AgencyPhone,
		GtfsAgencyEmail:    &a.AgencyEmail,
	}

	// Check if the agency already exists
	db := database.GetDB()
	var existingAgency entities.Agency
	db.Where("gtfs_agency_id = ?", a.AgencyID).First(&existingAgency)

	if existingAgency.ID != 0 { // Agency already exists
		log.Printf("Agency %s already exists\n", a.AgencyName)
		return
	} else {
		log.Printf("Agency %s doesn't exist, creating it\n", a.AgencyName)
		db.Save(&agency)
	}
}
