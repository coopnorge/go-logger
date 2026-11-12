# go-logger

This Go package is used to offer a unified logging interface among projects.

## Install

```shell
go get github.com/coopnorge/go-logger
```

## Import

```go
import "github.com/coopnorge/go-logger"
```

## Default behavior

By default, all logs will include:

- Log level
- Full path to file which called the logger, and line number
- Signature of the function that called the logger
- Timestamp of the log entry

Example log entry with default settings:

```go
package main

import "github.com/coopnorge/go-logger"

func main() {
	logger.Warn("something went wrong")
	// Output:
	// {"file":"/Users/anonymous/Projects/my-project/main.go:7","function":"main.main","level":"warning","msg":"something went wrong","time":"2022-02-17T15:04:06+01:00"}
}
```

## Example usage

See [logger_examples_test.go](logger_examples_test.go) for more examples.

### Using global logger

```go
package main

import "github.com/coopnorge/go-logger"

func main() {
	logger.Info("this won't be logged because the default log level is higher than info")
	logger.Warn("but this will be logged")
	// Output:
	// {"level":"warning","msg":"but this will be logged","time":"2022-02-17T11:01:28+01:00"}
}
```

### Setting log level

```go
package main

import "github.com/coopnorge/go-logger"

func main() {
	// global logger
	logger.Info("this won't be logged because the default log level is higher than info")
	logger.ConfigureGlobalLogger(logger.WithLevel(logger.LevelInfo))
	logger.Info("now this will be logged")
	// Output:
	// {"level":"info","msg":"now this will be logged","time":"2022-02-17T10:54:54+01:00"}

	// logger instance
	prodLogger := logger.New(logger.WithLevel(logger.LevelWarn))
	prodLogger.Info("this won't be logged because prodLogger's level is set to Warn...")
	prodLogger.Error("...but this will, because Error >= Warn")
	// Output:
	// {"level":"error","msg":"...but this will, because Error >= Warn","time":"2022-02-17T10:54:54+01:00"}

	debugLogger := logger.New(logger.WithLevel(logger.LevelDebug))
	debugLogger.Debug("this logger will log anything as Debug is the lowest available level")
	debugLogger.Warn("and this will be logged too")
	// Output:
	// {"level":"debug","msg":"this logger will log anything as Debug is the lowest available level","time":"2022-02-17T10:54:54+01:00"}
	// {"level":"warning","msg":"and this will be logged too","time":"2022-02-17T10:54:54+01:00"}
}
```

## Adapters

Adapters are separate Go modules. The core module `github.com/coopnorge/go-logger`
depends only on Logrus, so an application only downloads the libraries of the
adapters it uses. Every adapter is released together with the core module and
has the same version number. Install the adapter at the same version as the core
module:

```shell
go get github.com/coopnorge/go-logger@vX.Y.Z
go get github.com/coopnorge/go-logger/adapter/<name>@vX.Y.Z
```

The import path of an adapter is `github.com/coopnorge/go-logger/adapter/<name>`.
When the adapted library has a major version above 1, each major version has its
own module with a dotted suffix: `adapter/echo.v4` for Echo v4,
`adapter/goose.v3` for Goose v3 and `adapter/kratos.v2` for Kratos v2. Gorm has
no major version in its import path, so `adapter/gorm` has no suffix. The
package names are `echo`, `goose` and `logs`, so set an import alias.

### Echo

```shell
go get github.com/coopnorge/go-logger/adapter/echo.v4@vX.Y.Z
```

To ensure that Echo outputs logs in the correct format, set the logger of the
Echo instance to the wrapped logger.

```go
package main

import (
	echoLogger "github.com/coopnorge/go-logger/adapter/echo.v4"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Logger = echoLogger.NewWrappedEchoLogger()
}
```

### Gorm

```shell
go get github.com/coopnorge/go-logger/adapter/gorm@vX.Y.Z
```

