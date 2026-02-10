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
			expected: "service: failed",
		},
		{
			name:     "with cause",
			err:      &Error{Message: "failed", Cause: errors.New("root cause")},
			expected: "failed: root cause",
		},
		{
			name:     "with context and cause",
			err:      &Error{Context: "docker", Message: "connection failed", Cause: errors.New("timeout")},
			expected: "docker: connection failed: timeout",
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

func TestNew(t *testing.T) {
	cause := errors.New("cause")
	err := New(ErrCodeNotFound, "service", "not found", cause)

	assert.Equal(t, ErrCodeNotFound, err.Code)
	assert.Equal(t, "service", err.Context)
	assert.Equal(t, "not found", err.Message)
	assert.Equal(t, cause, err.Cause)
}

func TestNewf(t *testing.T) {
	err := Newf(ErrCodeInvalid, "field", "invalid value: %s", "test")

	assert.Equal(t, ErrCodeInvalid, err.Code)
	assert.Equal(t, "field", err.Context)
	assert.Equal(t, "invalid value: test", err.Message)
	assert.Nil(t, err.Cause)
}

func TestGetErrorCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: "",
		},
		{
			name:     "custom error",
			err:      &Error{Code: ErrCodeNotFound},
			expected: ErrCodeNotFound,
		},
		{
			name:     "standard error",
			err:      errors.New("standard"),
			expected: ErrCodeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, GetErrorCode(tt.err))
		})
	}
}

func TestIsNotFound(t *testing.T) {
	assert.True(t, IsNotFound(&Error{Code: ErrCodeNotFound}))
	assert.False(t, IsNotFound(&Error{Code: ErrCodeInvalid}))
	assert.False(t, IsNotFound(nil))
}

func TestIsAlreadyExists(t *testing.T) {
	assert.True(t, IsAlreadyExists(&Error{Code: ErrCodeAlreadyExists}))
	assert.False(t, IsAlreadyExists(&Error{Code: ErrCodeNotFound}))
}

func TestIsRetryable(t *testing.T) {
	assert.True(t, IsRetryable(&Error{Code: ErrCodeTimeout}))
	assert.True(t, IsRetryable(&Error{Code: ErrCodeUnavailable}))
	assert.False(t, IsRetryable(&Error{Code: ErrCodeNotFound}))
}

func TestIsPermissionDenied(t *testing.T) {
	assert.True(t, IsPermissionDenied(&Error{Code: ErrCodePermission}))
	assert.False(t, IsPermissionDenied(&Error{Code: ErrCodeNotFound}))
}

func TestIsInvalid(t *testing.T) {
	assert.True(t, IsInvalid(&Error{Code: ErrCodeInvalid}))
	assert.False(t, IsInvalid(&Error{Code: ErrCodeNotFound}))
}

func TestFormatForUser(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		contains string
	}{
		{
			name:     "nil error",
			err:      nil,
			contains: "",
		},
		{
			name:     "not found with hint",
			err:      &Error{Code: ErrCodeNotFound, Message: "service not found"},
			contains: "Hint: Check the name",
		},
		{
			name:     "unavailable with hint",
			err:      &Error{Code: ErrCodeUnavailable, Message: "service unavailable"},
			contains: "Hint: Ensure required services",
		},
		{
			name:     "permission with hint",
			err:      &Error{Code: ErrCodePermission, Message: "permission denied"},
			contains: "Hint: Check file permissions",
		},
		{
			name:     "other error",
			err:      &Error{Code: ErrCodeInternal, Message: "internal error"},
			contains: "internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatForUser(tt.err)
			if tt.contains != "" {
				assert.Contains(t, result, tt.contains)
			} else {
				assert.Empty(t, result)
			}
		})
	}
}

func TestGetExitCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: ExitSuccess,
		},
		{
			name:     "not found",
			err:      &Error{Code: ErrCodeNotFound},
			expected: ExitNotFound,
		},
		{
			name:     "invalid",
			err:      &Error{Code: ErrCodeInvalid},
			expected: ExitInvalidInput,
		},
		{
			name:     "permission",
			err:      &Error{Code: ErrCodePermission},
			expected: ExitPermissionError,
		},
		{
			name:     "already exists",
			err:      &Error{Code: ErrCodeAlreadyExists},
			expected: ExitAlreadyExists,
		},
		{
			name:     "timeout",
			err:      &Error{Code: ErrCodeTimeout},
			expected: ExitTimeout,
		},
		{
			name:     "unavailable",
			err:      &Error{Code: ErrCodeUnavailable},
			expected: ExitUnavailable,
		},
		{
			name:     "unknown",
			err:      &Error{Code: ErrCodeUnknown},
			expected: ExitGeneralError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, GetExitCode(tt.err))
		})
	}
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
