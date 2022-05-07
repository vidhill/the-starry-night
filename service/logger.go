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
	repo domain.LoggerRepository
}

func (l DefaultLoggerService) Debug(v ...interface{}) {
	l.repo.Debug(v...)
}

func (l DefaultLoggerService) Info(v ...interface{}) {
	l.repo.Info(v...)
}

func (l DefaultLoggerService) Warn(v ...interface{}) {
	l.repo.Warn(v...)
}

func (l DefaultLoggerService) Error(v ...interface{}) {
	l.repo.Error(v...)
}

func NewLoggerService(repository domain.LoggerRepository) LoggerService {
	return DefaultLoggerService{
		repo: repository,
	}
}
