package goose

import (
	"testing"

	log "github.com/coopnorge/go-logger"
)

func TestForwardersDelegateToV3(_ *testing.T) {
	var l *Logger = New(log.New()) //nolint:staticcheck
	l.Printf("hello %s", "world")
}
