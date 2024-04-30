package stops

import (
	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {
	stops := r.PathPrefix("/stops").Subrouter()

	stops.HandleFunc("", index).Methods("GET")
}
