package entities

import "gorm.io/gorm"

// Line represents a line of a transport agency.
// Eg: Line 1 of Renfe, Line 2 of EMT, etc.
// GTFS: Routes.txt
type Line struct {
	gorm.Model
	Name        string     `json:"name"` // GTFS: route_long_name o route_short_name
	AgencyID    uint       `json:"agency_id"`
	Description *string    `json:"description"` // GTFS: route_desc
	RouteType   int        `json:"route_type"`  // GTFS: route_type
	Schedule    []Schedule `json:"schedules"`

	GtfsRouteId        string  `json:"gtfs_route_id"`
	GtfsRouteShortName *string `json:"gtfs_route_short_name"`
	GtfsRouteLongName  *string `json:"gtfs_route_long_name"`
	GtfsRouteUrl       *string `json:"gtfs_route_url"`
	GtfsRouteColor     *string `json:"gtfs_route_color"`
	GtfsRouteTextColor *string `json:"gtfs_route_text_color"`
}
