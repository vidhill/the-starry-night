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
	w.Header().Add("Content-Type", "application/json")

	lat := req.URL.Query().Get("lat")
	long := req.URL.Query().Get("long")

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
