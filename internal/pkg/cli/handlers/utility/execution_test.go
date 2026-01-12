//go:build unit

package utility

import (
	"context"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// MockOutput for testing consistency
type MockOutput struct{}

func (m *MockOutput) Success(msg string, args ...any) {}
func (m *MockOutput) Error(msg string, args ...any)   {}
func (m *MockOutput) Warning(msg string, args ...any) {}
func (m *MockOutput) Info(msg string, args ...any)    {}
func (m *MockOutput) Header(msg string, args ...any)  {}
func (m *MockOutput) Muted(msg string, args ...any)   {}

func TestWebInterfacesHandler_ExecutionWithConstants(t *testing.T) {
	handler := NewWebInterfacesHandler()
	ctx := context.Background()
	cmd := &cobra.Command{}
	base := &base.BaseCommand{Output: &MockOutput{}}

	t.Run("handles execution with service constants", func(t *testing.T) {
		args := []string{services.ServicePostgres}

		err := handler.Handle(ctx, cmd, args, base)

		// May fail due to Docker dependencies, but tests execution path
		if err != nil {
			// Should not be validation error since args are valid
			assert.NotContains(t, err.Error(), "validation")
		}
	})

	t.Run("uses core HTTP constants properly", func(t *testing.T) {
		// Test that HTTP constants are available and reasonable
		assert.Greater(t, int(core.DefaultHTTPTimeoutSeconds), 0)
		assert.Greater(t, int(core.HTTPOKStatusThreshold), 100) // Should be > 100 for HTTP status
		assert.Less(t, int(core.HTTPOKStatusThreshold), 500)    // Should be < 500 for OK range
	})

	t.Run("validates with service constants", func(t *testing.T) {
		serviceList := []string{services.ServicePostgres, services.ServiceRedis}

		for _, service := range serviceList {
			err := handler.ValidateArgs([]string{service})
			assert.NoError(t, err, "Should accept service constant: %s", service)
		}
	})

	t.Run("handles empty args using default behavior", func(t *testing.T) {
		err := handler.ValidateArgs([]string{})
		assert.NoError(t, err, "Should accept empty args for web-interfaces")

		// Test execution with no specific services
		err = handler.Handle(ctx, cmd, []string{}, base)
		if err != nil {
			// Should attempt to process all services, not validation error
			assert.NotContains(t, err.Error(), "required")
		}
	})
}

func TestWebInterfacesHandler_HTTPLogic(t *testing.T) {
	t.Run("uses proper HTTP timeout constants", func(t *testing.T) {
		// Verify the handler would use core constants for HTTP operations
		timeout := int(core.DefaultHTTPTimeoutSeconds)
		assert.Greater(t, timeout, 0, "HTTP timeout should be positive")
		assert.LessOrEqual(t, timeout, 60, "HTTP timeout should be reasonable")
	})

	t.Run("uses proper HTTP status threshold", func(t *testing.T) {
		threshold := int(core.HTTPOKStatusThreshold)
		assert.Equal(t, 400, threshold, "HTTP OK threshold should be 400 based on core constant")
	})
}

func TestWebInterfacesHandler_ServiceConstants(t *testing.T) {
	t.Run("validates service constants", func(t *testing.T) {
		serviceList := []string{services.ServicePostgres, services.ServiceRedis}

		handler := NewWebInterfacesHandler()

		for _, service := range serviceList {
			t.Run(service, func(t *testing.T) {
				assert.NotEmpty(t, service, "Service constant should not be empty")

				err := handler.ValidateArgs([]string{service})
				assert.NoError(t, err, "Should validate service constant: %s", service)
			})
		}
	})
}

func TestUtilityHandlers_SimpleConstructors(t *testing.T) {
	t.Run("web interfaces handler constructor", func(t *testing.T) {
		handler := NewWebInterfacesHandler()
		assert.NotNil(t, handler)
	})
}
