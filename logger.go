package logger

import (
	"os"
	"sync"
)

var (
	logOnce     sync.Once
	errOnce     sync.Once
	logInstance *Logger
	errInstance *Logger
)

// StructuredLogger is an interface representing our logger
type StructuredLogger interface {
	LogWithFields(fields Fields) Entry
	LogPrint(args ...interface{})
	LogError(args ...interface{})
	LogTracef(format string, args ...interface{})
	LogDebugf(format string, args ...interface{})
	LogInfof(format string, args ...interface{})
	LogPrintf(format string, args ...interface{})
	LogWarnf(format string, args ...interface{})
	LogErrorf(format string, args ...interface{})
	LogFatalf(format string, args ...interface{})
	LogPanicf(format string, args ...interface{})
}

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

// Logger is our logger with the needed structured logger we use
type Logger struct {
	structuredLogger StructuredLogger
}

// Err returns a single isntance of a Logger using the Singleton Pattern with Gos sync.Once
func err() *Logger {
	errOnce.Do(func() {
		errInstance = &Logger{
			structuredLogger: Logrus(os.Stderr),
		}
	})

	return errInstance
}

// Log returns a single instance of a Logger using the Singleton Pattern with Gos sync.Once
func log() *Logger {
	logOnce.Do(func() {
		logInstance = &Logger{
			structuredLogger: Logrus(os.Stdout),
		}
	})

	return logInstance
}

// WithFields forwards a logging call with fields
func (*Logger) WithFields(fields Fields) Entry {
	return log().structuredLogger.LogWithFields(fields)
}

// Print forwards a standard print logging call
func (*Logger) Print(args ...interface{}) {
	log().structuredLogger.LogPrint(args...)
}

// Error forwards an error logging call
func (*Logger) Error(args ...interface{}) {
	err().structuredLogger.LogError(args...)
}

// Tracef forwards a tracing logging call
func (*Logger) Tracef(format string, args ...interface{}) {
	log().structuredLogger.LogTracef(format, args...)
}

// Debugf forwards a debugging logging call
func (*Logger) Debugf(format string, args ...interface{}) {
	log().structuredLogger.LogDebugf(format, args...)
}

// Infof forwards a logging call in the (format, args) format
func (*Logger) Infof(format string, args ...interface{}) {
	log().structuredLogger.LogInfof(format, args...)
}

// Printf forwards a standard printf logging call
func (*Logger) Printf(format string, args ...interface{}) {
	log().structuredLogger.LogPrintf(format, args...)
}

// Warnf forwards a warning logging call
func (*Logger) Warnf(format string, args ...interface{}) {
	err().structuredLogger.LogWarnf(format, args...)
}

// Errorf forwards an error logging call
func (*Logger) Errorf(format string, args ...interface{}) {
	err().structuredLogger.LogErrorf(format, args...)
}

// Fatalf forwards a fatal logging call
func (*Logger) Fatalf(format string, args ...interface{}) {
	err().structuredLogger.LogFatalf(format, args...)
}

// Panicf forwards a panic logging call
func (*Logger) Panicf(format string, args ...interface{}) {
	err().structuredLogger.LogPanicf(format, args...)
}
