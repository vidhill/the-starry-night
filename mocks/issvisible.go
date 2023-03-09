package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/model"
)

type ISSVisibleService struct {
	mock.Mock
}

func (mock *ISSVisibleService) GetISSVisible(now time.Time, coordinates model.Coordinates) (domain.ISSVisibleResult, error) {
	args := mock.Called(now, coordinates)
	result := args.Get(0)
	err := args.Error(1)

	return result.(domain.ISSVisibleResult), err
}

func NewISSVisibleService() ISSVisibleService {
	return ISSVisibleService{}
}
