package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/vidhill/the-starry-night/model"
	"github.com/vidhill/the-starry-night/service"
)

type ISSVisibleService struct {
	mock.Mock
}

func (mock ISSVisibleService) GetISSVisible(now time.Time, coordinates model.Coordinates) (service.ISSVisibleResult, error) {
	args := mock.Called(now, coordinates)
	result := args.Get(0)
	err := args.Error(1)

	return result.(service.ISSVisibleResult), err
}

func NewISSVisibleService() ISSVisibleService {
	return ISSVisibleService{}
}
