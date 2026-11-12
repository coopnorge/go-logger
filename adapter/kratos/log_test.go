package logs

import (
	"testing"

	coopLog "github.com/coopnorge/go-logger"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
)

func TestForwardersDelegateToV2(t *testing.T) {
	var a *LoggerKratosAdapter = NewLoggerKratosAdapter(coopLog.New()) //nolint:staticcheck
	assert.NoError(t, a.Log(log.LevelInfo, "msg", "hello"))
}
