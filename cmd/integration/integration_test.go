package main

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"testing"

	"github.com/jgroeneveld/schema"
	"github.com/stretchr/testify/assert"
	assertUtils "github.com/vidhill/the-starry-night/utils/assert"
)

var (
	baseUrl        string
	defaultBaseUrl = "http://localhost:8080"
	re             = regexp.MustCompile("http[s]?://.+")
)

func TestMain(m *testing.M) {
	// use base url from env variable if present
	baseUrl = setBaseUrl(defaultBaseUrl)

	// convert baseurl string into url.URL
	host, err := url.Parse(baseUrl)

	if err != nil {
		log.Printf("\n\n\tUnable to parse baseUrl, value was: \"%s\"\n\n\n", baseUrl)
		os.Exit(1)
	}

	if !rawTCPConnect(host.Host) {
		log.Printf("Integration tests not run because unable to connect to host at: %s\n Is the application running?", host.Host)
		os.Exit(1)
	}

	// execute tests
	code := m.Run()
	os.Exit(code)
}

func Test(t *testing.T) {

	t.Run("valid_request", func(t *testing.T) {

		response, err := http.Get(baseUrl + "/iss-position?lat=51.89764968941597&long=-8.46828736406348")

		if err != nil {
			log.Println("Error fetching", err.Error())
			assert.FailNow(t, err.Error())
		}

		assertStatusCode(t, http.StatusOK, response)

		contentTypeHeaders := response.Header.Values("Content-type")

		assert.Len(t, contentTypeHeaders, 1)
		assert.Equal(t, "application/json", contentTypeHeaders[0])

		validResponseSchema := schema.Map{
			"iss_overhead": schema.IsBool,
		}

		assertUtils.MatchesJSONSchema(t, validResponseSchema, response.Body)

	})

	t.Run("invalid_request", func(t *testing.T) {
		response, err := http.Get(baseUrl + "/iss-position")

		if err != nil {
			log.Println("Error fetching", err.Error())
			assert.FailNow(t, err.Error())
		}

		assertStatusCode(t, http.StatusBadRequest, response)

	})

	t.Run("Test_invalid_request_out_range_lat_long", func(t *testing.T) {

		response, err := http.Get(baseUrl + "/iss-position?lat=100&long=-8")

		if err != nil {
			log.Println("Error fetching", err.Error())
			assert.FailNow(t, err.Error())
		}

		assertStatusCode(t, http.StatusBadRequest, response)

	})
}

func assertStatusCode(t *testing.T, expected int, response *http.Response) {
	t.Helper()
	if response == nil {
		assert.FailNow(t, "no http.Response pointer was passed to assertion")
		return
	}

	actual := response.StatusCode
	assert.Equalf(t, expected, actual, `expected to return a response with the status code %v, actual response status was %v`, expected, actual)
}

func setBaseUrl(defaultBaseUrl string) string {
	host := os.Getenv("INTEGRATION_TEST_HOSTNAME") //nolint:forbidigo // using getEnv here as do not want to import config just for single value
	if host == "" {
		return defaultBaseUrl
	}

	if !re.MatchString(host) {
		log.Printf("\n\n\tInvalid base path passed as env variable, value was: \"%s\"\n\n\n", host)
		os.Exit(1)
	}

	return host
}

func rawTCPConnect(host string) bool {

	timeout := time.Second

	conn, err := net.DialTimeout("tcp", host, timeout)

	if err != nil {
		return false
	}

	defer conn.Close()

	return conn != nil
}
