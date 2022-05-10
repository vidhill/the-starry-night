package service

import "github.com/vidhill/the-starry-night/domain"

type ConfigService interface {
	GetString(string) string
	GetBool(string) bool
	GetInt(string) int
}

type DefaultConfigService struct {
	Repo domain.ConfigRepository
}

func (s DefaultConfigService) GetString(id string) string {
	return s.Repo.GetString(id)
}

func (s DefaultConfigService) GetBool(id string) bool {
	return s.Repo.GetBool(id)
}

func (s DefaultConfigService) GetInt(id string) int {
	return s.Repo.GetInt(id)
}

func NewConfigService(repository domain.ConfigRepository) ConfigService {
	return DefaultConfigService{
		Repo: repository,
	}
}
