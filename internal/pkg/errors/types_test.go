//go:build unit

package errors

import (
	"errors"
	"testing"
)

func TestValidationError(t *testing.T) {
	t.Run("creates validation error without cause", func(t *testing.T) {
		err := NewValidationError("field", "message", nil)
		expected := "validation failed for field: message"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("creates validation error with cause", func(t *testing.T) {
		cause := errors.New("underlying error")
		err := NewValidationError("field", "message", cause)
		expected := "validation failed for field: message: underlying error"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
		if err.Unwrap() != cause {
			t.Errorf("expected unwrap to return cause")
		}
	})
}

func TestServiceError(t *testing.T) {
	t.Run("creates service error", func(t *testing.T) {
		cause := errors.New("connection failed")
		err := NewServiceError("postgres", "start", cause)
		expected := "service postgres failed to start: connection failed"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
		if err.Unwrap() != cause {
			t.Errorf("expected unwrap to return cause")
		}
	})
}

func TestConfigError(t *testing.T) {
	t.Run("creates config error with path", func(t *testing.T) {
		cause := errors.New("invalid yaml")
		err := NewConfigError("/path/to/config.yaml", "parse failed", cause)
		expected := "config error in /path/to/config.yaml: parse failed: invalid yaml"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("creates config error without path", func(t *testing.T) {
		cause := errors.New("invalid yaml")
		err := NewConfigError("", "parse failed", cause)
		expected := "config error: parse failed: invalid yaml"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})
}

func TestDockerError(t *testing.T) {
	t.Run("creates docker error with container", func(t *testing.T) {
		cause := errors.New("container not found")
		err := NewDockerError("start", "postgres-1", cause)
		expected := "docker start failed for container postgres-1: container not found"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("creates docker error without container", func(t *testing.T) {
		cause := errors.New("daemon not running")
		err := NewDockerError("connect", "", cause)
		expected := "docker connect failed: daemon not running"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})
}
