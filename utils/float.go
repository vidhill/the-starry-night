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

// returns a function that checks if a float
// 	is not
//     greater than the value, or
//     less than the negative
func MakeCheckFloatInRange(i int) func(float64) bool {
	bounds := float64(i)
	lowerBound := -1 * bounds
	return func(f float64) bool {
		if f < lowerBound || f > bounds {
			return false
		}
		return true
	}
}
