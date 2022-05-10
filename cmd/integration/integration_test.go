package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

// SMOKE TEST
func Test_valid_request(t *testing.T) {
	response, err := http.Get("http://localhost:8080/iss-position?lat=51.89764968941597&long=-8.46828736406348")

	if err != nil {
		log.Println("Error fetching", err.Error())
		assert.FailNow(t, err.Error())
	}

	assertStatusCode(t, http.StatusOK, response)

	contentTypeHeaders := response.Header.Values("Content-type")

	fmt.Println()
	// dumping raw response to stdout
	io.Copy(os.Stdout, response.Body)
	fmt.Println()
	fmt.Println()

	assert.Len(t, contentTypeHeaders, 1)
	assert.Equal(t, "application/json", contentTypeHeaders[0])

}

func Test_invalid_request(t *testing.T) {
	response, err := http.Get("http://localhost:8080/iss-position")

	if err != nil {
		log.Println("Error fetching", err.Error())
		assert.FailNow(t, err.Error())
	}

	assertStatusCode(t, http.StatusBadRequest, response)
}

func Test_invalid_request_out_range_lat_long(t *testing.T) {
	response, err := http.Get("http://localhost:8080/iss-position?lat=100&long=-8")

	if err != nil {
		log.Println("Error fetching", err.Error())
		assert.FailNow(t, err.Error())
	}

	assertStatusCode(t, http.StatusBadRequest, response)

}

func assertStatusCode(t *testing.T, expected int, response *http.Response) {
	if response == nil {
		assert.FailNow(t, "no http.Response pointer was passed to assertion")
		return
	}

	actual := response.StatusCode
	assert.Equalf(t, expected, actual, `expected to return a response with the status code %v, actual response status was %v`, expected, actual)
}
