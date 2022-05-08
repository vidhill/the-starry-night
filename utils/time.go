package utils

import "time"

func DetermineIsNight(currentTime, sunrise, sunset time.Time) bool {
	isBeforeSunrise := currentTime.Before(sunrise)
	isAfterSunset := currentTime.After(sunset)

	return isBeforeSunrise || isAfterSunset
}
