package project

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
)

// DepsHandler handles the deps command
type DepsHandler struct{}

// NewDepsHandler creates a new deps handler
func NewDepsHandler() *DepsHandler {
	return &DepsHandler{}
}

// Handle executes the deps command
func (h *DepsHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	base.Output.Header("%s Service Dependencies", ui.IconInfo)

	if len(args) == 0 {
		base.Output.Warning(messages.InfoNoServicesSpecified)
		return nil
	}

	allServices, err := h.loadServices()
	if err != nil {
		return err
	}

	rows := h.buildDependencyRows(args, allServices)

	headers := []string{display.HeaderService, display.HeaderDependencies}
	display.RenderTable(base.Output.Writer(), headers, rows)

	base.Output.Success(messages.SuccessDependenciesDisplayed)
	return nil
}

func (h *DepsHandler) loadServices() (map[string]servicetypes.ServiceConfig, error) {
	manager, err := services.New()
	if err != nil {
		return nil, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentServices, messages.ErrorsServiceManagerCreateFailed, err)
	}
	return manager.GetAllServices(), nil
}

func (h *DepsHandler) buildDependencyRows(serviceNames []string, allServices map[string]servicetypes.ServiceConfig) [][]string {
	rows := make([][]string, 0, len(serviceNames))

	for _, name := range serviceNames {
		service, exists := allServices[name]
		if !exists {
			continue
		}

		depsStr := h.formatDependencies(service.Service.Dependencies.Required)
		rows = append(rows, []string{name, depsStr})
	}

	return rows
}

func (h *DepsHandler) formatDependencies(deps []string) string {
	if len(deps) == 0 {
		return messages.InfoProjectsNone
	}

	var result strings.Builder
	for i, dep := range deps {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(dep)
	}
	return result.String()
}

// ValidateArgs validates the command arguments
func (h *DepsHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *DepsHandler) GetRequiredFlags() []string {
	return []string{}
}
