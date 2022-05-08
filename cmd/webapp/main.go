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
	httpService := service.NewHttpService(domain.NewDefaultHttpClient(loggerService))

	ISSRepository := domain.NewISSRepositoryRest(configService, httpService, loggerService)

	ISSService := service.NewISSLocationService(ISSRepository)

	dh := handlers.NewHandlers(ISSService)

	mux := http.NewServeMux()

	//
	// route handlers
	//
	// health endpoint for kubernetes liveness probe
	mux.HandleFunc("/health", dh.Health)

	mux.HandleFunc("/iss-position", dh.ISSPosition)

	// start server
	port := configService.GetString("SERVER_PORT")

	loggerService.Info("listening on port", port)

	err := http.ListenAndServe(":"+port, mux)

	if err != nil {
		loggerService.Error("Error starting server", err.Error())
	}
}
