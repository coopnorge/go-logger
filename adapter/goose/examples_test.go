package goose_test

import (
	"github.com/coopnorge/go-logger"
	gooseLogger "github.com/coopnorge/go-logger/adapter/goose"
	"github.com/pressly/goose/v3"
)

func ExampleNew() {
	goose.SetLogger(gooseLogger.New(logger.Global()))
}
