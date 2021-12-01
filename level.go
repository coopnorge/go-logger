package logger

import "github.com/sirupsen/logrus"

type Level uint8

const (
	// LevelPanic is the highest log level, intended for unrecoverable errors that make the service unsuable. After logging, the app will be shut down.
	LevelPanic Level = iota
	// LevelFatal is to be used to log predictable errors that make the service unsuable, eg misconfiguration. After logging, the app will be shut down. This is true even if log level is set to Panic.
	LevelFatal
	// LevelError  isto be used for recoverable errors that limit the service's functionality, eg timeouts.
	LevelError
	// LevelWarn is to be used for non-critical errors that may require some attention.
	LevelWarn
	// LevelInfo is to be used for monitoring successful interactions, eg run time or job result.
	LevelInfo
	// LevelDebug should only be used in dev/test environments.
	LevelDebug
	// LevelTrace is to be used to log minute details. Should follow the same rules as LevelDebug.
	LevelTrace
)

func mapLevelToLogrusLevel(l Level) logrus.Level {
	return logrus.Level(l)
}
