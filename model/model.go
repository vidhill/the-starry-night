package model

import "time"

//
// Models that are not domain specific
//

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

type DaylightTimes struct {
	Observation time.Time
	Sunrise     time.Time
	Sunset      time.Time
}
