package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func assertLogEntryContains(t *testing.T, logReader io.Reader, key string, expectedValue interface{}) {
	log := make(map[string]interface{})
	err := json.NewDecoder(logReader).Decode(&log)
	if err != nil {
		t.Fatalf("cannot decode log entry: %v", err)
	}
	v, ok := log[key]
	if !ok {
		t.Fatalf("no value found for key %v", key)
	}
	if v != expectedValue {
		t.Fatalf("expected log[%v] to have value: %v, got: %v", key, expectedValue, v)
	}
}

func TestLogLevels(t *testing.T) {
	type testCase struct {
		logFunc          func(args ...interface{})
		expectedLogLevel string
	}
	buf := &bytes.Buffer{}
	testLogger := New(WithOutput(buf), WithLevel(LevelTrace))
	testLogger.logrusLogger.ExitFunc = func(int) {} // prevent .Fatal() from shutting down test runner
	testCases := map[string]testCase{
		"logger.Print() should log with level info": {
			logFunc:          testLogger.Print,
			expectedLogLevel: "info",
		},
		"logger.Info() should log with level info": {
			logFunc:          testLogger.Info,
			expectedLogLevel: "info",
		},
		"logger.Error() should log with level error": {
			logFunc:          testLogger.Error,
			expectedLogLevel: "error",
		},
		"logger.Trace() should log with level trace": {
			logFunc:          testLogger.Trace,
			expectedLogLevel: "trace",
		},
		"logger.Debug() should log with level debug": {
			logFunc:          testLogger.Debug,
			expectedLogLevel: "debug",
		},
		"logger.Warn() should log with level warning": {
			logFunc:          testLogger.Warn,
			expectedLogLevel: "warning",
		},
		"logger.Fatal() should log with level fatal": {
			logFunc:          testLogger.Fatal,
			expectedLogLevel: "fatal",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.logFunc("foobar")
			assertLogEntryContains(t, buf, "level", tc.expectedLogLevel)
		})
	}
}

func TestLogLevelsInFormatFuncs(t *testing.T) {
	type testCase struct {
		logFunc          func(format string, args ...interface{})
		expectedLogLevel string
	}
	buf := &bytes.Buffer{}
	testLogger := New(WithOutput(buf), WithLevel(LevelTrace))
	testLogger.logrusLogger.ExitFunc = func(int) {} // prevent .Fatal() from shutting down test runner
	testCases := map[string]testCase{
		"logger.Printf() should log with level info": {
			logFunc:          testLogger.Printf,
			expectedLogLevel: "info",
		},
		"logger.Infof() should log with level info": {
			logFunc:          testLogger.Infof,
			expectedLogLevel: "info",
		},
		"logger.Errorf() should log with level error": {
			logFunc:          testLogger.Errorf,
			expectedLogLevel: "error",
		},
		"logger.Tracef() should log with level trace": {
			logFunc:          testLogger.Tracef,
			expectedLogLevel: "trace",
		},
		"logger.Debugf() should log with level debug": {
			logFunc:          testLogger.Debugf,
			expectedLogLevel: "debug",
		},
		"logger.Warnf() should log with level warning": {
			logFunc:          testLogger.Warnf,
			expectedLogLevel: "warning",
		},
		"logger.Fatalf() should log with level fatal": {
			logFunc:          testLogger.Fatalf,
			expectedLogLevel: "fatal",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.logFunc("foobar")
			assertLogEntryContains(t, buf, "level", tc.expectedLogLevel)
		})
	}
}

