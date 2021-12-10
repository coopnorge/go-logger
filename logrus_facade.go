package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

// LogrusFacade abstracts logrus as a dependency
type LogrusFacade struct {
	logger *logrus.Logger
}

// Logrus represents the structure logger service we use
func Logrus(out io.Writer) *LogrusFacade {
	facade := &LogrusFacade{
		logger: logrus.New(),
	}

	facade.logger.Out = out
	facade.logger.Formatter = &logrus.JSONFormatter{}

	return facade
}

// LogWithFields accepts a map[string]interface{} and transforms it to what the logger expects
func (f *LogrusFacade) LogWithFields(fields Fields) Entry {
	logrusFields := make(logrus.Fields, len(fields))
	for k, v := range fields {
		logrusFields[k] = v
	}
	return f.logger.WithFields(logrusFields)
}

// LogPrint forwards a Print log call
func (f *LogrusFacade) LogPrint(args ...interface{}) {
	f.logger.Print(args...)
}

// LogError forwards a Error log call
func (f *LogrusFacade) LogError(args ...interface{}) {
	f.logger.Error(args...)
}

// LogTracef forawrds a Tracef log call
func (f *LogrusFacade) LogTracef(format string, args ...interface{}) {
	f.logger.Tracef(format, args...)
}

// LogDebugf forwards a Debugf log call
func (f *LogrusFacade) LogDebugf(format string, args ...interface{}) {
	f.logger.Debugf(format, args...)
}

// LogInfof forwards a Infof log call
func (f *LogrusFacade) LogInfof(format string, args ...interface{}) {
	f.logger.Infof(format, args...)
}

// LogPrintf forwards a Printf log call
func (f *LogrusFacade) LogPrintf(format string, args ...interface{}) {
	f.logger.Printf(format, args...)
}

// LogWarnf forwards a Warnf log call
func (f *LogrusFacade) LogWarnf(format string, args ...interface{}) {
	f.logger.Warnf(format, args...)
}

// LogErrorf forwards a Errorf log call
func (f *LogrusFacade) LogErrorf(format string, args ...interface{}) {
	f.logger.Errorf(format, args...)
}

// LogFatalf forwards a Fatalf log call
func (f *LogrusFacade) LogFatalf(format string, args ...interface{}) {
	f.logger.Fatalf(format, args...)
}

// LogPanicf forwards a Panicf log call
func (f *LogrusFacade) LogPanicf(format string, args ...interface{}) {
	f.logger.Panicf(format, args...)
}
