package entities

import (
	"way2go/infraestructure/database"

	"gorm.io/gorm"
)

// Stop represents a stop in a line. It can be a bus stop, a train station, etc.
type Stop struct {
	gorm.Model
	Name               string  `json:"name"`                // GTFS: stop_name
	Description        *string `json:"description"`         // GTFS: stop_desc
	Latitude           float64 `json:"latitude"`            // GTFS: stop_lat
	Longitude          float64 `json:"longitude"`           // GTFS: stop_lon
	WheelchairBoarding *bool   `json:"wheelchair_boarding"` // GTFS: wheelchair_boarding

	StopTimes []ScheduleStop `json:"stop_times"`

	GtfsStopId       string  `json:"gtfs_stop_id"`
	GtfsStopCode     *string `json:"gtfs_stop_code"`
	GtfsLocationType *int    `json:"gtfs_location_type"`
	GtfsStopTimezone *string `json:"gtfs_stop_timezone"`
}

func GetStops(page int) (stops []Stop, total int64, err error) {
	db := database.GetDB()
	err = db.Offset((page - 1) * 10).Limit(10).Order("name").Find(&stops).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Model(&Stop{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return stops, total, nil
}