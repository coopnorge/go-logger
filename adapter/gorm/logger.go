package gorm

import (
	"context"
	"errors"
	"time"

	coopLogger "github.com/coopnorge/go-logger"
	"gorm.io/gorm/logger"
)

// Logger is a logging adapter between Gorm an go-logger, do not create
// this directly, use NewLogger()
type Logger struct {
	instance     *coopLogger.Logger
	traceEnabled bool
}

// NewLogger creates a new Gorm logger that passes message to go-logger
//
// To inject the logger into Gorm use
//
//	package main
//
//	import (
//		gormLogger "github.com/coopnorge/go-logger/adapter/gorm"
//
//		"gorm.io/driver/postgres"
//		"gorm.io/gorm"
//	)
//
//	func main() {
//		l, err := gormLogger.NewLogger(gormLogger.WithGlobalLogger())
//		if err != nil {
//			panic(err)
//		}
//		db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
//			Logger: l,
//		})
//		if err != nil {
//			panic(err)
//		}
//	}
func NewLogger(opts ...LoggerOption) (*Logger, error) {
	logger := &Logger{
		traceEnabled: false,
	}
	for _, opt := range opts {
		opt.Apply(logger)
	}
	if logger.instance == nil {
		return nil, errors.New("no go-logger instance provided, use WithGlobalLogger() or WithLogger() to configure the logger")
	}
	return logger, nil
}

// LogMode sets the log level, this is ignored and filtering of logs is left to
// go-logger
func (l *Logger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

// Info writes info level statements to the log
func (l *Logger) Info(ctx context.Context, msg string, data ...any) {
	l.instance.WithContext(ctx).Infof(msg, data...)
}

// Warn writes warn level statements to the log
func (l *Logger) Warn(ctx context.Context, msg string, data ...any) {
	l.instance.WithContext(ctx).Warnf(msg, data...)
}

// Error writes error level statements to the log
func (l *Logger) Error(ctx context.Context, msg string, data ...any) {
	l.instance.WithContext(ctx).Errorf(msg, data...)
}

// Trace write SQL trace to the log
func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if !l.traceEnabled {
		return
	}
	elapsed := time.Since(begin)
	_, rows := fc() // sql is dropped for privacy, use tracing in Datadog instead.
	e := l.instance.WithContext(ctx).WithField("rows", rows).WithField("elapsed", elapsed)
	if err != nil {
		e.WithError(err).Error()
		return
	}
	e.Debug()
}
