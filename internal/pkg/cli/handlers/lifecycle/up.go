package lifecycle

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/command"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/shared"
	"github.com/otto-nation/otto-stack/internal/pkg/validation"
)

// UpHandler handles the up command
type UpHandler struct {
	stateManager *StateManager
}

// NewUpHandler creates a new up handler
func NewUpHandler() *UpHandler {
	return &UpHandler{
		stateManager: NewStateManager(),
	}
}

// Handle executes the up command
func (h *UpHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	// Process flags and build CLI context
	cliCtx, err := h.buildContext(cmd, args)
	if err != nil {
		return err
	}

	// Create command and middleware chain
	upCommand := NewServiceCommand(core.CommandUp, h.stateManager)
	validationMiddleware, loggingMiddleware := shared.CreateStandardMiddlewareChain()

	handler := command.NewHandler(upCommand, loggingMiddleware, validationMiddleware)

	// Execute through command pattern
	return handler.Execute(ctx, cliCtx, base)
}

// buildContext processes flags and arguments to build CLI context
func (h *UpHandler) buildContext(cmd *cobra.Command, args []string) (clicontext.Context, error) {
	return BuildStackContext(cmd, args)
}

// Legacy methods - will be moved to UpCommand.Execute gradually

// ValidateArgs validates the command arguments
func (h *UpHandler) ValidateArgs(args []string) error {
	return validation.ValidateUpArgs(args)
}

// GetRequiredFlags returns required flags for this command
func (h *UpHandler) GetRequiredFlags() []string {
	// No flags are strictly required for the up command
	return []string{}
}
