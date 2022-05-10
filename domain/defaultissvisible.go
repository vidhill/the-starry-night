package domain

import (
	"time"

	"github.com/vidhill/the-starry-night/model"
)

type DefaultISSVisible struct {
	config  ConfigRepository
	logger  LoggerRepository
	iss     ISSLocationRepository
	weather WeatherRepository
}

func (s DefaultISSVisible) GetISSVisible(now time.Time, coordinates model.Coordinates) (ISSVisibleResult, error) {
	return ISSVisibleResult{}, nil
}

func (h DefaultISSVisible) CallAPIsParallel(coordinates model.Coordinates) (model.Coordinates, WeatherResult, error) {
	logger := h.logger

	coordinatesChan := make(chan model.Coordinates, 1)
	weatherChan := make(chan WeatherResult, 1)
	errorsChan := make(chan error, 1)
	errorsChan1 := make(chan error, 1)

	go func() {
		logger.Info("requesting from ISS endpoint")
		ISSlocation, err := h.iss.GetCurrentLocation()
		coordinatesChan <- ISSlocation
		errorsChan <- err
		logger.Info("response from ISS endpoint")
	}()

	go func() {
		logger.Info("requesting from weather endpoint")
		weatherResult, err := h.weather.GetCurrent(coordinates)
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
		config:  config,
		logger:  logger,
		iss:     iss,
		weather: weather,
	}
}
