package project

import (
	"context"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/spf13/cobra"
)

// DepsHandler handles the deps command
type DepsHandler struct{}

// NewDepsHandler creates a new deps handler
func NewDepsHandler() *DepsHandler {
	return &DepsHandler{}
}

// Handle executes the deps command
func (h *DepsHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	base.Output.Header("%s", core.MsgDependencies_header)

	format, _ := cmd.Flags().GetString(core.FlagFormat)

	dependencies := make(map[string][]string) // No dependencies to load since functionality was deprecated

	if len(dependencies) == 0 {
		base.Output.Info("%s", core.MsgDependencies_none_found)
		return nil
	}

	services := h.transformDependenciesToServices(dependencies)
	formatter := display.New(cmd.OutOrStdout(), base.Output)

	if err := formatter.FormatStatus(services, display.Options{Format: format}); err != nil {
		return pkgerrors.NewServiceError(ComponentFormatter, ActionFormatOutput, err)
	}

	return nil
}

// transformDependenciesToServices converts dependencies map to ServiceStatus slice
func (h *DepsHandler) transformDependenciesToServices(dependencies map[string][]string) []display.ServiceStatus {
	var services []display.ServiceStatus
	for serviceName, deps := range dependencies {
		if len(deps) == 0 {
			services = append(services, display.ServiceStatus{
				Name:  serviceName,
				State: "None",
			})
		} else {
			for _, dep := range deps {
				services = append(services, display.ServiceStatus{
					Name:  serviceName,
					State: dep,
				})
			}
		}
	}
	return services
}

// ValidateArgs validates the command arguments
func (h *DepsHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *DepsHandler) GetRequiredFlags() []string {
	return []string{}
}
