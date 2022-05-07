package domain

import "log"

type DefaultLogger struct{}

func (s DefaultLogger) Debug(v ...interface{}) {
	log.Println(v...)
}

func (s DefaultLogger) Info(v ...interface{}) {
	log.Println(v...)
}

func (s DefaultLogger) Warn(v ...interface{}) {
	log.Println(v...)
}

func (s DefaultLogger) Error(v ...interface{}) {
	log.Println(v...)
}

func NewDefaultLogger() DefaultLogger {
	return DefaultLogger{}
}
