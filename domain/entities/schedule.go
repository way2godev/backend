package entities

import (
	"time"

	"gorm.io/gorm"
)

// Schedule represents a schedule for a line variant.
// Eg: The one that departs at 8:00am, at 9:00am, etc. those are different schedules.
type Schedule struct {
	gorm.Model
	LineVariantID uint   `json:"line_variant_id"`
	Days          string `json:"days" gorm:"type:varchar(7);default:'1111111'"`

	ScheduleStops []ScheduleStop `json:"schedule_stops"`
}

// ScheduleStop represents a stop in a schedule.
// Eg: Inside a schedule, the first stop is at 8:00am, the second at 8:10am, etc.
type ScheduleStop struct {
	gorm.Model
	ScheduleID uint      `json:"schedule_id"`
	StopID     uint      `json:"stop_id"`
	Order      uint      `json:"order"`
	Time       time.Time `json:"time"`
}
