package service

import (
	"time"

	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/model"
)

type ISSVisibleService interface {
	GetISSVisible(time.Time, model.Coordinates) (domain.ISSVisibleResult, error)
}

type DefaultISSVisibleService struct {
	Repo domain.ISSVisibleRepository
}

func (s DefaultISSVisibleService) GetISSVisible(now time.Time, coordinates model.Coordinates) (domain.ISSVisibleResult, error) {
	return s.Repo.GetISSVisible(now, coordinates)
}

func NewISSVisibleService(repository domain.ISSVisibleRepository) ISSVisibleService {
	return DefaultISSVisibleService{
		Repo: repository,
	}
}
