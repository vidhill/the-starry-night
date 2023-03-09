package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/utils"
)

var (
	okMessage               = []byte("ok")
	isValidLatitude         = utils.MakeCheckFloatInRange(90)
	isValidLongitude        = utils.MakeCheckFloatInRange(180)
	invalidLatitudeMessage  = makeOutsideBoundsMessage("latitude", 90)
	invalidLongitudeMessage = makeOutsideBoundsMessage("longitude", 180)
)

type Handlers struct {
	Logger            domain.LogProvider
	ISSVisibleService domain.ISSVisibleProvider
}

// swagger:parameters ISSRequest
type ISSRequest struct {
	// required: true
	// example: 51.89764968941597
	// min: -90
	// max: 90
	// In: query
	Latitude float64 `json:"lat"`
	// required: true
	// example: -8.46828736406348
	// min: -180
	// max: 180
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
	logger := h.Logger

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

	if !isValidLatitude(coordinates.Latitude) {
		handleInvalidRequest(w, req, invalidLatitudeMessage)
		return
	}

	if !isValidLongitude(coordinates.Longitude) {
		handleInvalidRequest(w, req, invalidLongitudeMessage)
		return
	}

	res, err := h.ISSVisibleService.GetISSVisible(time.Now(), coordinates)

	if err != nil {
		logger.Error(err.Error())
		handleInternalServerError(w, req, "failed")
		return
	}

	bs, err := json.Marshal(res)

	if err != nil {
		logger.Error(err.Error())
		handleInternalServerError(w, req, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func NewHandlers(logger domain.LogProvider, ISSVisible domain.ISSVisibleProvider) Handlers {

	return Handlers{
		Logger:            logger,
		ISSVisibleService: ISSVisible,
	}
}

func makeOutsideBoundsMessage(name string, i int) string {
	return fmt.Sprintf("Invalid values for %s, value should value should not be greater/less than %v", name, i)
}
