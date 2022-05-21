package handlers_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vidhill/the-starry-night/handlers"
	"github.com/vidhill/the-starry-night/model"
	"github.com/vidhill/the-starry-night/service"
	"github.com/vidhill/the-starry-night/stubrepository"
)

type MockISSService struct {
	mock.Mock
}

func (s MockISSService) GetISSVisible(now time.Time, coordinates model.Coordinates) (service.ISSVisibleResult, error) {
	return service.ISSVisibleResult{}, nil
}

func Test_ISSPosition_happyPath(t *testing.T) {
	h := initHandler()

	w, req := makeRequestRecorder("/iss-position?lat=51.89764968941597&long=-8.46828736406348")

	h.ISSPosition(w, req)

	res, data := getRecordedResponse(t, w)

	expected := `
	{
		"iss_overhead": false
	}
	`

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, expected, string(data))

}

func Test_ISSPosition_missingQueryParam(t *testing.T) {
	h := initHandler()

	w, req := makeRequestRecorder("/iss-position?long=-8.46828736406348")

	h.ISSPosition(w, req)

	res := w.Result()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

}

// Test utils
//

func initHandler() handlers.Handlers {
	stubLogger := stubrepository.NewStubLogger()
	mockISSService := MockISSService{}
	return handlers.NewHandlers(stubLogger, mockISSService)
}

func makeRequestRecorder(url string) (*httptest.ResponseRecorder, *http.Request) {
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	return w, req
}

func getRecordedResponse(t *testing.T, w *httptest.ResponseRecorder) (*http.Response, []byte) {
	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		assert.FailNow(t, "failed request")
		return res, []byte{}
	}

	return res, data
}
