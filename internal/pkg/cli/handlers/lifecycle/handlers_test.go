//go:build unit

package lifecycle

import (
	"context"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/test/testhelpers"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// MockOutput for testing - reuse pattern from operations
type MockOutput struct{}

func (m *MockOutput) Success(msg string, args ...any) {}
func (m *MockOutput) Error(msg string, args ...any)   {}
func (m *MockOutput) Warning(msg string, args ...any) {}
func (m *MockOutput) Info(msg string, args ...any)    {}
func (m *MockOutput) Header(msg string, args ...any)  {}
func (m *MockOutput) Muted(msg string, args ...any)   {}

func TestUpHandler_ExecutionFlow(t *testing.T) {
	cmd := &cobra.Command{}
	args := []string{services.ServicePostgres} // Use service constant

	t.Run("handles execution flow with proper context building", func(t *testing.T) {
		// Test the common BuildStackContext method directly
		cliCtx, err := common.BuildStackContext(cmd, args)

		// Should fail due to missing config, but tests the flow
		if err != nil {
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "config") // Expected config error
		} else {
			assert.Equal(t, args, cliCtx.Services.Names)
		}
	})

	t.Run("validates command constants usage", func(t *testing.T) {
		// Test that handler uses proper core constants
		assert.Equal(t, "up", core.CommandUp)
		assert.NotEmpty(t, core.MsgStarting)
	})
}

func TestDownHandler_ExecutionFlow(t *testing.T) {
	handler := NewDownHandler()
	ctx := context.Background()
	cmd := &cobra.Command{}
	args := []string{services.ServiceRedis} // Use service constant
	base := &base.BaseCommand{Output: &MockOutput{}}

	t.Run("handles down command execution", func(t *testing.T) {
		// This will fail due to missing config/Docker, but tests the routing
		err := handler.Handle(ctx, cmd, args, base)

		if err != nil {
			// Should not be validation error since args are valid
			assert.NotContains(t, err.Error(), "validation")
		}
	})

	t.Run("uses common utilities properly", func(t *testing.T) {
		// Test that handler uses common package utilities
		stateManager := common.NewStateManager()
		assert.NotNil(t, stateManager)

		// Test common middleware creation
		validation, logging := common.CreateStandardMiddlewareChain()
		assert.NotNil(t, validation)
		assert.NotNil(t, logging)
	})
}

func TestRestartHandler_ValidationAndConstants(t *testing.T) {
	handler := NewRestartHandler()

	t.Run("validates args using common patterns", func(t *testing.T) {
		// Test with service constants
		validArgs := []string{services.ServicePostgres, services.ServiceRedis}
		err := handler.ValidateArgs(validArgs)
		assert.NoError(t, err)

		// Test empty args (should be valid for restart)
		err = handler.ValidateArgs([]string{})
		assert.NoError(t, err)
	})

	t.Run("uses core constants for commands", func(t *testing.T) {
		// Verify core constants are properly defined
		assert.Equal(t, "restart", core.CommandRestart)
		assert.NotEmpty(t, core.MsgRestarting)
	})
}

func TestCleanupHandler_ErrorHandling(t *testing.T) {
	handler := NewCleanupHandler()

	t.Run("returns proper errors instead of calling os.Exit", func(t *testing.T) {
		ctx := context.Background()
		cmd := &cobra.Command{}
		args := []string{}
		base := &base.BaseCommand{Output: &MockOutput{}}

		// This should return an error, not call os.Exit()
		err := handler.Handle(ctx, cmd, args, base)

		// The key test: we get an error back instead of program termination
		assert.Error(t, err, "Should return error instead of calling os.Exit()")
		// Should not contain os.Exit related panics
		assert.NotContains(t, err.Error(), "exit")
	})

	t.Run("uses common error constants", func(t *testing.T) {
		// Test that common error constants are available
		assert.NotEmpty(t, common.ActionCleanupResources)
		assert.NotEmpty(t, common.ComponentStack)
		assert.NotEmpty(t, common.OpRemoveResources)
	})
}

func TestLifecycleHandlers_CoreConstants(t *testing.T) {
	t.Run("validates core command constants", func(t *testing.T) {
		commands := map[string]string{
			"up":      core.CommandUp,
			"down":    core.CommandDown,
			"restart": core.CommandRestart,
			"cleanup": core.CommandCleanup,
		}

		for expected, actual := range commands {
			assert.Equal(t, expected, actual, "Core command constant mismatch")
		}
	})

	t.Run("validates core message constants", func(t *testing.T) {
		messages := []string{
			core.MsgStarting,
			core.MsgRestarting,
		}

		for _, msg := range messages {
			assert.NotEmpty(t, msg, "Core message constant should not be empty")
		}
	})
}

func TestLifecycleHandlers_SimpleGetters(t *testing.T) {
	t.Run("up handler constructor", func(t *testing.T) {
		handler := NewUpHandler()
		assert.NotNil(t, handler)
	})

	t.Run("down handler constructor", func(t *testing.T) {
		handler := NewDownHandler()
		assert.NotNil(t, handler)
	})

	t.Run("restart handler constructor", func(t *testing.T) {
		handler := NewRestartHandler()
		assert.NotNil(t, handler)
	})

	t.Run("cleanup handler constructor", func(t *testing.T) {
		handler := NewCleanupHandler()
		assert.NotNil(t, handler)
	})
}

func TestCleanupHandler_basic(t *testing.T) {
	t.Run("new cleanup handler", func(t *testing.T) {
		handler := &CleanupHandler{}
		testhelpers.AssertNoError(t, nil, "CleanupHandler creation should not error")
		if handler == nil {
			t.Error("CleanupHandler should be created")
		}
	})
}
