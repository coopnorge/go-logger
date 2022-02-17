# go-logger

This Go package is used to offer a unified logging interface among projects.

## Import

```
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
logger.Warn("something went wrong")
// Output:
// {"file":"/Users/anonymous/Projects/my-project/main.go:7","function":"main.main","level":"warning","msg":"something went wrong","time":"2022-02-17T15:04:06+01:00"}
```

## Example usage

See `logger_examples_test.go`
