package domain

import "net/http"

type HttpRepository interface {
	Get(url string) (*http.Response, error)
}
