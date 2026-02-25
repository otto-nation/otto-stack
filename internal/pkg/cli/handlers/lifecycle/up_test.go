//go:build unit

package lifecycle

import (
	"context"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"github.com/otto-nation/otto-stack/test/testhelpers"
)

// Test constants - following DRY principles
const (
	// Using existing core constants
	expectedDefaultTimeout = core.DefaultStartTimeoutSeconds
)

func TestNewUpHandler(t *testing.T) {
	handler := NewUpHandler()

	assert.NotNil(t, handler)
	assert.IsType(t, &UpHandler{}, handler)
}

func TestUpHandler_ValidateArgs_NoArgs(t *testing.T) {
	handler := NewUpHandler()
	err := handler.ValidateArgs([]string{})
	assert.NoError(t, err, "Up command should accept no arguments")
}

func TestUpHandler_ValidateArgs_SingleService(t *testing.T) {
	handler := NewUpHandler()
	err := handler.ValidateArgs([]string{testhelpers.TestServiceName})
	assert.NoError(t, err, "Up command should accept service names")
}

func TestUpHandler_ValidateArgs_MultipleServices(t *testing.T) {
	handler := NewUpHandler()
	err := handler.ValidateArgs([]string{testhelpers.TestServiceName, "service2", "service3"})
	assert.NoError(t, err, "Up command should accept multiple service names")
}

func TestUpHandler_GetRequiredFlags(t *testing.T) {
	handler := NewUpHandler()
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags, "Up command should have no required flags")
}

func TestUpHandler_Handle(t *testing.T) {
	handler := NewUpHandler()
	cmd := &cobra.Command{
		Use: core.CommandUp,
	}

	cmd.Flags().Bool("build", false, "Build images before starting")
	cmd.Flags().Bool("detach", false, "Run in detached mode")

	base := &base.BaseCommand{
		Output: ui.NewOutput(),
	}

	ctx := context.Background()
	args := []string{}

	err := handler.Handle(ctx, cmd, args, base)
	assert.Error(t, err, "Expected error due to missing project initialization in test environment")
}

func TestDefaultTimeoutConstant(t *testing.T) {
	assert.Equal(t, expectedDefaultTimeout, core.DefaultStartTimeoutSeconds,
		"Should use core.DefaultStartTimeoutSeconds directly")
}

func TestNewDownHandler(t *testing.T) {
	handler := NewDownHandler()

	assert.NotNil(t, handler)
	assert.IsType(t, &DownHandler{}, handler)
}

func TestDownHandler_ValidateArgs_NoArgs(t *testing.T) {
	handler := NewDownHandler()
	err := handler.ValidateArgs([]string{})
	assert.NoError(t, err, "Down command should accept no arguments")
}

func TestDownHandler_ValidateArgs_SingleService(t *testing.T) {
	handler := NewDownHandler()
	err := handler.ValidateArgs([]string{testhelpers.TestServiceName})
	assert.NoError(t, err, "Down command should accept service names")
}

func TestDownHandler_ValidateArgs_MultipleServices(t *testing.T) {
	handler := NewDownHandler()
	err := handler.ValidateArgs([]string{testhelpers.TestServiceName, "service2"})
	assert.NoError(t, err, "Down command should accept multiple service names")
}

func TestDownHandler_GetRequiredFlags(t *testing.T) {
	handler := NewDownHandler()
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags, "Down command should have no required flags")
}

func TestDownHandler_Handle(t *testing.T) {
	// Redirect HOME to a fresh temp dir so DetectExecutionContext uses an empty
	// registry instead of the developer's real ~/.otto-stack/shared directory.
	t.Setenv("HOME", t.TempDir())

	handler := NewDownHandler()
	cmd := &cobra.Command{
		Use: core.CommandDown,
	}

	cmd.Flags().Bool("remove", false, "Remove containers")
	cmd.Flags().Bool("volumes", false, "Remove volumes")

	base := &base.BaseCommand{
		Output: ui.NewOutput(),
	}

	ctx := context.Background()
	args := []string{}

	err := handler.Handle(ctx, cmd, args, base)
	assert.NoError(t, err, "Down handler should handle global context without error")
}

func TestUpHandler_validateShareableServices_AllShareable(t *testing.T) {
	handler := NewUpHandler()
	serviceConfigs := []types.ServiceConfig{
		{Name: "redis", Shareable: true},
		{Name: "postgres", Shareable: true},
	}
	err := handler.validateShareableServices(serviceConfigs)
	require.NoError(t, err)
}

func TestUpHandler_validateShareableServices_NonShareable(t *testing.T) {
	handler := NewUpHandler()
	serviceConfigs := []types.ServiceConfig{
		{Name: "localstack-sqs", Shareable: false},
	}
	err := handler.validateShareableServices(serviceConfigs)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "localstack-sqs")
}

func TestUpHandler_validateShareableServices_Mixed(t *testing.T) {
	handler := NewUpHandler()
	serviceConfigs := []types.ServiceConfig{
		{Name: "redis", Shareable: true},
		{Name: "localstack-sqs", Shareable: false},
	}
	err := handler.validateShareableServices(serviceConfigs)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "localstack-sqs")
}

func TestUpHandler_validateShareableServices_Empty(t *testing.T) {
	handler := NewUpHandler()
	err := handler.validateShareableServices([]types.ServiceConfig{})
	require.NoError(t, err)
}

func TestFilterStatusQueryNames(t *testing.T) {
	tests := []struct {
		name     string
		configs  []types.ServiceConfig
		expected []string
	}{
		{
			name:     "empty",
			configs:  []types.ServiceConfig{},
			expected: []string{},
		},
		{
			name: "all persistent",
			configs: []types.ServiceConfig{
				{Name: "postgres"},
				{Name: "redis"},
			},
			expected: []string{"postgres", "redis"},
		},
		{
			name: "hidden provider included",
			configs: []types.ServiceConfig{
				{Name: "redis"},
				{Name: "localstack", Hidden: true},
			},
			expected: []string{"redis", "localstack"},
		},
		{
			name: "init container excluded",
			configs: []types.ServiceConfig{
				{Name: "postgres"},
				{Name: "init-job", Container: types.ContainerSpec{Restart: types.RestartPolicyNo}},
				{Name: "redis"},
			},
			expected: []string{"postgres", "redis"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterStatusQueryNames(tt.configs)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TODO: Add unit tests for buildContext method with various flag combinations
// TODO: Add tests for error handling scenarios
// TODO: Add tests for middleware chain execution
// TODO: Add E2E tests for full lifecycle up workflow
// TODO: Consider extracting common test utilities to reduce duplication across handler tests
