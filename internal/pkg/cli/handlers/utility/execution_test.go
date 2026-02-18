//go:build unit

package utility

import (
	"context"
	"io"
	"os"
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
func (m *MockOutput) Writer() io.Writer               { return os.Stdout }

func TestWebInterfacesHandler_ExecutionWithConstants(t *testing.T) {
	handler := NewWebInterfacesHandler()
	ctx := context.Background()
	cmd := &cobra.Command{}
	base := &base.BaseCommand{Output: &MockOutput{}}

	args := []string{services.ServicePostgres}
	err := handler.Handle(ctx, cmd, args, base)
	if err != nil {
		assert.NotContains(t, err.Error(), "validation")
	}
}

func TestWebInterfacesHandler_HTTPConstants(t *testing.T) {
	assert.Greater(t, int(core.DefaultHTTPTimeoutSeconds), 0)
	assert.Greater(t, int(core.HTTPOKStatusThreshold), 100)
	assert.Less(t, int(core.HTTPOKStatusThreshold), 500)
}

func TestWebInterfacesHandler_ValidatesServiceConstants(t *testing.T) {
	handler := NewWebInterfacesHandler()
	serviceList := []string{services.ServicePostgres, services.ServiceRedis}

	for _, service := range serviceList {
		err := handler.ValidateArgs([]string{service})
		assert.NoError(t, err, "Should accept service constant: %s", service)
	}
}

func TestWebInterfacesHandler_EmptyArgs(t *testing.T) {
	handler := NewWebInterfacesHandler()
	ctx := context.Background()
	cmd := &cobra.Command{}
	base := &base.BaseCommand{Output: &MockOutput{}}

	err := handler.ValidateArgs([]string{})
	assert.NoError(t, err, "Should accept empty args for web-interfaces")

	err = handler.Handle(ctx, cmd, []string{}, base)
	if err != nil {
		assert.NotContains(t, err.Error(), "required")
	}
}

func TestWebInterfacesHandler_HTTPTimeout(t *testing.T) {
	timeout := int(core.DefaultHTTPTimeoutSeconds)
	assert.Greater(t, timeout, 0, "HTTP timeout should be positive")
	assert.LessOrEqual(t, timeout, 60, "HTTP timeout should be reasonable")
}

func TestWebInterfacesHandler_HTTPStatusThreshold(t *testing.T) {
	threshold := int(core.HTTPOKStatusThreshold)
	assert.Equal(t, 400, threshold, "HTTP OK threshold should be 400 based on core constant")
}

func TestWebInterfacesHandler_Constructor(t *testing.T) {
	handler := NewWebInterfacesHandler()
	assert.NotNil(t, handler)
}
