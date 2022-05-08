package main

import (
	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/service"
)

func main() {
	configService := service.NewConfigService(domain.NewViperConfig())
	loggerService := service.NewLoggerService(domain.NewStandardLogger())
	httpService := service.NewHttpService(domain.NewDefaultHttpClient(loggerService))

	ISSRepository := domain.NewISSRepositoryRest(configService, httpService, loggerService)

	issService := service.NewISSLocationService(ISSRepository)

	res, err := issService.GetCurrentLocation()
	if err != nil {
		loggerService.Error(err.Error())
		return
	}
	loggerService.Info(res)
}
