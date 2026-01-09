package stack

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/command"
)

// DownHandler handles the down command
type DownHandler struct {
	stateManager *StateManager
}

// NewDownHandler creates a new down handler
func NewDownHandler() *DownHandler {
	return &DownHandler{
		stateManager: NewStateManager(),
	}
}

// Handle executes the down command
func (h *DownHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	// Build CLI context from flags and args
	cliCtx, err := BuildStackContext(cmd, args)
	if err != nil {
		return err
	}

	// Create command and middleware chain
	downCommand := NewDownCommand(h.stateManager)
	validationMiddleware, loggingMiddleware := CreateStandardMiddlewareChain()

	handler := command.NewHandler(downCommand, loggingMiddleware, validationMiddleware)

	// Execute through command pattern
	return handler.Execute(ctx, cliCtx, base)
}

// ValidateArgs validates the command arguments
func (h *DownHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *DownHandler) GetRequiredFlags() []string {
	return []string{}
}
