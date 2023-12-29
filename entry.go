package logger

import (
	"context"
	"fmt"
	"os"
	"sort"

	"go.uber.org/zap/zapcore"
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
func (e *Entry) WithField(key string, value any) *Entry {
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
func (e *Entry) Info(msg string) {
	msg, level := e.fireHooks(msg, LevelInfo)
	zapFields := e.getZapFields()
	zapLevel := mapLevelToZapLevel(level)
	e.logger.zapLogger.Log(zapLevel, msg, zapFields...)
}

// Infof forwards a logging call in the (format, args) format
func (e *Entry) Infof(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	msg, level := e.fireHooks(msg, LevelInfo)
	zapFields := e.getZapFields()
	zapLevel := mapLevelToZapLevel(level)
	e.logger.zapLogger.Log(zapLevel, msg, zapFields...)
}

// Error forwards an error logging call
func (e *Entry) Error(msg string) {
	msg, level := e.fireHooks(msg, LevelError)
	zapFields := e.getZapFields()
	zapLevel := mapLevelToZapLevel(level)
	e.logger.zapLogger.Log(zapLevel, msg, zapFields...)
}

// Errorf forwards an error logging call
func (e *Entry) Errorf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	msg, level := e.fireHooks(msg, LevelError)
	zapFields := e.getZapFields()
	zapLevel := mapLevelToZapLevel(level)
	e.logger.zapLogger.Log(zapLevel, msg, zapFields...)
}

// Debug forwards a debugging logging call
func (e *Entry) Debug(msg string) {
	msg, level := e.fireHooks(msg, LevelDebug)
	zapFields := e.getZapFields()
	zapLevel := mapLevelToZapLevel(level)
	e.logger.zapLogger.Log(zapLevel, msg, zapFields...)
}

// Debugf forwards a debugging logging call
func (e *Entry) Debugf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	msg, level := e.fireHooks(msg, LevelDebug)
	zapFields := e.getZapFields()
	zapLevel := mapLevelToZapLevel(level)
	e.logger.zapLogger.Log(zapLevel, msg, zapFields...)
}

// Warn forwards a warning logging call
func (e *Entry) Warn(msg string) {
	msg, level := e.fireHooks(msg, LevelWarn)
	zapFields := e.getZapFields()
	zapLevel := mapLevelToZapLevel(level)
	e.logger.zapLogger.Log(zapLevel, msg, zapFields...)
}

// Warnf forwards a warning logging call
func (e *Entry) Warnf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	msg, level := e.fireHooks(msg, LevelWarn)
	zapFields := e.getZapFields()
	zapLevel := mapLevelToZapLevel(level)
	e.logger.zapLogger.Log(zapLevel, msg, zapFields...)
}

// Fatal forwards a fatal logging call
func (e *Entry) Fatal(msg string) {
	msg, level := e.fireHooks(msg, LevelFatal)
	zapFields := e.getZapFields()
	zapLevel := mapLevelToZapLevel(level)
	e.logger.zapLogger.Log(zapLevel, msg, zapFields...)
}

// Fatalf forwards a fatal logging call
func (e *Entry) Fatalf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	msg, level := e.fireHooks(msg, LevelFatal)
	zapFields := e.getZapFields()
	zapLevel := mapLevelToZapLevel(level)
	e.logger.zapLogger.Log(zapLevel, msg, zapFields...)
}

func (e *Entry) fireHooks(msg string, level Level) (string, Level) {
	if len(e.logger.hooks) > 0 {
		for _, hook := range e.logger.hooks {
			he := &HookEntry{
				Fields:  e.fields,
				Level:   level,
				Message: msg,
				Context: e.context,
			}
			changed, err := hook.Fire(he)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to fire hook: %v\n", err)
				continue
			}
			if !changed {
				continue
			}
			e.fields = he.Fields
			level = he.Level
			msg = he.Message
			e.context = he.Context

		}
	}
	return msg, level
}

func (e *Entry) getZapFields() []zapcore.Field {
	res := make([]zapcore.Field, 0, len(e.fields)+1)
	for k, v := range e.fields {
		f := zapcore.Field{
			Key: k,
		}
		switch v := v.(type) {
		case error:
			f.Type = zapcore.ErrorType
			f.Interface = v
		case string:
			f.Type = zapcore.StringType
			f.String = v
		default:
			f.Type = zapcore.ReflectType
			f.Interface = v
		}
		res = append(res, f)
	}
	if e.logger.attemptConsistentOrdering {
		sort.Slice(res, func(i, j int) bool {
			return res[i].Key < res[j].Key
		})
	}
	return res
}
