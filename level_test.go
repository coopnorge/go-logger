package logger

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestMapLevelToZapLevel(t *testing.T) {
	type testCase struct {
		input          Level
		expectedOutput zapcore.Level
	}
	testCases := map[string]testCase{
		"map LevelFatal": {
			input:          LevelFatal,
			expectedOutput: zapcore.FatalLevel,
		},

		"map LevelError": {
			input:          LevelError,
			expectedOutput: zapcore.ErrorLevel,
		},
		"map LevelWarn": {
			input:          LevelWarn,
			expectedOutput: zapcore.WarnLevel,
		},
		"map LevelInfo": {
			input:          LevelInfo,
			expectedOutput: zapcore.InfoLevel,
		},
		"map levelDebug": {
			input:          LevelDebug,
			expectedOutput: zapcore.DebugLevel,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			res := mapLevelToZapLevel(tc.input)
			if res != tc.expectedOutput {
				t.Fatalf("expected %v to map to %v, got: %v", tc.input, tc.expectedOutput, res)
			}
		})
	}
}
