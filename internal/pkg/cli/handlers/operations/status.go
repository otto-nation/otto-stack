package operations

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/registry"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// StatusHandler handles the status command
type StatusHandler struct {
	logger *slog.Logger
}

// NewStatusHandler creates a new status handler
func NewStatusHandler() *StatusHandler {
	return &StatusHandler{
		logger: logger.GetLogger(),
	}
}

// Handle executes the status command
func (h *StatusHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	execCtx, err := common.DetectExecutionContext()
	if err != nil {
		return err
	}

	showAll, _ := cmd.Flags().GetBool(docker.FlagAll)
	showShared, _ := cmd.Flags().GetBool(docker.FlagShared)
	projectName, _ := cmd.Flags().GetString(docker.FlagProject)

	if projectName != "" {
		switch mode := execCtx.(type) {
		case *clicontext.ProjectMode:
			return h.handleProjectSharedStatus(ctx, cmd, args, base, mode, projectName)
		case *clicontext.SharedMode:
			return h.handleProjectSharedStatus(ctx, cmd, args, base, mode, projectName)
		}
	}

	switch mode := execCtx.(type) {
	case *clicontext.ProjectMode:
		if showAll || showShared {
			return h.handleSharedStatus(ctx, cmd, args, base, mode)
		}
		return h.handleProjectStatus(ctx, cmd, args, base)
	case *clicontext.SharedMode:
		return h.handleSharedStatus(ctx, cmd, args, base, mode)
	default:
		return pkgerrors.NewSystemErrorf(pkgerrors.ErrCodeInternal, messages.ErrorsContextUnknownMode, execCtx)
	}
}

func (h *StatusHandler) handleProjectStatus(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	ciFlags := ci.GetFlags(cmd)
	statusFlags, err := core.ParseStatusFlags(cmd)
	if err != nil {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldFlags, messages.ValidationFailedParseFlags, err)
	}

	if !ciFlags.Quiet {
		base.Output.Header(messages.LifecycleStatus)
	}

	setup, cleanup, err := common.SetupCoreCommand(ctx, base)
	if err != nil {
		return ci.FormatError(ciFlags, err)
	}
	defer cleanup()

	serviceConfigs, err := h.resolveServices(args, setup, &ciFlags)
	if err != nil {
		return err
	}

	statuses, err := h.getServiceStatuses(ctx, setup.Config.Project.Name, serviceConfigs, &ciFlags)
	if err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsStatusGetStatusesFailed, err)
	}

	if ciFlags.JSON || statusFlags.Format == "json" {
		h.outputJSON(&ciFlags, statuses)
		return nil
	}

	if statusFlags.Format == "yaml" {
		return h.outputYAML(statuses, base)
	}

	h.displayStatus(base, cmd, statuses, serviceConfigs)
	return nil
}

func (h *StatusHandler) handleSharedStatus(ctx context.Context, cmd *cobra.Command, _ []string, base *base.BaseCommand, mode clicontext.ExecutionMode) error {
	ciFlags := ci.GetFlags(cmd)
	showAll, _ := cmd.Flags().GetBool(docker.FlagAll)

	if !ciFlags.Quiet {
		if showAll {
			base.Output.Header(messages.InfoStatusAllProjects)
		} else {
			base.Output.Header(messages.InfoSharedContainersStatus)
		}
	}

	var sharedRoot string
	switch m := mode.(type) {
	case *clicontext.ProjectMode:
		sharedRoot = m.Shared.Root
	case *clicontext.SharedMode:
		sharedRoot = m.Shared.Root
	}

	reg := registry.NewManager(sharedRoot)
	_, err := reg.Load()
	if err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsRegistryLoadFailed, err)
	}

	sharedContainers, err := reg.List()
	if err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsStatusListSharedFailed, err)
	}

	if len(sharedContainers) == 0 {
		base.Output.Info(messages.InfoNoSharedContainers)
		return nil
	}

	dockerClient, err := docker.NewClient(h.logger)
	if err != nil {
		return pkgerrors.NewDockerError(pkgerrors.ErrCodeOperationFail, messages.ErrorsDockerClientCreateFailed, err)
	}
	defer func() { _ = dockerClient.Close() }()

	statuses := h.buildSharedStatuses(ctx, sharedContainers, dockerClient)

	if ciFlags.JSON {
		ci.OutputResult(ciFlags, display.SharedStatusResponse{
			SharedContainers: statuses,
			Count:            len(statuses),
		}, core.ExitSuccess)
		return nil
	}

	h.displaySharedStatus(base, cmd, statuses)
	return nil
}

// statusRequest encapsulates parameters for status operations
type statusRequest struct {
	ctx         context.Context
	cmd         *cobra.Command
	base        *base.BaseCommand
	mode        clicontext.ExecutionMode
	projectName string
}

func (h *StatusHandler) handleProjectSharedStatus(ctx context.Context, cmd *cobra.Command, _ []string, base *base.BaseCommand, mode clicontext.ExecutionMode, projectName string) error {
	req := statusRequest{
		ctx:         ctx,
		cmd:         cmd,
		base:        base,
		mode:        mode,
		projectName: projectName,
	}
	return h.handleProjectSharedStatusWithRequest(req)
}

