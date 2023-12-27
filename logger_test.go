package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
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
	switch expected := expectedValue.(type) {
	case *regexp.Regexp:
		vString, ok := v.(string)
		if !ok {
			t.Fatalf("cannot match regexp against non-string value")
		}
		if !expected.MatchString(vString) {
			t.Fatalf("log[%v]=%v does not match regexp: %v", key, expected, vString)
		}
	default:
		if v != expected {
			t.Fatalf("expected log[%v] to have value: %v, got: %v", key, expected, v)
		}
	}
}

func assertLogEntryHasKey(t *testing.T, logReader io.Reader, key string) {
	log := make(map[string]interface{})
	err := json.NewDecoder(logReader).Decode(&log)
	if err != nil {
		t.Fatalf("cannot decode log entry: %v", err)
	}
	_, ok := log[key]
	if !ok {
		t.Fatalf("key not found: %v", key)
	}
}

func assertLogEntryDoesNotHaveKey(t *testing.T, logReader io.Reader, key string) {
	log := make(map[string]interface{})
	err := json.NewDecoder(logReader).Decode(&log)
	if err != nil {
		t.Fatalf("cannot decode log entry: %v", err)
	}
	_, ok := log[key]
	if ok {
		t.Fatalf("unexpected key found: %v", key)
	}
}

func TestLogLevels(t *testing.T) {
	type testCase struct {
		logFunc          func(args ...interface{})
		expectedLogLevel string
	}
	buf := &bytes.Buffer{}
	testLogger := New(WithOutput(buf), WithLevel(LevelDebug))
	testLogger.logrusLogger.ExitFunc = func(int) {} // prevent .Fatal() from shutting down test runner
	testCases := map[string]testCase{
		"logger.Info() should log with level info": {
			logFunc:          testLogger.Info,
			expectedLogLevel: "info",
		},
		"logger.Error() should log with level error": {
			logFunc:          testLogger.Error,
			expectedLogLevel: "error",
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
	testLogger := New(WithOutput(buf), WithLevel(LevelDebug))
	testLogger.logrusLogger.ExitFunc = func(int) {} // prevent .Fatal() from shutting down test runner
	testCases := map[string]testCase{
		"logger.Infof() should log with level info": {
			logFunc:          testLogger.Infof,
			expectedLogLevel: "info",
		},
		"logger.Errorf() should log with level error": {
			logFunc:          testLogger.Errorf,
			expectedLogLevel: "error",
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
		f := func(t *testing.T, useSingle bool) {
			buf := &bytes.Buffer{}
			logger := New(WithLevel(LevelDebug), WithOutput(buf))
			if useSingle {
				logger.WithField("customField", tc.customFieldValue).Warnf("blabla")
			} else {
				logger.WithFields(Fields{
					"customField": tc.customFieldValue,
				}).Warnf("blabla")
			}
			log, err := io.ReadAll(buf)
			if err != nil {
				t.Fatalf("cannot read buffer: %v", err)
			}
			expectedLoggedValue := fmt.Sprintf(`"customField":%v`, tc.expectedLoggedValue)
			if !strings.Contains(string(log), expectedLoggedValue) {
				t.Fatalf("expected to find %v in log (%v)", expectedLoggedValue, string(log))
			}
		}
		t.Run(name+"_single", func(t *testing.T) { f(t, true) })
		t.Run(name+"_multiple", func(t *testing.T) { f(t, false) })
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
	b, err := io.ReadAll(logReader)
	if err != nil && err != io.EOF {
		t.Fatalf("cannot read log entry: %v", err)
	}
	return !(err == io.EOF || len(b) == 0)
}

func TestSettingLogLevel(t *testing.T) {
	type testCase struct {
		logLevel             Level
		expectedLoggedLevels []Level
	}

	testCases := map[string]testCase{
		"log fatal and above": {
			logLevel:             LevelFatal,
			expectedLoggedLevels: []Level{LevelFatal},
		},
		"log error and above": {
			logLevel:             LevelError,
			expectedLoggedLevels: []Level{LevelError, LevelFatal},
		},
		"log warn and above": {
			logLevel:             LevelWarn,
			expectedLoggedLevels: []Level{LevelWarn, LevelError, LevelFatal},
		},
		"log info and above": {
			logLevel:             LevelInfo,
			expectedLoggedLevels: []Level{LevelInfo, LevelWarn, LevelError, LevelFatal},
		},
		"log debug and above": {
			logLevel:             LevelDebug,
			expectedLoggedLevels: []Level{LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := New(WithOutput(buf), WithLevel(tc.logLevel))
			logger.logrusLogger.ExitFunc = func(int) {}

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
		})
	}
}

func TestWithLevelName(t *testing.T) {
	for name, lvl := range nameMapping {
		logger := New(WithLevelName(name))
		if logger.level != lvl {
			t.Fatalf("expected level %v, got %v", lvl, logger.logrusLogger.Level)
		}
	}
}

func TestBadLevelName(t *testing.T) {
	_, ok := LevelNameToLevel("invalid")
	if ok {
		t.Fatalf("Invalid level should not be ok")
	}
}

func TestReporingCaller(t *testing.T) {
	buf := &bytes.Buffer{}
	tee := io.TeeReader(buf, buf)
	logger := New(WithOutput(buf))
	logger.Error("foobar")
	assertLogEntryContains(t, tee, "file", regexp.MustCompile(`.*\.go:\d+$`))
	assertLogEntryHasKey(t, tee, "function")
}

func TestDisableReportingCaller(t *testing.T) {
	buf := &bytes.Buffer{}
	tee := io.TeeReader(buf, buf)
	logger := New(WithOutput(buf), WithReportCaller(false))
	logger.Error("foobar")
	assertLogEntryDoesNotHaveKey(t, tee, "file")
	assertLogEntryDoesNotHaveKey(t, tee, "function")
}

type myCtxKey struct{}

// Fire - implement Hook interface fire the entry
func testHook(entry *HookEntry) (bool, error) {
	ctx := entry.Context
	if ctx == nil {
		return false, nil
	}

	val := ctx.Value(myCtxKey{})
	if val == nil {
		return false, nil
	}

	str, ok := val.(string)
	if !ok || str == "" {
		return false, nil
	}

	entry.Data["my-custom-log-key"] = str
	return true, nil
}

func TestHookWithContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), myCtxKey{}, "my-custom-ctx-value")

	buf := &bytes.Buffer{}
	tee := io.TeeReader(buf, buf)
	logger := New(WithOutput(buf), WithHookFunc(testHook))
	logger.WithContext(ctx).Error("foobar")
	assertLogEntryContains(t, tee, "my-custom-log-key", "my-custom-ctx-value")
}

func TestHookWithoutContext(t *testing.T) {
	buf := &bytes.Buffer{}
	tee := io.TeeReader(buf, buf)
	logger := New(WithOutput(buf), WithHookFunc(testHook))
	logger.Error("foobar")
	assertLogEntryDoesNotHaveKey(t, tee, "my-custom-log-key")
}

func TestHookWithContext2(t *testing.T) {
	ctx := context.WithValue(context.Background(), myCtxKey{}, "my-custom-ctx-value")

	buf := &bytes.Buffer{}
	tee := io.TeeReader(buf, buf)
	logger := New(WithOutput(buf), WithHookFunc(testHook))
	logger.WithContext(ctx).Error("foobar")
	assertLogEntryContains(t, tee, "my-custom-log-key", "my-custom-ctx-value")
}
