package utils

import (
	"strconv"

	"github.com/vidhill/the-starry-night/model"
)

func MakeCoordinatesFromString(lat, long string) (model.Coordinates, error) {
	latitude, err := parseFloat(lat)
	if err != nil {
		return model.Coordinates{}, err
	}

	longitude, err := parseFloat(long)
	if err != nil {
		return model.Coordinates{}, err
	}

	result := model.Coordinates{
		Latitude:  latitude,
		Longitude: longitude,
	}

	return result, nil
}

func parseFloat(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}