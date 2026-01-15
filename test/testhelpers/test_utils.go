package testhelpers

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

// MockContext creates a basic context for testing
func MockContext() context.Context {
	return context.Background()
}

// MockLogger creates a basic logger for testing
func MockLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

// CreateTempDir creates a temporary directory for testing
func CreateTempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})
	return dir
}

// CreateTempFile creates a temporary file with content
func CreateTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	const fileMode = 0o644
	err := os.WriteFile(path, []byte(content), fileMode)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	return path
}

// AssertNoError asserts no error occurred
func AssertNoError(t *testing.T, err error, operation string) {
	t.Helper()
	if err != nil {
		t.Errorf("%s should not return error, got: %v", operation, err)
	}
}

// AssertError asserts an error occurred
func AssertError(t *testing.T, err error, operation string) {
	t.Helper()
	if err == nil {
		t.Errorf("%s should return error, got nil", operation)
	}
}
