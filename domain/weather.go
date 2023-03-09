package domain

import (
	"github.com/vidhill/the-starry-night/model"
)

type WeatherResult struct {
	CloudCover    int
	DaylightTimes model.DaylightTimes
}

type WeatherProvider interface {
	GetCurrent(position model.Coordinates) (WeatherResult, error)
}
