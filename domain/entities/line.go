package entities

import "gorm.io/gorm"

// Line represents a line of a transport agency.
// Eg: Line 1 of Renfe, Line 2 of EMT, etc.
type Line struct {
	gorm.Model
	Name         string        `json:"name"`
	AgencyID     uint          `json:"agency_id"`
	RawRouteId   *string       `json:"raw_route_id"`
	LineVariants []LineVariant `json:"line_variants"`
}

// LineVariant represents a variant of a line.
// Eg: Line 1 of Renfe has a variant that goes from Madrid to Figueres,
// but there is a variant that does not stop on Calatayud or Lleida, etc
type LineVariant struct {
	gorm.Model
	Tag        string `json:"tag"`
	LineID     uint   `json:"line_id"`
	OriginStop Stop   `json:"origin_stop" gorm:"foreignKey:ID"`
	FinalStop  Stop   `json:"final_stop" gorm:"foreignKey:ID"`
	
	// The default variant is the one that includes all stops and is from origin to final stop.
	IsDefault  bool   `json:"is_default" gorm:"default:false"` 
}
