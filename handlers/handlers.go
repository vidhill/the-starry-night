package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vidhill/the-starry-night/service"
)

var okMessage = []byte("ok")

type Handlers struct {
	ISSService service.ISSLocationService
	Logger     service.LoggerService
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

	location, err := h.ISSService.GetCurrentLocation()

	if err != nil {
		handleInternalServerError(w, req, "Error calling ISS API")
		return
	}

	bs, err := json.Marshal(location)

	if err != nil {
		handleInternalServerError(w, req, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func NewHandlers(logger service.LoggerService, issService service.ISSLocationService) Handlers {
	return Handlers{
		Logger:     logger,
		ISSService: issService,
	}
}

func getLatLongQueryParams(req *http.Request) (string, string) {
	lat := getQueryParam(req, "lat")
	long := getQueryParam(req, "long")

	return lat, long
}
