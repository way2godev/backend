package entities

import "gorm.io/gorm"

// Datastore represents a datasource of gtfs data.
type Datasource struct {
	gorm.Model
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Url         string  `json:"url"`
	Comments    *string `json:"other"`
}
