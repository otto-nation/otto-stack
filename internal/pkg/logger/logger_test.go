//go:build unit

package logger

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	// Clear env var for test
	os.Unsetenv(EnvLogLevel)
	config := DefaultConfig()
	assert.Equal(t, LogLevelWarn, config.Level)
}

func TestDefaultConfigWithEnv(t *testing.T) {
	os.Setenv(EnvLogLevel, "debug")
	defer os.Unsetenv(EnvLogLevel)
	config := DefaultConfig()
	assert.Equal(t, LogLevelDebug, config.Level)
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected slog.Level
	}{
		{"debug", LogLevelDebug, slog.LevelDebug},
		{"info", LogLevelInfo, slog.LevelInfo},
		{"warn", LogLevelWarn, slog.LevelWarn},
		{"warning", LogLevelWarning, slog.LevelWarn},
		{"error", LogLevelError, slog.LevelError},
		{"invalid", "invalid", slog.LevelInfo},
		{"uppercase", "DEBUG", slog.LevelDebug},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseLogLevel(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInit(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{"default config", DefaultConfig()},
		{"debug level", Config{Level: LogLevelDebug}},
		{"error level", Config{Level: LogLevelError}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Init(tt.config)
			assert.NoError(t, err)
			assert.NotNil(t, defaultLogger)

			// Clean up
			Close()
		})
	}
}

func TestGetLogger(t *testing.T) {
	// Clean state
	Close()

	logger := GetLogger()
	assert.NotNil(t, logger)
	assert.NotNil(t, defaultLogger)
}

func TestLoggingFunctions(t *testing.T) {
	err := Init(Config{Level: LogLevelDebug})
	assert.NoError(t, err)

	// Test that functions don't panic
	assert.NotPanics(t, func() {
		Debug("debug message", "key", "value")
		Info("info message", "key", "value")
		Error("error message", "key", "value")
	})

	Close()
}

func TestClose(t *testing.T) {
	err := Init(DefaultConfig())
	assert.NoError(t, err)
	assert.NotNil(t, defaultLogger)

	Close()
	assert.Nil(t, defaultLogger)
}

func TestMultipleInit(t *testing.T) {
	err1 := Init(Config{Level: LogLevelDebug})
	assert.NoError(t, err1)

	err2 := Init(Config{Level: LogLevelError})
	assert.NoError(t, err2)

	// Should not panic
	logger := GetLogger()
	assert.NotNil(t, logger)

	Close()
}

func TestLoggerPersistence(t *testing.T) {
	err := Init(DefaultConfig())
	assert.NoError(t, err)

	logger1 := GetLogger()
	logger2 := GetLogger()

	// Should return the same instance
	assert.Equal(t, logger1, logger2)

	Close()
}
