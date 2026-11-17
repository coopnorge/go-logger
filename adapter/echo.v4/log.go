package echo

import (
	"encoding/json"
	"io"

	log "github.com/coopnorge/go-logger"
	labstackecho "github.com/labstack/echo/v4"
	gommonlog "github.com/labstack/gommon/log"
)

// Ensure WrappedEchoLogger implements the Echo v4 Logger interface.
var _ labstackecho.Logger = (*WrappedEchoLogger)(nil)

// WrappedEchoLogger that can be passed to Echo middleware for Datadog integration
// implements Echo Logger from github.com/labstack/echo/v4/log.go
type WrappedEchoLogger struct {
	log *log.Logger

	// prefix for logs
	prefix string
	level  log.Level
}

// NewWrappedEchoLogger instance
func NewWrappedEchoLogger() *WrappedEchoLogger {
	return &WrappedEchoLogger{
		log:    log.New(log.WithLevel(log.LevelInfo)),
		prefix: "echo",
		level:  log.LevelInfo,
	}
}

// Output not supported to access in Coop logger, so it's returns only stub
func (wel *WrappedEchoLogger) Output() io.Writer {
	return wel.log.OutputHandler()
}

// SetOutput not supported to change in Coop logger, accepting only stub
func (wel *WrappedEchoLogger) SetOutput(_ io.Writer) {}

// Prefix returns the current log prefix
func (wel *WrappedEchoLogger) Prefix() string {
	return wel.prefix
}

// SetPrefix sets a string that will be prefixed to each log line
func (wel *WrappedEchoLogger) SetPrefix(p string) {
	wel.prefix = p
}

// Level returns a mapped Echo log level
func (wel *WrappedEchoLogger) Level() gommonlog.Lvl {
	switch wel.level {
	case log.LevelDebug:
		return gommonlog.DEBUG
	case log.LevelInfo:
		return gommonlog.INFO
	case log.LevelWarn:
		return gommonlog.WARN
	case log.LevelError:
		return gommonlog.ERROR
	case log.LevelFatal:
		return gommonlog.ERROR
	default:
		return gommonlog.OFF
	}
}

// SetLevel maps the log level from an Echo log level
func (wel *WrappedEchoLogger) SetLevel(v gommonlog.Lvl) {
	switch v {
	case gommonlog.DEBUG:
		wel.level = log.LevelDebug
	case gommonlog.INFO:
		wel.level = log.LevelInfo
	case gommonlog.WARN:
		wel.level = log.LevelWarn
	case gommonlog.ERROR:
		wel.level = log.LevelError
	case gommonlog.OFF:
		return // Ignore Coop logger cannot be disabled yet
	}

	log.SetLevel(wel.level)
}

// SetHeader not supported
func (wel *WrappedEchoLogger) SetHeader(_ string) {}

// Print logs at the info level
func (wel *WrappedEchoLogger) Print(i ...any) {
	wel.log.Info(i...)
}

// Printf logs formatted output at the info level
func (wel *WrappedEchoLogger) Printf(format string, args ...any) {
	wel.log.Infof(format, args...)
}

// Printj marshals a map to JSON and and logs it at the info level
func (wel *WrappedEchoLogger) Printj(j gommonlog.JSON) {
	wel.log.Info(wel.jsonToString(j))
}

// Debug logs at the debug level
func (wel *WrappedEchoLogger) Debug(i ...any) {
	wel.log.Debug(i...)
}

// Debugf logs formatted output at the debug level
func (wel *WrappedEchoLogger) Debugf(format string, args ...any) {
	wel.log.Debugf(format, args...)
}

// Debugj marshals a map to JSON and and logs it at the debug level
func (wel *WrappedEchoLogger) Debugj(j gommonlog.JSON) {
	wel.log.Debug(wel.jsonToString(j))
}

// Info logs at the info level
func (wel *WrappedEchoLogger) Info(i ...any) {
	wel.log.Info(i...)
}

// Infof logs formatted output at the info level
func (wel *WrappedEchoLogger) Infof(format string, args ...any) {
	wel.log.Infof(format, args...)
}

// Infoj marshals a map to JSON and and logs it at the info level
func (wel *WrappedEchoLogger) Infoj(j gommonlog.JSON) {
	wel.log.Info(wel.jsonToString(j))
}

// Warn logs at the warn level
func (wel *WrappedEchoLogger) Warn(i ...any) {
	wel.log.Warn(i...)
}

// Warnf logs formatted output at the warn level
func (wel *WrappedEchoLogger) Warnf(format string, args ...any) {
	wel.log.Warnf(format, args...)
}

// Warnj marshals a map to JSON and and logs it at the warn level
func (wel *WrappedEchoLogger) Warnj(j gommonlog.JSON) {
	wel.log.Warn(wel.jsonToString(j))
}

// Error logs at the error level
func (wel *WrappedEchoLogger) Error(i ...any) {
	wel.log.Error(i...)
}

// Errorf logs formatted output at the error level
func (wel *WrappedEchoLogger) Errorf(format string, args ...any) {
	wel.log.Errorf(format, args...)
}

// Errorj marshals a map to JSON and and logs it at the error level
func (wel *WrappedEchoLogger) Errorj(j gommonlog.JSON) {
	wel.log.Error(wel.jsonToString(j))
}

// Fatal logs fatally
func (wel *WrappedEchoLogger) Fatal(i ...any) {
	wel.log.Fatal(i...)
}

// Fatalf logs formatted output fatally
func (wel *WrappedEchoLogger) Fatalf(format string, args ...any) {
	wel.log.Fatalf(format, args...)
}

// Fatalj marshals a map to JSON and and logs it fatally
func (wel *WrappedEchoLogger) Fatalj(j gommonlog.JSON) {
	wel.log.Fatal(wel.jsonToString(j))
}

// Panic wraps a call to Fatal
func (wel *WrappedEchoLogger) Panic(i ...any) {
	wel.log.Fatal(i...)
}

// Panicf wraps a call to Fatalf
func (wel *WrappedEchoLogger) Panicf(format string, args ...any) {
	wel.log.Fatalf(format, args...)
}

// Panicj wraps a call to Fatalj
func (wel *WrappedEchoLogger) Panicj(j gommonlog.JSON) {
	wel.log.Fatal(wel.jsonToString(j))
}

func (wel *WrappedEchoLogger) jsonToString(j gommonlog.JSON) string {
	jb, marshalErr := json.Marshal(j)
	if marshalErr != nil {
		wel.log.Errorf("unable to marshal json data for logs, marshal error: :%v", marshalErr)
		return ""
	}

	return string(jb)
}
