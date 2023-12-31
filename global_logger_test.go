package logger

import (
	"bytes"
	"context"
	"fmt"
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

func TestChainingSetup(t *testing.T) {
	buf := &bytes.Buffer{}
	oldOutput := globalLogger.output
	oldLevel := globalLogger.level
	oldNowFunc := globalLogger.now
	defer func() {
		// bring global logger to original state after tests
		NewGlobalLogger(WithOutput(oldOutput), WithLevel(oldLevel), WithNowFunc(oldNowFunc))
	}()
	NewGlobalLogger(WithOutput(buf), WithLevel(LevelDebug), WithNowFunc(mockNowFunc), WithHookFunc(testHook))

	ctx := context.WithValue(context.Background(), myCtxKey{}, "my-custom-ctx-value")
	err := fmt.Errorf("some error")

	{
		// Start with global WithField
		WithField("foo", "bar").WithContext(ctx).WithFields(Fields{"baz": "quoo", "number": 42}).WithError(err).Infof("hello %s", "world")
		b := buf.Bytes() // get bytes for multiple-reads
		buf.Reset()      // Prepare for next log message

		assertLogEntryContains(t, bytes.NewReader(b), "msg", "hello world")
		assertLogEntryContains(t, bytes.NewReader(b), "foo", "bar")
		assertLogEntryContains(t, bytes.NewReader(b), "my-custom-log-key", "my-custom-ctx-value")
		assertLogEntryContains(t, bytes.NewReader(b), "baz", "quoo")
		assertLogEntryContains(t, bytes.NewReader(b), "number", float64(42))
	}

	{
		// Start with global WithFields
		WithFields(Fields{"baz": "quoo", "number": 42}).WithField("foo", "bar").WithContext(ctx).WithError(err).Infof("hello %s", "world")
		b := buf.Bytes() // get bytes for multiple-reads
		buf.Reset()      // Prepare for next log message

		assertLogEntryContains(t, bytes.NewReader(b), "msg", "hello world")
		assertLogEntryContains(t, bytes.NewReader(b), "foo", "bar")
		assertLogEntryContains(t, bytes.NewReader(b), "my-custom-log-key", "my-custom-ctx-value")
		assertLogEntryContains(t, bytes.NewReader(b), "baz", "quoo")
		assertLogEntryContains(t, bytes.NewReader(b), "number", float64(42))
	}

	{
		// Start with global WithError
		WithError(err).WithFields(Fields{"baz": "quoo", "number": 42}).WithField("foo", "bar").WithContext(ctx).Infof("hello %s", "world")
		b := buf.Bytes() // get bytes for multiple-reads
		buf.Reset()      // Prepare for next log message

		assertLogEntryContains(t, bytes.NewReader(b), "msg", "hello world")
		assertLogEntryContains(t, bytes.NewReader(b), "foo", "bar")
		assertLogEntryContains(t, bytes.NewReader(b), "my-custom-log-key", "my-custom-ctx-value")
		assertLogEntryContains(t, bytes.NewReader(b), "baz", "quoo")
		assertLogEntryContains(t, bytes.NewReader(b), "number", float64(42))
	}

	{
		// Start with global WithContext
		WithContext(ctx).WithError(err).WithFields(Fields{"baz": "quoo", "number": 42}).WithField("foo", "bar").Infof("hello %s", "world")
		b := buf.Bytes() // get bytes for multiple-reads
		buf.Reset()      // Prepare for next log message

		assertLogEntryContains(t, bytes.NewReader(b), "msg", "hello world")
		assertLogEntryContains(t, bytes.NewReader(b), "foo", "bar")
		assertLogEntryContains(t, bytes.NewReader(b), "my-custom-log-key", "my-custom-ctx-value")
		assertLogEntryContains(t, bytes.NewReader(b), "baz", "quoo")
		assertLogEntryContains(t, bytes.NewReader(b), "number", float64(42))
	}
}
