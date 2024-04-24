package entities

import (
	"gorm.io/gorm"
)

// Stop represents a stop in a line. It can be a bus stop, a train station, etc.
type Stop struct {
	gorm.Model
	Name              string  `json:"name"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	RawRenfeCode      *string `json:"raw_renfe_code"`
	Adress            *string `json:"adress"`
	City              *string `json:"city"`
	PostalCode        *string `json:"postal_code"`
	Province          *string `json:"province"`
	Country           *string `json:"country"`
	WeelchairBoarding *bool   `json:"weelchair_boarding"`
}
