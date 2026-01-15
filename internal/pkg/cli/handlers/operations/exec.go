package operations

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/command"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
)

// ExecHandler handles the exec command
type ExecHandler struct {
	stateManager *common.StateManager
}

// NewExecHandler creates a new exec handler
func NewExecHandler() *ExecHandler {
	return &ExecHandler{
		stateManager: common.NewStateManager(),
	}
}

// Handle executes the exec command
func (h *ExecHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	// Get verbose flag
	verbose := base.GetVerbose(cmd)

	// Create command and middleware chain
	execCommand := NewServiceCommand(core.CommandExec, h.stateManager)
	execCommand.SetVerbose(verbose)
	validationMiddleware, loggingMiddleware := CreateStandardMiddlewareChain()

	handler := command.NewHandler(execCommand, loggingMiddleware, validationMiddleware)

	// For now, create empty context - will be enhanced with flag processing
	cliCtx := clicontext.Context{}

	// Execute through command pattern
	return handler.Execute(ctx, cliCtx, base)
}

// ValidateArgs validates the command arguments
func (h *ExecHandler) ValidateArgs(args []string) error {
	if len(args) < core.MinArgumentCount {
		return fmt.Errorf("%s", core.MsgErrors_requires_service_and_command)
	}
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *ExecHandler) GetRequiredFlags() []string {
	return []string{}
}
