package lifecycle

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
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

	execCtx, err := detector.DetectContext()
	if err != nil {
		return err
	}

	switch mode := execCtx.(type) {
	case *clicontext.ProjectMode:
		return h.handleProjectContext(ctx, cmd, args, base, mode)
	case *clicontext.SharedMode:
		return h.handleGlobalContext(ctx, cmd, args, base, mode)
	default:
		return fmt.Errorf("unknown execution mode: %T", execCtx)
	}
}

func (h *UpHandler) handleProjectContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand, execCtx *clicontext.ProjectMode) error {
	base.Output.Header("%s", messages.LifecycleStarting)

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
		if err := h.registerSharedContainersForProject(sharedConfigs, setup.Config.Project.Name, execCtx.Shared.Root, base); err != nil {
			return err
		}
	}

	service, err := common.NewServiceManager(false)
	if err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, messages.ErrorsStackStartFailed, err)
	}

	upFlags, _ := core.ParseUpFlags(cmd)
	force, _ := cmd.Flags().GetBool(pkgerrors.FieldFlags)

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
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, messages.ErrorsStackStartFailed, err)
	}

	base.Output.Success(messages.SuccessServicesStarted)
	base.Output.Info(messages.InfoProjectInfo, setup.Config.Project.Name)
	for _, svc := range serviceConfigs {
		base.Output.Info("  %s %s", ui.IconSuccess, svc.Name)
	}

	return nil
}

func (h *UpHandler) handleGlobalContext(_ context.Context, _ *cobra.Command, args []string, base *base.BaseCommand, execCtx *clicontext.SharedMode) error {
	base.Output.Header(messages.SharedStarting)

	if len(args) == 0 {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationNoServicesSelected, nil)
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
			return nil, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentService, messages.ErrorsServiceLoadFailed, err)
		}
		if cfg != nil {
			configs = append(configs, *cfg)
		}
	}

	return configs, nil
}

func (h *UpHandler) registerSharedContainers(serviceConfigs []types.ServiceConfig, execCtx *clicontext.SharedMode, base *base.BaseCommand) error {
	return h.registerSharedContainersForProject(serviceConfigs, "global", execCtx.Shared.Root, base)
}

func (h *UpHandler) registerSharedContainersForProject(serviceConfigs []types.ServiceConfig, projectName string, sharedRoot string, base *base.BaseCommand) error {
	reg := registry.NewManager(sharedRoot)

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

	var shared []types.ServiceConfig
	for _, svc := range serviceConfigs {
		// Skip services marked as non-shareable
		if !svc.Shareable {
			continue
		}

		// If no specific services listed, include all shareable services
		if len(cfg.Sharing.Services) == 0 {
			shared = append(shared, svc)
			continue
		}

		// Include only if explicitly listed in sharing config
		if cfg.Sharing.Services[svc.Name] {
			shared = append(shared, svc)
		}
	}
	return shared
}

func (h *UpHandler) displaySuccess(base *base.BaseCommand, serviceConfigs []types.ServiceConfig) {
	base.Output.Success(messages.SharedRegistered)
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
