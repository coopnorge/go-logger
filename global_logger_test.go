package logger

import (
	"bytes"
	"testing"
)

func TestGlobalLoggerLogLevels(t *testing.T) {
	type testCase struct {
		logFunc          func(msg string)
		expectedLogLevel string
	}

	buf := &bytes.Buffer{}
	oldOutput := globalLogger.output
	oldLevel := globalLogger.level
	defer func() {
		// bring global logger to original state after tests
		NewGlobalLogger(WithOutput(oldOutput), WithLevel(oldLevel))
	}()
	NewGlobalLogger(WithOutput(buf), WithLevel(LevelDebug))

	testCases := map[string]testCase{
		"Info() should log with level info": {
			logFunc:          Info,
			expectedLogLevel: "info",
		},
		"Error() should log with level error": {
			logFunc:          Error,
			expectedLogLevel: "error",
		},
		"Debug() should log with level debug": {
			logFunc:          Debug,
			expectedLogLevel: "debug",
		},
		"Warn() should log with level warning": {
			logFunc:          Warn,
			expectedLogLevel: "warning",
		},
		// "Fatal() should log with level fatal": {
		// 	logFunc:          Fatal,
		// 	expectedLogLevel: "fatal",
		// },
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.logFunc("foobar")
			assertLogEntryContains(t, buf, "level", tc.expectedLogLevel)
		})
	}
}

func TestGlobalLoggerLogLevelsInFormatFuncs(t *testing.T) {
	type testCase struct {
		logFunc          func(format string, args ...interface{})
		expectedLogLevel string
	}

	buf := &bytes.Buffer{}
	oldOutput := globalLogger.output
	oldLevel := globalLogger.level
	defer func() {
		// bring global logger to original state after tests
		NewGlobalLogger(WithOutput(oldOutput), WithLevel(oldLevel))
	}()
	NewGlobalLogger(WithOutput(buf), WithLevel(LevelDebug))

	testCases := map[string]testCase{
		"Info() should log with level info": {
			logFunc:          Infof,
			expectedLogLevel: "info",
		},
		"Error() should log with level error": {
			logFunc:          Errorf,
			expectedLogLevel: "error",
		},
		"Debug() should log with level debug": {
			logFunc:          Debugf,
			expectedLogLevel: "debug",
		},
		"Warn() should log with level warning": {
			logFunc:          Warnf,
			expectedLogLevel: "warning",
		},
		// "Fatal() should log with level fatal": {
		// 	logFunc:          Fatalf,
		// 	expectedLogLevel: "fatal",
		// },
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.logFunc("foobar")
			assertLogEntryContains(t, buf, "level", tc.expectedLogLevel)
		})
	}
}
