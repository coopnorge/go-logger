// Package labstack_logger forwards to github.com/coopnorge/go-logger/adapter/echo.v4.
//
// Deprecated: use github.com/coopnorge/go-logger/adapter/echo.v4.
package labstack_logger //nolint:all

import (
	echov4 "github.com/coopnorge/go-logger/adapter/echo.v4"
)

// WrappedEchoLogger is an Echo v4 Logger backed by go-logger.
//
// Deprecated: use github.com/coopnorge/go-logger/adapter/echo.v4.WrappedEchoLogger.
//
//go:fix inline
type WrappedEchoLogger = echov4.WrappedEchoLogger

// NewWrappedEchoLogger creates a WrappedEchoLogger.
//
// Deprecated: use github.com/coopnorge/go-logger/adapter/echo.v4.NewWrappedEchoLogger.
//
//go:fix inline
func NewWrappedEchoLogger() *echov4.WrappedEchoLogger {
	return echov4.NewWrappedEchoLogger()
}
