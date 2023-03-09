package restapirepository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vidhill/the-starry-night/utils"
)

// TODO move these tests elsewhere as are testing two different domains

func Test_DetermineIsNight(t *testing.T) {

	t.Run("A", func(t *testing.T) {

		r, err := determineTimes("2017-08-28 16:45", "10:44", "23:47")
		assert.Nil(t, err)

		res := utils.DetermineIsNight(r.Observation, r)
		assert.False(t, res)
	})

	t.Run("B", func(t *testing.T) {
		r, err := determineTimes("2017-08-28 23:49", "10:44", "23:47")
		assert.Nil(t, err)

		res := utils.DetermineIsNight(r.Observation, r)
		assert.True(t, res)
	})
	t.Run("C", func(t *testing.T) {
		r, err := determineTimes("2017-08-28 08:49", "10:44", "23:47")
		assert.Nil(t, err)

		res := utils.DetermineIsNight(r.Observation, r)
		assert.True(t, res)
	})

	t.Run("D", func(t *testing.T) {
		r, err := determineTimes("2017-08-28 10:43", "10:44", "23:47")
		assert.Nil(t, err)

		res := utils.DetermineIsNight(r.Observation, r)
		assert.True(t, res)
	})
}

func Test_extractDateString(t *testing.T) {
	res := extractDateString("2017-08-28 10:43")

	assert.Equal(t, res, "2017-08-28")
}
