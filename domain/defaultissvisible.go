package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/vidhill/the-starry-night/model"
	"github.com/vidhill/the-starry-night/utils"
)

type DefaultISSVisible struct {
	Config  ConfigRepository
	Logger  LoggerRepository
	ISS     ISSLocationRepository
	Weather WeatherRepository
}

func (s DefaultISSVisible) GetISSVisible(now time.Time, coordinates model.Coordinates) (ISSVisibleResult, error) {
	ISSlocation, weatherResult, err := s.CallAPIsParallel(coordinates)

	if err != nil {
		return ISSVisibleResult{}, err
	}

	cloudCoverThreshold := s.Config.GetInt("CLOUD_COVER_THRESHOLD")
	accuracyNumDecimalPlaces := uint(s.Config.GetInt("ACCURACY_NUM_DECIMAL_PLACES"))

	res := ISSVisibleResult{
		ISSOverhead: CheckISSVisible(
			coordinates,
			ISSlocation,
			weatherResult,
			cloudCoverThreshold,
			accuracyNumDecimalPlaces),
	}

	return res, nil
}

func (h DefaultISSVisible) CallAPIsParallel(coordinates model.Coordinates) (model.Coordinates, WeatherResult, error) {
	logger := h.Logger

	coordinatesChan := make(chan model.Coordinates, 1)
	weatherChan := make(chan WeatherResult, 1)
	errorsChan := make(chan error, 1)
	errorsChan1 := make(chan error, 1)

	go func() {
		logger.Info("requesting from ISS endpoint")
		ISSlocation, err := h.ISS.GetCurrentLocation()
		coordinatesChan <- ISSlocation
		errorsChan <- err
		logger.Info("response from ISS endpoint")
	}()

	go func() {
		logger.Info("requesting from weather endpoint")
		weatherResult, err := h.Weather.GetCurrent(coordinates)
		weatherChan <- weatherResult
		errorsChan1 <- err
		logger.Info("response from weather endpoint")
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
			return ISSlocation, weatherResult, fmt.Errorf("failed to get iss location %w", issErr)
		}

		if weatherErr != nil {
			return ISSlocation, weatherResult, fmt.Errorf("failed to get weather info %w", issErr)
		}

	}

	return ISSlocation, weatherResult, nil
}

func NewDefaultISSVisible(
	config ConfigRepository,
	logger LoggerRepository,
	iss ISSLocationRepository,
	weather WeatherRepository,
) DefaultISSVisible {

	return DefaultISSVisible{
		Config:  config,
		Logger:  logger,
		ISS:     iss,
		Weather: weather,
	}
}

func CheckISSVisible(position, ISSPosition model.Coordinates, weatherResult WeatherResult, cloudCoverThreshold int, precision uint) bool {

	if weatherResult.CloudCover >= cloudCoverThreshold {
		return false
	}

	positionsMatch := MakeCoordinatesMatch(precision)

	return positionsMatch(position, ISSPosition)
}

func MakeCoordinatesMatch(precision uint) func(model.Coordinates, model.Coordinates) bool {
	positionsMatch := MakePositionMatch(precision)

	return func(a, b model.Coordinates) bool {
		if !positionsMatch(a.Latitude, b.Latitude) {
			return false
		}
		return positionsMatch(a.Longitude, b.Longitude)
	}
}

func MakePositionMatch(precision uint) func(a, b float64) bool {
	roundToPrecision := utils.MakeRoundToNPlaces(precision)
	return func(a, b float64) bool {
		return roundToPrecision(a) == roundToPrecision(b)
	}
}
