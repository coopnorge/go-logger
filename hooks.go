package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

// Hook defines the interface a custom Hook needs to implement
type Hook interface {
	Fire(*HookEntry) (changed bool, err error)
}

// HookFunc can be used to convert a simple function to implement the Hook interface.
type HookFunc func(*HookEntry) (changed bool, err error)

// Fire redirects a function call to the function receiver
func (hf HookFunc) Fire(he *HookEntry) (changed bool, err error) {
	return hf(he)
}

// HookEntry contains the fields provided for mutation in a hook.
type HookEntry struct {
	// Contains all the fields set by the user.
	Data Fields

	// Level the log entry was logged at: Trace, Debug, Info, Warn, Error, Fatal or Panic
	// This field will be set on entry firing and the value will be equal to the one in Logger struct field.
	Level Level

	// Message passed to Trace, Debug, Info, Warn, Error, Fatal or Panic
	Message string

	// Contains the context set by the user. Useful for hook processing etc.
	Context context.Context
}

type logrusHook struct {
	hook Hook
}

// Levels implements the logrus.Hook interface.
func (h *logrusHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire implements the logrus.Hook interface.
func (h *logrusHook) Fire(entry *logrus.Entry) error {
	// Provide all entry-data so the hook can mutate them.
	hookEntry := &HookEntry{
		Data:    Fields(entry.Data),
		Level:   mapLogrusLevelToLevel(entry.Level),
		Message: entry.Message,
		Context: entry.Context,
	}
	changed, err := h.hook.Fire(hookEntry)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}

	// Mutate the actual logrus entry with the mutations done in the hook.
	entry.Data = logrus.Fields(hookEntry.Data)
	entry.Level = mapLevelToLogrusLevel(hookEntry.Level)
	entry.Message = hookEntry.Message
	entry.Context = hookEntry.Context

	return nil
}
