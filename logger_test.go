package logger_test

import (
	"testing"

	"github.com/coopnorge/go-logger"
)

func TestLogger(t *testing.T) {
	logger.Log("foobar")
}
