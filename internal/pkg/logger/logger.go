package logger

import (
	"log/slog"
	"os"
	"strings"
)

var defaultLogger *slog.Logger

// Config represents minimal logger configuration
type Config struct {
	Level string `yaml:"level" json:"level"`
}

// DefaultConfig returns default logger configuration
func DefaultConfig() Config {
	return Config{
		Level: LogLevelInfo,
	}
}

// Init initializes the logger with the given configuration
func Init(config Config) error {
	level := parseLogLevel(config.Level)
	opts := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	defaultLogger = slog.New(handler)
	return nil
}

// Close closes the logger (simplified - no cleanup needed)
func Close() {
	defaultLogger = nil
}

// parseLogLevel converts string to slog.Level
func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case LogLevelDebug:
		return slog.LevelDebug
	case LogLevelInfo:
		return slog.LevelInfo
	case LogLevelWarn, "warning":
		return slog.LevelWarn
	case LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// GetLogger returns the default logger instance
func GetLogger() *slog.Logger {
	if defaultLogger == nil {
		if err := Init(DefaultConfig()); err != nil {
			defaultLogger = slog.Default()
		}
	}
	return defaultLogger
}

// Debug logs a debug message
func Debug(msg string, args ...any) {
	GetLogger().Debug(msg, args...)
}

// Info logs an info message
func Info(msg string, args ...any) {
	GetLogger().Info(msg, args...)
}

// Warn logs a warning message
func Warn(msg string, args ...any) {
	GetLogger().Warn(msg, args...)
}

// Error logs an error message
func Error(msg string, args ...any) {
	GetLogger().Error(msg, args...)
}

// With returns a logger with the given attributes
func With(args ...any) *slog.Logger {
	return GetLogger().With(args...)
}
