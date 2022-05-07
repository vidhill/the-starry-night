package service

import (
	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/model"
)

type ISSService interface {
	GetCurrentLocation() (model.Coordinates, error)
}

type DefaultISSService struct {
	repo domain.ISSRepository
}

func (s DefaultISSService) GetCurrentLocation() (model.Coordinates, error) {
	return s.repo.GetCurrentLocation()
}

func NewISSService(repository domain.ISSRepository) ISSService {
	return DefaultISSService{
		repo: repository,
	}
}
