package stack

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
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
func (h *StatusHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	// Check initialization first
	if err := utils.CheckInitialization(); err != nil {
		return err
	}

	// Get CI-friendly flags
	ciFlags := utils.GetCIFlags(cmd)

	if !ciFlags.Quiet {
		base.Output.Header(constants.MsgStatus)
	}

	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		utils.HandleError(ciFlags, err)
		return nil
	}
	defer cleanup()

	// Determine services to check
	serviceNames := args
	if len(serviceNames) == 0 {
		serviceNames = setup.Config.Stack.Enabled
	}

	// Apply same service resolution as up command
	serviceUtils := utils.NewServiceUtils()
	resolvedServices, err := serviceUtils.ResolveServices(serviceNames)
	if err != nil {
		utils.HandleError(ciFlags, fmt.Errorf(constants.MsgStack_failed_resolve_services, err))
		return nil
	}

	// Get service status
	statuses, err := setup.DockerClient.GetServiceStatus(ctx, setup.Config.Project.Name, resolvedServices)
	if err != nil {
		utils.HandleError(ciFlags, fmt.Errorf(constants.MsgStack_failed_get_service_status, err))
		return nil
	}

	// Handle CI-friendly output
	if ciFlags.JSON {
		utils.OutputResult(ciFlags, map[string]any{
			"services": statuses,
			"count":    len(statuses),
		}, constants.ExitSuccess)
		return nil
	}

	// Display user-friendly status
	if len(statuses) == 0 {
		// Restart operation
		return nil
	}

	h.logger.Info("Displaying service status", "service_count", len(statuses))
	fmt.Printf("%-20s %-12s %s\n", constants.StatusHeaderService, constants.StatusHeaderState, constants.StatusHeaderHealth)
	fmt.Println(strings.Repeat(constants.StatusSeparator, constants.StatusSeparatorLength))
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
