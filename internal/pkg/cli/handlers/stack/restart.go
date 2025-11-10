package stack

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
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

	base.Output.Header(core.MsgRestarting)

	ciFlags := ci.GetFlags(cmd)

	if ciFlags.DryRun {
		base.Output.Info("%s", core.MsgDry_run_showing_what_would_happen)
		base.Output.Info(core.MsgDry_run_would_restart_services, fmt.Sprintf("%v", args))
		return nil
	}

	// Check if otto-stack is initialized
	configPath := filepath.Join(core.OttoStackDir, core.ConfigFileName)
	if !func() bool { _, err := os.Stat(configPath); return err == nil }() {
		return errors.New(core.MsgErrors_not_initialized)
	}

	// Load project configuration
	cfg, err := LoadProjectConfig(configPath)
	if err != nil {
		return fmt.Errorf(core.MsgStack_failed_load_config, err)
	}

	// Create Docker client
	dockerClient, err := docker.NewClient(nil)
	if err != nil {
		return fmt.Errorf(core.MsgStack_failed_create_docker_client, err)
	}
	defer func() {
		_ = dockerClient.Close()
	}()

	// Parse all flags with validation - single line!
	flags, err := core.ParseRestartFlags(cmd)
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
		return fmt.Errorf(core.MsgStack_failed_stop_services, err)
	}

	// Start services
	startOptions := docker.StartOptions{
		Detach: true,
	}
	if err := dockerClient.ComposeUp(ctx, cfg.Project.Name, serviceNames, startOptions); err != nil {
		return fmt.Errorf(core.MsgStack_failed_start_services, err)
	}

	base.Output.Success(core.MsgRestartSuccess)
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