To ensure that Gorm outputs logs in the correct format Gorm must be configured
with a [custom logger](https://gorm.io/docs/logger.html#Customize-Logger).

```go
package main

import (
	gormLogger "github.com/coopnorge/go-logger/adapter/gorm"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	l, err := gormLogger.NewLogger(gormLogger.WithGlobalLogger())
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: l,
	})
	if err != nil {
		panic(err)
	}
}
```

### Kratos

```shell
go get github.com/coopnorge/go-logger/adapter/kratos.v2@vX.Y.Z
```

To ensure that Kratos outputs logs in the correct format Kratos must be
configured with a custom logger.

```go
package main

import (
	"github.com/coopnorge/go-logger"
	logs "github.com/coopnorge/go-logger/adapter/kratos.v2"
	"github.com/go-kratos/kratos/v2/log"
)

func main() {
	log.SetLogger(logs.NewLoggerKratosAdapter(logger.Global()))
}
```

### Goose

```shell
go get github.com/coopnorge/go-logger/adapter/goose.v3@vX.Y.Z
```

To ensure that Goose outputs logs in the correct format, Goose must be
configured with a custom logger.

```go
package main

import (
	"github.com/coopnorge/go-logger"
	gooseLogger "github.com/coopnorge/go-logger/adapter/goose.v3"
	"github.com/pressly/goose/v3"
)

func main() {
	goose.SetLogger(gooseLogger.New(logger.Global()))
}
```

## Migrating to v0.16.0

Before v0.16.0 the adapters were packages in the core module, so the core
module pulled in Gorm, Kratos, Goose, Gommon and SQLite for every consumer.
From v0.16.0 each adapter is its own module and the core module only depends on
Logrus.

The Echo, Goose and Kratos adapters have new module paths, see below. If your
application imports the Gorm adapter, run `go get` for it at the same version as
the core module:

```shell
go get github.com/coopnorge/go-logger@v0.16.0
go get github.com/coopnorge/go-logger/adapter/gorm@v0.16.0
```

Applications that only import `github.com/coopnorge/go-logger` need no change
and get fewer indirect dependencies. If `go build` reports an ambiguous import
after the upgrade, the adapter's `go get` is missing.

## Migrating to versioned adapter modules

The Echo, Goose and Kratos adapters moved to modules that match the major
version of the library they adapt:

| Old import path                                 | New import path                                    |
| ----------------------------------------------- | -------------------------------------------------- |
| `github.com/coopnorge/go-logger/adapter/echo`   | `github.com/coopnorge/go-logger/adapter/echo.v4`   |
| `github.com/coopnorge/go-logger/adapter/goose`  | `github.com/coopnorge/go-logger/adapter/goose.v3`  |
| `github.com/coopnorge/go-logger/adapter/kratos` | `github.com/coopnorge/go-logger/adapter/kratos.v2` |

The old modules are deprecated. They still work and forward every call to the
new modules. The package name of `adapter/echo.v4` is `echo`, so alias the
import when you also import `github.com/labstack/echo/v4`.

To migrate, get the new module and run the `inline` analyzer of `go fix`. It
rewrites the old imports and calls to the new ones:

```shell
go get github.com/coopnorge/go-logger/adapter/<name>.vN@vX.Y.Z
go fix -inline ./...
go mod tidy
```

`go mod tidy` drops the deprecated module from your `go.mod` once nothing
imports it. If the rewrite leaves two imports of the same new path under
different names, remove one by hand or run `goimports`. Without `go fix`,
replace the import path and keep the identifiers, which are the same in the new
modules.

## Hooks

Hooks are functions that are triggered on all log-entries and allow for data to
be collected and added to log-entry.

To configure a new Hook for the global logger call
`logger.ConfigureGlobalLogger(opts ...logger.LoggerOption)` and pass the
`LoggerOption` created by `logger.WithHook(hook logger.Hook)` or
`logger.WithHookFunc(hook logger.HookFunc)`.

The function `logger.WithHook` takes a pointer to a struct that implement the
`logger.Hook` interface, the function `logger.WithHookFunc` takes a function
with a signature that matches the `logger.HookFunc`.

```go
package logger

// Hook defines the interface a custom Hook needs to implement
type Hook interface {
	Fire(*HookEntry) (changed bool, err error)
}

// HookFunc can be used to convert a simple function to implement the Hook interface.
type HookFunc func(*HookEntry) (changed bool, err error)
```

The struct `logger.HookEntry` contains the fields provided for mutation in a
Hook.

A typical use case for a Hook is to extract data from the `context.Context` set
by the user using `logger.WithContext(ctx context.Context)`. The data in the
context set by the user may have been set further up in the call stack of where
the log-entry is created.

### Example - username logging hook

```go title="app/userhook/hook.go"
package userhook

import (
	"github.com/coopnorge/go-logger"
)

type UserContextLogHook struct{}

type UserKey struct{}

// Fire implements logger.Hook interface
func (u *UserContextLogHook) Fire(he *logger.HookEntry) (bool, error) {
	ctx := he.Context
	if ctx == nil {
		return false, nil
	}

	value, ok := ctx.Value(UserKey{}).(string)
	if !ok || value == "" {
		return false, nil
	}

	he.Data["user"] = value

	return true, nil
}

func NewHook() logger.Hook {
	return &UserContextLogHook{}
}
```

```go title="app/main.go"
package main

import (
	"github.com/coopnorge/app/userhook"

	"github.com/coopnorge/go-telemetry-lib/loghook"
	"github.com/coopnorge/go-logger"
)

func main() {
	ctx := context.Background()
	logger.ConfigureGlobalLogger(
		logger.WithHook(userhook.NewHook()),
		logger.WithHook(loghook.NewHook()),
	)
	a(ctx)
}

func a(ctx context.Context) {
	username := "peter"
	ctx := context.WithValue(ctx, userhook.UserKey{}, username)
	b(ctx)
}

func b(ctx context.Context) {
	logger.WithContext(ctx).Warn("Hello")
	// Output:
	// {"user": "peter", "level":"warning","msg":"hello","time":"2024-09-16T09:09:00+01:00"}
}
```

### Known Hooks

- `github.com/coopnorge/go-telemetry-lib/loghook.Hook`: relates log entries to a
  span by using `context.Context`.
  [Documentation](https://inventory.internal.coop/docs/default/component/go-telemetry-lib/getting-started/#step-2-configure-the-log-hook-required)
