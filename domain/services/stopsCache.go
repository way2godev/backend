package services

import (
	"strings"
	"sync"
	"way2go/domain/entities"
	"way2go/infraestructure/database"

	"gorm.io/gorm"
)

type StopSearchService struct {
	// stopSearchSharededCache is a cache that stores the stops in a sharded way.
	stopSearchSharededCache [][]entities.Stop
}

// StopSearchServiceInstance is the instance of the StopSearchService.
var StopSearchServiceInstance StopSearchService = StopSearchService{
	stopSearchSharededCache: [][]entities.Stop{},
}

// Search searches for stops by name and city.
func (s *StopSearchService) Search(search string) []entities.Stop {
	results := make(chan []entities.Stop)
	defer close(results)

	var wg sync.WaitGroup

	for _, dataChunk := range s.stopSearchSharededCache {
		wg.Add(1)
		go func(chunk []entities.Stop) {
			for _, stop := range chunk {
				city := strings.ToLower(*stop.City)
				name := strings.ToLower(stop.Name)
				search = strings.ToLower(search)
				if strings.Contains(city, search) || strings.Contains(name, search) {
					results <- append(<-results, stop)
				}
			}
			wg.Done()
		}(dataChunk)
	} 

	if len(s.stopSearchSharededCache) == 0 {
		return []entities.Stop{}
	}

	wg.Wait()
	return <-results
}

// Update updates the search cache.
// It should be called every time a stop is created, updated, deleted.
// It is used by the jobs that update the cache.
func (s *StopSearchService) Update() {
	db := database.GetDB()

	// Get all stops
	var allStops []entities.Stop
	db.FindInBatches(&allStops, 100, func(tx *gorm.DB, batch int) error {
		return nil
	})
	chardedStops := [][]entities.Stop{}

	// Divide the stops in chunks of 100
	for i := 0; i < len(allStops); i += 100 {
		end := i + 100
		if end > len(allStops) {
			end = len(allStops)
		}
		chardedStops = append(chardedStops, allStops[i:end])
	}

	s.stopSearchSharededCache = chardedStops
}