// func TestGlobalLoggerLogLevels(t *testing.T) {
func GlobalLoggerLogLevels(t *testing.T) {
	type testCase struct {
		logFunc          func(args ...interface{})
		expectedLogLevel string
	}

	buf := &bytes.Buffer{}
	oldOutput := globalLogger.output
	oldLevel := globalLogger.level
	oldExitFunc := globalLogger.logrusLogger.ExitFunc
	defer func() {
		// bring global logger to original state after tests
		ConfigureGlobalLogger(WithOutput(oldOutput), WithLevel(oldLevel))
		globalLogger.logrusLogger.ExitFunc = oldExitFunc
	}()
	ConfigureGlobalLogger(WithOutput(buf), WithLevel(LevelTrace))
	globalLogger.logrusLogger.ExitFunc = func(int) {}

	testCases := map[string]testCase{
		"Print() should log with level info": {
			logFunc:          Print,
			expectedLogLevel: "info",
		},
		"Info() should log with level info": {
			logFunc:          Info,
			expectedLogLevel: "info",
		},
		"Error() should log with level error": {
			logFunc:          Error,
			expectedLogLevel: "error",
		},
		"Trace() should log with level trace": {
			logFunc:          Trace,
			expectedLogLevel: "trace",
		},
		"Debug() should log with level debug": {
			logFunc:          Debug,
			expectedLogLevel: "debug",
		},
		"Warn() should log with level warning": {
			logFunc:          Warn,
			expectedLogLevel: "warning",
		},
		"Fatal() should log with level fatal": {
			logFunc:          Fatal,
			expectedLogLevel: "fatal",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.logFunc("foobar")
			assertLogEntryContains(t, buf, "level", tc.expectedLogLevel)
		})
	}
}

// func TestGlobalLoggerLogLevelsInFormatFuncs(t *testing.T) {
func GlobalLoggerLogLevelsInFormatFuncs(t *testing.T) {
	type testCase struct {
		logFunc          func(format string, args ...interface{})
		expectedLogLevel string
	}

	buf := &bytes.Buffer{}
	oldOutput := globalLogger.output
	oldLevel := globalLogger.level
	oldExitFunc := globalLogger.logrusLogger.ExitFunc
	defer func() {
		// bring global logger to original state after tests
		ConfigureGlobalLogger(WithOutput(oldOutput), WithLevel(oldLevel))
		globalLogger.logrusLogger.ExitFunc = oldExitFunc
	}()
	ConfigureGlobalLogger(WithOutput(buf), WithLevel(LevelTrace))
	globalLogger.logrusLogger.ExitFunc = func(int) {}

	testCases := map[string]testCase{
		"Print() should log with level info": {
			logFunc:          Printf,
			expectedLogLevel: "info",
		},
		"Info() should log with level info": {
			logFunc:          Infof,
			expectedLogLevel: "info",
		},
		"Error() should log with level error": {
			logFunc:          Errorf,
			expectedLogLevel: "error",
		},
		"Trace() should log with level trace": {
			logFunc:          Tracef,
			expectedLogLevel: "trace",
		},
		"Debug() should log with level debug": {
			logFunc:          Debugf,
			expectedLogLevel: "debug",
		},
		"Warn() should log with level warning": {
			logFunc:          Warnf,
			expectedLogLevel: "warning",
		},
		"Fatal() should log with level fatal": {
			logFunc:          Fatalf,
			expectedLogLevel: "fatal",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.logFunc("foobar")
			assertLogEntryContains(t, buf, "level", tc.expectedLogLevel)
		})
	}
}

func TestLoggingCustomFields(t *testing.T) {
	type testCase struct {
		customFieldValue    interface{}
		expectedLoggedValue interface{}
	}
	testCases := map[string]testCase{
		"string": {
			customFieldValue:    "foobar",
			expectedLoggedValue: `"foobar"`,
		},
		"int": {
			customFieldValue:    3,
			expectedLoggedValue: "3",
		},
		"slice of ints": {
			customFieldValue:    []int{1, 2, 3},
			expectedLoggedValue: "[1,2,3]",
		},
		"error": {
			customFieldValue:    fmt.Errorf("foobar"),
			expectedLoggedValue: `"foobar"`,
		},
		"map": {
			customFieldValue: map[string]int{
				"foo": 1,
			},
			expectedLoggedValue: `{"foo":1}`,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := New(WithLevel(LevelTrace), WithOutput(buf))
			logger.WithFields(Fields{
				"customField": tc.customFieldValue,
			}).Warnf("blabla")
			log, err := ioutil.ReadAll(buf)
			if err != nil {
				t.Fatalf("cannot read buffer: %v", err)
			}
			expectedLoggedValue := fmt.Sprintf(`"customField":%v`, tc.expectedLoggedValue)
			if !strings.Contains(string(log), expectedLoggedValue) {
				t.Fatalf("expected to find %v in log (%v)", expectedLoggedValue, string(log))
			}
		})
	}

}

