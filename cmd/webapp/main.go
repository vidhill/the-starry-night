package main

import (
	"log"
	"net/http"

	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/handlers"
	"github.com/vidhill/the-starry-night/service"
)

func main() {

	// wiring
	configService := service.NewConfigService(domain.NevViperConfig())
	dh := handlers.NewHandlers()

	mux := http.NewServeMux()

	//
	// route handlers
	//
	// health endpoint for kubernetes liveness probe
	mux.HandleFunc("/health", dh.Health)

	mux.HandleFunc("/hello", dh.Hello)
	mux.HandleFunc("/foo", dh.GetFoo)

	// start server
	port := configService.GetString("PORT")

	err := http.ListenAndServe(":"+port, mux)

	if err != nil {
		log.Println("Error starting server", err.Error())
	}
}
