package stack

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
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

	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		ci.HandleError(ciFlags, err)
		return nil
	}
	defer cleanup()

	// Determine services to check
	serviceNames := h.resolveServiceNames(args, setup.Config.Stack.Enabled)

	// Apply same service resolution as up command
	manager, err := GetServicesManager()
	if err != nil {
		ci.HandleError(ciFlags, fmt.Errorf("failed to create service manager: %w", err))
		return nil
	}

	// Validate services exist
	if err := manager.ValidateServices(serviceNames); err != nil {
		ci.HandleError(ciFlags, fmt.Errorf(core.MsgStack_failed_resolve_services, err))
		return nil
	}

	// Resolve services to actual container services
	resolvedServices, err := manager.ResolveServices(serviceNames)
	if err != nil {
		ci.HandleError(ciFlags, fmt.Errorf(core.MsgStack_failed_resolve_services, err))
		return nil
	}

	// Filter out init containers (restart: "no") from status display
	filteredServices := filterInitContainers(manager, resolvedServices)

	// Get service status
	statuses, err := setup.DockerClient.GetDockerServiceStatus(ctx, setup.Config.Project.Name, filteredServices)
	if err != nil {
		ci.HandleError(ciFlags, fmt.Errorf(core.MsgStack_failed_get_service_status, err))
		return nil
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

	fmt.Printf("%-20s %-12s %s\n", ui.StatusHeaderService, ui.StatusHeaderState, ui.StatusHeaderHealth)
	fmt.Println(ui.StatusSeparator)
	for _, status := range statuses {
		fmt.Printf("%-20s %-12s %s\n", status.Name, status.State, status.Health)
	}

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

// filterInitContainers removes init containers (restart: "no") from status display
func filterInitContainers(manager *services.Manager, serviceNames []string) []string {
	var filtered []string
	for _, serviceName := range serviceNames {
		if service, err := manager.GetService(serviceName); err == nil && service.Container.Restart != services.RestartPolicyNo {
			filtered = append(filtered, serviceName)
		}
	}
	return filtered
}

// resolveServiceNames determines which services to check status for
func (h *StatusHandler) resolveServiceNames(args, enabledServices []string) []string {
	if len(args) > 0 {
		return args
	}
	return enabledServices
}
