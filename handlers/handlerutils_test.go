package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_makeErrorResponse(t *testing.T) {
	b := makeErrorResponse(http.StatusNotFound, "Not found")

	expected := `
	{
		"error_code": 404,
		"message": "Not found"
	}
	`

	assert.JSONEq(t, expected, string(b))
}
