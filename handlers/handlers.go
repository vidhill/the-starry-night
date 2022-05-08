package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vidhill/the-starry-night/service"
)

var okMessage = []byte("ok")

type Handlers struct {
	ISSService service.ISSLocationService
}

func (s Handlers) Health(w http.ResponseWriter, req *http.Request) {
	w.Write(okMessage)
}

func (h Handlers) ISSPosition(w http.ResponseWriter, req *http.Request) {
	location, _ := h.ISSService.GetCurrentLocation()
	w.Header().Add("Content-Type", "application/json")

	bs, err := json.Marshal(location)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func NewHandlers(issService service.ISSLocationService) Handlers {
	return Handlers{
		ISSService: issService,
	}
}
