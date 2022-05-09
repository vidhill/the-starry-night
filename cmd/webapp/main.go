package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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
	weatherRepository := rest_api_repository.NewWeatherbitRepository(configService, httpService, loggerService)

	ISSService := service.NewISSLocationService(ISSRepository)
	weatherService := service.NewWeatherService(weatherRepository)

	dh := handlers.NewHandlers(configService, loggerService, ISSService, weatherService)

	mux := chi.NewRouter()

	//
	// route handlers
	//
	//

	// swagger:route GET /health api
	//
	// health endpoint for kubernetes liveness probe
	//
	// Produces:
	// - text/plain
	// Responses:
	// 				200:
	mux.Get("/health", dh.Health)

	// swagger:route GET /iss-position api ISSRequest
	//
	// Determines if ISS is overhead.
	//
	// Responses:
	// 				200: ISSResult
	// 				400: ErrorResponse
	// 				500: ErrorResponse
	mux.Get("/iss-position", composeHandler(addJsonHeader)(dh.ISSPosition))

	// start server
	port := configService.GetString("SERVER_PORT")

	loggerService.Info("listening on port", port)

	err := http.ListenAndServe(":"+port, mux)

	if err != nil {
		loggerService.Error("Error starting server", err.Error())
	}
}

// Composes multiple handler functions together
func composeHandler(manyHandlers ...func(http.HandlerFunc) http.HandlerFunc) func(http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		for _, v := range manyHandlers {
			h = v(h)
		}
		return h
	}
}

// Middleware function, adds Content-Type of json to any response
func addJsonHeader(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}
