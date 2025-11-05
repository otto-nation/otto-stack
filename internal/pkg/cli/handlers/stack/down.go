package stack

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/spf13/cobra"
)

// DownHandler handles the down command
type DownHandler struct{}

// NewDownHandler creates a new down handler
func NewDownHandler() *DownHandler {
	return &DownHandler{}
}

// Handle executes the down command
func (h *DownHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	// Check initialization first
	if err := utils.CheckInitialization(); err != nil {
		return err
	}

	// Start operation logging only after initialization check passes
	logger.Info(constants.LogMsgStartingOperation, constants.LogFieldOperation, constants.OperationStackDown, constants.LogFieldServices, args)
	defer func() {
		if r := recover(); r != nil {
			logger.Error(constants.LogMsgOperationFailed, constants.LogFieldOperation, constants.OperationStackDown, constants.LogFieldError, fmt.Errorf("panic: %v", r))
			panic(r)
		}
	}()

	base.Output.Header(constants.MsgStopping)
	logger.Info(constants.LogMsgServiceAction, constants.LogFieldAction, constants.ActionStop, constants.LogFieldService, "stack", constants.LogFieldServices, args)

	// Parse all flags with validation - single line!
	flags, err := constants.ParseDownFlags(cmd)
	if err != nil {
		logger.Error(constants.LogMsgOperationFailed, constants.LogFieldOperation, constants.OperationStackDown, constants.LogFieldError, err)
		return err
	}

	// Load project configuration
	configPath := filepath.Join(constants.OttoStackDir, constants.ConfigFileName)
	cfg, err := LoadProjectConfig(configPath)
	if err != nil {
		return fmt.Errorf(constants.MsgStack_failed_load_config, err)
	}

	// Create Docker client
	dockerLogger := base.Logger
	dockerClient, err := docker.NewClient(dockerLogger.SlogLogger())
	if err != nil {
		return fmt.Errorf(constants.MsgStack_failed_create_docker_client, err)
	}
	defer func() {
		if err := dockerClient.Close(); err != nil {
			base.Logger.Error("Failed to close Docker client", "error", err)
		}
	}()

	// Determine services to stop
	// Convert CLI options to internal options
	internalOptions := types.StopOptions{
		Timeout:       flags.Timeout,
		Remove:        true,
		RemoveVolumes: flags.Volumes,
	}

	// Stop services
	if err := dockerClient.ComposeDown(ctx, cfg.Project.Name, internalOptions); err != nil {
		logger.Error(constants.LogMsgOperationFailed, constants.LogFieldOperation, constants.OperationStackDown, constants.LogFieldError, err)
		return fmt.Errorf(constants.MsgStack_failed_stop_services, err)
	}

	base.Output.Success(constants.MsgStopSuccess)
	logger.Info(constants.LogMsgOperationCompleted, constants.LogFieldOperation, constants.OperationStackDown)
	return nil
}

// ValidateArgs validates the command arguments
func (h *DownHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *DownHandler) GetRequiredFlags() []string {
	return []string{}
}
