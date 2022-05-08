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

type ErrorResponse struct {
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
