package stack

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/utils"
	"github.com/spf13/cobra"
)

// RestartHandler handles the restart command
type RestartHandler struct{}

// NewRestartHandler creates a new restart handler
func NewRestartHandler() *RestartHandler {
	return &RestartHandler{}
}

// Handle executes the restart command
func (h *RestartHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	base.Output.Header(constants.MsgRestarting)

	// Check if otto-stack is initialized
	configPath := filepath.Join(constants.OttoStackDir, constants.ConfigFileName)
	if !utils.FileExists(configPath) {
		return errors.New(constants.Messages[constants.MsgErrors_not_initialized])
	}

	// Load project configuration
	cfg, err := LoadProjectConfig(configPath)
	if err != nil {
		return fmt.Errorf(constants.Messages[constants.MsgStack_failed_load_config], err)
	}

	// Create Docker client
	logger := base.Logger.(loggerAdapter)
	dockerClient, err := docker.NewClient(logger.SlogLogger())
	if err != nil {
		return fmt.Errorf(constants.Messages[constants.MsgStack_failed_create_docker_client], err)
	}
	defer func() {
		if err := dockerClient.Close(); err != nil {
			base.Logger.Error("Failed to close Docker client", "error", err)
		}
	}()

	// Parse all flags with validation - single line!
	flags, err := constants.ParseRestartFlags(cmd)
	if err != nil {
		return err
	}

	// Determine services to restart
	serviceNames := args
	if len(serviceNames) == 0 {
		serviceNames = cfg.Stack.Enabled
	}

	// Stop services first
	// Restart operation
	// Stop services
	stopOptions := types.StopOptions{
		Timeout: flags.Timeout,
	}
	if err := dockerClient.Containers().Stop(ctx, cfg.Project.Name, serviceNames, stopOptions); err != nil {
		return fmt.Errorf(constants.Messages[constants.MsgStack_failed_stop_services], err)
	}

	// Start services
	startOptions := types.StartOptions{
		Detach: true,
	}
	if err := dockerClient.Containers().Start(ctx, cfg.Project.Name, serviceNames, startOptions); err != nil {
		return fmt.Errorf(constants.Messages[constants.MsgStack_failed_start_services], err)
	}

	base.Output.Success(constants.MsgRestartSuccess)
	// Restart operation
	return nil
}

// ValidateArgs validates the command arguments
func (h *RestartHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *RestartHandler) GetRequiredFlags() []string {
	return []string{}
}
