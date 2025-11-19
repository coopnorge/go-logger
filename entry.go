package logger

import (
	"context"
	"fmt"

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

func addCallerFields(logrusFields logrus.Fields, reportCaller bool) {
	if reportCaller {
		frame := getCaller()
		logrusFields["file"] = fmt.Sprintf("%s:%v", frame.File, frame.Line)
		logrusFields["function"] = frame.Function
	}
}

// Info forwards a logging call in the (format, args) format
func (e *Entry) Info(args ...interface{}) {
	e.Log(LevelInfo, args...)
}

// Infof forwards a logging call in the (format, args) format
func (e *Entry) Infof(format string, args ...interface{}) {
	e.Logf(LevelInfo, format, args...)
}

// Error forwards an error logging call
func (e *Entry) Error(args ...interface{}) {
	e.Log(LevelError, args...)
}

// Errorf forwards an error logging call
func (e *Entry) Errorf(format string, args ...interface{}) {
	e.Logf(LevelError, format, args...)
}

// Debug forwards a debugging logging call
func (e *Entry) Debug(args ...interface{}) {
	e.Log(LevelDebug, args...)
}

// Debugf forwards a debugging logging call
func (e *Entry) Debugf(format string, args ...interface{}) {
	e.Logf(LevelDebug, format, args...)
}

// Warn forwards a warning logging call
func (e *Entry) Warn(args ...interface{}) {
	e.Log(LevelWarn, args...)
}

// Warnf forwards a warning logging call
func (e *Entry) Warnf(format string, args ...interface{}) {
	e.Logf(LevelWarn, format, args...)
}

// Fatal forwards a fatal logging call
func (e *Entry) Fatal(args ...interface{}) {
	e.Log(LevelFatal, args...)
}

// Fatalf forwards a fatal logging call
func (e *Entry) Fatalf(format string, args ...interface{}) {
	e.Logf(LevelFatal, format, args...)
}

// Logf forwards a logging call
func (e *Entry) Logf(level Level, format string, args ...any) {
	logrusFields := logrus.Fields(e.fields)
	addCallerFields(logrusFields, e.logger.reportCaller)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.logger.now()).WithFields(logrusFields).Logf(mapLevelToLogrusLevel(level), format, args...)

	// This ensures that logging with level Fatal results in Exit regardless if using .Fatalf or .Logf(LevelFatal, ...)
	if level == LevelFatal {
		e.logger.logrusLogger.Exit(1)
	}
}

// Log forwards a logging call
func (e *Entry) Log(level Level, args ...any) {
	logrusFields := logrus.Fields(e.fields)
	addCallerFields(logrusFields, e.logger.reportCaller)
	e.logger.logrusLogger.WithContext(e.context).WithTime(e.logger.now()).WithFields(logrusFields).Log(mapLevelToLogrusLevel(level), args...)

	// This ensures that logging with level Fatal results in Exit regardless if using .Fatal or .Log(LevelFatal, ...)
	if level == LevelFatal {
		e.logger.logrusLogger.Exit(1)
	}
}
