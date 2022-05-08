package domain

import "github.com/vidhill/the-starry-night/model"

type WeatherResult struct {
	AfterSunset bool
	CloudCover  float32
}

type WeatherRepository interface {
	GetCurrent(position model.Coordinates) (WeatherResult, error)
}
