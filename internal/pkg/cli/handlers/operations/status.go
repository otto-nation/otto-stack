package operations

import (
	"context"
	"fmt"
	"log/slog"
	"os"

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
	// Check initialization first

	// Get CI-friendly flags
	ciFlags := ci.GetFlags(cmd)

	if !ciFlags.Quiet {
		base.Output.Header(core.MsgStatus)
	}

	setup, cleanup, err := common.SetupCoreCommand(ctx, base)
	if err != nil {
		return ci.FormatError(ciFlags, err)
	}
	defer cleanup()

	// Resolve services to ServiceConfigs
	serviceConfigs, err := ResolveServiceConfigs(args, setup)
	if err != nil {
		return ci.FormatError(ciFlags, fmt.Errorf(core.MsgStack_failed_resolve_services, err))
	}

	// Filter out init containers (restart: "no") from status display
	filteredServices := filterInitContainers(serviceConfigs)

	// Get service status using StackService
	stackService, err := common.NewServiceManager(false)
	if err != nil {
		return ci.FormatError(ciFlags, pkgerrors.NewServiceError("stack", "create service", err))
	}

	statuses, err := stackService.DockerClient.GetServiceStatus(ctx, setup.Config.Project.Name, filteredServices)
	if err != nil {
		return ci.FormatError(ciFlags, fmt.Errorf(core.MsgStack_failed_get_service_status, err))
	}

	// Handle CI-friendly output
	if ciFlags.JSON {
		ci.OutputResult(ciFlags, map[string]any{
			"services": statuses,
			"count":    len(statuses),
		}, core.ExitSuccess)
		return nil
	}

	// Display user-friendly status
	if len(statuses) == 0 {
		// Restart operation
		return nil
	}

	// Map services to their container names
	serviceToContainer := make(map[string]string)
	for _, config := range serviceConfigs {
		containerName := getContainerName(config)
		serviceToContainer[config.Name] = containerName
	}

	// Convert statuses with inheritance and display
	serviceStatuses := convertToDisplayStatuses(statuses, serviceConfigs, serviceToContainer)

	formatter := display.NewStatusFormatter(os.Stdout)
	_ = formatter.FormatTable(serviceStatuses, display.Options{Compact: true})

	return nil
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
	// If service is hidden, it's the actual container
	if config.Hidden {
		return config.Name
	}

	// Check if service has dependencies that provide the actual container
	if len(config.Service.Dependencies.Required) > 0 {
		// Return the first required dependency as the container name
		return config.Service.Dependencies.Required[0]
	}

	// If no dependencies, service name is the container name
	return config.Name
}

// convertToDisplayStatuses creates display service statuses with health inheritance
func convertToDisplayStatuses(containerStatuses []docker.ContainerStatus, serviceConfigs []types.ServiceConfig, serviceToContainer map[string]string) []display.ServiceStatus {
	containerMap := make(map[string]docker.ContainerStatus)
	for _, status := range containerStatuses {
		containerMap[status.Name] = status
	}

	var result []display.ServiceStatus
	for _, config := range serviceConfigs {
		if config.Container.Restart == types.RestartPolicyNo || config.Hidden {
			continue // Skip init containers and hidden services
		}

		provider := serviceToContainer[config.Name]
		providerName := ""
		if provider != config.Name {
			providerName = provider
		}

		if containerStatus, exists := containerMap[provider]; exists {
			result = append(result, display.ServiceStatus{
				Name:     config.Name,
				Provider: providerName,
				State:    containerStatus.State,
				Health:   containerStatus.Health,
			})
		} else {
			result = append(result, display.ServiceStatus{
				Name:     config.Name,
				Provider: providerName,
				State:    "not found",
				Health:   "unknown",
			})
		}
	}

	return result
}

// filterInitContainers removes init containers (restart: "no") from status display
func filterInitContainers(serviceConfigs []types.ServiceConfig) []string {
	var filtered []string
	for _, config := range serviceConfigs {
		if config.Container.Restart != types.RestartPolicyNo {
			filtered = append(filtered, config.Name)
		}
	}
	return filtered
}
