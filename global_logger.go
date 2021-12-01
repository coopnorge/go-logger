package logger

var globalLogger = New()

func ConfigureGlobalLogger(opts ...LoggerOption) {
	globalLogger.applyOptions(opts...)
}

func Global() *Logger {
	return globalLogger
}

// WithFields creates log entry using global logger
func WithFields(fields Fields) Entry {
	return globalLogger.WithFields(fields)
}

// Print uses global logger to log payload on "info" level
func Print(args ...interface{}) {
	globalLogger.Print(args...)
}

// Printf uses global logger to log payload on "info" level
func Printf(format string, args ...interface{}) {
	globalLogger.Printf(format, args...)
}

// Infof uses global logger to log payload on "info" level
func Info(args ...interface{}) {
	globalLogger.Info(args...)
}

// Infof uses global logger to log payload on "info" level
func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

// Error uses global logger to log payload on "error" level
func Error(args ...interface{}) {
	globalLogger.Error(args...)
}

// Errorf uses global logger to log payload on "error" level
func Errorf(format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}

// Trace uses global logger to log payload on "trace" level
func Trace(args ...interface{}) {
	globalLogger.Trace(args...)
}

// Tracef uses global logger to log payload on "trace" level
func Tracef(format string, args ...interface{}) {
	globalLogger.Tracef(format, args...)
}

// Debug uses global logger to log payload on "debug" level
func Debug(args ...interface{}) {
	globalLogger.Debug(args...)
}

// Debugf uses global logger to log payload on "debug" level
func Debugf(format string, args ...interface{}) {
	globalLogger.Debugf(format, args...)
}

// Warn uses global logger to log payload on "warn" level
func Warn(args ...interface{}) {
	globalLogger.Warn(args...)
}

// Warnf uses global logger to log payload on "warn" level
func Warnf(format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}

// Fatal uses global logger to log payload on "fatal" level
func Fatal(args ...interface{}) {
	globalLogger.Fatal(args...)
}

// Fatalf uses global logger to log payload on "fatal" level
func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatalf(format, args...)
}

// Panic uses global logger to log payload on "panic" level
func Panic(args ...interface{}) {
	globalLogger.Panic(args...)
}

// Panicf uses global logger to log payload on "panic" level
func Panicf(format string, args ...interface{}) {
	globalLogger.Panicf(format, args...)
}
