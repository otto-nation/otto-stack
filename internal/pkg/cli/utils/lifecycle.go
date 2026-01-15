package utils

import (
	"context"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/command"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/operations"
	"github.com/spf13/cobra"
)

// ExecuteLifecycleCommand provides common execution pattern for lifecycle commands
// Eliminates duplication between UpHandler and DownHandler
func ExecuteLifecycleCommand(ctx context.Context, cmd *cobra.Command, args []string,
	base *base.BaseCommand, commandType string, stateManager *common.StateManager) error {

	// Build CLI context from flags and args
	cliCtx, err := common.BuildStackContext(cmd, args)
	if err != nil {
		return err
	}

	// Get verbose flag from command
	verbose := base.GetVerbose(cmd)

	// Create command and middleware chain
	serviceCommand := operations.NewServiceCommand(commandType, stateManager)
	serviceCommand.SetVerbose(verbose)
	validationMiddleware, loggingMiddleware := common.CreateStandardMiddlewareChain()
	handler := command.NewHandler(serviceCommand, loggingMiddleware, validationMiddleware)

	// Execute through command pattern
	return handler.Execute(ctx, cliCtx, base)
}
