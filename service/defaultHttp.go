package service

import (
	"net/http"

	"github.com/vidhill/the-starry-night/domain"
)

//
// Wrapper around the standard http service
//

type DefaultHttpClient struct {
	Logger domain.LogProvider
}

func (s DefaultHttpClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

func NewDefaultHttpClient(logger domain.LogProvider) domain.HttpProvider {
	return DefaultHttpClient{
		Logger: logger,
	}
}
