package stack

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/spf13/cobra"
)

// RestartHandler handles the restart command
type RestartHandler struct{}

// NewRestartHandler creates a new restart handler
func NewRestartHandler() *RestartHandler {
	return &RestartHandler{}
}

// Handle executes the restart command
func (h *RestartHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	// Check initialization first
	if err := utils.CheckInitialization(); err != nil {
		return err
	}

	base.Output.Header(constants.MsgRestarting)

	// Check if otto-stack is initialized
	configPath := filepath.Join(constants.OttoStackDir, constants.ConfigFileName)
	if !func() bool { _, err := os.Stat(configPath); return err == nil }() {
		return errors.New(constants.MsgErrors_not_initialized)
	}

	// Load project configuration
	cfg, err := LoadProjectConfig(configPath)
	if err != nil {
		return fmt.Errorf(constants.MsgStack_failed_load_config, err)
	}

	// Create Docker client
	dockerClient, err := docker.NewClient(nil)
	if err != nil {
		return fmt.Errorf(constants.MsgStack_failed_create_docker_client, err)
	}
	defer func() {
		_ = dockerClient.Close()
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
	stopOptions := docker.StopOptions{
		Timeout: flags.Timeout,
	}
	if err := dockerClient.ComposeDown(ctx, cfg.Project.Name, stopOptions); err != nil {
		return fmt.Errorf(constants.MsgStack_failed_stop_services, err)
	}

	// Start services
	startOptions := docker.StartOptions{
		Detach: true,
	}
	if err := dockerClient.ComposeUp(ctx, cfg.Project.Name, serviceNames, startOptions); err != nil {
		return fmt.Errorf(constants.MsgStack_failed_start_services, err)
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
