package gorm

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	coopLogger "github.com/coopnorge/go-logger"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestGormInit(t *testing.T) {
	logger, err := NewLogger(WithGlobalLogger())
	require.NoError(t, err)

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger,
	})
	require.NoError(t, err)
	require.NotNil(t, db)
}

func TestGlobalLogger(t *testing.T) {
	ctx := context.Background()

	output := &strings.Builder{}
	coopLogger.ConfigureGlobalLogger(coopLogger.WithLevel(coopLogger.LevelDebug), coopLogger.WithOutput(output))

	logger, err := NewLogger(WithGlobalLogger())
	require.NoError(t, err)

	output.Reset()
	logger.Info(ctx, "This is a test")
	assert.Contains(t, output.String(), "This is a test")

	output.Reset()
	logger.Warn(ctx, "This is a test")
	assert.Contains(t, output.String(), "This is a test")

	output.Reset()
	logger.Error(ctx, "This is a test")
	assert.Contains(t, output.String(), "This is a test")

	output.Reset()
	fc := func() (sql string, rowsAffected int64) {
		return "SQL STATEMENT", 10
	}
	logger.Trace(ctx, time.Now(), fc, nil)
	assert.Contains(t, output.String(), "\"rows\":10")

	output.Reset()
	fc = func() (sql string, rowsAffected int64) {
		return "SQL STATEMENT", 0
	}
	logger.Trace(ctx, time.Now(), fc, errors.New("This is a test"))
	assert.Contains(t, output.String(), "\"rows\":0")
	assert.Contains(t, output.String(), "This is a test")
}

func TestCustomLogger(t *testing.T) {
	ctx := context.Background()

	output := &strings.Builder{}
	logger, err := NewLogger(WithLogger(coopLogger.New(coopLogger.WithLevel(coopLogger.LevelDebug), coopLogger.WithOutput(output))))
	require.NoError(t, err)

	output.Reset()
	logger.Info(ctx, "This is a test")
	assert.Contains(t, output.String(), "This is a test")

	output.Reset()
	logger.Warn(ctx, "This is a test")
	assert.Contains(t, output.String(), "This is a test")

	output.Reset()
	logger.Error(ctx, "This is a test")
	assert.Contains(t, output.String(), "This is a test")

	output.Reset()
	fc := func() (sql string, rowsAffected int64) {
		return "SQL STATEMENT", 10
	}
	logger.Trace(ctx, time.Now(), fc, nil)
	assert.Contains(t, output.String(), "\"rows\":10")

	output.Reset()
	fc = func() (sql string, rowsAffected int64) {
		return "SQL STATEMENT", 0
	}
	logger.Trace(ctx, time.Now(), fc, errors.New("This is a test"))
	assert.Contains(t, output.String(), "\"rows\":0")
	assert.Contains(t, output.String(), "This is a test")
}
