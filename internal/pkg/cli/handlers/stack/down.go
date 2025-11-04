package stack

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/utils"
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
	finishOp := logger.StartOperation("stack_down", "services", args)
	defer func() {
		if r := recover(); r != nil {
			finishOp(fmt.Errorf("panic: %v", r))
			panic(r)
		}
	}()

	base.Output.Header(constants.MsgStopping)
	logger.LogServiceAction("stop", "stack", "services", args)

	// Parse all flags with validation - single line!
	flags, err := constants.ParseDownFlags(cmd)
	if err != nil {
		finishOp(err)
		return err
	}

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

	// Determine services to stop
	// Convert CLI options to internal options
	internalOptions := types.StopOptions{
		Timeout:       flags.Timeout,
		Remove:        true,
		RemoveVolumes: flags.Volumes,
	}

	// Stop services
	if err := dockerClient.ComposeDown(ctx, cfg.Project.Name, internalOptions); err != nil {
		finishOp(err)
		return fmt.Errorf(constants.Messages[constants.MsgStack_failed_stop_services], err)
	}

	base.Output.Success(constants.MsgStopSuccess)
	finishOp(nil)
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
