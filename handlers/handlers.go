package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vidhill/the-starry-night/service"
	"github.com/vidhill/the-starry-night/utils"
)

var okMessage = []byte("ok")

type Handlers struct {
	Config         service.ConfigService
	ISSService     service.ISSLocationService
	Logger         service.LoggerService
	WeatherService service.WeatherService
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

func NewHandlers(
	config service.ConfigService,
	logger service.LoggerService,
	issService service.ISSLocationService,
	weatherService service.WeatherService,
) Handlers {

	return Handlers{
		Config:         config,
		Logger:         logger,
		ISSService:     issService,
		WeatherService: weatherService,
	}
}
