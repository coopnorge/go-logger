package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Fields type, used to pass to `WithFields`.
type Fields map[string]any

// NowFunc is a typedef for a function which returns the current time
type NowFunc func() time.Time

// Logger is our logger with the needed structured logger we use
type Logger struct {
	zapLogger                 *zap.Logger
	now                       NowFunc
	output                    io.Writer
	level                     Level
	reportCaller              bool
	attemptConsistentOrdering bool
	hooks                     []Hook
}

// New creates and returns a new logger with supplied options. Before it is discarded or the application exits, Flush should be called.
func New(opts ...LoggerOption) *Logger {
	logger := &Logger{
		// now:          NowFunc(time.Now),
		output:       os.Stdout,
		level:        LevelWarn,
		reportCaller: true,
	}
	for _, opt := range opts {
		opt.Apply(logger)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = customLevelEncoder()

	encoder := zapcore.NewJSONEncoder(encoderConfig)

	writeSyncer := zapcore.AddSync(logger.output)

	levelEnabler := zap.NewAtomicLevelAt(mapLevelToZapLevel(logger.level))

	core := zapcore.NewCore(encoder, writeSyncer, levelEnabler)

	zapOpts := make([]zap.Option, 0, 1)
	if logger.now != nil {
		zapOpts = append(zapOpts, zap.WithClock(functionClock{now: logger.now}))
	}

	zapLogger := zap.New(core, zapOpts...)
	logger.zapLogger = zapLogger
	return logger
}

func customLevelEncoder() zapcore.LevelEncoder {
	// This LevelEncoder is consistent with what we had with logrus.
	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		str := l.String()
		if str == "warn" {
			str = "warning"
		}
		enc.AppendString(str)
	}
}

func (logger *Logger) Flush() error {
	return logger.zapLogger.Sync()
}

func (logger *Logger) entry() *Entry {
	fields := Fields{}
	if logger.reportCaller {
		frame := getCaller()
		fields["file"] = fmt.Sprintf("%s:%v", frame.File, frame.Line)
		fields["function"] = frame.Function
	}

	return &Entry{logger: logger, fields: fields}
}

const errorKey = "error"

// WithError is a convenience wrapper for WithField("error", err)
func (logger *Logger) WithError(err error) *Entry {
	return logger.entry().WithError(err)
}

// WithField forwards a logging call with a field
func (logger *Logger) WithField(key string, value any) *Entry {
	return logger.entry().WithField(key, value)
}

// WithFields forwards a logging call with fields
func (logger *Logger) WithFields(fields Fields) *Entry {
	return logger.entry().WithFields(fields)
}

// WithContext forwards a logging call with fields
func (logger *Logger) WithContext(ctx context.Context) *Entry {
	return logger.entry().WithContext(ctx)
}

// OutputHandler returns logger output handler
func (logger *Logger) OutputHandler() io.Writer {
	return logger.output
}

// Info forwards a logging call in the (format, args) format
func (logger *Logger) Info(msg string) {
	logger.entry().Info(msg)
}

// Infof forwards a logging call in the (format, args) format
func (logger *Logger) Infof(format string, args ...any) {
	logger.entry().Infof(format, args...)
}

// Error forwards an error logging call
func (logger *Logger) Error(msg string) {
	logger.entry().Error(msg)
}

// Errorf forwards an error logging call
func (logger *Logger) Errorf(format string, args ...any) {
	logger.entry().Errorf(format, args...)
}

// Debug forwards a debugging logging call
func (logger *Logger) Debug(msg string) {
	logger.entry().Debug(msg)
}

// Debugf forwards a debugging logging call
func (logger *Logger) Debugf(format string, args ...any) {
	logger.entry().Debugf(format, args...)
}

// Warn forwards a warning logging call
func (logger *Logger) Warn(msg string) {
	logger.entry().Warn(msg)
}

// Warnf forwards a warning logging call
func (logger *Logger) Warnf(format string, args ...any) {
	logger.entry().Warnf(format, args...)
}

// Fatal forwards a fatal logging call
func (logger *Logger) Fatal(msg string) {
	logger.entry().Fatal(msg)
}

// Fatalf forwards a fatal logging call
func (logger *Logger) Fatalf(format string, args ...any) {
	logger.entry().Fatalf(format, args...)
}
