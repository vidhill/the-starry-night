package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/model"
	"github.com/vidhill/the-starry-night/utils"
)

type ISSVisibleProvider interface {
	GetISSVisible(now time.Time, coordinates model.Coordinates) (domain.ISSVisibleResult, error)
}

type DefaultISSVisibleService struct {
	Config      domain.ConfigProvider
	Logger      domain.LogProvider
	ISSLocation domain.ISSLocationProvider
	Weather     domain.WeatherProvider
}

func (s DefaultISSVisibleService) GetISSVisible(now time.Time, coordinates model.Coordinates) (domain.ISSVisibleResult, error) {
	ISSlocation, weatherResult, err := s.CallAPIsParallel(coordinates)

	if err != nil {
		return domain.ISSVisibleResult{}, err
	}

	cloudCoverThreshold := s.Config.GetInt("CLOUD_COVER_THRESHOLD")
	accuracyNumDecimalPlaces := uint(s.Config.GetInt("ACCURACY_NUM_DECIMAL_PLACES"))

	issOverhead := CheckISSVisible(coordinates, ISSlocation, weatherResult, cloudCoverThreshold, accuracyNumDecimalPlaces)

	res := domain.ISSVisibleResult{
		ISSOverhead: issOverhead,
	}

	return res, nil
}

func (h DefaultISSVisibleService) CallAPIsParallel(coordinates model.Coordinates) (model.Coordinates, domain.WeatherResult, error) {

	coordinatesChan := make(chan model.Coordinates, 1)
	weatherChan := make(chan domain.WeatherResult, 1)
	errorsChan := make(chan error, 1)
	errorsChan1 := make(chan error, 1)

	go func() {
		ISSlocation, err := h.ISSLocation.GetCurrentLocation()
		coordinatesChan <- ISSlocation
		errorsChan <- err
	}()

	go func() {
		weatherResult, err := h.Weather.GetCurrent(coordinates)
		weatherChan <- weatherResult
		errorsChan1 <- err
	}()

	ISSlocation := <-coordinatesChan
	issErr := <-errorsChan

	weatherResult := <-weatherChan
	weatherErr := <-errorsChan1

	close(coordinatesChan)
	close(weatherChan)
	close(errorsChan)
	close(errorsChan1)

	if utils.AnyErrorNotNil(issErr, weatherErr) {

		if utils.AllErrorsPresent(issErr, weatherErr) {
			return ISSlocation, weatherResult, errors.New("both services returned errors")
		}

		/// todo tidy, use error folding ?
		if issErr != nil {
			return ISSlocation, weatherResult, fmt.Errorf("failed to get iss location, %w", issErr)
		}

		if weatherErr != nil {
			return ISSlocation, weatherResult, fmt.Errorf("failed to get weather info, %w", weatherErr)
		}

	}

	return ISSlocation, weatherResult, nil
}

func NewISSVisibleService(config domain.ConfigProvider, logger domain.LogProvider, iss domain.ISSLocationProvider, weather domain.WeatherProvider) ISSVisibleProvider {
	return DefaultISSVisibleService{
		Config:      config,
		Logger:      logger,
		ISSLocation: iss,
		Weather:     weather,
	}
}

func CheckISSVisible(position, ISSPosition model.Coordinates, weatherResult domain.WeatherResult, cloudCoverThreshold int, precision uint) bool {

	// if it's not night will not be visible
	if !utils.DetermineIsNight(weatherResult.DaylightTimes.Observation, weatherResult.DaylightTimes) {
		return false
	}

	if weatherResult.CloudCover >= cloudCoverThreshold {
		return false
	}

	positionsMatch := MakeCoordinatesMatch(precision)

	return positionsMatch(position, ISSPosition)
}

func MakeCoordinatesMatch(precision uint) func(model.Coordinates, model.Coordinates) bool {
	positionsMatch := MakePositionMatch(precision)

	return func(a, b model.Coordinates) bool {
		// just first check latitude first
		// if they don't match then must not match
		if !positionsMatch(a.Latitude, b.Latitude) {
			return false
		}
		// latitudes match, so now check longitude
		return positionsMatch(a.Longitude, b.Longitude)
	}
}

func MakePositionMatch(precision uint) func(a, b float64) bool {
	roundToPrecision := utils.MakeRoundToNPlaces(precision)
	return func(a, b float64) bool {
		return roundToPrecision(a) == roundToPrecision(b)
	}
}
