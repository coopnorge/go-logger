package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getCaller_initializes_vars(t *testing.T) {
	assert.Equal(t, "", goLoggerPackage)
	getCaller()
	assert.Equal(t, "github.com/coopnorge/go-logger", goLoggerPackage)
}

func Test_getPackageName(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "github.com/myorg/my-repo.something.func1",
			expected: "github.com/myorg/my-repo",
		},
		{
			input:    "myrepo.myfile.MyFunc",
			expected: "myrepo",
		},
		{
			input:    "package/subpackage.Function",
			expected: "package/subpackage",
		},
		{
			input:    "simplepkg.Function",
			expected: "simplepkg",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			assert.Equal(t, testCase.expected, getPackageName(testCase.input))
		})
	}
}
