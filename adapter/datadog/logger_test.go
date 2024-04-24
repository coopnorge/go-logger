package datadog

import (
	"fmt"
	"strings"
	"testing"

	coopLogger "github.com/coopnorge/go-logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
)

func TestSetup(t *testing.T) {
	logger, err := NewLogger(WithGlobalLogger())
	require.NoError(t, err)
	ddtrace.UseLogger(logger)
}

func TestGlobalLogger(t *testing.T) {
	output := &strings.Builder{}
	coopLogger.ConfigureGlobalLogger(coopLogger.WithLevel(coopLogger.LevelDebug), coopLogger.WithOutput(output))
	logger, err := NewLogger(WithGlobalLogger())
	require.NoError(t, err)

	tests := []struct {
		level    string
		input    string
		expected string
	}{
		{"error", "Datadog Tracer v1.63.0 ERROR: This is a test", "This is a test"},
		{"warning", "Datadog Tracer v1.63.0 WARN: This is a test", "This is a test"},
		{"info", "Datadog Tracer v1.63.0 INFO: This is a test", "This is a test"},
		{"debug", "Datadog Tracer v1.63.0 DEBUG: This is a test", "This is a test"},
		{"warning", "This is a test", "This is a test"},
	}
	for _, test := range tests {
		t.Run(test.level, func(t *testing.T) {
			output.Reset()
			logger.Log(test.input)
			assert.Contains(t, output.String(), fmt.Sprintf("\"level\":\"%v\"", test.level))
			assert.Contains(t, output.String(), fmt.Sprintf("\"msg\":\"%v\"", test.expected))
		})
	}
}

func TestCustomLogger(t *testing.T) {
	output := &strings.Builder{}
	logger, err := NewLogger(WithLogger(coopLogger.New(coopLogger.WithLevel(coopLogger.LevelDebug), coopLogger.WithOutput(output))))
	require.NoError(t, err)

	tests := []struct {
		level    string
		input    string
		expected string
	}{
		{"error", "Datadog Tracer v1.63.0 ERROR: This is a test", "This is a test"},
		{"warning", "Datadog Tracer v1.63.0 WARN: This is a test", "This is a test"},
		{"info", "Datadog Tracer v1.63.0 INFO: This is a test", "This is a test"},
		{"debug", "Datadog Tracer v1.63.0 DEBUG: This is a test", "This is a test"},
		{"warning", "This is a test", "This is a test"},
	}
	for _, test := range tests {
		t.Run(test.level, func(t *testing.T) {
			output.Reset()
			logger.Log(test.input)
			assert.Contains(t, output.String(), fmt.Sprintf("\"level\":\"%v\"", test.level))
			assert.Contains(t, output.String(), fmt.Sprintf("\"msg\":\"%v\"", test.expected))
		})
	}
}
