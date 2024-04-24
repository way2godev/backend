package cache

import (
	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {
	cache := r.PathPrefix("/cache").Subrouter()
	
	cache.HandleFunc("/stops", handleUpdateStopsCache).Methods("POST")
}
