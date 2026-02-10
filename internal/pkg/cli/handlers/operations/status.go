package operations

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

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
	detector, err := clicontext.NewDetector()
	if err != nil {
		return err
	}

	execCtx, err := detector.Detect()
	if err != nil {
		return err
	}

	showAll, _ := cmd.Flags().GetBool(docker.FlagAll)
	showShared, _ := cmd.Flags().GetBool(docker.FlagShared)
	projectName, _ := cmd.Flags().GetString(docker.FlagProject)

	if projectName != "" {
		return h.handleProjectSharedStatus(ctx, cmd, args, base, execCtx, projectName)
	}

	if execCtx.Type == clicontext.Shared || showAll || showShared {
		return h.handleSharedStatus(ctx, cmd, args, base, execCtx)
	}

	return h.handleProjectStatus(ctx, cmd, args, base, execCtx)
}

func (h *StatusHandler) handleProjectStatus(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand, _ *clicontext.ExecutionContext) error {
	ciFlags := ci.GetFlags(cmd)

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
		return err
	}

	if ciFlags.JSON {
		h.outputJSON(&ciFlags, statuses)
		return nil
	}

	h.displayStatus(base, cmd, statuses, serviceConfigs)
	return nil
}

func (h *StatusHandler) handleSharedStatus(ctx context.Context, cmd *cobra.Command, _ []string, base *base.BaseCommand, execCtx *clicontext.ExecutionContext) error {
	ciFlags := ci.GetFlags(cmd)
	showAll, _ := cmd.Flags().GetBool(docker.FlagAll)

	if !ciFlags.Quiet {
		if showAll {
			base.Output.Header(messages.InfoStatusAllProjects)
		} else {
			base.Output.Header(messages.InfoSharedContainersStatus)
		}
	}

	reg := registry.NewManager(execCtx.SharedContainers.Root)
	_, err := reg.Load()
	if err != nil {
		return err
	}

	sharedContainers, err := reg.List()
	if err != nil {
		return err
	}

	if len(sharedContainers) == 0 {
		base.Output.Info(messages.InfoNoSharedContainers)
		return nil
	}

	// Get docker client to check container states
	dockerClient, err := docker.NewClient(logger.GetLogger())
	if err != nil {
		return err
	}

	// Build display statuses with container state
	statuses := h.buildSharedStatuses(ctx, sharedContainers, dockerClient)

	if ciFlags.JSON {
		ci.OutputResult(ciFlags, display.SharedStatusResponse{
			SharedContainers: statuses,
			Count:            len(statuses),
		}, core.ExitSuccess)
		return nil
	}

	h.displaySharedStatus(base, statuses)
	return nil
}

func (h *StatusHandler) handleProjectSharedStatus(ctx context.Context, cmd *cobra.Command, _ []string, base *base.BaseCommand, execCtx *clicontext.ExecutionContext, projectName string) error {
	ciFlags := ci.GetFlags(cmd)

	if !ciFlags.Quiet {
		base.Output.Header(fmt.Sprintf(messages.InfoSharedContainersForProject, projectName))
	}

	reg := registry.NewManager(execCtx.SharedContainers.Root)
	if _, err := reg.Load(); err != nil {
		return err
	}

	containers, err := reg.List()
	if err != nil {
		return err
	}

	// Filter containers used by this project
	projectContainers := make([]*registry.ContainerInfo, 0, len(containers))
	for _, container := range containers {
		if h.containsProject(container.Projects, projectName) {
			projectContainers = append(projectContainers, container)
		}
	}

	if len(projectContainers) == 0 {
		base.Output.Info(fmt.Sprintf(messages.InfoNoSharedContainersForProject, projectName))
		return nil
	}

	dockerClient, err := docker.NewClient(h.logger)
	if err != nil {
		return err
	}

	statuses := h.buildSharedStatuses(ctx, projectContainers, dockerClient)

	if ciFlags.JSON {
		ci.OutputResult(ciFlags, display.ProjectSharedStatusResponse{
			Project:          projectName,
			SharedContainers: statuses,
			Count:            len(statuses),
		}, core.ExitSuccess)
		return nil
	}

	h.displaySharedStatus(base, statuses)
	return nil
}

func (h *StatusHandler) resolveServices(args []string, setup *common.CoreSetup, ciFlags *ci.Flags) ([]types.ServiceConfig, error) {
	serviceConfigs, err := common.ResolveServiceConfigs(args, setup)
	if err != nil {
		return nil, ci.FormatError(*ciFlags, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, messages.ErrorsStackResolveFailed, err))
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
	ci.OutputResult(*ciFlags, display.ServiceStatusResponse{
		Services: convertToInterfaceSlice(statuses),
		Count:    len(statuses),
	}, core.ExitSuccess)
}

func convertToInterfaceSlice(statuses []docker.ContainerStatus) []any {
	result := make([]any, len(statuses))
	for i, s := range statuses {
		result[i] = s
	}
	return result
}

func (h *StatusHandler) displayStatus(base *base.BaseCommand, cmd *cobra.Command, statuses []docker.ContainerStatus, serviceConfigs []types.ServiceConfig) {
	if len(statuses) == 0 {
		return
	}

	serviceToContainer := h.buildServiceContainerMap(serviceConfigs)
	serviceStatuses := convertToDisplayStatuses(statuses, serviceConfigs, serviceToContainer)
	verbose := base.GetVerbose(cmd)

	formatter := display.NewStatusFormatter(os.Stdout)
	_ = formatter.FormatTable(serviceStatuses, display.Options{
		Compact: !verbose,
		Verbose: verbose,
	})
}

