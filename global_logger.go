package logger

import (
	"context"
)

var globalLogger = New()

func SetGlobalLogger(logger *Logger) {
	globalLogger = logger
}

// NewGlobalLogger is a shortcut for SetGlobalLogger(New(...))
func NewGlobalLogger(opts ...LoggerOption) *Logger {
	logger := New(opts...)
	SetGlobalLogger(logger)
	return logger
}

// Global logger that can be accessed without prior instantiation
func Global() *Logger {
	return globalLogger
}

// WithError is a convenience wrapper for WithField("err", err)
func WithError(err error) *Entry {
	return globalLogger.WithError(err)
}

// WithField creates log entry using global logger
func WithField(key string, value any) *Entry {
	return globalLogger.WithField(key, value)
}

// WithFields creates log entry using global logger
func WithFields(fields Fields) *Entry {
	return globalLogger.WithFields(fields)
}

// WithContext creates log entry using global logger
func WithContext(ctx context.Context) *Entry {
	return globalLogger.WithContext(ctx)
}

// Info uses global logger to log payload on "info" level
func Info(msg string) {
	globalLogger.Info(msg)
}

// Infof uses global logger to log payload on "info" level
func Infof(format string, args ...any) {
	globalLogger.Infof(format, args...)
}

// Error uses global logger to log payload on "error" level
func Error(msg string) {
	globalLogger.Error(msg)
}

// Errorf uses global logger to log payload on "error" level
func Errorf(format string, args ...any) {
	globalLogger.Errorf(format, args...)
}

// Debug uses global logger to log payload on "debug" level
func Debug(msg string) {
	globalLogger.Debug(msg)
}

// Debugf uses global logger to log payload on "debug" level
func Debugf(format string, args ...any) {
	globalLogger.Debugf(format, args...)
}

// Warn uses global logger to log payload on "warn" level
func Warn(msg string) {
	globalLogger.Warn(msg)
}

// Warnf uses global logger to log payload on "warn" level
func Warnf(format string, args ...any) {
	globalLogger.Warnf(format, args...)
}

// Fatal uses global logger to log payload on "fatal" level
func Fatal(msg string) {
	globalLogger.Fatal(msg)
}

// Fatalf uses global logger to log payload on "fatal" level
func Fatalf(format string, args ...any) {
	globalLogger.Fatalf(format, args...)
}
