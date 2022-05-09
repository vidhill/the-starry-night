package stubrepository

import (
	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/model"
)

type ISSLocationRepositoryStub struct {
	logger domain.LoggerRepository
}

func (s ISSLocationRepositoryStub) GetCurrentLocation() (model.Coordinates, error) {

	result := model.Coordinates{
		Latitude:  51.89764968941597,
		Longitude: -8.46828736406348,
	}

	return result, nil
}

//
// Repository 'Constructor' function
//
func NewISSRepositoryStub(logger domain.LoggerRepository) ISSLocationRepositoryStub {

	return ISSLocationRepositoryStub{
		logger: logger,
	}
}
