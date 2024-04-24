package datadog

import (
	coopLogger "github.com/coopnorge/go-logger"
)

// LoggerOption defines an applicator interface
type LoggerOption interface { //nolint:all
	Apply(l *Logger)
}

// LoggerOptionFunc defines a function which modifies a logger
type LoggerOptionFunc func(l *Logger) //nolint:all

// Apply redirects a function call to the function receiver
func (lof LoggerOptionFunc) Apply(l *Logger) {
	lof(l)
}

// WithGlobalLogger configures Grom to use our global logger
func WithGlobalLogger() LoggerOption {
	return LoggerOptionFunc(func(l *Logger) {
		l.instance = coopLogger.Global()
	})
}

// WithLogger configures Grom to use a logger instance
func WithLogger(logger *coopLogger.Logger) LoggerOption {
	return LoggerOptionFunc(func(l *Logger) {
		l.instance = logger
	})
}
