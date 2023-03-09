package stubrepository

import (
	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/model"
)

type WeatherStubRepository struct {
	logger     domain.LogProvider
	mockResult domain.WeatherResult
}

func (s WeatherStubRepository) GetCurrent(location model.Coordinates) (domain.WeatherResult, error) {
	return s.mockResult, nil
}

func NewStubWeatherRepository(logger domain.LogProvider, mockResult domain.WeatherResult) domain.WeatherProvider {
	return WeatherStubRepository{
		logger:     logger,
		mockResult: mockResult,
	}
}
