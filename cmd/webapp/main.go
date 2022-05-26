package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/handlers"
	"github.com/vidhill/the-starry-night/middleware"
	"github.com/vidhill/the-starry-night/service"
	stubRepository "github.com/vidhill/the-starry-night/stubrepository"
)

const SWAGGER_ROOT = "swagger-ui"

func main() {

	// wiring
	configService := service.NewConfigService(domain.NewViperConfig())
	loggerService := service.NewLoggerService(domain.NewStandardLogger(), configService.GetString("LOG_LEVEL"))

	configErr := validateEnvVariables(configService)
	if configErr != nil {
		loggerService.Error(configErr.Error())
		os.Exit(1)
	}

	// httpService := service.NewHttpService(domain.NewDefaultHttpClient(loggerService))

	// ISSRepository := rest_api_repository.NewISSRepositoryRest(configService, httpService, loggerService)
	ISSRepository := stubRepository.NewISSRepositoryStub(loggerService)

	// weatherRepository := rest_api_repository.NewWeatherbitRepository(configService, httpService, loggerService)
	weatherResult := makeMockClearNightResult()
	weatherRepository := stubRepository.NewStubWeatherRepository(loggerService, weatherResult)

	ISSService := service.NewISSLocationService(ISSRepository)
	weatherService := service.NewWeatherService(weatherRepository)

	ISSVisibleService := service.NewISSVisibleService(configService, loggerService, ISSService, weatherService)
	// custom middleware to log using the wrapped logger service
	customLogMiddleware := middleware.MakeMyLoggerMiddleware(loggerService)

	dh := handlers.NewHandlers(loggerService, ISSVisibleService)

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

func validateEnvVariables(config service.ConfigService) error {
	cloudCoverThreshold := config.GetInt("CLOUD_COVER_THRESHOLD")
	apiKey := config.GetString("WEATHER_BIT_API_KEY")

	if apiKey == "" {
		return errors.New("required environment varible for WEATHER_BIT_API_KEY is not set")
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

// helpers
//
//
func timeFromString(s string) time.Time {
	t, _ := time.Parse(time.RFC822, s)
	return t
}

func makeMockClearNightResult() domain.WeatherResult {
	return domain.WeatherResult{
		CloudCover:        10,
		ObserverationTime: timeFromString("02 Jan 06 06:04 MST"),
		Sunrise:           timeFromString("02 Jan 06 08:00 MST"),
		Sunset:            timeFromString("02 Jan 06 22:00 MST"),
	}
}
