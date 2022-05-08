package service

import (
	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/model"
)

type WeatherService interface {
	GetCurrent(position model.Coordinates) (domain.WeatherResult, error)
}

type DefaultWeatherService struct {
	repo domain.WeatherRepository
}

func (s DefaultWeatherService) GetCurrent(position model.Coordinates) (domain.WeatherResult, error) {
	return s.repo.GetCurrent(position)
}

func NewWeatherService(repository domain.WeatherRepository) WeatherService {
	return DefaultWeatherService{
		repo: repository,
	}
}
