package labstack_logger

import (
	"encoding/json"
	"io"

	log "github.com/coopnorge/go-logger"
	echo "github.com/labstack/gommon/log"
)

// WrappedEchoLogger that can be passed to Echo middleware for Datadog integration
// implements Echo Logger from vendor/github.com/labstack/echo/v4/log.go
type WrappedEchoLogger struct {
	log    log.Entry
	output io.Writer

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
	return wel.output
}

// SetOutput not supported to change in Coop logger, accepting only stub
func (wel *WrappedEchoLogger) SetOutput(w io.Writer) {
	wel.output = w
}

func (wel *WrappedEchoLogger) Prefix() string {
	return wel.prefix
}

func (wel *WrappedEchoLogger) SetPrefix(p string) {
	wel.prefix = p
}

func (wel *WrappedEchoLogger) Level() echo.Lvl {
	switch wel.level {
	case log.LevelDebug:
		return echo.DEBUG
	case log.LevelInfo:
		return echo.INFO
	case log.LevelWarn:
		return echo.WARN
	case log.LevelError:
		return echo.ERROR
	case log.LevelFatal:
		return echo.ERROR
	default:
		return echo.OFF
	}
}

func (wel *WrappedEchoLogger) SetLevel(v echo.Lvl) {
	switch v {
	case echo.DEBUG:
		wel.level = log.LevelDebug
	case echo.INFO:
		wel.level = log.LevelInfo
	case echo.WARN:
		wel.level = log.LevelWarn
	case echo.ERROR:
		wel.level = log.LevelError
	case echo.OFF:
		return // Ignore Coop logger cannot be disabled yet
	}

	log.SetLevel(wel.level)
}

// SetHeader not supported
func (wel *WrappedEchoLogger) SetHeader(_ string) {}

func (wel *WrappedEchoLogger) Print(i ...interface{}) {
	wel.log.Info(i...)
}

func (wel *WrappedEchoLogger) Printf(format string, args ...interface{}) {
	wel.log.Infof(format, args...)
}

func (wel *WrappedEchoLogger) Printj(j echo.JSON) {
	wel.log.Info(wel.jsonToString(j))
}

func (wel *WrappedEchoLogger) Debug(i ...interface{}) {
	wel.log.Debug(i...)
}

func (wel *WrappedEchoLogger) Debugf(format string, args ...interface{}) {
	wel.log.Debugf(format, args...)
}

func (wel *WrappedEchoLogger) Debugj(j echo.JSON) {
	wel.log.Debug(wel.jsonToString(j))
}

func (wel *WrappedEchoLogger) Info(i ...interface{}) {
	wel.log.Info(i...)
}

func (wel *WrappedEchoLogger) Infof(format string, args ...interface{}) {
	wel.log.Infof(format, args...)
}

func (wel *WrappedEchoLogger) Infoj(j echo.JSON) {
	wel.log.Info(wel.jsonToString(j))
}

func (wel *WrappedEchoLogger) Warn(i ...interface{}) {
	wel.log.Warn(i...)
}

func (wel *WrappedEchoLogger) Warnf(format string, args ...interface{}) {
	wel.log.Warnf(format, args...)
}

func (wel *WrappedEchoLogger) Warnj(j echo.JSON) {
	wel.log.Warn(wel.jsonToString(j))
}

func (wel *WrappedEchoLogger) Error(i ...interface{}) {
	wel.log.Error(i...)
}

func (wel *WrappedEchoLogger) Errorf(format string, args ...interface{}) {
	wel.log.Errorf(format, args...)
}

func (wel *WrappedEchoLogger) Errorj(j echo.JSON) {
	wel.log.Error(wel.jsonToString(j))
}

func (wel *WrappedEchoLogger) Fatal(i ...interface{}) {
	wel.log.Fatal(i...)
}

func (wel *WrappedEchoLogger) Fatalj(j echo.JSON) {
	wel.log.Fatal(wel.jsonToString(j))
}

func (wel *WrappedEchoLogger) Fatalf(format string, args ...interface{}) {
	wel.log.Fatalf(format, args...)
}

func (wel *WrappedEchoLogger) Panic(i ...interface{}) {
	wel.log.Fatal(i...)
}

func (wel *WrappedEchoLogger) Panicj(j echo.JSON) {
	wel.log.Fatal(wel.jsonToString(j))
}

func (wel *WrappedEchoLogger) Panicf(format string, args ...interface{}) {
	wel.log.Fatalf(format, args...)
}

func (wel *WrappedEchoLogger) jsonToString(j echo.JSON) string {
	jb, marshalErr := json.Marshal(j)
	if marshalErr != nil {
		wel.log.Errorf("unable to marshal json data for logs, marshal error: :%v", marshalErr)
		return ""
	}

	return string(jb)
}
