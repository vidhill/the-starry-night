//go:build integration

package main

import (
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// SMOKE TEST
func Test_foo(t *testing.T) {
	response, err := http.Get("http://localhost:8080/iss-position")

	if err != nil {
		log.Println("Error fetching", err.Error())
		assert.FailNow(t, err.Error())
	}

	if response.StatusCode != http.StatusOK {
		errMessage := "non success response code received from api"
		log.Println(errMessage)

		assert.FailNow(t, err.Error())
	}

	contentTypeHeaders := response.Header.Values("Content-type")

	assert.Len(t, contentTypeHeaders, 1)
	assert.Equal(t, "application/json", contentTypeHeaders[0])
	assert.FailNow(t, "dummy failure")
}
