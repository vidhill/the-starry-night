package main

import (
	"net/http"

	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/handlers"
	"github.com/vidhill/the-starry-night/service"
)

func main() {

	// wiring
	configService := service.NewConfigService(domain.NewViperConfig())
	loggerService := service.NewLoggerService(domain.NewStandardLogger())
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
	port := configService.GetString("SERVER_PORT")

	loggerService.Info("listening on port", port)

	err := http.ListenAndServe(":"+port, mux)

	if err != nil {
		loggerService.Error("Error starting server", err.Error())
	}
}
