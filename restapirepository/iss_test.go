package restapirepository

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	customErrors "github.com/vidhill/the-starry-night/customerrors"
	"github.com/vidhill/the-starry-night/mocks"
	"github.com/vidhill/the-starry-night/stubrepository"
)

func Test_GetCurrentLocation_happy_path(t *testing.T) {

	mockHttp := mocks.NewMockHTTP()
	config, logger, _ := createStubs()

	mockJSON := `
	{
		"message": "success",
		"iss_position": {
			"latitude": "19.2243",
			"longitude": "-32.4257"
		},
		"timestamp": 1652732569
	}	
	`

	mockResponse := http.Response{
		StatusCode: http.StatusOK,
		Body:       makeMockReadCloser(mockJSON),
	}

	mockHttp.On("Get", mock.AnythingOfType("string")).Return(&mockResponse, nil)

	instance := NewISSRepositoryRest(&config, &mockHttp, &logger)

	res, err := instance.GetCurrentLocation()

	mockHttp.AssertExpectations(t)

	assert.Nil(t, err)

	assert.Equal(t, float64(19.2243), res.Latitude)
	assert.Equal(t, float64(-32.4257), res.Longitude)

}

func Test_GetCurrentLocation_bad_request(t *testing.T) {

	mockHttp := mocks.NewMockHTTP()
	config, logger, _ := createStubs()

	mockResponse := http.Response{
		StatusCode: http.StatusBadRequest,
	}

	mockHttp.On("Get", mock.AnythingOfType("string")).Return(&mockResponse, nil)

	instance := NewISSRepositoryRest(&config, &mockHttp, &logger)

	_, err := instance.GetCurrentLocation()

	mockHttp.AssertExpectations(t)

	assert.NotNil(t, err)

	unexpectedResponseError := customErrors.UnexpectedResponseError{}
	assert.ErrorAs(t, err, &unexpectedResponseError)

}

func Test_SummarizeResponse(t *testing.T) {
	config := stubrepository.NewStubConfig()
	logger := stubrepository.NewStubLogger()
	http := mocks.NewMockHTTP()

	assert := assert.New(t)

	instance := NewISSRepositoryRest(&config, &http, &logger)

	mockJson := `
	{
		"timestamp": 1652652180,
		"message": "success",
		"iss_position": {
			"longitude": "150.1601",
			"latitude": "-37.9447"
		}
	}
	`
	resp, parseErr := parseMockJson(mockJson)
	if parseErr != nil {
		assert.FailNow("invalid mock JSON string passed to test")
		return
	}

	res, err := instance.SummarizeResponse(resp)

	assert.Nil(err)
	assert.Equal(float64(150.1601), res.Longitude)
	assert.Equal(float64(-37.9447), res.Latitude)

}

//
// Test utils
//

// parse a JSON string into an ApiResponse struct,
// fail if JSON string is invalid
func parseMockJson(s string) (ApiResponse, error) {
	response := ApiResponse{}

	err := json.Unmarshal([]byte(s), &response)

	return response, err
}

func createStubs() (stubrepository.StubConfig, stubrepository.StubLogger, stubrepository.StubHttp) {
	stubConfig := stubrepository.NewStubConfig()
	stubLogger := stubrepository.NewStubLogger()
	stubHttp := stubrepository.NewStubHttp()

	return stubConfig, stubLogger, stubHttp
}

func makeMockReadCloser(s string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(s))
}
