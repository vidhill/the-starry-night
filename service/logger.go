package service

import (
	"fmt"

	"github.com/vidhill/the-starry-night/domain"
)

const (
	ERROR = iota
	WARN
	INFO
	DEBUG
	TRACE
)

type LoggerService interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
}

type DefaultLoggerService struct {
	Repo     domain.LoggerRepository
	Loglevel int
}

func (l DefaultLoggerService) Debug(v ...interface{}) {
	if l.Loglevel > DEBUG {
		l.Repo.Debug(v...)
	}
}

func (l DefaultLoggerService) Info(v ...interface{}) {
	if l.Loglevel > INFO {
		l.Repo.Info(v...)
	}
}

func (l DefaultLoggerService) Warn(v ...interface{}) {
	if l.Loglevel > WARN {
		l.Repo.Warn(v...)
	}
}

func (l DefaultLoggerService) Error(v ...interface{}) {
	if l.Loglevel > ERROR {
		l.Repo.Error(v...)
	}
}

// configuring log levels at service level so log level is independent of implementation
func NewLoggerService(repository domain.LoggerRepository, loglevel string) LoggerService {

	// translating string to int as mini perf optimisation as logs get called often
	levels := map[string]int{
		"DEBUG": DEBUG,
		"INFO":  INFO,
		"WARN":  WARN,
		"ERROR": ERROR,
	}

	level, ok := levels[loglevel]

	// fallback to info level if misconfigured
	if !ok {
		level = INFO
	}

	return DefaultLoggerService{
		Repo:     repository,
		Loglevel: level + 1,
	}
}
