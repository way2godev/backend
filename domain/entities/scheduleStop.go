package entities

import (
	"time"

	"gorm.io/gorm"
)

// ScheduleStop represents a stop in a schedule.
// Eg: Inside a schedule, the first stop is at 8:00am, the second at 8:10am, etc.
type ScheduleStop struct {
	gorm.Model
	ScheduleID            uint      `json:"schedule_id"`
	StopID                uint      `json:"stop_id"`
	StopSequence          uint      `json:"stop_sequence"`
	ArrivalTime           time.Time `json:"arrival_time" gorm:"type:time without time zone"`
	DepartureTime         time.Time `json:"departure_time" gorm:"type:time without time zone"`
	ShapeDistanceTraveled *float64  `json:"shape_distance_traveled"`
}
