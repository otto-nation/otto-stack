package lifecycle

import (
	"context"
	"fmt"
	"time"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
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
	base.Output.Header(core.MsgLifecycle_restarting)

	ciFlags := ci.GetFlags(cmd)
	if ciFlags.DryRun {
		base.Output.Info("%s", core.MsgDry_run_showing_what_would_happen)
		base.Output.Info(core.MsgDry_run_would_restart_services, fmt.Sprintf("%v", args))
		return nil
	}

	setup, cleanup, err := common.SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	flags, err := core.ParseRestartFlags(cmd)
	if err != nil {
		return err
	}

	// Resolve services - inline logic since lifecycle uses different pattern than operations
	var serviceConfigs []types.ServiceConfig
	if len(args) > 0 {
		serviceConfigs, err = services.ResolveUpServices(args, setup.Config)
	} else {
		serviceConfigs, err = services.ResolveUpServices(setup.Config.Stack.Enabled, setup.Config)
	}
	if err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, "resolve services", err)
	}

	if err := h.restartServices(ctx, setup, serviceConfigs, flags); err != nil {
		return err
	}

	base.Output.Success(core.MsgLifecycle_restart_success)
	return nil
}

// restartServices performs the stop and start operations using new stack service
func (h *RestartHandler) restartServices(ctx context.Context, setup *common.CoreSetup, serviceConfigs []types.ServiceConfig, flags *core.RestartFlags) error {
	// Create stack service
	stackService, err := common.NewServiceManager(false)
	if err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, "create service", err)
	}

	// Stop services
	stopRequest := services.StopRequest{
		Project:        setup.Config.Project.Name,
		ServiceConfigs: serviceConfigs,
		Remove:         false, // Just stop, don't remove
		Timeout:        time.Duration(flags.Timeout) * time.Second,
	}
	if err := stackService.Stop(ctx, stopRequest); err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, "stop services", err)
	}

	// Start services
	startRequest := services.StartRequest{
		Project:        setup.Config.Project.Name,
		ServiceConfigs: serviceConfigs,
	}
	if err := stackService.Start(ctx, startRequest); err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, "start services", err)
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
