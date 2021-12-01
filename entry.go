package logger

// Entry represents a logging entry and all supported method we use
type Entry interface {
	Tracef(format string, args ...interface{})
	Trace(args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Printf(format string, args ...interface{})
	Print(args ...interface{})
	Warnf(format string, args ...interface{})
	Warn(args ...interface{})
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
	Panicf(format string, args ...interface{})
	Panic(args ...interface{})
}
