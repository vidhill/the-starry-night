package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vidhill/the-starry-night/utils"
)

//TODO move these tests elsewhere as are testing two different domains
//
func Test_determinIsNight_1(t *testing.T) {
	observerationTime, sunriseTime, sunsetTime, err := determineTimes("2017-08-28 16:45", "10:44", "23:47")
	assert.Nil(t, err)

	res := utils.DetermineIsNight(observerationTime, sunriseTime, sunsetTime)
	assert.False(t, res)
}

func Test_determinIsNight_2(t *testing.T) {
	observerationTime, sunriseTime, sunsetTime, err := determineTimes("2017-08-28 23:49", "10:44", "23:47")
	assert.Nil(t, err)

	res := utils.DetermineIsNight(observerationTime, sunriseTime, sunsetTime)
	assert.True(t, res)
}

func Test_determinIsNight_3(t *testing.T) {
	observerationTime, sunriseTime, sunsetTime, err := determineTimes("2017-08-28 08:49", "10:44", "23:47")
	assert.Nil(t, err)

	res := utils.DetermineIsNight(observerationTime, sunriseTime, sunsetTime)
	assert.True(t, res)
}

func Test_determinIsNight_4(t *testing.T) {
	observerationTime, sunriseTime, sunsetTime, err := determineTimes("2017-08-28 10:43", "10:44", "23:47")
	assert.Nil(t, err)

	res := utils.DetermineIsNight(observerationTime, sunriseTime, sunsetTime)
	assert.True(t, res)
}