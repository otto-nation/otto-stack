package stack

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/spf13/cobra"
)

// DownHandler handles the down command
type DownHandler struct{}

// NewDownHandler creates a new down handler
func NewDownHandler() *DownHandler {
	return &DownHandler{}
}

// Handle executes the down command
func (h *DownHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	// Check initialization first

	// Start operation logging only after initialization check passes
	logger.Info(logger.LogMsgStartingOperation, logger.LogFieldOperation, logger.OperationStackDown, logger.LogFieldServices, args)
	defer func() {
		if r := recover(); r != nil {
			logger.Error(logger.LogMsgOperationFailed, logger.LogFieldOperation, logger.OperationStackDown, logger.LogFieldError, fmt.Errorf("panic: %v", r))
			panic(r)
		}
	}()

	ciFlags := ci.GetFlags(cmd)

	if ciFlags.DryRun {
		base.Output.Info("%s", core.MsgDry_run_showing_what_would_happen)
		base.Output.Info(core.MsgDry_run_would_stop_services, fmt.Sprintf("%v", args))
		return nil
	}

	base.Output.Header(core.MsgStopping)
	logger.Info(logger.LogMsgServiceAction, logger.LogFieldAction, logger.ActionStop, logger.LogFieldService, "stack", logger.LogFieldServices, args)

	// Parse all flags with validation - single line!
	flags, err := core.ParseDownFlags(cmd)
	if err != nil {
		logger.Error(logger.LogMsgOperationFailed, logger.LogFieldOperation, logger.OperationStackDown, logger.LogFieldError, err)
		return err
	}

	// Load project configuration
	configPath := filepath.Join(core.OttoStackDir, core.ConfigFileName)
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

	// Determine services to stop
	// Convert CLI options to internal options
	internalOptions := docker.StopOptions{
		Timeout:       flags.Timeout,
		Remove:        flags.Remove,
		RemoveVolumes: flags.Volumes,
		RemoveOrphans: flags.RemoveOrphans,
	}

	// Stop services
	if err := dockerClient.ComposeDown(ctx, cfg.Project.Name, internalOptions); err != nil {
		logger.Error(logger.LogMsgOperationFailed, logger.LogFieldOperation, logger.OperationStackDown, logger.LogFieldError, err)
		return fmt.Errorf(core.MsgStack_failed_stop_services, err)
	}

	base.Output.Success(core.MsgStopSuccess)
	logger.Info(logger.LogMsgOperationCompleted, logger.LogFieldOperation, logger.OperationStackDown)
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
