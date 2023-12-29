package logger

import (
	"go.uber.org/zap/zapcore"
)

// Level is an integer representation of the logging level
type Level uint8

const (
	// LevelFatal is to be used to log predictable errors that make the service unusable, eg misconfiguration. After logging, the app will be shut down.
	LevelFatal = iota
	// LevelError  isto be used for recoverable errors that limit the service's functionality, eg timeouts.
	LevelError
	// LevelWarn is to be used for non-critical errors that may require some attention.
	LevelWarn
	// LevelInfo is to be used for monitoring successful interactions, eg run time or job result.
	LevelInfo
	// LevelDebug should only be used in dev/test environments.
	LevelDebug
)

func mapLevelToZapLevel(l Level) zapcore.Level {
	// Note that zapcore also has DPanicLevel and Panic level, but we don't expose those.
	switch l {
	case LevelFatal:
		return zapcore.FatalLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelDebug:
		return zapcore.DebugLevel
	}
	// should never get here
	return zapcore.DebugLevel
}
