package logger

import "errors"

// ErrInvalidNameInEnv to be used if a name is not found in the available levels
var ErrInvalidNameInEnv = errors.New("invalid name passed to available log levels")
