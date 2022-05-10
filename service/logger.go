package service

import (
	"github.com/vidhill/the-starry-night/domain"
)

type LoggerService interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
}

type DefaultLoggerService struct {
	Repo domain.LoggerRepository
}

func (l DefaultLoggerService) Debug(v ...interface{}) {
	l.Repo.Debug(v...)
}

func (l DefaultLoggerService) Info(v ...interface{}) {
	l.Repo.Info(v...)
}

func (l DefaultLoggerService) Warn(v ...interface{}) {
	l.Repo.Warn(v...)
}

func (l DefaultLoggerService) Error(v ...interface{}) {
	l.Repo.Error(v...)
}

func NewLoggerService(repository domain.LoggerRepository) LoggerService {
	return DefaultLoggerService{
		Repo: repository,
	}
}
