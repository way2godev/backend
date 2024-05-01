package entities

import (
	"gorm.io/gorm"
	"way2go/infraestructure/database"
)

// Stop represents a stop in a line. It can be a bus stop, a train station, etc.
type Stop struct {
	gorm.Model
	Name               string  `json:"name"`
	Description        *string `json:"description"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	WheelchairBoarding *bool   `json:"wheelchair_boarding"`

	StopTimes []ScheduleStop `json:"stop_times"`

	GtfsStopId       string  `json:"gtfs_stop_id"`
	GtfsStopCode     *string `json:"gtfs_stop_code"`
	GtfsLocationType *int    `json:"gtfs_location_type"`
	GtfsStopTimezone *string `json:"gtfs_stop_timezone"`
}

// GetStopsPaginatedJoinLines retrieves a paginated list of stops with their associated lines.
// It takes a page number as input and returns the stops parsed as a slice of maps, the total number of stops, and any error encountered.
// For each stop, it retrieves the lines that pass through the stop and constructs a parsed representation of each line.
// The parsed stops and lines are stored in maps and appended to their respective slices.
// Finally, the function returns the parsed stops, the total number of stops, and any error encountered during the process.
func GetStopsPaginatedJoinLines(page int) (stopsParsed []map[string]interface{}, total int64, err error) {
	db := database.GetDB()
	var stops []Stop
	err = db.Model(&Stop{}).
		Limit(10).Offset((page - 1) * 10).
		Group("stops.id").
		Find(&stops).Error

	if err != nil {
		return nil, 0, err
	}
	err = db.Model(&Stop{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	for _, stop := range stops {
		lines, err := stop.getLinesThatPassThrough()
		if err != nil {
			return nil, 0, err
		}
		parsedLines := []map[string]interface{}{}
		for _, line := range lines {
			parsedLine := map[string]interface{}{
				"id":                    line.ID,
				"name":                  line.Name,
				"gtfs_route_short_name": line.GtfsRouteShortName,
				"gtfs_route_long_name":  line.GtfsRouteLongName,
				"gtfs_route_color":      line.GtfsRouteColor,
				"gtfs_route_text_color": line.GtfsRouteTextColor,
			}
			parsedLines = append(parsedLines, parsedLine)
		}
		stopParsed := map[string]interface{}{
			"id":                  stop.ID,
			"name":                stop.Name,
			"description":         stop.Description,
			"latitude":            stop.Latitude,
			"longitude":           stop.Longitude,
			"wheelchair_boarding": stop.WheelchairBoarding,
			"lines":               parsedLines,
		}

		stopsParsed = append(stopsParsed, stopParsed)
	}

	return stopsParsed, total, nil
}

// getLinesThatPassThrough returns a list of lines that pass through the stop.
func (s *Stop) getLinesThatPassThrough() (lines []Line, err error) {
	db := database.GetDB()
	err = db.Model(&Line{}).
		Select("lines.id, lines.name, lines.gtfs_route_short_name, lines.gtfs_route_long_name, lines.gtfs_route_color, lines.gtfs_route_text_color").
		Joins("JOIN schedules ON lines.id = schedules.line_id").
		Joins("JOIN schedule_stops ON schedules.id = schedule_stops.schedule_id").
		Where("schedule_stops.stop_id = ?", s.ID).
		Distinct().
		Find(&lines).Error
	if err != nil {
		return nil, err
	}
	return lines, nil
}
