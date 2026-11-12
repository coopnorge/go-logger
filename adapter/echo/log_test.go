package labstack_logger //nolint:all

import (
	"testing"
)

func TestForwardersDelegateToV4(t *testing.T) {
	var l *WrappedEchoLogger = NewWrappedEchoLogger() //nolint:staticcheck
	if l.Prefix() != "echo" {
		t.Fatalf("unexpected prefix %q", l.Prefix())
	}
}
