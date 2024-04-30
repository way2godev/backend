package entities

import (
	"gorm.io/gorm"
)

// Agency represents a transport agency.
// Eg: Renfe, EMT, etc.
// GTFS: Agency.txt
type Agency struct {
	gorm.Model
	Name string `json:"name"` // GTFS: agency_name
	Line []Line `json:"lines"`

	GtfsAgencyId       string  `json:"gtfs_agency_id"`
	GtfsAgencyUrl      string  `json:"gtfs_agency_url"`
	GtfsAgencyTimezone string  `json:"gtfs_agency_timezone"`
	GtfsAgencyLang     *string `json:"gtfs_agency_lang"`
	GtfsAgencyPhone    *string `json:"gtfs_agency_phone"`
	GtfsAgencyEmail    *string `json:"gtfs_agency_email"`
}
