package service

import "github.com/vidhill/the-starry-night/domain"

type ConfigService interface {
	GetString(string) string
	GetBool(string) bool
}

type DefaultConfigService struct {
	repo domain.ConfigRepository
}

func (s DefaultConfigService) GetString(id string) string {
	return s.repo.GetString(id)
}

func (s DefaultConfigService) GetBool(id string) bool {
	return s.repo.GetBool(id)
}

func NewConfigService(repository domain.ConfigRepository) DefaultConfigService {
	return DefaultConfigService{
		repo: repository,
	}
}
