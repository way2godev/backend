package entities

import (
	"gorm.io/gorm"
)

// Schedule represents a schedule for a line variant.
// Eg: The one that departs at 8:00am, at 9:00am, etc. those are different schedules.
// GTFS: Trips.txt
type Schedule struct {
	gorm.Model
	LineID uint    `json:"line_id"`
	Name   *string `json:"name"` // GTFS: trip_headsign

	ScheduleStops []ScheduleStop `json:"schedule_stops"`
	Shape         Shape          `json:"shape"`
	Calendar      Calendar       `json:"calendar"`

	GtfsServiceId     string  `json:"gtfs_service_id"`
	GtfsTripId        string  `json:"gtfs_trip_id"`
	GtfsTripShortName *string `json:"gtfs_trip_short_name"`
	GtfsBikesAllowed  *int    `json:"gtfs_bikes_allowed"`
}
