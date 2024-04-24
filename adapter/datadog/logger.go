package datadog

import (
	"errors"
	"regexp"

	coopLogger "github.com/coopnorge/go-logger"
)

// Logger is a logging adapter between gopkg.in/DataDog/dd-trace-go.v1/ddtrace
// an go-logger, do not create this directly, use NewLogger()
type Logger struct {
	instance             *coopLogger.Logger
	msgPattern           *regexp.Regexp
	msgPatternGroupNames []string
}

// NewLogger creates a new Datadog logger that passes message to go-logger
//
// To inject the logger into ddtrace use
//
//	package main
//
//	import (
//		"github.com/coopnorge/go-logger/adapter/datadog"
//
//		"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
//	)
//
//	func main() {
//		l, err =: datadog.NewLogger(datadog.WithGlobalLogger())
//		ddtrace.UseLogger(l)
//	}
func NewLogger(opts ...LoggerOption) (*Logger, error) {
	msgPattern := regexp.MustCompile(`^(?P<source>Datadog\sTracer\sv\d+\.\d+\.\d+)\s(?P<level>(?:ERROR)|(?:WARN)|(?:INFO)|(?:DEBUG)):\s(?P<msg>.+)$`)
	logger := &Logger{
		msgPattern:           msgPattern,
		msgPatternGroupNames: msgPattern.SubexpNames(),
	}
	for _, opt := range opts {
		opt.Apply(logger)
	}
	if logger.instance == nil {
		return nil, errors.New("No go-logger instance provided, use WithGlobalLogger() or WithLogger() to configure the logger")
	}
	return logger, nil
}

// Log writes statements to the log
func (l *Logger) Log(msg string) {
	matchResult, err := l.parseMsg(msg)
	if err != nil {
		l.instance.WithField("log_parse_err", err).Warn(msg)
		return
	}

	e := l.instance.WithField("source", matchResult["source"])

	switch matchResult["level"] {
	case "ERROR":
		e.Error(matchResult["msg"])
	case "WARN":
		e.Warn(matchResult["msg"])
	case "INFO":
		e.Info(matchResult["msg"])
	case "DEBUG":
		e.Debug(matchResult["msg"])
	default:
		e.WithField("log_parse_err", "Could not resolve the level").WithField("parsed_level", matchResult["level"]).Warn(matchResult["msg"])
	}
}

func (l *Logger) parseMsg(msg string) (map[string]string, error) {
	matches := l.msgPattern.FindAllStringSubmatch(msg, -1)
	if len(matches) == 0 {
		return nil, errors.New("Datadog logger adapter could not match the log statement to the pattern, falling back to warning")
	}

	matchResult := map[string]string{}
	for i, n := range matches[0] {
		matchResult[l.msgPatternGroupNames[i]] = n
	}
	return matchResult, nil
}
