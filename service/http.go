package service

import (
	"net/http"

	"github.com/vidhill/the-starry-night/domain"
)

type HttpService interface {
	Get(url string) (*http.Response, error)
}

type DefaultHttpService struct {
	Repo domain.HttpRepository
}

func (s DefaultHttpService) Get(url string) (*http.Response, error) {
	return s.Repo.Get(url)
}

func NewHttpService(repository domain.HttpRepository) HttpService {
	return DefaultHttpService{
		Repo: repository,
	}
}
