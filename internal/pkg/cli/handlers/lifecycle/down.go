package lifecycle

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/registry"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
)

// DownHandler handles the down command
type DownHandler struct{}

// NewDownHandler creates a new down handler
func NewDownHandler() *DownHandler {
	return &DownHandler{}
}

// Handle executes the down command
func (h *DownHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	detector, err := clicontext.NewDetector()
	if err != nil {
		return err
	}

	execCtx, err := detector.Detect()
	if err != nil {
		return err
	}

	if execCtx.Type == clicontext.Global {
		return h.handleGlobalContext(ctx, cmd, args, base, execCtx)
	}

	return h.handleProjectContext(ctx, cmd, args, base, execCtx)
}

func (h *DownHandler) handleProjectContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand, execCtx *clicontext.ExecutionContext) error {
	base.Output.Header(core.MsgLifecycle_stopping)

	setup, cleanup, err := common.SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	serviceConfigs, err := common.ResolveServiceConfigs(args, setup)
	if err != nil {
		return err
	}

	serviceConfigs, err = h.filterSharedIfNeeded(serviceConfigs, execCtx, base)
	if err != nil {
		return err
	}

	if len(serviceConfigs) == 0 {
		base.Output.Info(core.MsgShared_no_services_to_stop)
		return nil
	}

	if err := h.stopServices(ctx, cmd, setup, serviceConfigs, base); err != nil {
		return err
	}

	// Unregister shared containers after stopping
	return h.unregisterSharedContainersForProject(serviceConfigs, setup.Config.Project.Name, execCtx, base)
}

func (h *DownHandler) filterSharedIfNeeded(serviceConfigs []types.ServiceConfig, execCtx *clicontext.ExecutionContext, base *base.BaseCommand) ([]types.ServiceConfig, error) {
	reg := registry.NewManager(execCtx.Shared.Root)
	_, err := reg.Load()
	if err != nil {
		return nil, pkgerrors.NewServiceError(common.ComponentRegistry, common.ActionLoadRegistry, err)
	}

	sharedServices := h.findSharedServices(serviceConfigs, reg)
	if len(sharedServices) == 0 {
		return serviceConfigs, nil
	}

	return h.promptAndFilterShared(base, reg, sharedServices, serviceConfigs), nil
}

func (h *DownHandler) findSharedServices(serviceConfigs []types.ServiceConfig, reg *registry.Manager) []string {
	var shared []string
	for _, svc := range serviceConfigs {
		isShared, err := reg.IsShared(svc.Name)
		if err == nil && isShared {
			shared = append(shared, svc.Name)
		}
	}
	return shared
}

func (h *DownHandler) promptAndFilterShared(base *base.BaseCommand, reg *registry.Manager, sharedServices []string, serviceConfigs []types.ServiceConfig) []types.ServiceConfig {
	base.Output.Warning(core.MsgShared_will_stop)
	for _, svc := range sharedServices {
		container, err := reg.Get(svc)
		if err == nil && container != nil {
			base.Output.Info("  - %s (used by: %v)", svc, container.Projects)
		}
	}

	if !h.promptStopShared(base) {
		base.Output.Info(core.MsgShared_skipping)
		return h.filterOutShared(sharedServices, serviceConfigs)
	}

	return serviceConfigs
}

func (h *DownHandler) filterOutShared(sharedServices []string, serviceConfigs []types.ServiceConfig) []types.ServiceConfig {
	var filtered []types.ServiceConfig
	for _, svc := range serviceConfigs {
		if !slices.Contains(sharedServices, svc.Name) {
			filtered = append(filtered, svc)
		}
	}
	return filtered
}

