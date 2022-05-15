package restapirepository

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockConfig struct{}
type MockLogger struct{}
type MockHttp struct{}

func (c MockConfig) GetBool(s string) bool     { return false }
func (c MockConfig) GetString(s string) string { return "" }
func (c MockConfig) GetInt(s string) int       { return 0 }

func (h MockHttp) Get(url string) (*http.Response, error) {
	return &http.Response{}, nil
}

func (l MockLogger) Debug(v ...interface{}) {}
func (l MockLogger) Info(v ...interface{})  {}
func (l MockLogger) Warn(v ...interface{})  {}
func (l MockLogger) Error(v ...interface{}) {}

func Test_SummarizeResponse(t *testing.T) {
	configService := MockConfig{}
	loggerService := MockLogger{}
	httpService := MockHttp{}

	assert := assert.New(t)

	instance := NewISSRepositoryRest(configService, httpService, loggerService)

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

	// if err := json.Unmarshal([]byte(s), &response); err != nil {
	// 	assert.FailNow(t, "invalid mock JSON string passed to test")
	// 	return response
	// }
	return response, err
}
