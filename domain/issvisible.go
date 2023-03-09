package domain

import (
	"time"

	"github.com/vidhill/the-starry-night/model"
)

type ISSVisibleProvider interface {
	GetISSVisible(now time.Time, coordinates model.Coordinates) (ISSVisibleResult, error)
}

// swagger:model ISSResult
type ISSVisibleResult struct {
	ISSOverhead bool `json:"iss_overhead"`
}
