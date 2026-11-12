package logs

import (
	"testing"

	coopLog "github.com/coopnorge/go-logger"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
)

func TestTestLoggerKratosAdapterEmptyLog(t *testing.T) {
	coopLogger := coopLog.New(coopLog.WithLevel(coopLog.LevelDebug))
	loggerAdapter := NewLoggerKratosAdapter(coopLogger)

	logErr := loggerAdapter.Log(log.LevelWarn)
	assert.NoError(t, logErr)
}

func TestLoggerKratosAdapterLevels(t *testing.T) {
	tests := []struct {
		level log.Level
	}{
		{log.LevelError},
		{log.LevelWarn},
		{log.LevelInfo},
		{log.LevelDebug},
	}

	for _, test := range tests {
		t.Run(test.level.String(), func(t *testing.T) {
			coopLogger := coopLog.New(coopLog.WithLevel(coopLog.LevelDebug))
			loggerAdapter := NewLoggerKratosAdapter(coopLogger)

			logErr := loggerAdapter.Log(test.level, "test message")
			assert.NoError(t, logErr)
		})
	}
}
