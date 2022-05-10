package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MakeRoundToNPlaces(t *testing.T) {

	var f float64 = 51.12345

	testCases := map[uint]float64{
		1:    51.1,
		2:    51.12,
		3:    51.123,
		4:    51.1235,
		5:    51.12345,
		1000: 51.12345,
	}

	for numPlaces, expected := range testCases {
		roundToNPlaces := MakeRoundToNPlaces(numPlaces)

		res := roundToNPlaces(f)
		assert.Equal(t, expected, res)
	}

}

func Test_MakeRoundToNPlaces_1(t *testing.T) {
	roundToNPlaces := MakeRoundToNPlaces(0)
	res := roundToNPlaces(51.12345)
	assert.Equal(t, float64(51), res)
}
