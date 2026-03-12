package lifecycle

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

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
)

// DownHandler handles the down command
type DownHandler struct{}

// NewDownHandler creates a new down handler
func NewDownHandler() *DownHandler {
	return &DownHandler{}
}

// Handle executes the down command
func (h *DownHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	if projectDir, _ := cmd.Flags().GetString(docker.FlagProject); projectDir != "" {
		mode, err := buildProjectMode(projectDir)
		if err != nil {
			return err
		}
		return h.handleProjectContext(ctx, cmd, args, base, mode)
	}

	execCtx, err := common.DetectExecutionContext()
	if err != nil {
		return err
	}

	showShared, _ := cmd.Flags().GetBool(docker.FlagShared)
	showAll, _ := cmd.Flags().GetBool(docker.FlagAll)

	if showShared {
		switch mode := execCtx.(type) {
		case *clicontext.ProjectMode:
			return h.handleGlobalContext(ctx, cmd, args, base, mode.Shared)
		case *clicontext.SharedMode:
			return h.handleGlobalContext(ctx, cmd, args, base, mode.Shared)
		}
	}

	switch mode := execCtx.(type) {
	case *clicontext.ProjectMode:
		if showAll {
			// --all: stop project containers first (without prompting about shared),
			// then stop all shared containers.
			if err := h.handleProjectContext(ctx, cmd, args, base, mode); err != nil {
				return err
			}
			return h.handleGlobalContext(ctx, cmd, args, base, mode.Shared)
		}
		return h.handleProjectContext(ctx, cmd, args, base, mode)
	case *clicontext.SharedMode:
		return h.handleGlobalContext(ctx, cmd, args, base, mode.Shared)
	default:
		return pkgerrors.NewSystemErrorf(pkgerrors.ErrCodeInternal, messages.ErrorsContextUnknownMode, execCtx)
	}
}

func (h *DownHandler) handleProjectContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand, execCtx *clicontext.ProjectMode) error {
	base.Output.Header(messages.LifecycleStopping)

	setup, cleanup, err := common.SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	serviceConfigs, err := common.ResolveServiceConfigs(args, setup)
	if err != nil {
		return err
	}

	stopAll, _ := cmd.Flags().GetBool(docker.FlagAll)
	ciFlags := ci.GetFlags(cmd)
	// When --all is set, include shared containers without prompting — the user's intent is explicit.
	if !stopAll {
		serviceConfigs, err = h.filterSharedIfNeeded(serviceConfigs, execCtx.Shared.Root, base, ciFlags.NonInteractive)
		if err != nil {
			return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsServiceFilterSharedFailed, err)
		}
	}

	if len(serviceConfigs) == 0 {
		base.Output.Info(messages.SharedNoServicesToStop)
		return nil
	}

	service, err := h.stopServices(ctx, cmd, setup, serviceConfigs)
	if err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentServices, messages.ErrorsStackStopFailed, err)
	}

	base.Output.Success(messages.SuccessServicesStopped)
	base.Output.Muted(messages.InfoProjectInfo, setup.Config.Project.Name)

	filteredNames := filterStatusQueryNames(serviceConfigs)
	if statuses, statusErr := service.Status(ctx, services.StatusRequest{
		Project:  setup.Config.Project.Name,
		Services: filteredNames,
	}); statusErr == nil {
		// Silent fallback: if Status() fails, the command already succeeded — skip the table
		_ = display.RenderStatusTable(base.Output.Writer(), statuses, serviceConfigs, true, base.Output.GetNoColor())
	}

	// Unregister shared containers after stopping
	return h.unregisterSharedContainersForProject(serviceConfigs, setup.Config.Project.Name, execCtx.Shared.Root, base)
}

func (h *DownHandler) filterSharedIfNeeded(serviceConfigs []types.ServiceConfig, sharedRoot string, base *base.BaseCommand, nonInteractive bool) ([]types.ServiceConfig, error) {
	reg := registry.NewManager(sharedRoot)
	_, err := reg.Load()
	if err != nil {
		return nil, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentRegistry, messages.ErrorsRegistryLoadFailed, err)
	}

	sharedServices := h.findSharedServices(serviceConfigs, reg)
	if len(sharedServices) == 0 {
		return serviceConfigs, nil
	}

	return h.promptAndFilterShared(base, reg, sharedServices, serviceConfigs, nonInteractive), nil
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

