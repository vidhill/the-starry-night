package stubrepository

import "net/http"

type StubHttp struct{}

func (mock StubHttp) Get(url string) (*http.Response, error) {
	return &http.Response{}, nil
}

func NewStubHttp() StubHttp {
	return StubHttp{}
}
