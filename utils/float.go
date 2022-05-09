package utils

import "math"

// todo perhaps should replace with more precise floating point implementation

// make function to round a float to n decimal places
func MakeRoundToNPlaces(places uint) func(f float64) float64 {
	placesF := float64(places)
	precisionF := math.Pow(10, placesF)

	return func(f float64) float64 {
		fl := float64(f)
		res := math.Round(fl*precisionF) / precisionF

		// above n decimal places rounding divides by zero
		if math.IsNaN(res) {
			return f
		}

		return res
	}
}
