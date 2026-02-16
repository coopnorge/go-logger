package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Fields type, used to pass to `WithFields`.
type Fields map[string]any

// NowFunc is a typedef for a function which returns the current time
type NowFunc func() time.Time

// Logger is our logger with the needed structured logger we use
type Logger struct {
	logrusLogger *logrus.Logger
	now          NowFunc
	output       io.Writer
	level        Level
	reportCaller bool
}

func (logger *Logger) applyOptions(opts ...LoggerOption) {
	for _, opt := range opts {
		opt.Apply(logger)
	}
	logger.logrusLogger.SetOutput(logger.output)
	logger.logrusLogger.SetLevel(mapLevelToLogrusLevel(logger.level))
}

// New creates and returns a new logger with supplied options
func New(opts ...LoggerOption) *Logger {
	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999Z07:00", // Adds milliseconds to the time-output
	})
	logger := &Logger{
		logrusLogger: logrusLogger,
		now:          NowFunc(time.Now),
		output:       os.Stdout,
		level:        LevelWarn,
		reportCaller: true,
	}
	logger.applyOptions(opts...)
	return logger
}

func (logger *Logger) entry() *Entry {
	return &Entry{logger: logger, fields: Fields{}}
}

const errorKey = "error"

// WithError is a convenience wrapper for WithField("error", err)
func (logger *Logger) WithError(err error) *Entry {
	return logger.WithField(errorKey, err)
}

// WithField forwards a logging call with a field
func (logger *Logger) WithField(key string, value any) *Entry {
	return logger.entry().WithField(key, value)
}

// WithFields forwards a logging call with fields
func (logger *Logger) WithFields(fields Fields) *Entry {
	return logger.entry().WithFields(fields)
}

// WithContext forwards a logging call with fields
func (logger *Logger) WithContext(ctx context.Context) *Entry {
	return logger.entry().WithContext(ctx)
}

// OutputHandler returns logger output handler
func (logger *Logger) OutputHandler() io.Writer {
	return logger.output
}

// Info forwards a logging call in the (format, args) format
func (logger *Logger) Info(args ...any) {
	logger.entry().Info(args...)
}

// Infof forwards a logging call in the (format, args) format
func (logger *Logger) Infof(format string, args ...any) {
	logger.entry().Infof(format, args...)
}

// Error forwards an error logging call
func (logger *Logger) Error(args ...any) {
	logger.entry().Error(args...)
}

// Errorf forwards an error logging call
func (logger *Logger) Errorf(format string, args ...any) {
	logger.entry().Errorf(format, args...)
}

// Debug forwards a debugging logging call
func (logger *Logger) Debug(args ...any) {
	logger.entry().Debug(args...)
}

// Debugf forwards a debugging logging call
func (logger *Logger) Debugf(format string, args ...any) {
	logger.entry().Debugf(format, args...)
}

// Warn forwards a warning logging call
func (logger *Logger) Warn(args ...any) {
	logger.entry().Warn(args...)
}

// Warnf forwards a warning logging call
func (logger *Logger) Warnf(format string, args ...any) {
	logger.entry().Warnf(format, args...)
}

// Fatal forwards a fatal logging call
func (logger *Logger) Fatal(args ...any) {
	logger.entry().Fatal(args...)
}

// Fatalf forwards a fatal logging call
func (logger *Logger) Fatalf(format string, args ...any) {
	logger.entry().Fatalf(format, args...)
}

// Log forwards a logging call
func (logger *Logger) Log(level Level, args ...any) {
	logger.entry().Log(level, args...)
}

// Logf forwards a logging call
func (logger *Logger) Logf(level Level, format string, args ...any) {
	logger.entry().Logf(level, format, args...)
}
