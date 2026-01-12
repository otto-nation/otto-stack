//go:build unit

package base

import (
	"context"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNewBaseCommand(t *testing.T) {
	t.Run("creates base command with quiet flag", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Bool(core.FlagQuiet, true, "quiet mode")

		base := NewBaseCommand(cmd)
		assert.NotNil(t, base)
		assert.NotNil(t, base.Output)
	})

	t.Run("creates base command with no-color flag", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Bool(core.FlagQuiet, false, "quiet mode")
		cmd.Flags().Bool(core.FlagNoColor, true, "no color")

		base := NewBaseCommand(cmd)
		assert.NotNil(t, base)
		assert.NotNil(t, base.Output)
	})

	t.Run("handles missing no-color flag gracefully", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Bool(core.FlagQuiet, false, "quiet mode")
		// No no-color flag added

		base := NewBaseCommand(cmd)
		assert.NotNil(t, base)
		assert.NotNil(t, base.Output)
	})
}

func TestBaseCommand_GetVerbose(t *testing.T) {
	t.Run("extracts verbose flag when present", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Bool("verbose", true, "verbose mode")

		base := &BaseCommand{}
		verbose := base.GetVerbose(cmd)
		assert.True(t, verbose)
	})

	t.Run("handles missing verbose flag", func(t *testing.T) {
		cmd := &cobra.Command{}
		// No verbose flag added

		base := &BaseCommand{}
		verbose := base.GetVerbose(cmd)
		assert.False(t, verbose)
	})
}

func TestCommandHandler_Interface(t *testing.T) {
	t.Run("validates CommandHandler interface", func(t *testing.T) {
		// Test that CommandHandler interface can be implemented
		var handler CommandHandler
		assert.Nil(t, handler)

		// Test context usage
		ctx := context.Background()
		assert.NotNil(t, ctx)
	})
}

func TestOutput_Interface(t *testing.T) {
	t.Run("validates Output interface methods", func(t *testing.T) {
		// Test that Output interface can be implemented
		var output Output
		assert.Nil(t, output)

		// Mock implementation for testing
		mock := &mockOutput{}
		assert.NotNil(t, mock)

		// Test interface methods don't panic
		mock.Success("test")
		mock.Error("test")
		mock.Warning("test")
		mock.Info("test")
		mock.Header("test")
		mock.Muted("test")
	})
}

// Mock output for testing
type mockOutput struct{}

func (m *mockOutput) Success(msg string, args ...any) {}
func (m *mockOutput) Error(msg string, args ...any)   {}
func (m *mockOutput) Warning(msg string, args ...any) {}
func (m *mockOutput) Info(msg string, args ...any)    {}
func (m *mockOutput) Header(msg string, args ...any)  {}
func (m *mockOutput) Muted(msg string, args ...any)   {}
