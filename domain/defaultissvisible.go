package domain

import (
	"time"

	"github.com/vidhill/the-starry-night/model"
)

type DefaultISSVisible struct {
	config  ConfigRepository
	logger  LoggerRepository
	iss     ISSLocationRepository
	weather WeatherRepository
}

func (s DefaultISSVisible) GetISSVisible(now time.Time, coordinates model.Coordinates) (ISSVisibleResult, error) {
	return ISSVisibleResult{}, nil
}

func NewDefaultISSVisible(
	config ConfigRepository,
	logger LoggerRepository,
	iss ISSLocationRepository,
	weather WeatherRepository,
) DefaultISSVisible {

	return DefaultISSVisible{
		config:  config,
		logger:  logger,
		iss:     iss,
		weather: weather,
	}
}
