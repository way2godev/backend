package entities

import (
	"gorm.io/gorm"
)

// Stop represents a stop in a line. It can be a bus stop, a train station, etc.
type Stop struct {
	gorm.Model
	Name              string  `json:"name"`               // GTFS: stop_name
	Description       *string `json:"description"`        // GTFS: stop_desc
	Latitude          float64 `json:"latitude"`           // GTFS: stop_lat
	Longitude         float64 `json:"longitude"`          // GTFS: stop_lon
	WeelchairBoarding *bool   `json:"weelchair_boarding"` // GTFS: weelchair_boarding

	StopTimes []ScheduleStop `json:"stop_times"`

	GtfsStopId       string  `json:"gtfs_stop_id"`
	GtfsStopCode     *string `json:"gtfs_stop_code"`
	GtfsLocationType *int    `json:"gtfs_location_type"`
	GtfsStopTimezone *string `json:"gtfs_stop_timezone"`
}
