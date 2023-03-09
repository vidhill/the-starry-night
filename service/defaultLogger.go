package service

import (
	"log"
	"os"

	"github.com/vidhill/the-starry-night/domain"
)

const (
	ERROR = iota
	WARN
	INFO
	DEBUG
	TRACE
)

//
// Wrapper around the standard lib logger; log
//

const flags = log.LstdFlags

type StandardLogger struct {
	logLevel    int
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
}

func (l StandardLogger) Debug(v ...interface{}) {
	if l.logLevel > DEBUG {
		l.DebugLogger.Println(v...)
	}
}

func (l StandardLogger) Info(v ...interface{}) {
	if l.logLevel > INFO {
		l.InfoLogger.Println(v...)
	}
}

func (l StandardLogger) Warn(v ...interface{}) {
	if l.logLevel > WARN {
		l.WarnLogger.Println(v...)
	}
}

func (l StandardLogger) Error(v ...interface{}) {
	if l.logLevel > ERROR {
		l.ErrorLogger.Println(v...)
	}
}

func makeLogger(prefix string) *log.Logger {
	prefixWithSeparator := prefix + ": "
	return log.New(os.Stdout, prefixWithSeparator, flags)
}

func NewStandardLogger(c domain.ConfigProvider) domain.LogProvider {

	loglevel := c.GetString("LOG_LEVEL")

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

	return StandardLogger{
		logLevel:    level + 1,
		DebugLogger: makeLogger("DEBUG"),
		InfoLogger:  makeLogger("INFO"),
		WarnLogger:  makeLogger("WARN"),
		ErrorLogger: makeLogger("ERROR"),
	}
}
