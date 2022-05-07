package domain

import "github.com/vidhill/the-starry-night/model"

type ISSRepository interface {
	GetCurrentLocation() (model.Coordinates, error)
}
