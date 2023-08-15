package logger

import "io"

// LoggerOption defines an applicator interface
type LoggerOption interface { //nolint:all
	Apply(l *Logger)
}

// LoggerOptionFunc defines a function which modifies a logger
type LoggerOptionFunc func(l *Logger) //nolint:all

// Apply redirects a function call to the function receiver
func (lof LoggerOptionFunc) Apply(l *Logger) {
	lof(l)
}

// WithNowFunc overrides default function used to determine current time.
// Intended to be used in tests only.
func WithNowFunc(nowFunc NowFunc) LoggerOption {
	return LoggerOptionFunc(func(l *Logger) {
		l.now = nowFunc
	})
}

// WithOutput overrides default output the logs are written to.
func WithOutput(output io.Writer) LoggerOption {
	return LoggerOptionFunc(func(l *Logger) {
		l.output = output
	})
}

// WithLevel sets minimum level for filtering logs
func WithLevel(level Level) LoggerOption {
	return LoggerOptionFunc(func(l *Logger) {
		l.level = level
	})
}

// WithReportCaller allows enabling/disabling including calling method in the log entry
func WithReportCaller(enable bool) LoggerOption {
	return LoggerOptionFunc(func(l *Logger) {
		l.reportCaller = enable
	})
}
