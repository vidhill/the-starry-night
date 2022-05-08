package domain

import "net/http"

type DefaultHttpClient struct{}

func (s DefaultHttpClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}
