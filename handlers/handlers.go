package handlers

import (
	"encoding/json"
	"math"
	"net/http"

	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/model"
	"github.com/vidhill/the-starry-night/service"
	"github.com/vidhill/the-starry-night/utils"
)

var okMessage = []byte("ok")

type Handlers struct {
	ISSService     service.ISSLocationService
	Logger         service.LoggerService
	WeatherService service.WeatherService
}

type Result struct {
	ISSOverhead bool `json:"iss_overhead"`
}

func (s Handlers) Health(w http.ResponseWriter, req *http.Request) {
	w.Write(okMessage)
}

func (h Handlers) ISSPosition(w http.ResponseWriter, req *http.Request) {
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

	ISSlocation, err := h.ISSService.GetCurrentLocation()

	if err != nil {
		handleInternalServerError(w, req, "Error calling ISS API")
		return
	}

	weatherResult, err := h.WeatherService.GetCurrent(coordinates)

	if err != nil {
		handleInternalServerError(w, req, "Error calling weather API")
		return
	}

	res := Result{
		ISSOverhead: CheckISSVisible(coordinates, ISSlocation, weatherResult, 30, 100),
	}

	bs, err := json.Marshal(res)

	if err != nil {
		handleInternalServerError(w, req, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func NewHandlers(logger service.LoggerService, issService service.ISSLocationService, weatherService service.WeatherService) Handlers {

	return Handlers{
		Logger:         logger,
		ISSService:     issService,
		WeatherService: weatherService,
	}
}

func getLatLongQueryParams(req *http.Request) (string, string) {
	lat := getQueryParam(req, "lat")
	long := getQueryParam(req, "long")

	return lat, long
}

func CheckISSVisible(position, ISSPosition model.Coordinates, weatherResult domain.WeatherResult, cloudCoverThreshold, precision int) bool {
	if weatherResult.CloudCover <= cloudCoverThreshold {
		return false
	}

	positionsMatch := MakePositionsMatch(precision)

	return positionsMatch(position, ISSPosition)
}

// todo perhaps should replace with more precise floating point implementation
func MakeRoundToNPlaces(precision int) func(f float32) float32 {
	precisionF := float64(precision)
	return func(f float32) float32 {
		fl := float64(f)
		res := math.Round(fl+precisionF) / precisionF
		return float32(res)
	}
}

func MakePositionsMatch(precision int) func(model.Coordinates, model.Coordinates) bool {
	roundToPrecision := MakeRoundToNPlaces(precision)
	return func(a, b model.Coordinates) bool {

		latsMatch := roundToPrecision(a.Latitude) == roundToPrecision(b.Latitude)
		if !latsMatch {
			return false
		}
		longsMatch := roundToPrecision(a.Longitude) == roundToPrecision(b.Longitude)
		return longsMatch
	}
}
