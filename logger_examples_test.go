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
	return time.Date(2020, 10, 10, 10, 10, 10, 1000000, time.UTC)
}

func ExampleInfo() {
	logger := New(WithNowFunc(mockNowFunc), WithReportCaller(false))
	logger.Warn("foobar")
	logger.Info("i won't be logged because the default log level is higher than info")
	logger.Error("foobar")
	// Output: {"level":"warning","msg":"foobar","time":"2020-10-10T10:10:10.001Z"}
	// {"level":"error","msg":"foobar","time":"2020-10-10T10:10:10.001Z"}
}

func ExampleWithLevel() {
	logger := New(WithNowFunc(mockNowFunc), WithLevel(LevelInfo), WithReportCaller(false))
	logger.Info("now log level is set to info or lower, I will be logged")
	// Output: {"level":"info","msg":"now log level is set to info or lower, I will be logged","time":"2020-10-10T10:10:10.001Z"}
}

func ExampleWithLevelName() {
	logger := New(WithNowFunc(mockNowFunc), WithLevelName("info"), WithReportCaller(false))
	logger.Info("now log level is set to info or lower, I will be logged")
	// Output: {"level":"info","msg":"now log level is set to info or lower, I will be logged","time":"2020-10-10T10:10:10.001Z"}
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
		ConfigureGlobalLogger(WithOutput(oldOutput), WithNowFunc(oldNowFunc), WithReportCaller(true))
	}()
	ConfigureGlobalLogger(WithOutput(os.Stdout), WithNowFunc(mockNowFunc), WithReportCaller(false))

	WithFields(Fields{
		"timeSpentOnConfiguration": 0,
		"defaultsLoaded":           true,
	}).Warn("use default logger with 0 configuration")
	WithField("singleField", true).Warn("example with a single field")
	// Output: {"defaultsLoaded":true,"level":"warning","msg":"use default logger with 0 configuration","time":"2020-10-10T10:10:10.001Z","timeSpentOnConfiguration":0}
	// {"level":"warning","msg":"example with a single field","singleField":true,"time":"2020-10-10T10:10:10.001Z"}
}

func ExampleWithError() {
	logger := New(WithNowFunc(mockNowFunc), WithReportCaller(false))

	err := errors.New("Test error")
	logger.WithError(err).Error("Operation failed")
	// Output: {"error":"Test error","level":"error","msg":"Operation failed","time":"2020-10-10T10:10:10.001Z"}
}

type warner interface {
	Warn(args ...any)
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
	// Output: {"level":"warning","msg":"foobar","time":"2020-10-10T10:10:10.001Z"}
}

// warnIfLastAttempt is a helper that demonstrates that you can use a function to provide the appropriate log-level
// based on data (deliveryCount in this instance)
func warnIfLastAttempt(deliveryCount int) Level {
	if deliveryCount < 5 {
		return LevelInfo
	}
	return LevelWarn
}

func ExampleLogger_Log_firstDelivery() {
	logger := New(WithNowFunc(mockNowFunc), WithReportCaller(false), WithLevel(LevelInfo))

	logger.Log(warnIfLastAttempt(1), "message delivery count: ", 1)
	// Output:
	// {"level":"info","msg":"message delivery count: 1","time":"2020-10-10T10:10:10.001Z"}
}

func ExampleLogger_Log_thirdDelivery() {
	logger := New(WithNowFunc(mockNowFunc), WithReportCaller(false), WithLevel(LevelInfo))

	logger.Log(warnIfLastAttempt(3), "message delivery count: ", 3)
	// Output: {"level":"info","msg":"message delivery count: 3","time":"2020-10-10T10:10:10.001Z"}
}

func ExampleLogger_Log_lastDelivery() {
	logger := New(WithNowFunc(mockNowFunc), WithReportCaller(false), WithLevel(LevelInfo))

	logger.Log(warnIfLastAttempt(5), "message delivery count: ", 5)
	// Output: {"level":"warning","msg":"message delivery count: 5","time":"2020-10-10T10:10:10.001Z"}
}

func ExampleLogger_Logf() {
	logger := New(WithNowFunc(mockNowFunc), WithReportCaller(false))

	logger.Logf(LevelWarn, "DeliveryCount is %d; level chosen based on value", 5)
	// Output:
	// {"level":"warning","msg":"DeliveryCount is 5; level chosen based on value","time":"2020-10-10T10:10:10.001Z"}
}
