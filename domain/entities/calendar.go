package entities

import "gorm.io/gorm"

// Calendar represents the days of the week in which a schedule is active.
// GTFS: Calendar.txt
type Calendar struct {
	gorm.Model
	ScheduleID uint   `json:"schedule_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Monday     bool   `json:"monday"`
	Tuesday    bool   `json:"tuesday"`
	Wednesday  bool   `json:"wednesday"`
	Thursday   bool   `json:"thursday"`
	Friday     bool   `json:"friday"`
	Saturday   bool   `json:"saturday"`
	Sunday     bool   `json:"sunday"`
}

// CalendarException represents an exception to the regular schedule.
// GTFS: CalendarDates.txt
type CalendarException struct {
	gorm.Model
	CalendarID    uint   `json:"calendar_id"`
	Date          string `json:"date"`
	ExceptionType int    `json:"exception_type"` // 1: Added, 2: Removed
}
