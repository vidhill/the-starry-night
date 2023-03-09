package domain

import "net/http"

type HttpProvider interface {
	Get(url string) (*http.Response, error)
}
