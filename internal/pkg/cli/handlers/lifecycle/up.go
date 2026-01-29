package lifecycle

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/registry"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"github.com/otto-nation/otto-stack/internal/pkg/validation"
)

// UpHandler handles the up command
type UpHandler struct{}

// NewUpHandler creates a new up handler
func NewUpHandler() *UpHandler {
	return &UpHandler{}
}

// Handle executes the up command
func (h *UpHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
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

func (h *UpHandler) handleProjectContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand, execCtx *clicontext.ExecutionContext) error {
	base.Output.Header("%s", core.MsgLifecycle_starting)

	// Validate flags
	if err := validation.ValidateUpFlags(cmd); err != nil {
		return err
	}

	setup, cleanup, err := common.SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	serviceConfigs, err := common.ResolveServiceConfigs(args, setup)
	if err != nil {
		return err
	}

	// Register shared containers before starting
	sharedConfigs := h.filterSharedServices(serviceConfigs, setup.Config)
	if len(sharedConfigs) > 0 {
		if err := h.registerSharedContainersForProject(sharedConfigs, setup.Config.Project.Name, execCtx, base); err != nil {
			return err
		}
	}

	service, err := common.NewServiceManager(false)
	if err != nil {
		return pkgerrors.NewServiceError(common.ComponentStack, common.ActionCreateService, err)
	}

	upFlags, _ := core.ParseUpFlags(cmd)
	force, _ := cmd.Flags().GetBool(core.FlagForce)

	const defaultTimeout = 30 * time.Second
	timeout, err := time.ParseDuration(upFlags.Timeout)
	if err != nil {
		timeout = defaultTimeout
	}

	startRequest := services.StartRequest{
		Project:        setup.Config.Project.Name,
		ServiceConfigs: serviceConfigs,
		Build:          upFlags.Build || force, // force implies build
		ForceRecreate:  upFlags.ForceRecreate,
		Detach:         upFlags.Detach,
		NoDeps:         upFlags.NoDeps,
		Timeout:        timeout,
	}

	if err = service.Start(ctx, startRequest); err != nil {
		return pkgerrors.NewServiceError(common.ComponentStack, common.ActionStartServices, err)
	}

	base.Output.Success("Services started successfully")
	base.Output.Info("Project: %s", setup.Config.Project.Name)
	for _, svc := range serviceConfigs {
		base.Output.Info("  %s %s", ui.IconSuccess, svc.Name)
	}

	return nil
}

func (h *UpHandler) handleGlobalContext(_ context.Context, _ *cobra.Command, args []string, base *base.BaseCommand, execCtx *clicontext.ExecutionContext) error {
	base.Output.Header(core.MsgShared_starting)

	if len(args) == 0 {
		return pkgerrors.NewValidationError("", core.MsgShared_no_services_specified, nil)
	}

	serviceConfigs, err := h.loadServiceConfigs(args)
	if err != nil {
		return err
	}

	if err := h.registerSharedContainers(serviceConfigs, execCtx, base); err != nil {
		return err
	}

	h.displaySuccess(base, serviceConfigs)
	return nil
}

func (h *UpHandler) loadServiceConfigs(serviceNames []string) ([]types.ServiceConfig, error) {
	serviceUtils := services.NewServiceUtils()
	var configs []types.ServiceConfig

	for _, name := range serviceNames {
		cfg, err := serviceUtils.LoadServiceConfig(name)
		if err != nil {
			return nil, pkgerrors.NewServiceError(common.ComponentService, common.ActionLoadService, err)
		}
		if cfg != nil {
			configs = append(configs, *cfg)
		}
	}

	return configs, nil
}

func (h *UpHandler) registerSharedContainers(serviceConfigs []types.ServiceConfig, execCtx *clicontext.ExecutionContext, base *base.BaseCommand) error {
	return h.registerSharedContainersForProject(serviceConfigs, common.ContextGlobal, execCtx, base)
}

func (h *UpHandler) registerSharedContainersForProject(serviceConfigs []types.ServiceConfig, projectName string, execCtx *clicontext.ExecutionContext, base *base.BaseCommand) error {
	reg := registry.NewManager(execCtx.Shared.Root)

	for _, svc := range serviceConfigs {
		containerName := core.SharedContainerPrefix + svc.Name
		if err := reg.Register(svc.Name, containerName, projectName); err != nil {
			base.Output.Warning("Failed to register %s: %v", svc.Name, err)
		}
	}

	return nil
}

func (h *UpHandler) filterSharedServices(serviceConfigs []types.ServiceConfig, cfg *config.Config) []types.ServiceConfig {
	if !cfg.Sharing.Enabled {
		return nil
	}

	// If no specific services listed, all are shared
	if len(cfg.Sharing.Services) == 0 {
		return serviceConfigs
	}

	// Filter to only services marked as shared
	var shared []types.ServiceConfig
	for _, svc := range serviceConfigs {
		if cfg.Sharing.Services[svc.Name] {
			shared = append(shared, svc)
		}
	}
	return shared
}

func (h *UpHandler) displaySuccess(base *base.BaseCommand, serviceConfigs []types.ServiceConfig) {
	base.Output.Success(core.MsgShared_registered)
	for _, svc := range serviceConfigs {
		base.Output.Info("  %s %s (shared)", ui.IconSuccess, svc.Name)
	}
}

// ValidateArgs validates the command arguments
func (h *UpHandler) ValidateArgs(args []string) error {
	return validation.ValidateUpArgs(args)
}

// GetRequiredFlags returns required flags for this command
func (h *UpHandler) GetRequiredFlags() []string {
	return []string{}
}
