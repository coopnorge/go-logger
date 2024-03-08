package logs

import (
	"github.com/coopnorge/go-logger"
	"github.com/go-kratos/kratos/v2/log"
)

// Ensure LoggerKratosAdapter implements the log.Logger interface.
var _ log.Logger = (*LoggerKratosAdapter)(nil)

// LoggerKratosAdapter Adapter for Go-Kratos.
type LoggerKratosAdapter struct {
	log *logger.Logger
}

// NewLoggerKratosAdapter Coop logger adapter for Go-Kratos.
func NewLoggerKratosAdapter(coopLog *logger.Logger) *LoggerKratosAdapter {
	return &LoggerKratosAdapter{log: coopLog}
}

// Log print the keyValPairs to log.
func (l *LoggerKratosAdapter) Log(level log.Level, keyValPairs ...interface{}) error {
	if len(keyValPairs) == 0 {
		return nil
	}

	switch level {
	case log.LevelFatal:
		logger.Fatal(keyValPairs...)
	case log.LevelError:
		logger.Error(keyValPairs...)
	case log.LevelWarn:
		logger.Warn(keyValPairs...)
	case log.LevelInfo:
		logger.Info(keyValPairs...)
	default:
		logger.Debug(keyValPairs...)
	}

	return nil
}
