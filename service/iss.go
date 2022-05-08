package service

import (
	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/model"
)

type ISSLocationService interface {
	GetCurrentLocation() (model.Coordinates, error)
}

type DefaultISSService struct {
	repo domain.ISSLocationRepository
}

func (s DefaultISSService) GetCurrentLocation() (model.Coordinates, error) {
	return s.repo.GetCurrentLocation()
}

func NewISSLocationService(repository domain.ISSLocationRepository) ISSLocationService {
	return DefaultISSService{
		repo: repository,
	}
}
