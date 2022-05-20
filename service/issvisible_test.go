package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/mocks"
	"github.com/vidhill/the-starry-night/model"
	"github.com/vidhill/the-starry-night/service"
	"github.com/vidhill/the-starry-night/stubrepository"
)

// using the same threshold for unit tests
const testCloudCoverThreshold = 30

func Test_GetISSVisible_happy_path(t *testing.T) {

	// create mocks
	//
	mockConfig, mockISS, mockWeather := initMocks()

	setupMockConfigValues(&mockConfig)

	mockISSPosition := model.Coordinates{
		Latitude:  0.123456,
		Longitude: 0.2,
	}
	mockISS.On("GetCurrentLocation").Return(mockISSPosition, nil)

	mockWeatherResult := makeMockClearNightResult()
	mockWeather.On("GetCurrent", mock.AnythingOfType("Coordinates")).Return(mockWeatherResult, nil)

	//
	// mocking end

	ISSVisibleService := makeService(mockConfig, mockISS, mockWeather)

	mockISSCoordinates := model.Coordinates{
		Latitude:  0.12,
		Longitude: 0.2,
	}

	mockNow := time.Unix(0, 0)

	res, err := ISSVisibleService.GetISSVisible(mockNow, mockISSCoordinates)

	mockISS.AssertExpectations(t)
	mockWeather.AssertExpectations(t)

	assert.Nil(t, err)
	assert.True(t, res.ISSOverhead)
}

func Test_GetISSVisible_error(t *testing.T) {

	// create mocks
	//
	mockConfig, mockISS, mockWeather := initMocks()

	setupMockConfigValues(&mockConfig)

	mockISSPosition := model.Coordinates{}
	mockISS.On("GetCurrentLocation").Return(mockISSPosition, errors.New("Dummy Error"))

	mockWeatherResult := makeMockClearNightResult()
	mockWeather.On("GetCurrent", mock.AnythingOfType("Coordinates")).Return(mockWeatherResult, nil)

	//
	// mocking end

	ISSVisibleService := makeService(mockConfig, mockISS, mockWeather)

	mockISSCoordinates := model.Coordinates{
		Latitude:  0.12,
		Longitude: 0.2,
	}

	mockNow := time.Unix(0, 0)

	res, err := ISSVisibleService.GetISSVisible(mockNow, mockISSCoordinates)

	mockISS.AssertExpectations(t)
	mockWeather.AssertExpectations(t)

	assert.NotNil(t, err)
	assert.False(t, res.ISSOverhead)
}

// exact position match scenarios
//
func Test_CheckISSVisible(t *testing.T) {

	loc := model.Coordinates{
		Latitude:  51.89764968941597,
		Longitude: -8.46828736406348,
	}

	w := makeMockClearNightResult()

	res := service.CheckISSVisible(loc, loc, w, testCloudCoverThreshold, 100)

	assert.True(t, res)
}

func Test_CheckISSVisible_daytime(t *testing.T) {

	loc := model.Coordinates{
		Latitude:  51.89764968941597,
		Longitude: -8.46828736406348,
	}

	w := makeMockClearDayResult()

	res := service.CheckISSVisible(loc, loc, w, testCloudCoverThreshold, 100)

	assert.False(t, res)
}

func Test_CheckISSVisible_clear_night(t *testing.T) {

	loc := model.Coordinates{
		Latitude:  51.89764968941597,
		Longitude: -8.46828736406348,
	}

	w := makeMockClearNightResult()

	res := service.CheckISSVisible(loc, loc, w, testCloudCoverThreshold, 100)

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

	res := service.CheckISSVisible(loc, ISSloc, w, testCloudCoverThreshold, 100)

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

	res := service.CheckISSVisible(loc, ISSloc, w, testCloudCoverThreshold, 2)

	assert.True(t, res)
}

func Test_CheckISSVisible_overcast_night(t *testing.T) {

	loc := model.Coordinates{
		Latitude:  51.89764968941597,
		Longitude: -8.46828736406348,
	}

	w := makeMockOvercastNightResult()

	res := service.CheckISSVisible(loc, loc, w, testCloudCoverThreshold, 100)

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

func makeService(
	mockConfig mocks.Config,
	mockISS mocks.ISS,
	mockWeather mocks.Weather,
) service.ISSVisibleService {

	stubLogger := service.NewLoggerService(stubrepository.NewStubLogger(), "INFO")

	configService := service.NewConfigService(&mockConfig)
	ISSService := service.NewISSLocationService(&mockISS)
	weatherService := service.NewWeatherService(&mockWeather)

	ISSVisibleService := service.NewISSVisibleService(configService, stubLogger, ISSService, weatherService)

	return ISSVisibleService
}

func initMocks() (mocks.Config, mocks.ISS, mocks.Weather) {
	mockConfig := mocks.NewMockConfig()
	mockISS := mocks.NewMockISSRepository()
	mockWeather := mocks.NewMockWeatherRepository()

	return mockConfig, mockISS, mockWeather
}

func setupMockConfigValues(mockConfig *mocks.Config) {
	mockConfig.On("GetInt", "CLOUD_COVER_THRESHOLD").Return(testCloudCoverThreshold)
	mockConfig.On("GetInt", "ACCURACY_NUM_DECIMAL_PLACES").Return(2)
}
