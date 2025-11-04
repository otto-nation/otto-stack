package project

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/spf13/cobra"
)

// DepsHandler handles the deps command
type DepsHandler struct{}

// NewDepsHandler creates a new deps handler
func NewDepsHandler() *DepsHandler {
	return &DepsHandler{}
}

// Handle executes the deps command
func (h *DepsHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	base.Output.Header("%s", constants.Messages[constants.MsgDependencies_header])

	// Get output format
	format, _ := cmd.Flags().GetString(constants.FlagFormat)

	// Load service dependencies
	serviceUtils := utils.NewServiceUtils()
	dependencies, err := serviceUtils.LoadAllServiceDependencies()
	if err != nil {
		return fmt.Errorf("failed to load dependencies: %w", err)
	}

	if len(dependencies) == 0 {
		base.Output.Info("%s", constants.Messages[constants.MsgDependencies_none_found])
		return nil
	}

	// Create display data
	var displayData []map[string]any
	for serviceName, deps := range dependencies {
		if len(deps) == 0 {
			displayData = append(displayData, map[string]any{
				"Service":      serviceName,
				"Dependencies": "None",
			})
		} else {
			for _, dep := range deps {
				displayData = append(displayData, map[string]any{
					"Service":      serviceName,
					"Dependencies": dep,
				})
			}
		}
	}

	// Display results
	formatter := display.New(cmd.OutOrStdout(), base.Output)

	// Convert to ServiceStatus format for display
	var services []display.ServiceStatus
	for _, item := range displayData {
		services = append(services, display.ServiceStatus{
			Name:  item["Service"].(string),
			State: item["Dependencies"].(string),
		})
	}

	if err := formatter.FormatStatus(services, display.Options{Format: format}); err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	return nil
}

// ValidateArgs validates the command arguments
func (h *DepsHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *DepsHandler) GetRequiredFlags() []string {
	return []string{}
}
