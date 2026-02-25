package lifecycle

import (
	"context"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
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
	execCtx, err := common.DetectExecutionContext()
	if err != nil {
		return err
	}

	switch mode := execCtx.(type) {
	case *clicontext.ProjectMode:
		return h.handleProjectContext(ctx, cmd, args, base, mode)
	case *clicontext.SharedMode:
		return h.handleGlobalContext(ctx, cmd, args, base, mode)
	default:
		return pkgerrors.NewSystemErrorf(pkgerrors.ErrCodeInternal, messages.ErrorsContextUnknownMode, execCtx)
	}
}

func (h *UpHandler) handleProjectContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand, execCtx *clicontext.ProjectMode) error {
	base.Output.Header("%s", messages.LifecycleStarting)

	// Validate flags
	if err := validation.ValidateUpFlags(cmd); err != nil {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldFlags, messages.ValidationFailed, err)
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
		configDir, _ := filepath.Abs(core.OttoStackDir)
		project := registry.ProjectRef{Name: setup.Config.Project.Name, ConfigDir: configDir}
		if err := h.registerSharedContainersForProject(sharedConfigs, project, execCtx.Shared.Root, base); err != nil {
			return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsServiceRegisterSharedFailed, err)
		}
	}

	service, err := common.NewServiceManager(false)
	if err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, messages.ErrorsStackStartFailed, err)
	}

	upFlags, err := core.ParseUpFlags(cmd)
	if err != nil {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldFlags, messages.ValidationFailedParseFlags, err)
	}

	const defaultTimeout = 30 * time.Second
	timeout, err := time.ParseDuration(upFlags.Timeout)
	if err != nil {
		timeout = defaultTimeout
	}

	startRequest := services.StartRequest{
		Project:        setup.Config.Project.Name,
		ServiceConfigs: serviceConfigs,
		Build:          upFlags.Build,
		ForceRecreate:  upFlags.ForceRecreate,
		Detach:         upFlags.Detach,
		NoDeps:         upFlags.NoDeps,
		Timeout:        timeout,
	}

	if err = service.Start(ctx, startRequest); err != nil {
		return err
	}

	base.Output.Success(messages.SuccessServicesStarted)
	base.Output.Muted(messages.InfoProjectInfo, setup.Config.Project.Name)

	filteredNames := filterStatusQueryNames(serviceConfigs)
	if statuses, err := service.Status(ctx, services.StatusRequest{
		Project:  setup.Config.Project.Name,
		Services: filteredNames,
	}); err == nil {
		// Silent fallback: if Status() fails, the command already succeeded — skip the table
		_ = display.RenderStatusTable(base.Output.Writer(), statuses, serviceConfigs, true, base.Output.GetNoColor())
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

	if err := h.validateShareableServices(serviceConfigs); err != nil {
		return err
	}

	if err := h.registerSharedContainers(serviceConfigs, execCtx, base); err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsServiceRegisterSharedFailed, err)
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
			return nil, err
		}
		if cfg != nil {
			configs = append(configs, *cfg)
		}
	}

	return configs, nil
}

func (h *UpHandler) registerSharedContainers(serviceConfigs []types.ServiceConfig, execCtx *clicontext.SharedMode, base *base.BaseCommand) error {
	project := registry.ProjectRef{Name: "global", ConfigDir: execCtx.Shared.Root}
	return h.registerSharedContainersForProject(serviceConfigs, project, execCtx.Shared.Root, base)
}

func (h *UpHandler) registerSharedContainersForProject(serviceConfigs []types.ServiceConfig, project registry.ProjectRef, sharedRoot string, base *base.BaseCommand) error {
	reg := registry.NewManager(sharedRoot)

	// Auto-heal: purge any non-shareable entries from previous bugs
	if shareableMap, err := h.buildShareableMap(); err == nil {
		if err := reg.PurgeNonShareable(shareableMap); err != nil {
			base.Output.Warning(messages.WarningsRegistryCleanFailed, err)
		}
	}

	for _, svc := range serviceConfigs {
		if !svc.Shareable {
			continue // defense in depth: validateShareableServices should catch this first
		}
		containerName := core.SharedContainerPrefix + svc.Name
		if err := reg.Register(svc.Name, containerName, project); err != nil {
			base.Output.Warning(messages.WarningsRegistryRegisterFailed, svc.Name, err)
		}
	}

	return nil
}

func (h *UpHandler) validateShareableServices(serviceConfigs []types.ServiceConfig) error {
	for _, svc := range serviceConfigs {
		if !svc.Shareable {
			return pkgerrors.NewValidationErrorf(
				pkgerrors.ErrCodeInvalid,
				pkgerrors.FieldServiceName,
				messages.ValidationServiceNotShareable,
				svc.Name,
			)
		}
	}
	return nil
}

func (h *UpHandler) buildShareableMap() (map[string]bool, error) {
	mgr, err := services.New()
	if err != nil {
		return nil, err
	}
	allServices := mgr.GetAllServices()
	shareableMap := make(map[string]bool, len(allServices))
	for name, cfg := range allServices {
		shareableMap[name] = cfg.Shareable
	}
	return shareableMap, nil
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
		base.Output.Info(messages.SharedRegisteredItem, ui.IconOK, svc.Name)
	}
}

// filterStatusQueryNames returns service names to include in a Docker status query.
// It excludes init containers (RestartPolicy: no) which have no persistent running containers,
// but includes hidden services since they may act as providers for other services.
func filterStatusQueryNames(configs []types.ServiceConfig) []string {
	names := make([]string, 0, len(configs))
	for _, cfg := range configs {
		if cfg.Container.Restart != types.RestartPolicyNo {
			names = append(names, cfg.Name)
		}
	}
	return names
}

// ValidateArgs validates the command arguments
func (h *UpHandler) ValidateArgs(args []string) error {
	return validation.ValidateUpArgs(args)
}

// GetRequiredFlags returns required flags for this command
func (h *UpHandler) GetRequiredFlags() []string {
	return []string{}
}
