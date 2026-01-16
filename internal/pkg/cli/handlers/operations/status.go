package operations

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
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
	ciFlags := ci.GetFlags(cmd)

	if !ciFlags.Quiet {
		base.Output.Header(core.MsgLifecycle_status)
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

func (h *StatusHandler) resolveServices(args []string, setup *common.CoreSetup, ciFlags *ci.Flags) ([]types.ServiceConfig, error) {
	serviceConfigs, err := common.ResolveServiceConfigs(args, setup)
	if err != nil {
		return nil, ci.FormatError(*ciFlags, fmt.Errorf(core.MsgStack_failed_resolve_services, err))
	}
	return serviceConfigs, nil
}

func (h *StatusHandler) getServiceStatuses(ctx context.Context, projectName string, serviceConfigs []types.ServiceConfig, ciFlags *ci.Flags) ([]docker.ContainerStatus, error) {
	filteredServices := filterInitContainers(serviceConfigs)

	stackService, err := common.NewServiceManager(false)
	if err != nil {
		return nil, ci.FormatError(*ciFlags, pkgerrors.NewServiceError("stack", "create service", err))
	}

	statuses, err := stackService.DockerClient.GetServiceStatus(ctx, projectName, filteredServices)
	if err != nil {
		return nil, ci.FormatError(*ciFlags, fmt.Errorf(core.MsgStack_failed_get_service_status, err))
	}

	return statuses, nil
}

func (h *StatusHandler) outputJSON(ciFlags *ci.Flags, statuses []docker.ContainerStatus) {
	ci.OutputResult(*ciFlags, map[string]any{
		"services": statuses,
		"count":    len(statuses),
	}, core.ExitSuccess)
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
		State:    "not found",
		Health:   "unknown",
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
