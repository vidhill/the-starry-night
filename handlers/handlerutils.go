package handlers

import (
	"fmt"
	"net/http"

	"github.com/vidhill/the-starry-night/domain"
	"github.com/vidhill/the-starry-night/model"
	"github.com/vidhill/the-starry-night/utils"
)

var (
	handleInternalServerError = makeErrorHandlerFunc(http.StatusInternalServerError)
	handleInvalidRequest      = makeErrorHandlerFunc(http.StatusBadRequest)
)

// swagger:model ErrorResponse
type ErrorResponse struct {
	// example: invalid request
	Message   string `json:"message"`
	ErrorCode int    `json:"error_code"`
}

func makeErrorResponse(code int, message string) []byte {
	e := ErrorResponse{
		Message:   message,
		ErrorCode: code,
	}
	return utils.MarshalIgnoreError(e)
}

func makeErrorHandlerFunc(statusCode int) func(w http.ResponseWriter, r *http.Request, message string) {
	return func(w http.ResponseWriter, r *http.Request, message string) {
		resp := makeErrorResponse(statusCode, message)

		w.WriteHeader(statusCode)
		w.Write(resp)
	}
}

func handleInvalidMissingQueryParm(w http.ResponseWriter, r *http.Request, missing string) {
	handleInvalidRequest(w, r, fmt.Sprintf(`Invalid request, query param "%s" is required`, missing))
}

func getQueryParam(req *http.Request, id string) string {
	return req.URL.Query().Get(id)
}

func getLatLongQueryParams(req *http.Request) (string, string) {
	lat := getQueryParam(req, "lat")
	long := getQueryParam(req, "long")

	return lat, long
}

func CheckISSVisible(position, ISSPosition model.Coordinates, weatherResult domain.WeatherResult, cloudCoverThreshold int, precision uint) bool {

	if weatherResult.CloudCover >= cloudCoverThreshold {
		return false
	}

	positionsMatch := MakeCoordinatesMatch(precision)

	return positionsMatch(position, ISSPosition)
}

func MakeCoordinatesMatch(precision uint) func(model.Coordinates, model.Coordinates) bool {
	positionsMatch := MakePositionMatch(precision)

	return func(a, b model.Coordinates) bool {
		if !positionsMatch(a.Latitude, b.Latitude) {
			return false
		}
		return positionsMatch(a.Longitude, b.Longitude)
	}
}

func MakePositionMatch(precision uint) func(a, b float64) bool {
	roundToPrecision := utils.MakeRoundToNPlaces(precision)
	return func(a, b float64) bool {
		return roundToPrecision(a) == roundToPrecision(b)
	}
}
