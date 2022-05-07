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

	//
	// route handlers
	//
	// health endpoint for kubernetes liveness probe
	http.HandleFunc("/health", dh.Health)

	http.HandleFunc("/hello", dh.Hello)
	http.HandleFunc("/foo", dh.GetFoo)

	// start server
	port := configService.GetString("PORT")

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Println("Error starting server", err.Error())
	}
}
