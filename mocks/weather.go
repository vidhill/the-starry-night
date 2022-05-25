package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/model"
)

type Weather struct {
	mock.Mock
}

func (mock *Weather) GetCurrent(location model.Coordinates) (domain.WeatherResult, error) {
	args := mock.Called(location)
	result := args.Get(0)
	err := args.Error(1)
	return result.(domain.WeatherResult), err
}

func NewMockWeatherRepository() Weather {
	return Weather{}
}
