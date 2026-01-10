package lifecycle

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/command"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/operations"
)

// DownHandler handles the down command
type DownHandler struct {
	stateManager *common.StateManager
}

// NewDownHandler creates a new down handler
func NewDownHandler() *DownHandler {
	return &DownHandler{
		stateManager: common.NewStateManager(),
	}
}

// Handle executes the down command
// TODO: Refactor - this method has significant duplication with UpHandler.Handle()
// Consider extracting common handler execution pattern to shared utility
func (h *DownHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	// Build CLI context from flags and args
	cliCtx, err := common.BuildStackContext(cmd, args)
	if err != nil {
		return err
	}

	// Create command and middleware chain
	downCommand := operations.NewServiceCommand(core.CommandDown, h.stateManager)
	validationMiddleware, loggingMiddleware := common.CreateStandardMiddlewareChain()

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
