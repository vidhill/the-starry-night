//go:build integration

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

	if response.StatusCode != http.StatusOK {
		errMessage := "non success response code received from api"
		log.Println(errMessage)

		assert.FailNow(t, errMessage)
		return
	}

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

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	contentTypeHeaders := response.Header.Values("Content-type")

	assert.Len(t, contentTypeHeaders, 1)
	assert.Equal(t, "application/json", contentTypeHeaders[0])

}
