package stubrepository

import (
	"fmt"
	"testing"

	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/model"
)

func Test_(t *testing.T) {

	logger := domain.NewStandardLogger()

	mockResult := makeMockClearNightResult()

	weather := NewStubWeatherRepository(logger, mockResult)

	res, _ := weather.GetCurrent(model.Coordinates{})

	fmt.Println("now", res.ObserverationTime)
	fmt.Println("Sunrise", res.Sunrise)
	fmt.Println("Sunset", res.Sunset)

}

func makeMockClearNightResult() domain.WeatherResult {

	t := time.Unix(0, 0).Add(time.Hour * 5)

	return domain.WeatherResult{
		CloudCover:        10,
		ObserverationTime: t,
		Sunrise:           t.Add(time.Hour * 2),
		Sunset:            t.Add(time.Hour * 15),
	}
}
