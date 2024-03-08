package logs

import (
	coopLog "github.com/coopnorge/go-logger"
	"github.com/go-kratos/kratos/v2/log"
)

// Ensure LoggerKratosAdapter implements the log.Logger interface.
var _ log.Logger = (*LoggerKratosAdapter)(nil)

// LoggerKratosAdapter adapter for Go-Kratos.
type LoggerKratosAdapter struct {
	log *coopLog.Logger
}

// NewLoggerKratosAdapter constructor that accepting Coop logger adapter for Go-Kratos.
//
// Example:
//
//	package main
//
//	import (
//		"github.com/coopnorge/go-logger"
//		"github.com/go-kratos/kratos/v2/log"
//	)
//
//	func main() {
//		// Create a Coop logger.
//		logger.ConfigureGlobalLogger(logger.WithLevel(logger.LevelDebug), logger.WithHook(tracelogger.NewHook()))
//
//		// Create a LoggerKratosAdapter and pass it to Go-Kratos so it will know what adapter of logger to use.
//		log.SetLogger(logs.NewLoggerKratosAdapter(logger.Global()))
//	}
func NewLoggerKratosAdapter(coopLog *coopLog.Logger) *LoggerKratosAdapter {
	return &LoggerKratosAdapter{log: coopLog}
}

// Log prints the keyValPairs to the log.
//
// Example of logs:
//
// {"file":"/project/internal/pkg/logs/kratos.go:36","function":"xxx","level":"debug","msg":"msgconfig loaded: MY_ENV format: ","time":"2000-03-08T10:49:12Z"}
//
// {"file":"/project/internal/pkg/logs/kratos.go:34","function":"xxx","level":"info","msg":"I'm godoc example message","time":"2000-03-08T10:49:12Z"}
func (l *LoggerKratosAdapter) Log(level log.Level, keyValPairs ...interface{}) error {
	if len(keyValPairs) == 0 {
		return nil
	}

	switch level {
	case log.LevelFatal:
		coopLog.Fatal(keyValPairs...)
	case log.LevelError:
		coopLog.Error(keyValPairs...)
	case log.LevelWarn:
		coopLog.Warn(keyValPairs...)
	case log.LevelInfo:
		coopLog.Info(keyValPairs...)
	default:
		coopLog.Debug(keyValPairs...)
	}

	return nil
}