func (h *DownHandler) stopServices(ctx context.Context, cmd *cobra.Command, setup *common.CoreSetup, serviceConfigs []types.ServiceConfig, base *base.BaseCommand) error {
	service, err := common.NewServiceManager(false)
	if err != nil {
		return pkgerrors.NewServiceError(common.ComponentStack, common.ActionCreateService, err)
	}

	downFlags, _ := core.ParseDownFlags(cmd)

	timeout := time.Duration(downFlags.Timeout) * time.Second

	stopRequest := services.StopRequest{
		Project:        setup.Config.Project.Name,
		ServiceConfigs: serviceConfigs,
		Remove:         true,
		RemoveVolumes:  downFlags.Volumes,
		RemoveOrphans:  downFlags.RemoveOrphans,
		Timeout:        timeout,
	}

	if err = service.Stop(ctx, stopRequest); err != nil {
		return pkgerrors.NewServiceError(common.ComponentStack, common.ActionStopServices, err)
	}

	h.displayStopSuccess(base, setup, serviceConfigs)
	return nil
}

func (h *DownHandler) displayStopSuccess(base *base.BaseCommand, setup *common.CoreSetup, serviceConfigs []types.ServiceConfig) {
	base.Output.Success("Services stopped successfully")
	base.Output.Info("Project: %s", setup.Config.Project.Name)
	for _, svc := range serviceConfigs {
		base.Output.Info("  %s %s", ui.IconSuccess, svc.Name)
	}
}

func (h *DownHandler) handleGlobalContext(_ context.Context, _ *cobra.Command, args []string, base *base.BaseCommand, execCtx *clicontext.ExecutionContext) error {
	base.Output.Header(core.MsgShared_stopping)

	servicesToStop, err := h.determineServicesToStop(args, execCtx, base)
	if err != nil {
		return err
	}

	if servicesToStop == nil {
		return nil
	}

	return h.unregisterSharedContainers(servicesToStop, execCtx, base)
}

func (h *DownHandler) determineServicesToStop(args []string, execCtx *clicontext.ExecutionContext, base *base.BaseCommand) ([]string, error) {
	if len(args) > 0 {
		return args, nil
	}

	reg := registry.NewManager(execCtx.Shared.Root)
	_, err := reg.Load()
	if err != nil {
		return nil, pkgerrors.NewServiceError(common.ComponentRegistry, common.ActionLoadRegistry, err)
	}

	containers, err := reg.List()
	if err != nil {
		return nil, err
	}

	if len(containers) == 0 {
		base.Output.Info(core.MsgShared_no_containers)
		return nil, nil
	}

	return h.promptStopAll(containers, base), nil
}

func (h *DownHandler) promptStopAll(containers []*registry.ContainerInfo, base *base.BaseCommand) []string {
	base.Output.Warning(core.MsgShared_stop_all_prompt)
	var services []string
	for _, container := range containers {
		base.Output.Info("  - %s (used by: %v)", container.Name, container.Projects)
		services = append(services, container.Name)
	}

	if !h.promptStopShared(base) {
		base.Output.Info(core.MsgShared_cancelled)
		return nil
	}

	return services
}

func (h *DownHandler) unregisterSharedContainers(servicesToStop []string, execCtx *clicontext.ExecutionContext, base *base.BaseCommand) error {
	return h.unregisterSharedContainersForProject(h.serviceNamesToConfigs(servicesToStop), common.ContextGlobal, execCtx, base)
}

func (h *DownHandler) unregisterSharedContainersForProject(serviceConfigs []types.ServiceConfig, projectName string, execCtx *clicontext.ExecutionContext, base *base.BaseCommand) error {
	reg := registry.NewManager(execCtx.Shared.Root)

	for _, svc := range serviceConfigs {
		if err := reg.Unregister(svc.Name, projectName); err != nil {
			base.Output.Warning("Failed to unregister %s: %v", svc.Name, err)
		}
	}

	return nil
}

func (h *DownHandler) serviceNamesToConfigs(serviceNames []string) []types.ServiceConfig {
	configs := make([]types.ServiceConfig, len(serviceNames))
	for i, name := range serviceNames {
		configs[i] = types.ServiceConfig{Name: name}
	}
	return configs
}

func (h *DownHandler) promptStopShared(_ *base.BaseCommand) bool {
	var response string
	fmt.Print("⚠️  Stop shared containers? (y/N): ")
	_, _ = fmt.Scanln(&response)
	return response == "y" || response == "Y"
}

// ValidateArgs validates the command arguments
func (h *DownHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *DownHandler) GetRequiredFlags() []string {
	return []string{}
}
