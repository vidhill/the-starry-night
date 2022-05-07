package handlers

import (
	"net/http"

	"github.com/vidhill/the-starry-night/utils"
)

func handleInternalServerError(w http.ResponseWriter, r *http.Request) {
	statusCode := http.StatusInternalServerError
	resp := makeErrorResponse(statusCode, "Internal server error")

	w.WriteHeader(statusCode)
	w.Write(resp)
}

func makeErrorResponse(code int, message string) []byte {
	e := ErrorResponse{
		Message:   message,
		ErrorCode: code,
	}
	return utils.MarshalIgnoreError(e)
}
