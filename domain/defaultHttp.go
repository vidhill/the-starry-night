package domain

import "net/http"

//
// Wrapper around the standard http service
//

type DefaultHttpClient struct {
	logger LoggerRepository
}

func (s DefaultHttpClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

func NewDefaultHttpClient(logger LoggerRepository) DefaultHttpClient {
	return DefaultHttpClient{
		logger: logger,
	}
}
