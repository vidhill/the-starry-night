package main

import (
	"net/http"

	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/handlers"
	rest_api_repository "github.com/vidhill/the-starry-night/restapirepository"
	"github.com/vidhill/the-starry-night/service"
)

func main() {

	// wiring
	configService := service.NewConfigService(domain.NewViperConfig())
	loggerService := service.NewLoggerService(domain.NewStandardLogger())
	httpService := service.NewHttpService(domain.NewDefaultHttpClient(loggerService))

	ISSRepository := rest_api_repository.NewISSRepositoryRest(configService, httpService, loggerService)

	ISSService := service.NewISSLocationService(ISSRepository)

	dh := handlers.NewHandlers(loggerService, ISSService)

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
