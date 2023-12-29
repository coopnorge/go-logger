// All examples disable reporing caller as file path depends on workdir
// which is different for every run environment and therefore it's impossible
// to match exact output.
package logger

import (
	"errors"
	"os"
	"time"
)

func mockNowFunc() time.Time {
	return time.Date(2020, 10, 10, 10, 10, 10, 10, time.UTC)
}

func ExampleInfo() {
	logger := New(WithNowFunc(mockNowFunc), WithReportCaller(false))
	logger.Warn("foobar")
	logger.Info("i won't be logged because the default log level is higher than info")
	logger.Error("foobar")
	// Output: {"level":"warning","time":"2020-10-10T10:10:10Z","msg":"foobar"}
	// {"level":"error","time":"2020-10-10T10:10:10Z","msg":"foobar"}
}

func ExampleWithLevel() {
	logger := New(WithNowFunc(mockNowFunc), WithLevel(LevelInfo), WithReportCaller(false))
	logger.Info("now log level is set to info or lower, I will be logged")
	// Output: {"level":"info","time":"2020-10-10T10:10:10Z","msg":"now log level is set to info or lower, I will be logged"}
}

func ExampleWithFields() {
	// Example runner replaces os.Stdout to catch output and compares it with desired output.
	// Global logger instance sets default output to os.Stdin before example runner has a chance to overwrite it.
	// Within this function, os.Stdin is already replaced by example runner.
	// So we need to tell the global logger instance to use modified os.Stdout as its output,
	// otherwise example runner will fail as logs would be written to real stdout
	// and nothing would get written to example runnner's buffer.
	oldOutput := globalLogger.output
	oldNowFunc := globalLogger.now
	defer func() {
		NewGlobalLogger(WithOutput(oldOutput), WithNowFunc(oldNowFunc), WithReportCaller(true), withAttemptedConsistentOrdering(false))
	}()
	NewGlobalLogger(WithOutput(os.Stdout), WithNowFunc(mockNowFunc), WithReportCaller(false), withAttemptedConsistentOrdering(true))

	WithFields(Fields{
		"timeSpentOnConfiguration": 0,
		"defaultsLoaded":           true,
	}).Warn("use default logger with 0 configuration")
	WithField("singleField", true).Warn("example with a single field")
	// Output: {"level":"warning","time":"2020-10-10T10:10:10Z","msg":"use default logger with 0 configuration","defaultsLoaded":true,"timeSpentOnConfiguration":0}
	// {"level":"warning","time":"2020-10-10T10:10:10Z","msg":"example with a single field","singleField":true}
}

func ExampleWithError() {
	logger := New(WithNowFunc(mockNowFunc), WithReportCaller(false))

	err := errors.New("Test error")
	logger.WithError(err).Error("Operation failed")
	// Output: {"level":"error","time":"2020-10-10T10:10:10Z","msg":"Operation failed","error":"Test error"}
}

type warner interface {
	Warn(msg string)
}

func funcThatAcceptsInterface(warner warner) {
	warner.Warn("foobar")
}

func ExampleGlobal() {
	oldOutput := globalLogger.output
	oldNowFunc := globalLogger.now
	defer func() {
		NewGlobalLogger(WithOutput(oldOutput), WithNowFunc(oldNowFunc), WithReportCaller(true))
	}()
	NewGlobalLogger(WithOutput(os.Stdout), WithNowFunc(mockNowFunc), WithReportCaller(false))

	funcThatAcceptsInterface(Global())
	// Output: {"level":"warning","time":"2020-10-10T10:10:10Z","msg":"foobar"}
}
