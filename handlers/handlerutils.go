package handlers

import (
	"fmt"
	"net/http"

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

// Composes multiple handler functions together
func ComposeHandlers(manyHandlers ...func(http.HandlerFunc) http.HandlerFunc) func(http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		for _, v := range manyHandlers {
			h = v(h)
		}
		return h
	}
}

// Middleware function, adds Content-Type of json to any response
func AddJsonHeader(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}
