package logger

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

// Entry represents a logging entry and all supported method we use
type Entry struct {
	logger    *Logger
	timestamp time.Time
	fields    Fields
	context   context.Context
}

// WithError is a convenience wrapper for WithField("error", err)
func (e *Entry) WithError(err error) *Entry {
	return e.WithField(errorKey, err)
}

// WithField forwards a logging call with a field
func (e *Entry) WithField(key string, value interface{}) *Entry {
	e.fields[key] = value
	return e
}

// WithFields forwards a logging call with fields
func (e *Entry) WithFields(fields Fields) *Entry {
	for k, v := range fields {
		e.fields[k] = v
	}
	return e
}

// WithContext sets the context for the log-message. Useful when using hooks.
func (e *Entry) WithContext(ctx context.Context) *Entry {
	e.context = ctx
	return e
}

// Info forwards a logging call in the (format, args) format
func (e *Entry) Info(args ...interface{}) {
	logrusFields := logrus.Fields(e.fields)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.timestamp).WithFields(logrusFields).Info(args...)
}

// Infof forwards a logging call in the (format, args) format
func (e *Entry) Infof(format string, args ...interface{}) {
	logrusFields := logrus.Fields(e.fields)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.timestamp).WithFields(logrusFields).Infof(format, args...)
}

// Error forwards an error logging call
func (e *Entry) Error(args ...interface{}) {
	logrusFields := logrus.Fields(e.fields)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.timestamp).WithFields(logrusFields).Error(args...)
}

// Errorf forwards an error logging call
func (e *Entry) Errorf(format string, args ...interface{}) {
	logrusFields := logrus.Fields(e.fields)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.timestamp).WithFields(logrusFields).Errorf(format, args...)
}

// Debug forwards a debugging logging call
func (e *Entry) Debug(args ...interface{}) {
	logrusFields := logrus.Fields(e.fields)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.timestamp).WithFields(logrusFields).Debug(args...)
}

// Debugf forwards a debugging logging call
func (e *Entry) Debugf(format string, args ...interface{}) {
	logrusFields := logrus.Fields(e.fields)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.timestamp).WithFields(logrusFields).Debugf(format, args...)
}

// Warn forwards a warning logging call
func (e *Entry) Warn(args ...interface{}) {
	logrusFields := logrus.Fields(e.fields)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.timestamp).WithFields(logrusFields).Warn(args...)
}

// Warnf forwards a warning logging call
func (e *Entry) Warnf(format string, args ...interface{}) {
	logrusFields := logrus.Fields(e.fields)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.timestamp).WithFields(logrusFields).Warnf(format, args...)
}

// Fatal forwards a fatal logging call
func (e *Entry) Fatal(args ...interface{}) {
	logrusFields := logrus.Fields(e.fields)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.timestamp).WithFields(logrusFields).Fatal(args...)
}

// Fatalf forwards a fatal logging call
func (e *Entry) Fatalf(format string, args ...interface{}) {
	logrusFields := logrus.Fields(e.fields)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.timestamp).WithFields(logrusFields).Fatalf(format, args...)
}
