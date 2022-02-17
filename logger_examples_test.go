// All examples disable reporing caller as file path depends on workdir
// which is different for every run environment and therefore it's impossible
// to match exact output.
package logger

import (
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
	// Output: {"level":"warning","msg":"foobar","time":"2020-10-10T10:10:10Z"}
	// {"level":"error","msg":"foobar","time":"2020-10-10T10:10:10Z"}
}

func ExampleWithLevel() {
	logger := New(WithNowFunc(mockNowFunc), WithLevel(LevelInfo), WithReportCaller(false))
	logger.Info("now log level is set to info or lower, I will be logged")
	// Output: {"level":"info","msg":"now log level is set to info or lower, I will be logged","time":"2020-10-10T10:10:10Z"}
}

func ExampleWithFields() {
	// example runner replaces os.Stdout to catch output and compare it with desired output, so we need to overwrite default output, which points to actual system standard output, since it is initialized before tests and examples run
	oldOutput := globalLogger.output
	oldNowFunc := globalLogger.now
	defer func() {
		ConfigureGlobalLogger(WithOutput(oldOutput), WithNowFunc(oldNowFunc), WithReportCaller(true))
	}()
	ConfigureGlobalLogger(WithOutput(os.Stdout), WithNowFunc(mockNowFunc), WithReportCaller(false))

	WithFields(Fields{
		"timeSpentOnConfiguration": 0,
	}).Warn("use default logger with 0 configuration")
	// Output: {"level":"warning","msg":"use default logger with 0 configuration","time":"2020-10-10T10:10:10Z","timeSpentOnConfiguration":0}
}

type warner interface {
	Warn(args ...interface{})
}

func funcThatAcceptsInterface(warner warner) {
	warner.Warn("foobar")
}

func ExampleGlobal() {
	oldOutput := globalLogger.output
	oldNowFunc := globalLogger.now
	defer func() {
		ConfigureGlobalLogger(WithOutput(oldOutput), WithNowFunc(oldNowFunc), WithReportCaller(true))
	}()
	ConfigureGlobalLogger(WithOutput(os.Stdout), WithNowFunc(mockNowFunc), WithReportCaller(false))

	funcThatAcceptsInterface(Global())
	// Output: {"level":"warning","msg":"foobar","time":"2020-10-10T10:10:10Z"}
}
