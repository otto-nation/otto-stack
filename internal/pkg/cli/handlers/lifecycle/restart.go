package lifecycle

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/registry"
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
	base.Output.Header(messages.LifecycleRestarting)

	ciFlags := ci.GetFlags(cmd)
	if ciFlags.DryRun {
		base.Output.Info("%s", messages.DryRunShowingWhatWouldHappen)
		base.Output.Info(messages.DryRunWouldRestartServices, fmt.Sprintf("%v", args))
		return nil
	}

	execCtx, err := h.detectContext()
	if err != nil {
		return err
	}

	switch mode := execCtx.(type) {
	case *clicontext.ProjectMode:
		return h.handleProjectContext(ctx, cmd, args, base)
	case *clicontext.SharedMode:
		return h.handleSharedContext(ctx, cmd, args, base, mode)
	default:
		return fmt.Errorf("unknown execution mode: %T", execCtx)
	}
}

func (h *RestartHandler) detectContext() (clicontext.ExecutionMode, error) {
	detector, err := clicontext.NewDetector()
	if err != nil {
		return nil, err
	}
	return detector.DetectContext()
}

func (h *RestartHandler) handleProjectContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	setup, cleanup, err := common.SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	flags, err := core.ParseRestartFlags(cmd)
	if err != nil {
		return err
	}

	var serviceConfigs []types.ServiceConfig
	if len(args) > 0 {
		serviceConfigs, err = services.ResolveUpServices(args, setup.Config)
	} else {
		serviceConfigs, err = services.ResolveUpServices(setup.Config.Stack.Enabled, setup.Config)
	}
	if err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, messages.ErrorsStackResolveFailed, err)
	}

	if err := h.restartServices(ctx, setup, serviceConfigs, flags); err != nil {
		return err
	}

	base.Output.Success(messages.LifecycleRestartSuccess)
	return nil
}

func (h *RestartHandler) handleSharedContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand, mode *clicontext.SharedMode) error {
	if len(args) == 0 {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationServiceNameRequired, nil)
	}

	reg, err := h.loadRegistry(mode.Shared.Root)
	if err != nil {
		return err
	}

	if err := h.verifyServicesInRegistry(args, reg); err != nil {
		return err
	}

	flags, err := core.ParseRestartFlags(cmd)
	if err != nil {
		return err
	}

	dockerClient, err := docker.NewClient(logger.GetLogger())
	if err != nil {
		return err
	}

	timeout := flags.Timeout
	stopOpts := container.StopOptions{Timeout: &timeout}

	for _, serviceName := range args {
		if err := dockerClient.GetDockerClient().ContainerRestart(ctx, serviceName, stopOpts); err != nil {
			return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentDocker, fmt.Sprintf("failed to restart %s", serviceName), err)
		}
		base.Output.Success("Restarted %s", serviceName)
	}

	return nil
}

func (h *RestartHandler) loadRegistry(sharedRoot string) (*registry.Registry, error) {
	reg := registry.NewManager(sharedRoot)
	return reg.Load()
}

func (h *RestartHandler) verifyServicesInRegistry(serviceNames []string, reg *registry.Registry) error {
	for _, serviceName := range serviceNames {
		if _, exists := reg.Containers[serviceName]; !exists {
			return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, fmt.Sprintf(messages.SharedServiceNotInRegistry, serviceName), nil)
		}
	}
	return nil
}

// restartServices performs the stop and start operations using new stack service
func (h *RestartHandler) restartServices(ctx context.Context, setup *common.CoreSetup, serviceConfigs []types.ServiceConfig, flags *core.RestartFlags) error {
	// Create stack service
	stackService, err := common.NewServiceManager(false)
	if err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, messages.ErrorsStackCreateFailed, err)
	}

	// Stop services
	stopRequest := services.StopRequest{
		Project:        setup.Config.Project.Name,
		ServiceConfigs: serviceConfigs,
		Remove:         false, // Just stop, don't remove
		Timeout:        time.Duration(flags.Timeout) * time.Second,
	}
	if err := stackService.Stop(ctx, stopRequest); err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, messages.ErrorsStackStopFailed, err)
	}

	// Start services
	startRequest := services.StartRequest{
		Project:        setup.Config.Project.Name,
		ServiceConfigs: serviceConfigs,
	}
	if err := stackService.Start(ctx, startRequest); err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, messages.ErrorsStackStartFailed, err)
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
