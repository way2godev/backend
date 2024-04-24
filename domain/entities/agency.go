package entities

import (
	"gorm.io/gorm"
)

// Agency represents a transport agency.
// Eg: Renfe, EMT, etc.
type Agency struct {
	gorm.Model
	Name        string  `json:"name"`
	RawId       *string `json:"raw_id"`
	Url         *string `json:"url"`
	PhoneNumber *string `json:"phone_number"`
	Line        []Line  `json:"lines"`
}
