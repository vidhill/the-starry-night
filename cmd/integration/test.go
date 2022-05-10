package main

import (
	"os"

	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/model"
	rest_api_repository "github.com/vidhill/the-starry-night/restapirepository"
	"github.com/vidhill/the-starry-night/service"
)

func main() {
	configService := service.NewConfigService(domain.NewViperConfig())
	loggerService := service.NewLoggerService(domain.NewStandardLogger(), configService.GetString("LOG_LEVEL"))
	httpService := service.NewHttpService(domain.NewDefaultHttpClient(loggerService))

	cliArgs := os.Args[1:]

	if len(cliArgs) == 0 {
		loggerService.Warn("No cli argument passed")

		os.Exit(2)
	}

	switch cliArgs[0] {
	case "iss":
		{
			ISSRepository := rest_api_repository.NewISSRepositoryRest(configService, httpService, loggerService)

			issService := service.NewISSLocationService(ISSRepository)

			res, err := issService.GetCurrentLocation()
			if err != nil {
				loggerService.Error(err.Error())
				return
			}
			loggerService.Info(res)
		}
	case "weather":
		{
			weatherRepository := rest_api_repository.NewWeatherbitRepository(configService, httpService, loggerService)
			weatherService := service.NewWeatherService(weatherRepository)

			res, err := weatherService.GetCurrent(model.Coordinates{Latitude: 51.89764968941597, Longitude: -8.46828736406348})
			if err != nil {
				loggerService.Error(err.Error())
				return
			}
			loggerService.Info(res)

		}
	default:
		{
			loggerService.Info("Default case")
		}
	}

}
