package stack

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/command"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/middleware"
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
	// Create command and middleware chain
	downCommand := NewDownCommand(h.stateManager)
	validationMiddleware := middleware.NewInitializationMiddleware()
	loggingMiddleware := middleware.NewLoggingMiddleware()

	handler := command.NewHandler(downCommand, loggingMiddleware, validationMiddleware)

	// For now, create empty context - will be enhanced with flag processing
	cliCtx := clicontext.Context{}

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
