//go:build integration

package project

import (
	"context"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitHandler_Integration(t *testing.T) {
	t.Run("handler creation and basic validation", func(t *testing.T) {
		// Test that we can create the handler without issues
		handler := NewInitHandler()
		require.NotNil(t, handler)

		// Test basic validation methods
		err := handler.ValidateArgs([]string{})
		assert.NoError(t, err)

		err = handler.ValidateArgs([]string{"extra", "args"})
		assert.NoError(t, err) // ValidateArgs currently always returns nil

		// Test required flags
		flags := handler.GetRequiredFlags()
		assert.NotNil(t, flags)
	})

	t.Run("handler with mock command", func(t *testing.T) {
		// Create handler
		handler := NewInitHandler()
		require.NotNil(t, handler)

		// Create a mock command with required flags
		cmd := &cobra.Command{
			Use: "init",
		}

		// Add the flags that the handler expects
		cmd.Flags().String("name", "", "Project name")
		cmd.Flags().Bool("force", false, "Force initialization")
		cmd.Flags().StringSlice("services", []string{}, "Services to include")

		// Set flag values to avoid prompting
		err := cmd.Flags().Set("name", "test-project")
		require.NoError(t, err)
		err = cmd.Flags().Set("force", "true")
		require.NoError(t, err)

		// Create base command with output that doesn't require user interaction
		output := ui.NewOutput()
		base := &base.BaseCommand{
			Output: output,
		}

		ctx := context.Background()

		// This should not panic, but may fail due to directory validation
		// We're mainly testing that the handler can be called without crashing
		err = handler.Handle(ctx, cmd, []string{}, base)

		// We expect this to fail in a controlled way (not panic)
		// The exact error depends on the validation logic
		if err != nil {
			t.Logf("Handler returned expected error: %v", err)
		}
	})
}
