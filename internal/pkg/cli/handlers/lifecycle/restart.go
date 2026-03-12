package lifecycle

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
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

	execCtx, err := common.DetectExecutionContext()
	if err != nil {
		return err
	}

	switch mode := execCtx.(type) {
	case *clicontext.ProjectMode:
		return h.handleProjectContext(ctx, cmd, args, base)
	case *clicontext.SharedMode:
		return h.handleSharedContext(ctx, cmd, args, base, mode)
	default:
		return pkgerrors.NewSystemErrorf(pkgerrors.ErrCodeInternal, messages.ErrorsContextUnknownMode, execCtx)
	}
}

func (h *RestartHandler) handleProjectContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	setup, cleanup, err := common.SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	flags, err := core.ParseRestartFlags(cmd)
	if err != nil {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldFlags, messages.ValidationFailedParseFlags, err)
	}

	var serviceConfigs []types.ServiceConfig
	if len(args) > 0 {
		serviceConfigs, err = services.ResolveUpServices(args, setup.Config)
	} else {
		serviceConfigs, err = services.ResolveUpServices(setup.Config.Stack.Enabled, setup.Config)
	}
	if err != nil {
		return err
	}

	if err := h.restartServices(ctx, setup, serviceConfigs, flags); err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentServices, messages.ErrorsServiceRestartFailed, err)
	}

	base.Output.Success(messages.LifecycleRestartSuccess)
	base.Output.Muted(messages.InfoProjectInfo, setup.Config.Project.Name)

	if svc, err := common.NewServiceManager(false); err == nil {
		filteredNames := filterStatusQueryNames(serviceConfigs)
		if statuses, err := svc.Status(ctx, services.StatusRequest{
			Project:  setup.Config.Project.Name,
			Services: filteredNames,
		}); err == nil {
			// Silent fallback: if Status() fails, the command already succeeded — skip the table
			_ = display.RenderStatusTable(base.Output.Writer(), statuses, serviceConfigs, true, base.Output.GetNoColor())
		}
	}

	return nil
}

func (h *RestartHandler) handleSharedContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand, mode *clicontext.SharedMode) error {
	if len(args) == 0 {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationServiceNameRequired, nil)
	}

	reg, err := h.loadRegistry(mode.Shared.Root)
	if err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsRegistryLoadFailed, err)
	}

	if err := common.VerifyServicesInRegistry(args, reg); err != nil {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ErrorsServiceNotInRegistry, err)
	}

	flags, err := core.ParseRestartFlags(cmd)
	if err != nil {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldFlags, messages.ValidationFailedParseFlags, err)
	}

	composePath := filepath.Join(mode.Shared.Root, core.GeneratedDir, docker.DockerComposeFileName)
	composeManager, err := docker.NewManager()
	if err != nil {
		return pkgerrors.NewDockerError(pkgerrors.ErrCodeOperationFail, messages.ErrorsDockerManagerCreateFailed, err)
	}

	proj, err := composeManager.LoadProject(ctx, []string{composePath}, "shared")
	if err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentDocker, messages.ErrorsDockerLoadProjectFailed, err)
	}

	timeout := time.Duration(flags.Timeout) * time.Second
	if err := composeManager.Stop(ctx, "shared", docker.StopOptions{Services: args, Timeout: &timeout}.ToSDK()); err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, messages.ErrorsServiceRestartFailed, err)
	}

	if err := composeManager.Up(ctx, proj, docker.UpOptions{Detach: true, Services: args}); err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, messages.ErrorsServiceRestartFailed, err)
	}

	for _, serviceName := range args {
		base.Output.Success(messages.SuccessRestartedService, serviceName)
	}

	return nil
}

func (h *RestartHandler) loadRegistry(sharedRoot string) (*registry.Registry, error) {
	reg := registry.NewManager(sharedRoot)
	return reg.Load()
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
		NoDeps:         flags.NoDeps,
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
