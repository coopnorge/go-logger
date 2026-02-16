package logger

import (
	"context"
	"io"
)

var globalLogger = New()

// ConfigureGlobalLogger applies supplied logger options to the global logger
func ConfigureGlobalLogger(opts ...LoggerOption) {
	globalLogger.applyOptions(opts...)
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
func Info(args ...any) {
	globalLogger.Info(args...)
}

// Infof uses global logger to log payload on "info" level
func Infof(format string, args ...any) {
	globalLogger.Infof(format, args...)
}

// Error uses global logger to log payload on "error" level
func Error(args ...any) {
	globalLogger.Error(args...)
}

// Errorf uses global logger to log payload on "error" level
func Errorf(format string, args ...any) {
	globalLogger.Errorf(format, args...)
}

// Debug uses global logger to log payload on "debug" level
func Debug(args ...any) {
	globalLogger.Debug(args...)
}

// Debugf uses global logger to log payload on "debug" level
func Debugf(format string, args ...any) {
	globalLogger.Debugf(format, args...)
}

// Warn uses global logger to log payload on "warn" level
func Warn(args ...any) {
	globalLogger.Warn(args...)
}

// Warnf uses global logger to log payload on "warn" level
func Warnf(format string, args ...any) {
	globalLogger.Warnf(format, args...)
}

// Fatal uses global logger to log payload on "fatal" level
func Fatal(args ...any) {
	globalLogger.Fatal(args...)
}

// Fatalf uses global logger to log payload on "fatal" level
func Fatalf(format string, args ...any) {
	globalLogger.Fatalf(format, args...)
}

// Log uses global logger to log payload on the level of the first parameter
func Log(level Level, args ...any) {
	globalLogger.Log(level, args...)
}

// Logf uses global logger to log payload on the level of the first parameter
func Logf(level Level, format string, args ...any) {
	globalLogger.Logf(level, format, args...)
}

// SetNowFunc sets `now` func user by global logger
func SetNowFunc(nowFunc NowFunc) {
	ConfigureGlobalLogger(WithNowFunc(nowFunc))
}

// SetOutput changes global logger's output
func SetOutput(output io.Writer) {
	ConfigureGlobalLogger(WithOutput(output))
}

// SetLevel sets minimum log level for global logger
func SetLevel(level Level) {
	ConfigureGlobalLogger(WithLevel(level))
}

// SetReportCaller allows controlling if caller info should be attached to logs by global logger
func SetReportCaller(enable bool) {
	ConfigureGlobalLogger(WithReportCaller(enable))
}
