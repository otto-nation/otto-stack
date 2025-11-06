package stack

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
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
	// Check initialization first
	if err := utils.CheckInitialization(); err != nil {
		return err
	}

	// Get CI-friendly flags
	ciFlags := utils.GetCIFlags(cmd)

	if !ciFlags.Quiet {
		base.Output.Header(core.MsgLogs)
	}

	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	// Parse all flags with validation - single line!
	flags, err := core.ParseLogsFlags(cmd)
	if err != nil {
		return err
	}

	// Determine services
	serviceNames := args
	if len(serviceNames) == 0 {
		serviceNames = setup.Config.Stack.Enabled
	}

	// Apply service resolution
	serviceUtils := services.NewServiceUtils()
	resolvedServices, err := serviceUtils.ResolveServices(serviceNames)
	if err != nil {
		utils.HandleError(ciFlags, fmt.Errorf("failed to resolve services: %w", err))
		return nil
	}

	// Create log options - clean usage with no repetitive error handling
	options := docker.LogOptions{
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
