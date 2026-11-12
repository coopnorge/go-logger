package goose

import (
	log "github.com/coopnorge/go-logger"
)

// Logger is a logging adapter between Goose and go-logger
type Logger struct {
	log *log.Logger
}

// New creates a new GooseLogger instance
func New(logger *log.Logger) *Logger {
	return &Logger{
		log: logger,
	}
}

// Printf writes info level statements to the log (goose Logger interface).
func (l *Logger) Printf(format string, v ...any) {
	l.log.Infof(format, v...)
}

// Fatalf writes fatal level statements to the log (goose Logger interface).
func (l *Logger) Fatalf(format string, v ...any) {
	l.log.Fatalf(format, v...)
}