func (h *StatusHandler) handleProjectSharedStatusWithRequest(req statusRequest) error {
	ciFlags := ci.GetFlags(req.cmd)

	if !ciFlags.Quiet {
		req.base.Output.Header(fmt.Sprintf(messages.InfoSharedContainersForProject, req.projectName))
	}

	var sharedRoot string
	switch m := req.mode.(type) {
	case *clicontext.ProjectMode:
		sharedRoot = m.Shared.Root
	case *clicontext.SharedMode:
		sharedRoot = m.Shared.Root
	}

	reg := registry.NewManager(sharedRoot)
	if _, err := reg.Load(); err != nil {
		return err
	}

	containers, err := reg.List()
	if err != nil {
		return err
	}

	projectContainers := make(map[string]*registry.ContainerInfo)
	for service, container := range containers {
		if h.containsProject(container.Projects, req.projectName) {
			projectContainers[service] = container
		}
	}

	if len(projectContainers) == 0 {
		req.base.Output.Info(fmt.Sprintf(messages.InfoNoSharedContainersForProject, req.projectName))
		return nil
	}

	dockerClient, err := docker.NewClient(h.logger)
	if err != nil {
		return err
	}
	defer func() { _ = dockerClient.Close() }()

	statuses := h.buildSharedStatuses(req.ctx, projectContainers, dockerClient)

	if ciFlags.JSON {
		ci.OutputResult(ciFlags, display.ProjectSharedStatusResponse{
			Project:          req.projectName,
			SharedContainers: statuses,
			Count:            len(statuses),
		}, core.ExitSuccess)
		return nil
	}

	h.displaySharedStatus(req.base, req.cmd, statuses)
	return nil
}

func (h *StatusHandler) resolveServices(args []string, setup *common.CoreSetup, ciFlags *ci.Flags) ([]types.ServiceConfig, error) {
	serviceConfigs, err := common.ResolveServiceConfigs(args, setup)
	if err != nil {
		return nil, ci.FormatError(*ciFlags, err)
	}
	return serviceConfigs, nil
}

func (h *StatusHandler) getServiceStatuses(ctx context.Context, projectName string, serviceConfigs []types.ServiceConfig, ciFlags *ci.Flags) ([]docker.ContainerStatus, error) {
	filteredServices := filterInitContainers(serviceConfigs)

	stackService, err := common.NewServiceManager(false)
	if err != nil {
		return nil, ci.FormatError(*ciFlags, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, messages.ErrorsStackCreateFailed, err))
	}

	statuses, err := stackService.Status(ctx, services.StatusRequest{
		Project:  projectName,
		Services: filteredServices,
	})
	if err != nil {
		return nil, ci.FormatError(*ciFlags, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, messages.ErrorsStackGetStatusFailed, err))
	}

	return statuses, nil
}

func (h *StatusHandler) outputJSON(ciFlags *ci.Flags, statuses []docker.ContainerStatus) {
	output := ci.StatusOutput{
		Services: make([]any, len(statuses)),
		Count:    len(statuses),
	}
	for i, s := range statuses {
		output.Services[i] = s
	}
	ci.OutputResult(*ciFlags, output, core.ExitSuccess)
}

func (h *StatusHandler) outputYAML(statuses []docker.ContainerStatus, base *base.BaseCommand) error {
	return yaml.NewEncoder(base.Output.Writer()).Encode(statuses)
}

func (h *StatusHandler) displayStatus(base *base.BaseCommand, cmd *cobra.Command, statuses []docker.ContainerStatus, serviceConfigs []types.ServiceConfig) {
	if len(statuses) == 0 {
		return
	}
	verbose := base.GetVerbose(cmd)
	_ = display.RenderStatusTable(base.Output.Writer(), statuses, serviceConfigs, !verbose, base.Output.GetNoColor())
}

// ValidateArgs validates the command arguments
func (h *StatusHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *StatusHandler) GetRequiredFlags() []string {
	return []string{}
}

// convertToDisplayStatuses creates display service statuses with health inheritance

// filterInitContainers removes init containers from status display
func filterInitContainers(serviceConfigs []types.ServiceConfig) []string {
	filtered := make([]string, 0, len(serviceConfigs))
	for _, config := range serviceConfigs {
		if config.Container.Restart != types.RestartPolicyNo {
			filtered = append(filtered, config.Name)
		}
	}
	return filtered
}

func (h *StatusHandler) buildSharedStatuses(ctx context.Context, containers map[string]*registry.ContainerInfo, dockerClient *docker.Client) []display.SharedContainerStatus {
	statuses := make([]display.SharedContainerStatus, 0, len(containers))
	now := time.Now()

	for service, container := range containers {
		cs := dockerClient.InspectContainer(ctx, container.Name)
		projectNames := make([]string, len(container.Projects))
		for i, ref := range container.Projects {
			projectNames[i] = ref.Name
		}
		uptime := time.Duration(0)
		if !cs.StartedAt.IsZero() {
			uptime = now.Sub(cs.StartedAt)
		}
		statuses = append(statuses, display.SharedContainerStatus{
			Name:      container.Name,
			Service:   service,
			State:     cs.State,
			Health:    cs.Health,
			Ports:     cs.Ports,
			Uptime:    uptime,
			Projects:  projectNames,
			CreatedAt: container.CreatedAt,
			UpdatedAt: container.UpdatedAt,
		})
	}

	return statuses
}

func (h *StatusHandler) displaySharedStatus(base *base.BaseCommand, cmd *cobra.Command, statuses []display.SharedContainerStatus) {
	if len(statuses) == 0 {
		base.Output.Info(messages.InfoNoSharedContainers)
		return
	}
	verbose := base.GetVerbose(cmd)
	_ = display.RenderSharedStatusTable(base.Output.Writer(), statuses, !verbose, base.Output.GetNoColor())
}

func (h *StatusHandler) containsProject(projects []registry.ProjectRef, projectName string) bool {
	return slices.ContainsFunc(projects, func(r registry.ProjectRef) bool { return r.Name == projectName })
}
