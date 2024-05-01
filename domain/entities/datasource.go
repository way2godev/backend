package entities

import "gorm.io/gorm"

// Datastore represents a datasource of gtfs data.
type Datasource struct {
	gorm.Model
	Name     string  `json:"name"`
	Provider string  `json:"provider"`
	Location string  `json:"location"`
	Url      string  `json:"url"`
	Expired  bool    `json:"expired"`
	Comments *string `json:"other"`
}
