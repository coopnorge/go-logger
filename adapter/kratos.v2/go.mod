module github.com/coopnorge/go-logger/adapter/kratos.v2

go 1.27.3

require (
	github.com/coopnorge/go-logger v0.15.10
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/stretchr/testify v1.12.1
)

require (
	github.com/sirupsen/logrus v1.10.2 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.49.0 // indirect
)

replace github.com/coopnorge/go-logger => ../..
