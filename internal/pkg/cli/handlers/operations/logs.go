package operations

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/command"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
)

// LogsHandler handles the logs command
type LogsHandler struct {
	stateManager *common.StateManager
}

// NewLogsHandler creates a new logs handler
func NewLogsHandler() *LogsHandler {
	return &LogsHandler{
		stateManager: common.NewStateManager(),
	}
}

// Handle executes the logs command
func (h *LogsHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	// Build CLI context from flags and args
	cliCtx, err := common.BuildStackContext(cmd, args)
	if err != nil {
		return err
	}

	// Get verbose flag
	verbose := base.GetVerbose(cmd)

	// Create command and middleware chain
	logsCommand := NewServiceCommand(core.CommandLogs, h.stateManager)
	logsCommand.SetVerbose(verbose)
	validationMiddleware, loggingMiddleware := CreateStandardMiddlewareChain()

	handler := command.NewHandler(logsCommand, loggingMiddleware, validationMiddleware)

	// Execute through command pattern
	return handler.Execute(ctx, cliCtx, base)
}

// ValidateArgs validates the command arguments
func (h *LogsHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *LogsHandler) GetRequiredFlags() []string {
	return []string{}
}
