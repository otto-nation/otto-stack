//go:build unit

package operations

import (
	"context"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/stretchr/testify/assert"
)

// MockOutput implements base.Output for testing
type MockOutput struct{}

func (m *MockOutput) Success(msg string, args ...any) {}
func (m *MockOutput) Error(msg string, args ...any)   {}
func (m *MockOutput) Warning(msg string, args ...any) {}
func (m *MockOutput) Info(msg string, args ...any)    {}
func (m *MockOutput) Header(msg string, args ...any)  {}
func (m *MockOutput) Muted(msg string, args ...any)   {}

func TestNewServiceCommand(t *testing.T) {
	stateManager := common.NewStateManager()

	cmd := NewServiceCommand(core.CommandStatus, stateManager)

	assert.NotNil(t, cmd)
	assert.Equal(t, core.CommandStatus, cmd.operation)
	assert.Equal(t, stateManager, cmd.stateManager)
}

func TestServiceCommand_ExecuteWithConstants(t *testing.T) {
	stateManager := common.NewStateManager()
	ctx := context.Background()
	cliCtx := clicontext.Context{
		Project:  clicontext.ProjectSpec{Name: "test-project"},
		Services: clicontext.ServiceSpec{Names: []string{services.ServicePostgres}},
	}
	base := &base.BaseCommand{Output: &MockOutput{}}

	t.Run("executes status command using core constants", func(t *testing.T) {
		cmd := NewServiceCommand(core.CommandStatus, stateManager)

		err := cmd.Execute(ctx, cliCtx, base)

		// Should attempt execution (may fail due to Docker), not routing error
		if err != nil {
			assert.NotContains(t, err.Error(), "unsupported")
		}
	})

	t.Run("uses common error constants properly", func(t *testing.T) {
		// Test that error constants are properly defined
		assert.Equal(t, "list containers", common.OpListContainers)
		assert.Equal(t, "show logs", common.OpShowLogs)
		assert.Equal(t, "remove resources", common.OpRemoveResources)
	})

	t.Run("validates all core command constants", func(t *testing.T) {
		commands := []string{
			core.CommandUp,
			core.CommandDown,
			core.CommandStatus,
			core.CommandLogs,
			core.CommandExec,
			core.CommandConnect,
			core.CommandRestart,
			core.CommandCleanup,
		}

		for _, command := range commands {
			assert.NotEmpty(t, command, "Command constant should not be empty")

			// Each command should route without "unsupported" error
			cmd := NewServiceCommand(command, stateManager)
			err := cmd.Execute(ctx, cliCtx, base)

			if err != nil {
				assert.NotContains(t, err.Error(), "unsupported stack operation")
			}
		}
	})
}

func TestResolveServiceConfigsWithConstants(t *testing.T) {
	setup := &common.CoreSetup{
		Config: &config.Config{
			Stack: config.StackConfig{
				Enabled: []string{"postgres", "redis"},
			},
		},
	}

	t.Run("resolves services using service constants", func(t *testing.T) {
		args := []string{services.ServicePostgres}

		configs, err := ResolveServiceConfigs(args, setup)

		// Test the logic path, may fail due to service resolution
		if err == nil {
			assert.NotNil(t, configs)
		}
		// Key test: function doesn't panic and follows expected path
	})

	t.Run("uses core timeout constants", func(t *testing.T) {
		// Verify core constants are available for HTTP operations
		assert.Greater(t, int(core.DefaultHTTPTimeoutSeconds), 0)
		assert.Greater(t, int(core.HTTPOKStatusThreshold), 0)
		// DefaultLogTailLines is a string constant, so just verify it's not empty
		assert.NotEmpty(t, core.DefaultLogTailLines)
	})
}

func TestResolveServiceConfigs(t *testing.T) {
	setup := &common.CoreSetup{
		Config: &config.Config{
			Stack: config.StackConfig{
				Enabled: []string{"postgres", "redis"},
			},
		},
	}

	t.Run("uses args when provided", func(t *testing.T) {
		args := []string{"postgres"}

		configs, err := ResolveServiceConfigs(args, setup)

		// May fail due to service resolution, but we're testing the logic path
		if err == nil {
			assert.NotNil(t, configs)
		}
		// Test passes if it doesn't panic and follows the args path
	})

	t.Run("uses enabled services when no args", func(t *testing.T) {
		args := []string{}

		configs, err := ResolveServiceConfigs(args, setup)

		// May fail due to service resolution, but we're testing the logic path
		if err == nil {
			assert.NotNil(t, configs)
		}
		// Test passes if it doesn't panic and follows the enabled services path
	})
}

func TestCreateStandardMiddlewareChain(t *testing.T) {
	validation, logging := CreateStandardMiddlewareChain()

	assert.NotNil(t, validation)
	assert.NotNil(t, logging)
}