func (h *StatusHandler) buildServiceContainerMap(serviceConfigs []types.ServiceConfig) map[string]string {
	serviceToContainer := make(map[string]string)
	for _, config := range serviceConfigs {
		serviceToContainer[config.Name] = getContainerName(config)
	}
	return serviceToContainer
}

// ValidateArgs validates the command arguments
func (h *StatusHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *StatusHandler) GetRequiredFlags() []string {
	return []string{}
}

// getContainerName returns the actual container name for a service
func getContainerName(config types.ServiceConfig) string {
	if config.Hidden {
		return config.Name
	}

	if len(config.Service.Dependencies.Required) > 0 {
		return config.Service.Dependencies.Required[0]
	}

	return config.Name
}

// convertToDisplayStatuses creates display service statuses with health inheritance
func convertToDisplayStatuses(containerStatuses []docker.ContainerStatus, serviceConfigs []types.ServiceConfig, serviceToContainer map[string]string) []display.ServiceStatus {
	containerMap := buildContainerMap(containerStatuses)
	return buildDisplayStatuses(serviceConfigs, serviceToContainer, containerMap)
}

func buildContainerMap(containerStatuses []docker.ContainerStatus) map[string]docker.ContainerStatus {
	containerMap := make(map[string]docker.ContainerStatus, len(containerStatuses))
	for _, status := range containerStatuses {
		containerMap[status.Name] = status
	}
	return containerMap
}

func buildDisplayStatuses(serviceConfigs []types.ServiceConfig, serviceToContainer map[string]string, containerMap map[string]docker.ContainerStatus) []display.ServiceStatus {
	result := make([]display.ServiceStatus, 0, len(serviceConfigs))
	for _, config := range serviceConfigs {
		if shouldSkipService(config) {
			continue
		}
		result = append(result, createServiceStatus(config, serviceToContainer, containerMap))
	}
	return result
}

func shouldSkipService(config types.ServiceConfig) bool {
	return config.Container.Restart == types.RestartPolicyNo || config.Hidden
}

func createServiceStatus(config types.ServiceConfig, serviceToContainer map[string]string, containerMap map[string]docker.ContainerStatus) display.ServiceStatus {
	provider := serviceToContainer[config.Name]
	providerName := getProviderName(config.Name, provider)

	if containerStatus, exists := containerMap[provider]; exists {
		return buildFoundStatus(config.Name, providerName, containerStatus)
	}

	return buildNotFoundStatus(config.Name, providerName)
}

func getProviderName(serviceName, provider string) string {
	if provider == serviceName {
		return ""
	}
	return provider
}

func buildFoundStatus(name, provider string, containerStatus docker.ContainerStatus) display.ServiceStatus {
	uptime := time.Duration(0)
	if !containerStatus.StartedAt.IsZero() {
		uptime = time.Since(containerStatus.StartedAt)
	}

	return display.ServiceStatus{
		Name:      name,
		Provider:  provider,
		State:     containerStatus.State,
		Health:    containerStatus.Health,
		Ports:     containerStatus.Ports,
		CreatedAt: containerStatus.CreatedAt,
		UpdatedAt: containerStatus.StartedAt,
		Uptime:    uptime,
	}
}

func buildNotFoundStatus(name, provider string) display.ServiceStatus {
	return display.ServiceStatus{
		Name:     name,
		Provider: provider,
		State:    messages.InfoStateNotFound,
		Health:   messages.InfoHealthUnknown,
	}
}

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

func (h *StatusHandler) buildSharedStatuses(ctx context.Context, containers []*registry.ContainerInfo, dockerClient *docker.Client) []display.SharedContainerStatus {
	statuses := make([]display.SharedContainerStatus, 0, len(containers))

	for _, container := range containers {
		state := h.getContainerState(ctx, container.Name, dockerClient)
		statuses = append(statuses, display.SharedContainerStatus{
			Name:      container.Name,
			Service:   container.Service,
			State:     state,
			Projects:  container.Projects,
			CreatedAt: container.CreatedAt,
			UpdatedAt: container.UpdatedAt,
		})
	}

	return statuses
}

func (h *StatusHandler) getContainerState(ctx context.Context, containerName string, dockerClient *docker.Client) string {
	inspectResp, err := dockerClient.GetDockerClient().ContainerInspect(ctx, containerName)
	if err != nil {
		return messages.InfoStateNotFound
	}
	return inspectResp.State.Status
}

func (h *StatusHandler) displaySharedStatus(base *base.BaseCommand, statuses []display.SharedContainerStatus) {
	if len(statuses) == 0 {
		base.Output.Info(messages.InfoNoSharedContainers)
		return
	}

	tw := table.NewWriter()
	tw.SetOutputMirror(os.Stdout)
	tw.SetStyle(table.StyleLight)

	tw.AppendHeader(table.Row{
		messages.InfoHeaderContainer,
		messages.InfoHeaderService,
		messages.InfoHeaderState,
		messages.InfoHeaderUsedBy,
	})

	for _, status := range statuses {
		tw.AppendRow(table.Row{
			status.Name,
			status.Service,
			status.State,
			formatProjects(status.Projects),
		})
	}

	tw.Render()
}

func formatProjects(projects []string) string {
	if len(projects) == 0 {
		return messages.InfoProjectsNone
	}
	if len(projects) == 1 {
		return projects[0]
	}
	const maxShow = 3
	if len(projects) > maxShow {
		const projectsToShow = 2
		return fmt.Sprintf("%s, %s, +%d more", projects[0], projects[1], len(projects)-projectsToShow)
	}

	var b strings.Builder
	for i, p := range projects {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(p)
	}
	return b.String()
}

func (h *StatusHandler) containsProject(projects []string, projectName string) bool {
	return slices.Contains(projects, projectName)
}
