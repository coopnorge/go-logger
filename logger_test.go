package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertLogEntryContains(t *testing.T, logReader io.Reader, key string, expectedValue interface{}) {
	t.Helper()
	log := decodeLogToMap(t, logReader)
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

func getLineNumber(t *testing.T, logReader io.Reader) int {
	const key = "file"
	t.Helper()

	log := decodeLogToMap(t, logReader)

	vAny, ok := log[key]
	if !ok {
		t.Fatalf("no value found for key %v", key)
	}
	v, ok := vAny.(string)
	if !ok {
		t.Fatalf("value in key %v was not string: %v", key, vAny)
	}
	r := regexp.MustCompile(`.*\.go:(\d+)$`)
	matches := r.FindStringSubmatch(v)
	if len(matches) != 2 {
		// matches will first match the full string, and the second match is the line-number
		t.Fatalf("log[%v]=%v does not match regexp", key, v)
	}
	i, err := strconv.Atoi(matches[1])
	if err != nil {
		t.Fatalf("log[%v]=%v failed to parse to int", key, v)
	}
	return i
}

func assertLogEntryHasKey(t *testing.T, logReader io.Reader, key string) {
	t.Helper()
	log := decodeLogToMap(t, logReader)

	_, ok := log[key]
	if !ok {
		t.Fatalf("key not found: %v", key)
	}
}

func assertLogEntryDoesNotHaveKey(t *testing.T, logReader io.Reader, key string) {
	t.Helper()
	log := decodeLogToMap(t, logReader)

	_, ok := log[key]
	if ok {
		t.Fatalf("unexpected key found: %v", key)
	}
}

func decodeLogToMap(t *testing.T, logReader io.Reader) map[string]interface{} {
	log := make(map[string]interface{})
	err := json.NewDecoder(logReader).Decode(&log)
	if err != nil {
		t.Fatalf("cannot decode log entry: %v", err)
	}

	return log
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

func TestLogLevelsInLogFunc(t *testing.T) {
	type testCase struct {
		logLevel         Level
		expectedLogLevel string
	}
	buf := &bytes.Buffer{}
	testLogger := New(WithOutput(buf), WithLevel(LevelDebug))
	testLogger.logrusLogger.ExitFunc = func(int) {} // prevent .Fatal() from shutting down test runner
	testCases := map[string]testCase{
		"logger.Infof() should log with level info": {
			logLevel:         LevelInfo,
			expectedLogLevel: "info",
		},
		"logger.Errorf() should log with level error": {
			logLevel:         LevelError,
			expectedLogLevel: "error",
		},
		"logger.Debugf() should log with level debug": {
			logLevel:         LevelDebug,
			expectedLogLevel: "debug",
		},
		"logger.Warnf() should log with level warning": {
			logLevel:         LevelWarn,
			expectedLogLevel: "warning",
		},
		"logger.Fatalf() should log with level fatal": {
			logLevel:         LevelFatal,
			expectedLogLevel: "fatal",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			testLogger.Log(tc.logLevel, "foobar")
			assertLogEntryContains(t, buf, "level", tc.expectedLogLevel)
		})
	}
}

func TestLogLevelsInLogfFunc(t *testing.T) {
	type testCase struct {
		logLevel         Level
		expectedLogLevel string
	}
	buf := &bytes.Buffer{}
	testLogger := New(WithOutput(buf), WithLevel(LevelDebug))
	testLogger.logrusLogger.ExitFunc = func(int) {} // prevent .Fatal() from shutting down test runner
	testCases := map[string]testCase{
		"logger.Infof() should log with level info": {
			logLevel:         LevelInfo,
			expectedLogLevel: "info",
		},
		"logger.Errorf() should log with level error": {
			logLevel:         LevelError,
			expectedLogLevel: "error",
		},
		"logger.Debugf() should log with level debug": {
			logLevel:         LevelDebug,
			expectedLogLevel: "debug",
		},
		"logger.Warnf() should log with level warning": {
			logLevel:         LevelWarn,
			expectedLogLevel: "warning",
		},
		"logger.Fatalf() should log with level fatal": {
			logLevel:         LevelFatal,
			expectedLogLevel: "fatal",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			testLogger.Logf(tc.logLevel, "foobar")
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
	return len(b) != 0
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

func TestReportingCaller(t *testing.T) {
	builder := &strings.Builder{}
	logger := New(WithOutput(builder), WithLevel(LevelInfo))
	reportCallerInGoLoggerPackage = true

	// Simple log
	logger.Error("foobar")

	// 2 logs from same entry
	entry := logger.WithField("foo", "bar")
	entry.Info("first")
	entry.Error("second")

	result := strings.TrimSpace(builder.String()) // Trim trailing newline
	lines := strings.Split(result, "\n")

	if len(lines) != 3 {
		t.Fatalf("expected %d lines, got %d", 3, len(lines))
	}

	lineNumbers := make([]int, 0, 3)

	{
		line := lines[0]
		assertLogEntryContains(t, strings.NewReader(line), "file", regexp.MustCompile(`.*\.go:\d+$`))
		assertLogEntryHasKey(t, strings.NewReader(line), "function")
		lineNumbers = append(lineNumbers, getLineNumber(t, strings.NewReader(line)))
	}

	{
		line := lines[1]
		assertLogEntryContains(t, strings.NewReader(line), "file", regexp.MustCompile(`.*\.go:\d+$`))
		assertLogEntryHasKey(t, strings.NewReader(line), "function")
		lineNumbers = append(lineNumbers, getLineNumber(t, strings.NewReader(line)))
	}
	{
		line := lines[2]
		assertLogEntryContains(t, strings.NewReader(line), "file", regexp.MustCompile(`.*\.go:\d+$`))
		assertLogEntryHasKey(t, strings.NewReader(line), "function")
		lineNumbers = append(lineNumbers, getLineNumber(t, strings.NewReader(line)))
	}

	uniqueLineNumber := unique(lineNumbers)
	if len(lineNumbers) != len(uniqueLineNumber) {
		t.Fatalf("expected %d unique line-numbers, got %d. All line-numbers: %#v", len(lineNumbers), len(uniqueLineNumber), lineNumbers)
	}
}

// Get unique items in slice
func unique[T comparable](slice []T) []T {
	m := make(map[T]struct{}, len(slice))

	for _, v := range slice {
		m[v] = struct{}{}
	}

	u := make([]T, 0, len(m))
	for k := range m {
		u = append(u, k)
	}
	return u
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

func TestReuseEntry(t *testing.T) {
	builder := &strings.Builder{}
	logger := New(WithOutput(builder), WithLevel(LevelInfo))
	withFields := logger.WithField("foo", "bar")
	withFields.Info("quoo")
	withFields.Error("baz")

	str := builder.String()
	lines := strings.Split(str, "\n")
	if len(lines) != 3 {
		// Info + Error + empty newline
		t.Fatalf("expected %d lines, got %d", 3, len(lines))
	}
	if lines[2] != "" {
		t.Fatalf("expected last line to be empty, got %s", lines[2])
	}

	assertLogEntryContains(t, strings.NewReader(lines[0]), "foo", "bar")
	assertLogEntryContains(t, strings.NewReader(lines[0]), "msg", "quoo")
	assertLogEntryContains(t, strings.NewReader(lines[0]), "level", "info")

	assertLogEntryContains(t, strings.NewReader(lines[1]), "foo", "bar")
	assertLogEntryContains(t, strings.NewReader(lines[1]), "msg", "baz")
	assertLogEntryContains(t, strings.NewReader(lines[1]), "level", "error")
}

func TestReuseEntryWithFields(t *testing.T) {
	// This test asserts that entry.WithField(...) does mutate the entry itself
	builder := &strings.Builder{}
	logger := New(WithOutput(builder), WithLevel(LevelInfo))
	entryWithFields := logger.WithField("foo", "bar")
	entryWithFields.WithField("only-quoo", true).Info("quoo")
	entryWithFields.WithError(fmt.Errorf("only-baz")).Error("baz")
	entryWithFields.Info("final")

	str := builder.String()
	lines := strings.Split(str, "\n")
	if len(lines) != 4 {
		// Info + Error + empty newline
		t.Fatalf("expected %d lines, got %d", 4, len(lines))
	}
	if lines[3] != "" {
		t.Fatalf("expected last line to be empty, got %s", lines[3])
	}

	assertLogEntryContains(t, strings.NewReader(lines[0]), "foo", "bar")
	assertLogEntryContains(t, strings.NewReader(lines[0]), "msg", "quoo")
	assertLogEntryContains(t, strings.NewReader(lines[0]), "level", "info")
	assertLogEntryContains(t, strings.NewReader(lines[0]), "only-quoo", true)
	assertLogEntryDoesNotHaveKey(t, strings.NewReader(lines[0]), "error")

	assertLogEntryContains(t, strings.NewReader(lines[1]), "foo", "bar")
	assertLogEntryContains(t, strings.NewReader(lines[1]), "msg", "baz")
	assertLogEntryContains(t, strings.NewReader(lines[1]), "level", "error")
	assertLogEntryDoesNotHaveKey(t, strings.NewReader(lines[1]), "only-quoo")
	assertLogEntryContains(t, strings.NewReader(lines[1]), "error", "only-baz")

	assertLogEntryContains(t, strings.NewReader(lines[2]), "foo", "bar")
	assertLogEntryContains(t, strings.NewReader(lines[2]), "msg", "final")
	assertLogEntryContains(t, strings.NewReader(lines[2]), "level", "info")
	assertLogEntryDoesNotHaveKey(t, strings.NewReader(lines[2]), "only-quoo")
	assertLogEntryDoesNotHaveKey(t, strings.NewReader(lines[2]), "error")
}

func TestTimeFormat(t *testing.T) {
	osloLoc, err := time.LoadLocation("Europe/Oslo")
	require.NoError(t, err)
	testcases := []struct {
		name     string
		now      time.Time
		expected string
	}{
		{
			name:     "sample timestamp",
			now:      time.Date(2020, 2, 3, 4, 5, 6, 789012345, time.UTC),
			expected: "2020-02-03T04:05:06.789Z",
		},
		{
			name:     "removes trailing zeros in sub-seconds part",
			now:      time.Date(2020, 2, 3, 4, 5, 8, 123000000, time.UTC),
			expected: "2020-02-03T04:05:08.123Z",
		},
		{
			name:     "removes trailing zeros in sub-seconds part",
			now:      time.Date(2020, 2, 3, 4, 5, 8, 120000000, time.UTC),
			expected: "2020-02-03T04:05:08.12Z",
		},
		{
			name:     "removes trailing zeros in sub-seconds part",
			now:      time.Date(2020, 2, 3, 4, 5, 8, 100000000, time.UTC),
			expected: "2020-02-03T04:05:08.1Z",
		},
		{
			name:     "includes 1 millisecond in output",
			now:      time.Date(2020, 2, 3, 4, 5, 8, 1000000, time.UTC),
			expected: "2020-02-03T04:05:08.001Z",
		},
		{
			name:     "omits sub-second precision when nano is 0",
			now:      time.Date(2020, 2, 3, 4, 5, 8, 0, time.UTC),
			expected: "2020-02-03T04:05:08Z",
		},
		{
			name:     "omits sub-second precision when nano is rounded to 0",
			now:      time.Date(2020, 2, 3, 4, 5, 8, 999999, time.UTC),
			expected: "2020-02-03T04:05:08Z",
		},
		{
			name:     "sample timestamp Oslo (Winter-time)",
			now:      time.Date(2020, 2, 3, 4, 5, 6, 789012345, osloLoc),
			expected: "2020-02-03T04:05:06.789+01:00",
		},
		{
			name:     "sample timestamp Oslo (Summer-time)",
			now:      time.Date(2020, 6, 3, 4, 5, 6, 789012345, osloLoc),
			expected: "2020-06-03T04:05:06.789+02:00",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			builder := &strings.Builder{}
			nowFunc := func() time.Time { return tc.now }
			logger := New(WithOutput(builder), WithLevel(LevelInfo), WithNowFunc(nowFunc))
			logger.Info("test time")

			log := builder.String()
			assertLogEntryContains(t, strings.NewReader(log), "time", tc.expected)
		})
	}
}

func TestOnlyFatalExits(t *testing.T) {
	type testCase struct {
		logFunc          func(args ...interface{})
		logfFunc         func(format string, args ...interface{})
		logLevel         Level
		expectedLogLevel string
	}
	buf := &bytes.Buffer{}
	testLogger := New(WithOutput(buf), WithLevel(LevelDebug))
	testCases := map[string]testCase{
		"level info": {
			logFunc:          testLogger.Info,
			logfFunc:         testLogger.Infof,
			logLevel:         LevelInfo,
			expectedLogLevel: "info",
		},
		"level error": {
			logFunc:          testLogger.Error,
			logfFunc:         testLogger.Errorf,
			logLevel:         LevelError,
			expectedLogLevel: "error",
		},
		"level debug": {
			logFunc:          testLogger.Debug,
			logfFunc:         testLogger.Debugf,
			logLevel:         LevelDebug,
			expectedLogLevel: "debug",
		},
		"level warning": {
			logFunc:          testLogger.Warn,
			logfFunc:         testLogger.Warnf,
			logLevel:         LevelWarn,
			expectedLogLevel: "warning",
		},
		"level fatal": {
			logFunc:          testLogger.Fatal,
			logfFunc:         testLogger.Fatalf,
			logLevel:         LevelFatal,
			expectedLogLevel: "fatal",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			exitCount := 0
			// prevent .Fatal() from shutting down test runner, and also count times it happened
			testLogger.logrusLogger.ExitFunc = func(int) {
				exitCount++
			}

			// Log using all 4 ways to log at a level

			tc.logFunc("foobar")
			assertLogEntryContains(t, buf, "level", tc.expectedLogLevel)

			tc.logfFunc("foobar")
			assertLogEntryContains(t, buf, "level", tc.expectedLogLevel)

			testLogger.Log(tc.logLevel, "foobar")
			assertLogEntryContains(t, buf, "level", tc.expectedLogLevel)

			testLogger.Logf(tc.logLevel, "foobar")
			assertLogEntryContains(t, buf, "level", tc.expectedLogLevel)

			// Only LevelFatal should trigger exit
			if tc.logLevel == LevelFatal {
				assert.Equal(t, 4, exitCount)
			} else {
				assert.Equal(t, 0, exitCount)
			}
		})
	}
}
