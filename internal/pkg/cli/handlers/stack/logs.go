package stack

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/spf13/cobra"
)

// LogsHandler handles the logs command
type LogsHandler struct{}

// NewLogsHandler creates a new logs handler
func NewLogsHandler() *LogsHandler {
	return &LogsHandler{}
}

// Handle executes the logs command
func (h *LogsHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	// Get CI-friendly flags
	ciFlags := utils.GetCIFlags(cmd)

	if !ciFlags.Quiet {
		base.Output.Header(constants.MsgLogs)
	}

	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		utils.HandleError(ciFlags, err)
		return nil
	}
	defer cleanup()

	// Parse all flags with validation - single line!
	flags, err := constants.ParseLogsFlags(cmd)
	if err != nil {
		return err
	}

	// Determine services
	serviceNames := args
	if len(serviceNames) == 0 {
		serviceNames = setup.Config.Stack.Enabled
	}

	// Apply service resolution
	serviceUtils := utils.NewServiceUtils()
	resolvedServices, err := serviceUtils.ResolveServices(serviceNames)
	if err != nil {
		utils.HandleError(ciFlags, fmt.Errorf("failed to resolve services: %w", err))
		return nil
	}

	// Create log options - clean usage with no repetitive error handling
	options := types.LogOptions{
		Follow:     flags.Follow,
		Timestamps: flags.Timestamps,
		Tail:       flags.Tail,
	}

	// Get logs using Docker client
	return setup.DockerClient.ComposeLogs(ctx, setup.Config.Project.Name, resolvedServices, options)
}

// ValidateArgs validates the command arguments
func (h *LogsHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *LogsHandler) GetRequiredFlags() []string {
	return []string{}
}
