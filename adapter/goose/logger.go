// Package goose forwards to github.com/coopnorge/go-logger/adapter/goose.v3.
//
// Deprecated: use github.com/coopnorge/go-logger/adapter/goose.v3.
package goose

import (
	log "github.com/coopnorge/go-logger"
	goosev3 "github.com/coopnorge/go-logger/adapter/goose.v3"
)

// Logger is a logging adapter between Goose and go-logger.
//
// Deprecated: use github.com/coopnorge/go-logger/adapter/goose.v3.Logger.
//
//go:fix inline
type Logger = goosev3.Logger

// New creates a new Logger instance.
//
// Deprecated: use github.com/coopnorge/go-logger/adapter/goose.v3.New.
//
//go:fix inline
func New(logger *log.Logger) *goosev3.Logger {
	return goosev3.New(logger)
}
