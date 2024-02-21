package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

// Entry represents a logging entry and all supported method we use
type Entry struct {
	logger  *Logger
	fields  Fields
	context context.Context
}

// WithError is a convenience wrapper for WithField("error", err)
func (e *Entry) WithError(err error) *Entry {
	return e.WithField(errorKey, err)
}

// WithField forwards a logging call with a field
func (e *Entry) WithField(key string, value interface{}) *Entry {
	return e.WithFields(Fields{key: value})
}

// WithFields forwards a logging call with fields
func (e *Entry) WithFields(fields Fields) *Entry {
	// Make a copy, to prevent mutation of the old entry
	newFields := make(Fields, len(e.fields)+len(fields))
	// Copy old fields
	for k, v := range e.fields {
		newFields[k] = v
	}
	// Set new fields
	for k, v := range fields {
		newFields[k] = v
	}
	return &Entry{logger: e.logger, fields: newFields, context: e.context}
}

// WithContext sets the context for the log-message. Useful when using hooks.
func (e *Entry) WithContext(ctx context.Context) *Entry {
	// Make a copy, to prevent mutation of the old entry
	return &Entry{logger: e.logger, fields: e.fields, context: ctx}
}

// Info forwards a logging call in the (format, args) format
func (e *Entry) Info(args ...interface{}) {
	logrusFields := logrus.Fields(e.fields)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.logger.now()).WithFields(logrusFields).Info(args...)
}

// Infof forwards a logging call in the (format, args) format
func (e *Entry) Infof(format string, args ...interface{}) {
	logrusFields := logrus.Fields(e.fields)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.logger.now()).WithFields(logrusFields).Infof(format, args...)
}

// Error forwards an error logging call
func (e *Entry) Error(args ...interface{}) {
	logrusFields := logrus.Fields(e.fields)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.logger.now()).WithFields(logrusFields).Error(args...)
}

// Errorf forwards an error logging call
func (e *Entry) Errorf(format string, args ...interface{}) {
	logrusFields := logrus.Fields(e.fields)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.logger.now()).WithFields(logrusFields).Errorf(format, args...)
}

// Debug forwards a debugging logging call
func (e *Entry) Debug(args ...interface{}) {
	logrusFields := logrus.Fields(e.fields)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.logger.now()).WithFields(logrusFields).Debug(args...)
}

// Debugf forwards a debugging logging call
func (e *Entry) Debugf(format string, args ...interface{}) {
	logrusFields := logrus.Fields(e.fields)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.logger.now()).WithFields(logrusFields).Debugf(format, args...)
}

// Warn forwards a warning logging call
func (e *Entry) Warn(args ...interface{}) {
	logrusFields := logrus.Fields(e.fields)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.logger.now()).WithFields(logrusFields).Warn(args...)
}

// Warnf forwards a warning logging call
func (e *Entry) Warnf(format string, args ...interface{}) {
	logrusFields := logrus.Fields(e.fields)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.logger.now()).WithFields(logrusFields).Warnf(format, args...)
}

// Fatal forwards a fatal logging call
func (e *Entry) Fatal(args ...interface{}) {
	logrusFields := logrus.Fields(e.fields)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.logger.now()).WithFields(logrusFields).Fatal(args...)
}

// Fatalf forwards a fatal logging call
func (e *Entry) Fatalf(format string, args ...interface{}) {
	logrusFields := logrus.Fields(e.fields)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.logger.now()).WithFields(logrusFields).Fatalf(format, args...)
}
