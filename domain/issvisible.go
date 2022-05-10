package domain

import (
	"time"

	"github.com/vidhill/the-starry-night/model"
)

// swagger:model ISSResult
type ISSVisibleResult struct {
	ISSOverhead bool `json:"iss_overhead"`
}

type ISSVisibleRepository interface {
	GetISSVisible(time.Time, model.Coordinates) (ISSVisibleResult, error)
}
