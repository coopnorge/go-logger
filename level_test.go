package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestMapLevelToLogrusLevel(t *testing.T) {
	type testCase struct {
		input          Level
		expectedOutput logrus.Level
	}
	testCases := map[string]testCase{
		"map LevelFatal": {
			input:          LevelFatal,
			expectedOutput: logrus.FatalLevel,
		},

		"map LevelError": {
			input:          LevelError,
			expectedOutput: logrus.ErrorLevel,
		},
		"map LevelWarn": {
			input:          LevelWarn,
			expectedOutput: logrus.WarnLevel,
		},
		"map LevelInfo": {
			input:          LevelInfo,
			expectedOutput: logrus.InfoLevel,
		},
		"map levelDebug": {
			input:          LevelDebug,
			expectedOutput: logrus.DebugLevel,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			res := mapLevelToLogrusLevel(tc.input)
			if res != tc.expectedOutput {
				t.Fatalf("expected %v to map to %v, got: %v", tc.input, tc.expectedOutput, res)
			}
		})
	}
}
