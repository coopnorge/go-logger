// Package logs forwards to github.com/coopnorge/go-logger/adapter/kratos.v2.
//
// Deprecated: use github.com/coopnorge/go-logger/adapter/kratos.v2.
package logs

import (
	coopLog "github.com/coopnorge/go-logger"
	kratosv2 "github.com/coopnorge/go-logger/adapter/kratos.v2"
)

// LoggerKratosAdapter adapter for Go-Kratos.
//
// Deprecated: use github.com/coopnorge/go-logger/adapter/kratos.v2.LoggerKratosAdapter.
//
//go:fix inline
type LoggerKratosAdapter = kratosv2.LoggerKratosAdapter

// NewLoggerKratosAdapter creates a LoggerKratosAdapter.
//
// Deprecated: use github.com/coopnorge/go-logger/adapter/kratos.v2.NewLoggerKratosAdapter.
//
//go:fix inline
func NewLoggerKratosAdapter(coopLog *coopLog.Logger) *kratosv2.LoggerKratosAdapter {
	return kratosv2.NewLoggerKratosAdapter(coopLog)
}
