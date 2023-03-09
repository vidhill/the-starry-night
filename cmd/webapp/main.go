package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/handlers"
	"github.com/vidhill/the-starry-night/middleware"
	"github.com/vidhill/the-starry-night/restapirepository"
	"github.com/vidhill/the-starry-night/service"
)

const SWAGGER_ROOT = "swagger-ui"

func main() {

	// wiring
	configService := service.NewViperConfig()
	loggerService := service.NewStandardLogger(configService)

	configErr := validateEnvVariables(configService)
	if configErr != nil {
		loggerService.Error(configErr.Error())
		os.Exit(1)
	}

	httpService := service.NewDefaultHttpClient(loggerService)

	issService := restapirepository.NewISSRepositoryRest(configService, httpService, loggerService)
	weatherService := restapirepository.NewWeatherbitRepository(configService, httpService, loggerService)

	issVisibleService := service.NewISSVisibleService(configService, loggerService, issService, weatherService)
	// custom middleware to log using the wrapped logger service
	customLogMiddleware := middleware.MakeMyLoggerMiddleware(loggerService)

	dh := handlers.NewHandlers(loggerService, issVisibleService)

	mux := chi.NewRouter()

	mux.Use(customLogMiddleware)

	// Serve swagger-ui static files
	mux.Handle(MakeSwaggerStaticServe(SWAGGER_ROOT))

	//
	// route handlers
	//
	//

	// swagger:route GET /health api
	//
	// Health endpoint for kubernetes liveness probe.
	//
	// Produces:
	// - text/plain
	//
	// Responses:
	// 			200: healthResponse
	mux.Get("/health", dh.Health)

	// swagger:route GET /iss-position api ISSRequest
	//
	// Determines if ISS is overhead.
	//
	// Responses:
	// 				200: ISSResult
	// 				400: ErrorResponse
	// 				500: ErrorResponse
	mux.Get("/iss-position", handlers.ComposeHandlers(handlers.AddJsonHeader)(dh.ISSPosition))

	port := configService.GetString("SERVER_PORT")

	loggerService.Info("listening on port", port)

	// start server
	err := http.ListenAndServe(":"+port, mux)

	if err != nil {
		loggerService.Error("Error starting server", err.Error())
	}
}

func validateEnvVariables(config domain.ConfigProvider) error {
	cloudCoverThreshold := config.GetInt("CLOUD_COVER_THRESHOLD")
	apiKey := config.GetString("WEATHER_BIT_API_KEY")

	if apiKey == "" {
		return errors.New("required environment variable for WEATHER_BIT_API_KEY is not set")
	}

	if cloudCoverThreshold == 0 || cloudCoverThreshold < 0 {
		return errors.New("CLOUD_COVER_THRESHOLD is not set, value should be a positive int")
	}
	return nil
}

func MakeSwaggerStaticServe(root string) (string, http.Handler) {
	fs := http.FileServer(http.Dir("./" + root))

	basePath := fmt.Sprintf("/%s/", root)
	url := basePath + "*"

	handler := http.StripPrefix(basePath, fs)

	return url, handler
}
