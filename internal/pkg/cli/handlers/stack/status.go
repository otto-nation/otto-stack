package stack

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/output"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/validation"
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
	if err := validation.CheckInitialization(); err != nil {
		return err
	}

	// Get CI-friendly flags
	ciFlags := output.GetCIFlags(cmd)

	if !ciFlags.Quiet {
		base.Output.Header(core.MsgStatus)
	}

	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		output.HandleError(ciFlags, err)
		return nil
	}
	defer cleanup()

	// Determine services to check
	serviceNames := args
	if len(serviceNames) == 0 {
		serviceNames = setup.Config.Stack.Enabled
	}

	// Apply same service resolution as up command
	manager, err := services.New()
	if err != nil {
		output.HandleError(ciFlags, fmt.Errorf("failed to create service manager: %w", err))
		return nil
	}

	// Validate services exist
	if err := manager.ValidateServices(serviceNames); err != nil {
		output.HandleError(ciFlags, fmt.Errorf(core.MsgStack_failed_resolve_services, err))
		return nil
	}

	resolvedServices := serviceNames

	// Get service status
	statuses, err := setup.DockerClient.GetDockerServiceStatus(ctx, setup.Config.Project.Name, resolvedServices)
	if err != nil {
		output.HandleError(ciFlags, fmt.Errorf(core.MsgStack_failed_get_service_status, err))
		return nil
	}

	// Handle CI-friendly output
	if ciFlags.JSON {
		output.OutputResult(ciFlags, map[string]any{
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

	fmt.Printf("%-20s %-12s %s\n", core.StatusHeaderService, core.StatusHeaderState, core.StatusHeaderHealth)
	fmt.Println(strings.Repeat(core.StatusSeparator, core.StatusSeparatorLength))
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
