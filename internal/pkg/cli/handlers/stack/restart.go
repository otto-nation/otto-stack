package stack

import (
	"context"
	"fmt"
	"time"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
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
	base.Output.Header(core.MsgRestarting)

	ciFlags := ci.GetFlags(cmd)
	if ciFlags.DryRun {
		base.Output.Info("%s", core.MsgDry_run_showing_what_would_happen)
		base.Output.Info(core.MsgDry_run_would_restart_services, fmt.Sprintf("%v", args))
		return nil
	}

	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	flags, err := core.ParseRestartFlags(cmd)
	if err != nil {
		return err
	}

	serviceNames := h.resolveServiceNames(args, setup.Config.Stack.Enabled)

	if err := h.restartServices(ctx, setup, serviceNames, flags); err != nil {
		return err
	}

	base.Output.Success(core.MsgRestartSuccess)
	return nil
}

// resolveServiceNames determines which services to restart
func (h *RestartHandler) resolveServiceNames(args, enabledServices []string) []string {
	if len(args) > 0 {
		return args
	}
	return enabledServices
}

// restartServices performs the stop and start operations using new stack service
func (h *RestartHandler) restartServices(ctx context.Context, setup *CoreSetup, serviceNames []string, flags *core.RestartFlags) error {
	// Create stack service
	stackService, err := NewStackService(false)
	if err != nil {
		return fmt.Errorf("failed to create stack service: %w", err)
	}

	// Stop services
	stopRequest := services.StopRequest{
		Project: setup.Config.Project.Name,
		Remove:  false, // Just stop, don't remove
		Timeout: time.Duration(flags.Timeout) * time.Second,
	}
	if err := stackService.Stop(ctx, stopRequest); err != nil {
		return fmt.Errorf(core.MsgStack_failed_stop_services, err)
	}

	// Start services
	startRequest := services.StartRequest{
		Project:  setup.Config.Project.Name,
		Services: serviceNames,
	}
	if err := stackService.Start(ctx, startRequest); err != nil {
		return fmt.Errorf(core.MsgStack_failed_start_services, err)
	}

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
