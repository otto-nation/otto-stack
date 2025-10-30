package stack

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/types"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	pkgTypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
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
		ui.Header(constants.MsgLogs)
	}

	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		utils.HandleError(ciFlags, err)
		return nil
	}
	defer cleanup()

	// Get flags
	follow, _ := cmd.Flags().GetBool(constants.FlagFollow)
	tail, _ := cmd.Flags().GetString(constants.FlagTail)
	timestamps, _ := cmd.Flags().GetBool(constants.FlagTimestamps)

	// Determine services
	serviceNames := args
	if len(serviceNames) == 0 {
		serviceNames = setup.Config.Stack.Enabled
	}

	// Apply service resolution
	serviceUtils := utils.NewServiceUtils()
	resolvedServices, err := serviceUtils.ResolveServices(serviceNames)
	if err != nil {
		utils.HandleError(ciFlags, fmt.Errorf(constants.MsgFailedResolveServices.Content, err))
		return nil
	}

	// Create log options
	options := pkgTypes.LogOptions{
		Follow:     follow,
		Tail:       tail,
		Timestamps: timestamps,
	}

	// Get logs using Docker client
	return setup.DockerClient.Containers().Logs(ctx, setup.Config.Project.Name, resolvedServices, options)
}

// ValidateArgs validates the command arguments
func (h *LogsHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *LogsHandler) GetRequiredFlags() []string {
	return []string{}
}
