package general

import (
	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {
	general := r.PathPrefix("/general").Subrouter()
	
	general.HandleFunc("/ping", handlePing).Methods("GET")
}
