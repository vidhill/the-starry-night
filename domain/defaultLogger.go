package domain

import (
	"log"
	"os"
)

//
// Wrapper around the standard lib logger; log
//

const flags = log.LstdFlags

type StandardLogger struct {
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
}

func (l StandardLogger) Debug(v ...interface{}) {
	l.DebugLogger.Println(v...)
}

func (l StandardLogger) Info(v ...interface{}) {
	l.InfoLogger.Println(v...)
}

func (l StandardLogger) Warn(v ...interface{}) {
	l.WarnLogger.Println(v...)
}

func (l StandardLogger) Error(v ...interface{}) {
	l.ErrorLogger.Println(v...)
}

func makeLogger(prefix string) *log.Logger {
	prefixWithSeparator := prefix + ": "
	return log.New(os.Stdout, prefixWithSeparator, flags)
}

func NewStandardLogger() StandardLogger {
	return StandardLogger{
		DebugLogger: makeLogger("DEBUG"),
		InfoLogger:  makeLogger("INFO"),
		WarnLogger:  makeLogger("WARN"),
		ErrorLogger: makeLogger("ERROR"),
	}
}
