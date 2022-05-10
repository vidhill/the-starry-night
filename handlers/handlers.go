package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/model"
	"github.com/vidhill/the-starry-night/service"
	"github.com/vidhill/the-starry-night/utils"
)

var okMessage = []byte("ok")

type Handlers struct {
	Config            service.ConfigService
	ISSService        service.ISSLocationService
	Logger            service.LoggerService
	WeatherService    service.WeatherService
	ISSVisibleService service.ISSVisibleService
}

// swagger:parameters ISSRequest
type ISSRequest struct {
	// required: true
	// example: 51.89764968941597
	// In: query
	Latitude float64 `json:"lat"`
	// required: true
	// example: -8.46828736406348
	// In: query
	Longitude float64 `json:"long"`
}

// swagger:model ISSResult
type Result struct {
	ISSOverhead bool `json:"iss_overhead"`
}

func (s Handlers) Health(w http.ResponseWriter, req *http.Request) {
	w.Write(okMessage)
}

func (h Handlers) ISSPosition(w http.ResponseWriter, req *http.Request) {
	cloudCoverThreshold := h.Config.GetInt("CLOUD_COVER_THRESHOLD")
	accuracyNumDecimalPlaces := uint(h.Config.GetInt("ACCURACY_NUM_DECIMAL_PLACES"))

	lat, long := getLatLongQueryParams(req)

	if lat == "" {
		handleInvalidMissingQueryParm(w, req, "lat")
		return
	}

	if long == "" {
		handleInvalidMissingQueryParm(w, req, "long")
		return
	}

	// validates lat/long can be parsed into floats
	coordinates, err := utils.MakeCoordinatesFromString(lat, long)
	if err != nil {
		handleInvalidRequest(w, req, "Invalid float values for lat/long query params")
		return
	}

	ISSlocation, weatherResult, err := h.CallAPIsParallel(coordinates)

	if err != nil {
		handleInternalServerError(w, req, "failed")
		return
	}

	res := Result{
		ISSOverhead: CheckISSVisible(coordinates, ISSlocation, weatherResult, cloudCoverThreshold, accuracyNumDecimalPlaces),
	}

	bs, err := json.Marshal(res)

	if err != nil {
		handleInternalServerError(w, req, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func (h Handlers) CallAPIsParallel(coordinates model.Coordinates) (model.Coordinates, domain.WeatherResult, error) {
	logger := h.Logger

	coordinatesChan := make(chan model.Coordinates, 1)
	weatherChan := make(chan domain.WeatherResult, 1)
	errorsChan := make(chan error, 1)
	errorsChan1 := make(chan error, 1)

	go func() {
		logger.Info("requesting from ISS endpoint")
		ISSlocation, err := h.ISSService.GetCurrentLocation()
		coordinatesChan <- ISSlocation
		errorsChan <- err
		logger.Info("response from ISS endpoint")
	}()

	go func() {
		logger.Info("requesting from weather endpoint")
		weatherResult, err := h.WeatherService.GetCurrent(coordinates)
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

	if issErr != nil {
		logger.Error(issErr)
		return ISSlocation, weatherResult, issErr
	}

	if weatherErr != nil {
		logger.Error(weatherErr)
		return ISSlocation, weatherResult, weatherErr
	}

	return ISSlocation, weatherResult, nil
}

func NewHandlers(
	config service.ConfigService,
	logger service.LoggerService,
	issService service.ISSLocationService,
	weatherService service.WeatherService,
	ISSVisibleService service.ISSVisibleService,
) Handlers {

	return Handlers{
		Config:            config,
		Logger:            logger,
		ISSService:        issService,
		WeatherService:    weatherService,
		ISSVisibleService: ISSVisibleService,
	}
}
