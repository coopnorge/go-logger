// Deprecated: use github.com/coopnorge/go-logger/adapter/goose.v3
module github.com/coopnorge/go-logger/adapter/goose

go 1.27.3

require (
	github.com/coopnorge/go-logger v0.15.10
	github.com/coopnorge/go-logger/adapter/goose.v3 v0.15.10
)

require (
	github.com/sirupsen/logrus v1.10.2 // indirect
	golang.org/x/sys v0.49.0 // indirect
)

replace github.com/coopnorge/go-logger => ../..

replace github.com/coopnorge/go-logger/adapter/goose.v3 => ../goose.v3
