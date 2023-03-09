package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/vidhill/the-starry-night/model"
)

type ISS struct {
	mock.Mock
}

func (mock *ISS) GetCurrentLocation() (model.Coordinates, error) {
	args := mock.Called()
	result := args.Get(0)
	err := args.Error(1)
	return result.(model.Coordinates), err
}

// Repository 'Constructor' function
func NewMockISSRepository() ISS {
	return ISS{}
}
