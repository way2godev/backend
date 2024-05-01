package parsers

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
	"way2go/domain/entities"
	"way2go/infraestructure/database"
	"way2go/microservices/gtfs-parser/constants"
	"way2go/microservices/gtfs-parser/csv"
)

type gtfsRoute struct {
	RouteID        string
	AgencyID       string
	RouteShortName string
	RouteLongName  string
	RouteDesc      string
	RouteType      string
	RouteUrl       string
	RouteColor     string
	RouteTextColor string
	RouteSortOrder string
}

func Routes() {
	routes, err := csv.Read(fmt.Sprintf("%s/%s", constants.WORKDIR, constants.ROUTES_FILE))
	if err != nil {
		log.Fatalf("Error reading CSV: %v", err)
		return
	}

	log.Printf("Found in total %d routes\n", len(routes))
	var parsedRoutes []gtfsRoute
	for _, route := range routes {
		parsedRoute := gtfsRoute{
			RouteID:        route["route_id"],
			AgencyID:       route["agency_id"],
			RouteShortName: route["route_short_name"],
			RouteLongName:  route["route_long_name"],
			RouteDesc:      route["route_desc"],
			RouteType:      route["route_type"],
			RouteUrl:       route["route_url"],
			RouteColor:     route["route_color"],
			RouteTextColor: route["route_text_color"],
			RouteSortOrder: route["route_sort_order"],
		}
		parsedRoutes = append(parsedRoutes, parsedRoute)
	}
	log.Println("Routes parsed successfully")

	startTime := time.Now()
	var wg sync.WaitGroup
	chunkSize := len(parsedRoutes) / constants.PROCESSING_CHUNKS
	for chunk := 0; chunk < constants.PROCESSING_CHUNKS; chunk++ {
		wg.Add(1)
		go func(chunkIndex int) {
			defer wg.Done()
			startIndex := chunkIndex * chunkSize
			endIndex := (chunkIndex + 1) * chunkSize
			if endIndex > len(parsedRoutes) {
				endIndex = len(parsedRoutes)
			}
			for i := startIndex; i < endIndex; i++ {
				parsedRoutes[i].saveToDatabase()
			}
		}(chunk)
	}
	wg.Wait()
	log.Printf("Routes saved to database in %v\n", time.Since(startTime))

}

func (r *gtfsRoute) saveToDatabase() {
	db := database.GetDB()

	var agency entities.Agency
	err := db.Where("gtfs_agency_id = ?", r.AgencyID).First(&agency).Error
	if err != nil {
		log.Fatalf("Error getting agency: %v", err)
	}

	routeType, err := strconv.Atoi(r.RouteType)
	if err != nil {
		log.Printf("Error parsing route type: %v (value: %s)\n", err, r.RouteType)
		return
	}

	var name string
	if r.RouteLongName != "" {
		name = r.RouteLongName
	} else {
		name = r.RouteShortName
	}

	route := entities.Line{
		Name:               name,
		AgencyID:           agency.ID,
		Description:        &r.RouteDesc,
		RouteType:          routeType,
		GtfsRouteId:        r.RouteID,
		GtfsRouteShortName: &r.RouteShortName,
		GtfsRouteLongName:  &r.RouteLongName,
		GtfsRouteUrl:       &r.RouteUrl,
		GtfsRouteColor:     &r.RouteColor,
		GtfsRouteTextColor: &r.RouteTextColor,
	}

	// Check if the route already exists
	var existingRoute entities.Line
	db.Where("gtfs_route_id = ?", r.RouteID).First(&existingRoute)
	if existingRoute.ID != 0 {
		log.Printf("Route %s already exists\n", r.RouteID)
		return
	} else {
		db.Save(&route)
		log.Printf("Route %s saved\n", r.RouteID)
	}
}
