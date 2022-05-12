package service_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/model"
	"github.com/vidhill/the-starry-night/service"
)

// exact position match scenarios
//
func Test_CheckISSVisible(t *testing.T) {

	loc := model.Coordinates{
		Latitude:  51.89764968941597,
		Longitude: -8.46828736406348,
	}

	w := makeMockClearNightResult()

	res := service.CheckISSVisible(loc, loc, w, 30, 100)

	assert.True(t, res)
}

func Test_CheckISSVisible_daytime(t *testing.T) {

	loc := model.Coordinates{
		Latitude:  51.89764968941597,
		Longitude: -8.46828736406348,
	}

	w := makeMockClearDayResult()

	res := service.CheckISSVisible(loc, loc, w, 30, 100)

	assert.False(t, res)
}

func Test_CheckISSVisible_clear_night(t *testing.T) {

	loc := model.Coordinates{
		Latitude:  51.89764968941597,
		Longitude: -8.46828736406348,
	}

	w := makeMockClearNightResult()

	res := service.CheckISSVisible(loc, loc, w, 30, 100)

	assert.True(t, res)
}

// not directly overhead, with max accuracy
func Test_CheckISSVisible_clear_night_not_overhead(t *testing.T) {

	loc := model.Coordinates{
		Latitude:  51.89764968941597,
		Longitude: -8.46828736406348,
	}

	ISSloc := model.Coordinates{
		Latitude:  51.89764968941598,
		Longitude: -8.46828736406348,
	}

	w := makeMockClearNightResult()

	res := service.CheckISSVisible(loc, ISSloc, w, 30, 100)

	assert.False(t, res)
}

// positions not exactly the same but within accuracy (two decimal places) theshhold
func Test_CheckISSVisible_clear_night_not_exactly(t *testing.T) {

	loc := model.Coordinates{
		Latitude:  51.899,
		Longitude: -8.468,
	}

	ISSloc := model.Coordinates{
		Latitude:  51.897,
		Longitude: -8.468,
	}

	w := makeMockClearNightResult()

	res := service.CheckISSVisible(loc, ISSloc, w, 30, 2)

	assert.True(t, res)
}

func Test_CheckISSVisible_overcast_night(t *testing.T) {

	loc := model.Coordinates{
		Latitude:  51.89764968941597,
		Longitude: -8.46828736406348,
	}

	w := makeMockOvercastNightResult()

	res := service.CheckISSVisible(loc, loc, w, 30, 100)

	assert.False(t, res)
}

func makeMockClearNightResult() domain.WeatherResult {
	return domain.WeatherResult{
		CloudCover:        10,
		ObserverationTime: timeFromString("02 Jan 06 06:04 MST"),
		Sunrise:           timeFromString("02 Jan 06 08:00 MST"),
		Sunset:            timeFromString("02 Jan 06 22:00 MST"),
	}
}

func makeMockOvercastNightResult() domain.WeatherResult {
	return domain.WeatherResult{
		CloudCover:        100,
		ObserverationTime: timeFromString("02 Jan 06 06:04 MST"),
		Sunrise:           timeFromString("02 Jan 06 08:00 MST"),
		Sunset:            timeFromString("02 Jan 06 22:00 MST"),
	}
}

// helpers
//
//
func timeFromString(s string) time.Time {
	t, _ := time.Parse(time.RFC822, s)
	return t
}

func makeMockClearDayResult() domain.WeatherResult {

	return domain.WeatherResult{
		CloudCover:        10,
		ObserverationTime: timeFromString("02 Jan 06 10:00 MST"),
		Sunrise:           timeFromString("02 Jan 06 08:00 MST"),
		Sunset:            timeFromString("02 Jan 06 22:00 MST"),
	}
}
