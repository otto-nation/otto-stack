package stack

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/spf13/cobra"
)

// LogsHandler handles the logs command
type LogsHandler struct{}

// NewLogsHandler creates a new logs handler
func NewLogsHandler() *LogsHandler {
	return &LogsHandler{}
}

// Handle executes the logs command
func (h *LogsHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	ciFlags := ci.GetFlags(cmd)

	if !ciFlags.Quiet {
		base.Output.Header(core.MsgLogs)
	}

	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	flags, err := core.ParseLogsFlags(cmd)
	if err != nil {
		return err
	}

	serviceNames := h.resolveServiceNames(args, setup.Config.Stack.Enabled)
	resolvedServices, err := h.resolveServices(serviceNames)
	if err != nil {
		ci.HandleError(ciFlags, fmt.Errorf("failed to resolve services: %w", err))
		return nil
	}

	// Create stack service
	stackService, err := NewStackService(false)
	if err != nil {
		return fmt.Errorf("failed to create stack service: %w", err)
	}

	// Create logs request
	logRequest := services.LogRequest{
		Project:    setup.Config.Project.Name,
		Services:   resolvedServices,
		Follow:     flags.Follow,
		Timestamps: flags.Timestamps,
		Tail:       flags.Tail,
	}

	return stackService.Logs(ctx, logRequest)
}

// resolveServiceNames determines which services to get logs for
func (h *LogsHandler) resolveServiceNames(args, enabledServices []string) []string {
	if len(args) > 0 {
		return args
	}
	return enabledServices
}

// resolveServices resolves service names using service utils
func (h *LogsHandler) resolveServices(serviceNames []string) ([]string, error) {
	serviceUtils := services.NewServiceUtils()
	return serviceUtils.ResolveServices(serviceNames)
}

// ValidateArgs validates the command arguments
func (h *LogsHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *LogsHandler) GetRequiredFlags() []string {
	return []string{}
}
