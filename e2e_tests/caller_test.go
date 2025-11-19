package e2etests

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"testing"

	logger "github.com/coopnorge/go-logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCaller(t *testing.T) {
	// This function tests that the caller information is correctly added to log entries
	// It should contain the exact correct line number that the log was called from.

	builder := &strings.Builder{}
	logger := logger.New(logger.WithOutput(builder), logger.WithReportCaller(true), logger.WithLevel(logger.LevelDebug))

	logger.Info("Hello")
	_, _, line, ok := runtime.Caller(0) // Important: Run this IMMEDIATELY after the log line
	require.True(t, ok)
	expectedLine := line - 1 // We want the log to contain the line where logger.Info was called

	s := builder.String()
	logJSON := map[string]any{}
	err := json.Unmarshal([]byte(s), &logJSON)
	require.NoError(t, err)

	file, ok := logJSON["file"].(string)
	require.True(t, ok)
	assert.Contains(t, file, fmt.Sprintf("e2e_tests/caller_test.go:%d", expectedLine))

	function, ok := logJSON["function"].(string)
	require.True(t, ok)
	assert.Equal(t, "github.com/coopnorge/go-logger-e2e-tests.TestGetCaller", function)
}

func TestGetCallerFromPreviousEntry(t *testing.T) {
	// This function tests that the caller information is correctly added to log entries
	// It should contain the exact correct line number that the log was called from.

	builder := &strings.Builder{}
	logger := logger.New(logger.WithOutput(builder), logger.WithReportCaller(true), logger.WithLevel(logger.LevelDebug))

	entry := logger.WithField("key", "value")

	entry.Info("Hello")
	_, _, line, ok := runtime.Caller(0) // Important: Run this IMMEDIATELY after the log line
	require.True(t, ok)
	expectedLine := line - 1 // We want the log to contain the line where logger.Info was called, NOT where the entry was created

	s := builder.String()
	logJSON := map[string]any{}
	err := json.Unmarshal([]byte(s), &logJSON)
	require.NoError(t, err)

	file, ok := logJSON["file"].(string)
	require.True(t, ok)
	assert.Contains(t, file, fmt.Sprintf("e2e_tests/caller_test.go:%d", expectedLine))

	function, ok := logJSON["function"].(string)
	require.True(t, ok)
	assert.Equal(t, "github.com/coopnorge/go-logger-e2e-tests.TestGetCallerFromPreviousEntry", function)
}
