package stubrepository

import (
	"time"

	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/model"
)

type WeatherStubRepository struct {
	logger domain.LoggerRepository
}

func (s WeatherStubRepository) GetCurrent(location model.Coordinates) (domain.WeatherResult, error) {

	emptyResult := domain.WeatherResult{
		CloudCover:        10,
		ObserverationTime: time.Now(),
		Sunrise:           time.Now().Local().Add(time.Hour * -10),
		Sunset:            time.Now().Local().Add(time.Hour * -3),
	}

	return emptyResult, nil

}

func NewStubWeatherRepository(logger domain.LoggerRepository) WeatherStubRepository {

	return WeatherStubRepository{
		logger: logger,
	}
}
