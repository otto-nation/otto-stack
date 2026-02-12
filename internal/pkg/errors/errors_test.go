//go:build unit

package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected string
	}{
		{
			name:     "simple message",
			err:      &Error{Message: "test error"},
			expected: "test error",
		},
		{
			name:     "with context",
			err:      &Error{Context: "service", Message: "failed"},
			expected: "failed",
		},
		{
			name:     "with cause",
			err:      &Error{Message: "failed", Cause: errors.New("root cause")},
			expected: "failed: root cause",
		},
		{
			name:     "with context and cause",
			err:      &Error{Context: "docker", Message: "connection failed", Cause: errors.New("timeout")},
			expected: "connection failed: timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	cause := errors.New("root cause")
	err := &Error{Message: "wrapped", Cause: cause}
	assert.Equal(t, cause, err.Unwrap())
}

func TestNewSystemError(t *testing.T) {
	cause := errors.New("cause")
	err := NewSystemError(ErrCodeNotFound, "not found", cause)

	assert.Equal(t, ErrCodeNotFound, err.Code)
	assert.Equal(t, "", err.Context)
	assert.Equal(t, "not found", err.Message)
	assert.Equal(t, cause, err.Cause)
}

func TestNewSystemErrorf(t *testing.T) {
	err := NewSystemErrorf(ErrCodeInvalid, "invalid value: %s", "test")

	assert.Equal(t, ErrCodeInvalid, err.Code)
	assert.Equal(t, "", err.Context)
	assert.Equal(t, "invalid value: test", err.Message)
	assert.Nil(t, err.Cause)
}

func TestLegacyAliases(t *testing.T) {
	t.Run("NewValidationError", func(t *testing.T) {
		err := NewValidationError(ErrCodeInvalid, "field", "invalid", nil)
		assert.Equal(t, ErrCodeInvalid, err.Code)
		assert.Equal(t, "field", err.Context)
	})

	t.Run("NewValidationErrorf", func(t *testing.T) {
		err := NewValidationErrorf(ErrCodeInvalid, "field", "invalid: %s", "test")
		assert.Equal(t, "invalid: test", err.Message)
	})

	t.Run("NewServiceError", func(t *testing.T) {
		err := NewServiceError(ErrCodeOperationFail, "docker", "start failed", nil)
		assert.Equal(t, ErrCodeOperationFail, err.Code)
		assert.Equal(t, "docker", err.Context)
	})

	t.Run("NewConfigError", func(t *testing.T) {
		err := NewConfigError(ErrCodeNotFound, "/path", "not found", nil)
		assert.Equal(t, ErrCodeNotFound, err.Code)
		assert.Equal(t, "/path", err.Context)
	})

	t.Run("NewConfigErrorf", func(t *testing.T) {
		err := NewConfigErrorf(ErrCodeInvalid, "/path", "invalid: %s", "yaml")
		assert.Equal(t, "invalid: yaml", err.Message)
	})
}
