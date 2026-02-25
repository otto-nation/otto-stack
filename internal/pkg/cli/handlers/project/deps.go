package project

import (
	"context"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/validation"
)

const depsDash = "—"

// DepsHandler handles the deps command
type DepsHandler struct{}

// NewDepsHandler creates a new deps handler
func NewDepsHandler() *DepsHandler {
	return &DepsHandler{}
}

// Handle executes the deps command
func (h *DepsHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	base.Output.Header(messages.DependenciesHeader)

	serviceConfigs, err := h.loadServices(args)
	if err != nil {
		return err
	}

	if len(serviceConfigs) == 0 {
		base.Output.Warning(messages.DependenciesNoEnabledServices)
		return nil
	}

	headers, rows := h.buildTable(serviceConfigs)
	display.RenderTable(base.Output.Writer(), headers, rows)

	return nil
}

func (h *DepsHandler) loadServices(args []string) ([]types.ServiceConfig, error) {
	if err := validation.CheckInitialization(); err != nil {
		return nil, err
	}

	configPath := filepath.Join(core.OttoStackDir, core.ConfigFileName)
	cfg, err := common.LoadProjectConfig(configPath)
	if err != nil {
		return nil, pkgerrors.NewConfigError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentConfig, messages.ErrorsConfigLoadFailed, err)
	}

	if len(args) > 0 {
		return services.ResolveUpServices(args, cfg)
	}
	return services.ResolveUpServices(cfg.Stack.Enabled, cfg)
}

// buildTable constructs headers and rows, collapsing columns that have no data.
func (h *DepsHandler) buildTable(configs []types.ServiceConfig) ([]string, [][]string) {
	hasRequired := h.anyHasField(configs, func(c types.ServiceConfig) bool { return len(c.Service.Dependencies.Required) > 0 })
	hasSoft := h.anyHasField(configs, func(c types.ServiceConfig) bool { return len(c.Service.Dependencies.Soft) > 0 })
	hasConflicts := h.anyHasField(configs, func(c types.ServiceConfig) bool { return len(c.Service.Dependencies.Conflicts) > 0 })
	hasProvides := h.anyHasField(configs, func(c types.ServiceConfig) bool { return len(c.Service.Dependencies.Provides) > 0 })

	headers := []string{display.HeaderService}
	if hasRequired {
		headers = append(headers, display.HeaderRequired)
	}
	if hasSoft {
		headers = append(headers, display.HeaderSoft)
	}
	if hasConflicts {
		headers = append(headers, display.HeaderConflicts)
	}
	if hasProvides {
		headers = append(headers, display.HeaderProvides)
	}

	rows := make([][]string, len(configs))
	for i, cfg := range configs {
		row := []string{cfg.Name}
		if hasRequired {
			row = append(row, joinOrDash(cfg.Service.Dependencies.Required))
		}
		if hasSoft {
			row = append(row, joinOrDash(cfg.Service.Dependencies.Soft))
		}
		if hasConflicts {
			row = append(row, joinOrDash(cfg.Service.Dependencies.Conflicts))
		}
		if hasProvides {
			row = append(row, joinOrDash(cfg.Service.Dependencies.Provides))
		}
		rows[i] = row
	}

	return headers, rows
}

func (h *DepsHandler) anyHasField(configs []types.ServiceConfig, check func(types.ServiceConfig) bool) bool {
	return slices.ContainsFunc(configs, check)
}

func joinOrDash(values []string) string {
	if len(values) == 0 {
		return depsDash
	}
	return strings.Join(values, ", ")
}

// ValidateArgs validates the command arguments
func (h *DepsHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *DepsHandler) GetRequiredFlags() []string {
	return []string{}
}
