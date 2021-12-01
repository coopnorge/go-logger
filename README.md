# go-logger

This Go package is used to offer a unified logging interface among projects.

## Import

```
import log "github.com/coopnorge/go-logger"
```

## Example usage

### Logger struct init

```
var logger *log.Logger = &log.Logger{}
```

Possible log levels are: Trace, Debug, Info, Print, Warn, Warning, Error, Fatal, Panic

### Log level: Info

```
logger.Infof("Logging string foo %s)", stringFoo)
```

### Log level: Error

```
logger.Errorf("Logging error value foo: %v", errorFoo)
```

### Log level: Fatal

```
logger.Fatalf("Log app-breaking error value and propagating fatal error: %v", badError)
```
