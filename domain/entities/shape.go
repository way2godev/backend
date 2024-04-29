package entities

import "gorm.io/gorm"

type Shape struct {
	gorm.Model
	ScheduleID uint `json:"schedule_id"`
	Elements   []ShapeElement
}

type ShapeElement struct {
	gorm.Model
	ShapeID          uint     `json:"shape_id"`
	Latitude         float64  `json:"latitude"`
	Longitude        float64  `json:"longitude"`
	Sequence         uint     `json:"sequence"`
	DistanceTraveled *float64 `json:"distance_traveled"`
}
