package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AllErrorsPresent_1(t *testing.T) {
	res := AllErrorsPresent(nil, nil)

	assert.False(t, res)
}

func Test_AllErrorsPresent_2(t *testing.T) {
	err := errors.New("dummy")
	res := AllErrorsPresent(err, err)

	assert.True(t, res)
}

func Test_AllErrorsPresent_3(t *testing.T) {
	err := errors.New("dummy")
	res := AllErrorsPresent(err, nil, err)

	assert.False(t, res)
}
