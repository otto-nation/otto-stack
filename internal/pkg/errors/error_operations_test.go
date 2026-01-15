//go:build unit

package errors

import (
	"errors"
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestError_formatters(t *testing.T) {
	t.Run("new validation errorf", func(t *testing.T) {
		err := NewValidationErrorf("field", "test %s", "message")
		testhelpers.AssertError(t, err, "NewValidationErrorf should return error")
		if err.Error() == "" {
			t.Error("NewValidationErrorf should return formatted error")
		}
	})

	t.Run("new service errorf", func(t *testing.T) {
		err := NewServiceErrorf("service", "test %s", "message")
		testhelpers.AssertError(t, err, "NewServiceErrorf should return error")
		if err.Error() == "" {
			t.Error("NewServiceErrorf should return formatted error")
		}
	})

	t.Run("new config errorf", func(t *testing.T) {
		err := NewConfigErrorf("key", "test %s", "message")
		testhelpers.AssertError(t, err, "NewConfigErrorf should return error")
		if err.Error() == "" {
			t.Error("NewConfigErrorf should return formatted error")
		}
	})

	t.Run("new docker errorf", func(t *testing.T) {
		err := NewDockerErrorf("operation", "container", "test %s", "message")
		testhelpers.AssertError(t, err, "NewDockerErrorf should return error")
		if err.Error() == "" {
			t.Error("NewDockerErrorf should return formatted error")
		}
	})
}

func TestError_unwrap(t *testing.T) {
	baseErr := errors.New("base error")

	t.Run("service error unwrap", func(t *testing.T) {
		err := NewServiceError("service", "message", baseErr)
		unwrapped := err.Unwrap()
		if unwrapped != baseErr {
			t.Error("ServiceError.Unwrap should return base error")
		}
	})

	t.Run("config error unwrap", func(t *testing.T) {
		err := NewConfigError("key", "message", baseErr)
		unwrapped := err.Unwrap()
		if unwrapped != baseErr {
			t.Error("ConfigError.Unwrap should return base error")
		}
	})
}
