package mocks

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type Http struct {
	mock.Mock
}

func (mock *Http) Get(url string) (*http.Response, error) {
	args := mock.Called(url)
	result := args.Get(0)
	err := args.Error(1)

	return result.(*http.Response), err
}

func NewMockHTTP() Http {
	return Http{}
}
