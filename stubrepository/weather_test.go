package stubrepository

import (
	"fmt"
	"testing"

	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/model"
)

func Test_(t *testing.T) {

	logger := domain.NewStandardLogger()
	weather := NewStubWeatherRepository(logger)

	res, _ := weather.GetCurrent(model.Coordinates{})

	fmt.Println("now", res.ObserverationTime)
	fmt.Println("Sunrise", res.Sunrise)
	fmt.Println("Sunset", res.Sunset)

}
