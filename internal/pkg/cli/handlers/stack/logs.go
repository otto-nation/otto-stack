package stack

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/command"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/middleware"
)

// LogsHandler handles the logs command
type LogsHandler struct {
	stateManager *StateManager
}

// NewLogsHandler creates a new logs handler
func NewLogsHandler() *LogsHandler {
	return &LogsHandler{
		stateManager: NewStateManager(),
	}
}

// Handle executes the logs command
func (h *LogsHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	// Create command and middleware chain
	logsCommand := NewLogsCommand(h.stateManager)
	validationMiddleware := middleware.NewInitializationMiddleware()
	loggingMiddleware := middleware.NewLoggingMiddleware()

	handler := command.NewHandler(logsCommand, loggingMiddleware, validationMiddleware)

	// For now, create empty context - will be enhanced with flag processing
	cliCtx := clicontext.Context{}

	// Execute through command pattern
	return handler.Execute(ctx, cliCtx, base)
}

// ResolveServiceNames resolves service names from args or config
func (h *LogsHandler) ResolveServiceNames(args []string, setup *CoreSetup) ([]string, error) {
	// Legacy method - will be moved to LogsCommand.Execute
	return args, nil
}

// ValidateArgs validates the command arguments
func (h *LogsHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *LogsHandler) GetRequiredFlags() []string {
	return []string{}
}
