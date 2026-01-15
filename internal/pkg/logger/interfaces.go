package logger

import "log/slog"

// Adapter interface for accessing underlying slog.Logger
type Adapter interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	SlogLogger() *slog.Logger
}
