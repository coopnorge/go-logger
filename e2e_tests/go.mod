module github.com/coopnorge/go-logger-e2e-tests

go 1.26.0

replace github.com/coopnorge/go-logger => ../

require (
	github.com/coopnorge/go-logger v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.12.0
)

require (
	github.com/sirupsen/logrus v1.9.4 // indirect
	golang.org/x/sys v0.44.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
