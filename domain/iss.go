package domain

import "github.com/vidhill/the-starry-night/model"

type ISSLocationProvider interface {
	GetCurrentLocation() (model.Coordinates, error)
}
