package domain

import "github.com/vidhill/the-starry-night/model"

type ISSLocationRepository interface {
	GetCurrentLocation() (model.Coordinates, error)
}