func (h *DownHandler) promptAndFilterShared(base *base.BaseCommand, reg *registry.Manager, sharedServices []string, serviceConfigs []types.ServiceConfig, nonInteractive bool) []types.ServiceConfig {
	base.Output.Warning(messages.SharedWillStop)
	for _, svc := range sharedServices {
		container, err := reg.Get(svc)
		if err == nil && container != nil {
			base.Output.Info(messages.InfoListItemWithUsers, svc, container.Projects)
		}
	}

	if nonInteractive {
		base.Output.Info(messages.SharedSkippingNonInteractive, len(sharedServices))
		return h.filterOutShared(sharedServices, serviceConfigs)
	}

	if !h.promptStopShared(base) {
		base.Output.Info(messages.SharedSkipping)
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

func (h *DownHandler) stopServices(ctx context.Context, cmd *cobra.Command, setup *common.CoreSetup, serviceConfigs []types.ServiceConfig) (*services.Service, error) {
	service, err := common.NewServiceManager(false)
	if err != nil {
		return nil, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, messages.ErrorsStackStopFailed, err)
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
		return nil, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, messages.ErrorsStackStopFailed, err)
	}

	return service, nil
}

func (h *DownHandler) handleGlobalContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand, sharedInfo *clicontext.SharedInfo) error {
	base.Output.Header(messages.SharedStopping)

	nonInteractive := ci.GetFlags(cmd).NonInteractive
	servicesToStop, err := h.determineServicesToStop(args, sharedInfo, base, nonInteractive)
	if err != nil {
		return err
	}

	if servicesToStop == nil {
		return nil
	}

	h.stopSharedContainersViaCompose(ctx, sharedInfo.Root, servicesToStop, base)

	return h.unregisterSharedContainers(servicesToStop, sharedInfo, base)
}

func (h *DownHandler) stopSharedContainersViaCompose(ctx context.Context, sharedRoot string, services []string, base *base.BaseCommand) {
	composePath := filepath.Join(sharedRoot, core.GeneratedDir, docker.DockerComposeFileName)
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		base.Output.Warning("%s", messages.SharedComposeFileNotFound)
		return
	}

	composeManager, err := docker.NewManager()
	if err != nil {
		base.Output.Warning(messages.ErrorsDockerManagerCreateFailed, err)
		return
	}

	proj, err := composeManager.LoadProject(ctx, []string{composePath}, "shared")
	if err != nil {
		base.Output.Warning(messages.ErrorsDockerLoadProjectFailed, err)
		return
	}

	// Ignore errors from down — containers may already be stopped
	_ = composeManager.Down(ctx, proj, docker.DownOptions{Services: services, RemoveOrphans: true}.ToSDK())
}

func (h *DownHandler) determineServicesToStop(args []string, sharedInfo *clicontext.SharedInfo, base *base.BaseCommand, nonInteractive bool) ([]string, error) {
	if len(args) > 0 {
		return args, nil
	}

	reg := registry.NewManager(sharedInfo.Root)
	_, err := reg.Load()
	if err != nil {
		return nil, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentRegistry, messages.ErrorsRegistryLoadFailed, err)
	}

	containers, err := reg.List()
	if err != nil {
		return nil, pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsStatusListSharedFailed, err)
	}

	if len(containers) == 0 {
		base.Output.Info(messages.SharedNoContainers)
		return nil, nil
	}

	return h.promptStopAll(containers, base, nonInteractive), nil
}

func (h *DownHandler) promptStopAll(containers map[string]*registry.ContainerInfo, base *base.BaseCommand, nonInteractive bool) []string {
	base.Output.Warning(messages.SharedStopAllPrompt)
	services := make([]string, 0, len(containers))
	for service := range containers {
		services = append(services, service)
	}
	sort.Strings(services)
	for _, service := range services {
		base.Output.Info(messages.InfoListItemWithUsers, service, containers[service].Projects)
	}

	if nonInteractive {
		base.Output.Info(messages.SharedStoppingNonInteractive)
		return services
	}

	if !h.promptStopShared(base) {
		base.Output.Info(messages.SharedCancelled)
		return nil
	}

	return services
}

func (h *DownHandler) unregisterSharedContainers(servicesToStop []string, sharedInfo *clicontext.SharedInfo, base *base.BaseCommand) error {
	return h.unregisterSharedContainersForProject(h.serviceNamesToConfigs(servicesToStop), "global", sharedInfo.Root, base)
}

func (h *DownHandler) unregisterSharedContainersForProject(serviceConfigs []types.ServiceConfig, projectName string, sharedRoot string, base *base.BaseCommand) error {
	reg := registry.NewManager(sharedRoot)

	for _, svc := range serviceConfigs {
		if err := reg.Unregister(svc.Name, projectName); err != nil {
			base.Output.Warning(messages.WarningsRegistryUnregisterFailed, svc.Name, err)
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
	prompt := &survey.Confirm{
		Message: messages.PromptsStopSharedContainers,
		Default: false,
	}
	var confirmed bool
	if err := survey.AskOne(prompt, &confirmed); err != nil {
		return false
	}
	return confirmed
}

// ValidateArgs validates the command arguments
func (h *DownHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *DownHandler) GetRequiredFlags() []string {
	return []string{}
}
