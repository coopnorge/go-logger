package logger

import (
	"strings"

	"github.com/sirupsen/logrus"
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

var nameMapping = map[string]Level{
	"fatal": LevelFatal,
	"error": LevelError,
	"warn":  LevelWarn,
	"info":  LevelInfo,
	"debug": LevelDebug,
}

// LevelNameToLevel converts a named log level to the Level type
func LevelNameToLevel(name string) (l Level, ok bool) {
	l, ok = nameMapping[strings.ToLower(name)]
	return
}

func mapLevelToLogrusLevel(l Level) logrus.Level {
	switch l {
	case LevelFatal:
		return logrus.FatalLevel
	case LevelError:
		return logrus.ErrorLevel
	case LevelWarn:
		return logrus.WarnLevel
	case LevelInfo:
		return logrus.InfoLevel
	case LevelDebug:
		return logrus.DebugLevel
	}
	// should never get here
	return logrus.DebugLevel
}

func mapLogrusLevelToLevel(l logrus.Level) Level {
	switch l {
	case logrus.FatalLevel:
		return LevelFatal
	case logrus.ErrorLevel:
		return LevelError
	case logrus.WarnLevel:
		return LevelWarn
	case logrus.InfoLevel:
		return LevelInfo
	case logrus.DebugLevel:
		return LevelDebug
	}
	// should never get here
	return LevelDebug
}
