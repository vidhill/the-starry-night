package domain

import (
	"time"

	"github.com/vidhill/the-starry-night/model"
)

type WeatherResult struct {
	CloudCover        int
	Sunrise           time.Time
	Sunset            time.Time
	ObserverationTime time.Time
}

type WeatherRepository interface {
	GetCurrent(position model.Coordinates) (WeatherResult, error)
}
