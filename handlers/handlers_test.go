package handlers_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jgroeneveld/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/handlers"
	"github.com/vidhill/the-starry-night/mocks"
	"github.com/vidhill/the-starry-night/stubrepository"
	assertUtils "github.com/vidhill/the-starry-night/utils/assert"
)

var (
	errorResponseJSONSchema = schema.Map{
		"message":    schema.IsString,
		"error_code": schema.IsInteger,
	}
)

func Test_ISSPosition(t *testing.T) {

	// case service returns a isOverhead response
	t.Run("happyPath", func(t *testing.T) {

		mockResponse := domain.ISSVisibleResult{
			ISSOverhead: true,
		}

		mockISSVisible := makeMockISSService(mockResponse, nil)

		h := initHandler(mockISSVisible)

		rr := httptest.NewRecorder()
		req := makeISSRequest("lat=51.89764968941597&long=-8.46828736406348")

		h.ISSPosition(rr, req)

		res, data := getRecordedResponse(t, rr)

		expected := `
		{
			"iss_overhead": true
		}
		`

		mockISSVisible.AssertExpectations(t)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.JSONEq(t, expected, string(data))

	})

	// case service returns an error
	t.Run("error", func(t *testing.T) {

		mockResponse := domain.ISSVisibleResult{}
		modkErr := errors.New("mock error")

		mockISSVisible := makeMockISSService(mockResponse, modkErr)

		h := initHandler(mockISSVisible)

		rr := httptest.NewRecorder()
		req := makeISSRequest("lat=51.89764968941597&long=-8.46828736406348")

		h.ISSPosition(rr, req)

		res := rr.Result()

		mockISSVisible.AssertExpectations(t)

		assertUtils.MatchesJSONSchema(t, errorResponseJSONSchema, res.Body)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	})

	t.Run("missingQueryParam", func(t *testing.T) {

		mockISSService := mocks.NewISSVisibleService()
		h := initHandler(&mockISSService)

		testCases := []string{
			"",                   // all query params missing
			"long=-8.4",          // lat param missing
			"lat=-8.4",           // long param missing
			"lat=91&long=-8.4",   // lat param outside range
			"lat=51.89&long=190", // long param outside range
		}

		for _, queryParam := range testCases {
			rr := httptest.NewRecorder()
			req := makeISSRequest(queryParam)

			h.ISSPosition(rr, req)

			res := rr.Result()

			assertUtils.MatchesJSONSchema(t, errorResponseJSONSchema, res.Body)

			assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		}
	})

}

// Test utils
//

func initHandler(mockISSService *mocks.ISSVisibleService) handlers.Handlers {
	stubLogger := stubrepository.NewStubLogger()
	return handlers.NewHandlers(stubLogger, mockISSService)
}

func makeISSRequest(queryParams string) *http.Request {
	path := "/iss-position?" + queryParams
	return httptest.NewRequest(http.MethodGet, path, nil)
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

func makeMockISSService(mockResponse domain.ISSVisibleResult, err error) *mocks.ISSVisibleService {
	mockISSVisibleService := mocks.NewISSVisibleService()

	mockISSVisibleService.On(
		"GetISSVisible",
		mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("model.Coordinates"),
	).Return(mockResponse, err)

	return &mockISSVisibleService
}