func contains(levels []Level, level Level) bool {
	for _, l := range levels {
		if l == level {
			return true
		}
	}
	return false
}

func wasLogged(t *testing.T, logReader io.Reader) bool {
	bytes, err := ioutil.ReadAll(logReader)
	if err != nil && err != io.EOF {
		t.Fatalf("cannot read log entry: %v", err)
	}
	return !(err == io.EOF || len(bytes) == 0)
}

func TestSettingLogLeve(t *testing.T) {
	type testCase struct {
		logLevel             Level
		expectedLoggedLevels []Level
	}

	testCases := map[string]testCase{
		"log panic": {
			logLevel:             LevelPanic,
			expectedLoggedLevels: []Level{LevelPanic},
		},
		"log fatal and above": {
			logLevel:             LevelFatal,
			expectedLoggedLevels: []Level{LevelFatal, LevelPanic},
		},
		"log error and above": {
			logLevel:             LevelError,
			expectedLoggedLevels: []Level{LevelError, LevelFatal, LevelPanic},
		},
		"log warn and above": {
			logLevel:             LevelWarn,
			expectedLoggedLevels: []Level{LevelWarn, LevelError, LevelFatal, LevelPanic},
		},
		"log info and above": {
			logLevel:             LevelInfo,
			expectedLoggedLevels: []Level{LevelInfo, LevelWarn, LevelError, LevelFatal, LevelPanic},
		},
		"log debug and above": {
			logLevel:             LevelDebug,
			expectedLoggedLevels: []Level{LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal, LevelPanic},
		},
		"log trace and above": {
			logLevel:             LevelTrace,
			expectedLoggedLevels: []Level{LevelTrace, LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal, LevelPanic},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatalf(".Panic() did not panic")
				}
			}()
			buf := &bytes.Buffer{}
			logger := New(WithOutput(buf), WithLevel(tc.logLevel))
			logger.logrusLogger.ExitFunc = func(int) {}

			logger.Trace("trace")
			if contains(tc.expectedLoggedLevels, LevelTrace) != wasLogged(t, buf) {
				t.Fatalf("trace level was incorrectly filtered")
			}

			logger.Debug("debug")
			if contains(tc.expectedLoggedLevels, LevelDebug) != wasLogged(t, buf) {
				t.Fatalf("debug level was incorrectly filtered")
			}

			logger.Info("info")
			if contains(tc.expectedLoggedLevels, LevelInfo) != wasLogged(t, buf) {
				t.Fatalf("info level was incorrectly filtered")
			}

			logger.Warn("warn")
			if contains(tc.expectedLoggedLevels, LevelWarn) != wasLogged(t, buf) {
				t.Fatalf("warn level was incorrectly filtered")
			}

			logger.Error("error")
			if contains(tc.expectedLoggedLevels, LevelError) != wasLogged(t, buf) {
				t.Fatalf("error level was incorrectly filtered")
			}

			logger.Fatal("fatal")
			if contains(tc.expectedLoggedLevels, LevelFatal) != wasLogged(t, buf) {
				t.Fatalf("fatal level was incorrectly filtered")
			}

			logger.Panic("panic")
			if contains(tc.expectedLoggedLevels, LevelPanic) != wasLogged(t, buf) {
				t.Fatalf("panic level was incorrectly filtered")
			}

		})
	}
}
