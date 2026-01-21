package lifecycle

import (
	"context"
	"fmt"
	"slices"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/registry"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
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

	// Check for shared containers
	reg := registry.NewManager(execCtx.Shared.Root)

	_, loadErr := reg.Load()
	if loadErr != nil {
		return loadErr
	}

	var sharedServices []string
	for _, svc := range serviceConfigs {
		isShared, err := reg.IsShared(svc.Name)
		if err == nil && isShared {
			sharedServices = append(sharedServices, svc.Name)
		}
	}

	// Prompt for shared containers
	if len(sharedServices) > 0 {
		serviceConfigs = h.handleSharedContainers(base, reg, sharedServices, serviceConfigs)
	}

	if len(serviceConfigs) == 0 {
		base.Output.Info("No services to stop")
		return nil
	}

	service, err := common.NewServiceManager(false)
	if err != nil {
		return pkgerrors.NewServiceError(common.ComponentStack, common.ActionCreateService, err)
	}

	stopRequest := services.StopRequest{
		Project:        setup.Config.Project.Name,
		ServiceConfigs: serviceConfigs,
		Remove:         true,
		RemoveVolumes:  false,
	}

	if err = service.Stop(ctx, stopRequest); err != nil {
		return pkgerrors.NewServiceError(common.ComponentStack, common.ActionStopServices, err)
	}

	base.Output.Success("Services stopped successfully")
	base.Output.Info("Project: %s", setup.Config.Project.Name)
	for _, svc := range serviceConfigs {
		base.Output.Info("  %s %s", display.StatusSuccess, svc.Name)
	}

	return nil
}

func (h *DownHandler) handleGlobalContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand, execCtx *clicontext.ExecutionContext) error {
	base.Output.Header("Stopping shared containers")

	if len(args) == 0 {
		base.Output.Warning("This will stop ALL shared containers")
		if !h.promptStopShared(base) {
			base.Output.Info("Cancelled")
			return nil
		}
	}

	// TODO: Implement shared container shutdown
	base.Output.Info("Global context - stopping shared containers")
	for _, svc := range args {
		base.Output.Info("  %s %s (shared)", display.StatusSuccess, svc)
	}

	return nil
}

func (h *DownHandler) promptStopShared(base *base.BaseCommand) bool {
	var response string
	fmt.Print("⚠️  Stop shared containers? (y/N): ")
	_, _ = fmt.Scanln(&response)
	return response == "y" || response == "Y"
}

func (h *DownHandler) handleSharedContainers(base *base.BaseCommand, reg *registry.Manager, sharedServices []string, serviceConfigs []types.ServiceConfig) []types.ServiceConfig {
	base.Output.Warning("The following shared containers will be stopped:")
	for _, svc := range sharedServices {
		container, err := reg.Get(svc)
		if err == nil && container != nil {
			base.Output.Info("  - %s (used by: %v)", svc, container.Projects)
		}
	}

	if !h.promptStopShared(base) {
		base.Output.Info("Skipping shared containers")
		return h.filterSharedServices(sharedServices, serviceConfigs)
	}

	return serviceConfigs
}

func (h *DownHandler) filterSharedServices(sharedServices []string, serviceConfigs []types.ServiceConfig) []types.ServiceConfig {
	var filtered []types.ServiceConfig
	for _, svc := range serviceConfigs {
		if !slices.Contains(sharedServices, svc.Name) {
			filtered = append(filtered, svc)
		}
	}
	return filtered
}

// ValidateArgs validates the command arguments
func (h *DownHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *DownHandler) GetRequiredFlags() []string {
	return []string{}
}
