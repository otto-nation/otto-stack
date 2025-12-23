package stack

import (
	"context"
	"fmt"
	"time"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
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

	// Parse flags directly from command
	timeout, _ := cmd.Flags().GetInt(core.FlagTimeout)
	volumes, _ := cmd.Flags().GetBool(core.FlagVolumes)
	_, _ = cmd.Flags().GetBool(core.FlagRemoveOrphans)  // Keep for future use
	_, _ = cmd.Flags().GetString(core.FlagRemoveImages) // Keep for future use

	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	// Use new stack service
	stackService, err := NewStackService(false) // Not verbose for down operations
	if err != nil {
		logger.Error(logger.LogMsgOperationFailed, logger.LogFieldOperation, logger.OperationStackDown, logger.LogFieldError, err)
		return err
	}

	// Create stop request
	stopRequest := services.StopRequest{
		Project:       setup.Config.Project.Name,
		Remove:        true, // Use down operation (remove containers)
		RemoveVolumes: volumes,
		Timeout:       time.Duration(timeout) * time.Second,
	}

	if err := stackService.Stop(ctx, stopRequest); err != nil {
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
