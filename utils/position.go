package utils

import (
	"strconv"

	"github.com/vidhill/the-starry-night/model"
)

func MakeCoordinatesFromString(lat, long string) (model.Coordinates, error) {
	emptyRes := model.Coordinates{}
	latitude, err := parseFloat(lat)
	if err != nil {
		return emptyRes, err
	}

	longitude, err := parseFloat(long)
	if err != nil {
		return emptyRes, err
	}

	result := model.Coordinates{
		Latitude:  latitude,
		Longitude: longitude,
	}

	return result, nil
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
