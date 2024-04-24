package agencies

import (
	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {
	agencies := r.PathPrefix("/agencies").Subrouter()
	
	agencies.HandleFunc("", indexHandler).Methods("GET")
	agencies.HandleFunc("", createHandler).Methods("POST")
}
