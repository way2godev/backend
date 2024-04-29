package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"way2go/api/agencies"
	"way2go/api/general"
	"way2go/api/middleware"
	"way2go/api/stops"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type server struct {
	router *mux.Router
}

var Server server

func (s *server) Run(started chan bool) {
	log.Print("Starting API server")
	s.router = mux.NewRouter()

	r := s.router
	r.Use(middleware.JsonMiddleware)

	v1 := r.PathPrefix("/api/v1").Subrouter()

	general.SetupRoutes(v1)
	agencies.SetupRoutes(v1)
	stops.SetupRoutes(v1)

	s.printRoutes()
	s.startHttpServer(started)
}

func (s *server) printRoutes() {
	fmt.Printf("\nAvailable routes:\n")
	s.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			log.Panic(err)
		}
		methods, _ := route.GetMethods()
		if len(methods) > 0 {
			fmt.Printf("%s (%s)\n", t, methods)
		}
		return nil
	})
	fmt.Printf("\n")
}

func (s *server) setupCors() http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	return c.Handler(s.router)
}

func (s *server) startHttpServer(started chan bool) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3050"
	}

	handler := s.setupCors()

	log.Printf("Starting server on port %s", port)
	started <- true
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return
	}
}
