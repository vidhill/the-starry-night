package utils

import (
	"time"

	"github.com/vidhill/the-starry-night/model"
)

func DetermineIsNight(currentTime time.Time, o model.DaylightTimes) bool {
	isBeforeSunrise := currentTime.Before(o.Sunrise)
	isAfterSunset := currentTime.After(o.Sunset)

	return isBeforeSunrise || isAfterSunset
}
