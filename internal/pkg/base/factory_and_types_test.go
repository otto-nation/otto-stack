//go:build unit

package base

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNewBaseCommand_WithQuiet(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool(core.FlagQuiet, true, "quiet mode")

	base := NewBaseCommand(cmd)
	assert.NotNil(t, base)
	assert.NotNil(t, base.Output)
}

func TestNewBaseCommand_WithNoColor(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool(core.FlagQuiet, false, "quiet mode")
	cmd.Flags().Bool(core.FlagNoColor, true, "no color")

	base := NewBaseCommand(cmd)
	assert.NotNil(t, base)
	assert.NotNil(t, base.Output)
}

func TestNewBaseCommand_MissingNoColor(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool(core.FlagQuiet, false, "quiet mode")

	base := NewBaseCommand(cmd)
	assert.NotNil(t, base)
	assert.NotNil(t, base.Output)
}

func TestBaseCommand_GetVerbose_Present(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("verbose", true, "verbose mode")

	base := &BaseCommand{}
	verbose := base.GetVerbose(cmd)
	assert.True(t, verbose)
}

func TestBaseCommand_GetVerbose_Missing(t *testing.T) {
	cmd := &cobra.Command{}

	base := &BaseCommand{}
	verbose := base.GetVerbose(cmd)
	assert.False(t, verbose)
}

func TestCommandHandler_Interface(t *testing.T) {
	var handler CommandHandler
	assert.Nil(t, handler)

	ctx := context.Background()
	assert.NotNil(t, ctx)
}

func TestOutput_Interface(t *testing.T) {
	var output Output
	assert.Nil(t, output)

	mock := &mockOutput{}
	assert.NotNil(t, mock)

	mock.Success("test")
	mock.Error("test")
	mock.Warning("test")
	mock.Info("test")
	mock.Header("test")
	mock.Muted("test")
}

// Mock output for testing
type mockOutput struct{}

func (m *mockOutput) Success(msg string, args ...any) {}
func (m *mockOutput) Error(msg string, args ...any)   {}
func (m *mockOutput) Warning(msg string, args ...any) {}
func (m *mockOutput) Info(msg string, args ...any)    {}
func (m *mockOutput) Header(msg string, args ...any)  {}
func (m *mockOutput) Muted(msg string, args ...any)   {}
func (m *mockOutput) Writer() io.Writer               { return os.Stdout }
