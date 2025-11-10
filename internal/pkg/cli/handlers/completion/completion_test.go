package completion

import (
	"context"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNewCompletionHandler(t *testing.T) {
	t.Run("creates handler successfully", func(t *testing.T) {
		handler := NewCompletionHandler()
		assert.NotNil(t, handler)
		assert.IsType(t, &CompletionHandler{}, handler)
	})
}

func TestCompletionHandler_ValidateArgs(t *testing.T) {
	handler := NewCompletionHandler()

	t.Run("accepts valid shell types", func(t *testing.T) {
		validShells := []string{"bash", "zsh", "fish", "powershell"}

		for _, shell := range validShells {
			err := handler.ValidateArgs([]string{shell})
			assert.NoError(t, err, "Should accept shell: %s", shell)
		}
	})

	t.Run("rejects invalid shell types", func(t *testing.T) {
		err := handler.ValidateArgs([]string{"invalid-shell"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported shell")
	})

	t.Run("rejects no arguments", func(t *testing.T) {
		err := handler.ValidateArgs([]string{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exactly one argument")
	})

	t.Run("rejects multiple arguments", func(t *testing.T) {
		err := handler.ValidateArgs([]string{"bash", "zsh"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exactly one argument")
	})
}

func TestCompletionHandler_GetRequiredFlags(t *testing.T) {
	t.Run("returns empty flags list", func(t *testing.T) {
		handler := NewCompletionHandler()
		flags := handler.GetRequiredFlags()

		assert.NotNil(t, flags)
		assert.Empty(t, flags)
		assert.IsType(t, []string{}, flags)
	})
}

func TestCompletionHandler_Handle(t *testing.T) {
	handler := NewCompletionHandler()
	ctx := context.Background()

	// Create a minimal root command for testing
	rootCmd := &cobra.Command{
		Use: "test",
	}

	// Create a subcommand to simulate the completion command
	completionCmd := &cobra.Command{
		Use: "completion",
	}
	rootCmd.AddCommand(completionCmd)

	baseCmd := &base.BaseCommand{}

	t.Run("handles bash completion", func(t *testing.T) {
		assert.NotPanics(t, func() {
			handler.Handle(ctx, completionCmd, []string{"bash"}, baseCmd)
		})
	})

	t.Run("handles zsh completion", func(t *testing.T) {
		assert.NotPanics(t, func() {
			handler.Handle(ctx, completionCmd, []string{"zsh"}, baseCmd)
		})
	})

	t.Run("handles fish completion", func(t *testing.T) {
		assert.NotPanics(t, func() {
			handler.Handle(ctx, completionCmd, []string{"fish"}, baseCmd)
		})
	})

	t.Run("handles powershell completion", func(t *testing.T) {
		assert.NotPanics(t, func() {
			handler.Handle(ctx, completionCmd, []string{"powershell"}, baseCmd)
		})
	})

	t.Run("handles invalid shell gracefully", func(t *testing.T) {
		// This should be caught by ValidateArgs, but test Handle directly
		err := handler.Handle(ctx, completionCmd, []string{"invalid"}, baseCmd)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported shell")
	})
}
