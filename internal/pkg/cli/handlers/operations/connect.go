package operations

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/command"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
)

// ConnectHandler handles the connect command
type ConnectHandler struct {
	stateManager *common.StateManager
}

// NewConnectHandler creates a new connect handler
func NewConnectHandler() *ConnectHandler {
	return &ConnectHandler{
		stateManager: common.NewStateManager(),
	}
}

// ValidateArgs validates the command arguments
func (h *ConnectHandler) ValidateArgs(args []string) error {
	if len(args) < 1 {
		return pkgerrors.NewValidationError(pkgerrors.FieldServiceName, "service name is required", nil)
	}
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *ConnectHandler) GetRequiredFlags() []string {
	return []string{}
}

// Handle executes the connect command
func (h *ConnectHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	// Get verbose flag
	verbose := base.GetVerbose(cmd)

	// Create command and middleware chain
	connectCommand := NewServiceCommand(core.CommandConnect, h.stateManager)
	connectCommand.SetVerbose(verbose)
	validationMiddleware, loggingMiddleware := CreateStandardMiddlewareChain()

	handler := command.NewHandler(connectCommand, loggingMiddleware, validationMiddleware)

	// For now, create empty context - will be enhanced with flag processing
	cliCtx := clicontext.Context{}

	// Execute through command pattern
	return handler.Execute(ctx, cliCtx, base)
}
