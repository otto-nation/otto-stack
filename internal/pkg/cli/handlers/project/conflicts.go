package project

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
	"slices"
	"sort"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/validation"
)

// ConflictsHandler handles the conflicts command
type ConflictsHandler struct{}

// NewConflictsHandler creates a new conflicts handler
func NewConflictsHandler() *ConflictsHandler {
	return &ConflictsHandler{}
}

type semanticConflict struct {
	serviceA   string
	serviceB   string
	capability string // non-empty when conflict is a provides overlap
}

type portConflict struct {
	service string
	port    int
}

// Handle executes the conflicts command
func (h *ConflictsHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	base.Output.Header(messages.ConflictsHeader)

	flags, err := core.ParseConflictsFlags(cmd)
	if err != nil {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldFlags, messages.ValidationFailedParseFlags, err)
	}

	serviceConfigs, err := h.loadEnabledServices()
	if err != nil {
		return err
	}

	hasConflicts := false

	semanticConflicts := h.detectSemanticConflicts(serviceConfigs)
	if len(semanticConflicts) > 0 {
		hasConflicts = true
		base.Output.Warning(messages.ConflictsFound, len(semanticConflicts))
		for _, c := range semanticConflicts {
			if c.capability != "" {
				base.Output.Info(messages.ConflictsProvidesOverlap, c.serviceA, c.serviceB, c.capability)
			} else {
				base.Output.Info(messages.ConflictsExplicitIncompatible, c.serviceA, c.serviceB)
			}
		}
	} else {
		base.Output.Success(messages.SuccessNoConflicts)
	}

	if flags.CheckPorts {
		portConflicts := h.detectPortConflicts(serviceConfigs)
		if len(portConflicts) > 0 {
			hasConflicts = true
			base.Output.Warning(messages.ConflictsPortFound, len(portConflicts))
			for _, p := range portConflicts {
				base.Output.Info(messages.ConflictsPortInUse, p.port, p.service)
			}
		} else {
			base.Output.Success(messages.ConflictsNoPortConflicts)
		}
	}

	if hasConflicts {
		return pkgerrors.ErrSilentExit
	}
	return nil
}

func (h *ConflictsHandler) loadEnabledServices() ([]types.ServiceConfig, error) {
	if err := validation.CheckInitialization(); err != nil {
		return nil, err
	}

	configPath := filepath.Join(core.OttoStackDir, core.ConfigFileName)
	cfg, err := common.LoadProjectConfig(configPath)
	if err != nil {
		return nil, pkgerrors.NewConfigError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentConfig, messages.ErrorsConfigLoadFailed, err)
	}

	return services.ResolveUpServices(cfg.Stack.Enabled, cfg)
}

func (h *ConflictsHandler) detectSemanticConflicts(configs []types.ServiceConfig) []semanticConflict {
	var conflicts []semanticConflict
	// seenPairs deduplicates (serviceA, serviceB) pairs regardless of order.
	seenPairs := make(map[string]bool)

	for i := range len(configs) {
		for j := i + 1; j < len(configs); j++ {
			a, b := configs[i], configs[j]
			key := pairKey(a.Name, b.Name)
			if seenPairs[key] {
				continue
			}
			seenPairs[key] = true

			if h.hasExplicitConflict(a, b) {
				conflicts = append(conflicts, semanticConflict{
					serviceA: a.Name,
					serviceB: b.Name,
				})
				continue
			}

			// Only report the first overlapping capability to avoid redundant rows.
			if overlaps := h.providesOverlap(a, b); len(overlaps) > 0 {
				conflicts = append(conflicts, semanticConflict{
					serviceA:   a.Name,
					serviceB:   b.Name,
					capability: overlaps[0],
				})
			}
		}
	}

	return conflicts
}

func (h *ConflictsHandler) hasExplicitConflict(a, b types.ServiceConfig) bool {
	return slices.Contains(a.Service.Dependencies.Conflicts, b.Name)
}

func (h *ConflictsHandler) providesOverlap(a, b types.ServiceConfig) []string {
	providesB := make(map[string]bool, len(b.Service.Dependencies.Provides))
	for _, p := range b.Service.Dependencies.Provides {
		providesB[p] = true
	}

	var overlap []string
	for _, p := range a.Service.Dependencies.Provides {
		if providesB[p] {
			overlap = append(overlap, p)
		}
	}
	sort.Strings(overlap)
	return overlap
}

func (h *ConflictsHandler) detectPortConflicts(configs []types.ServiceConfig) []portConflict {
	var conflicts []portConflict
	for _, cfg := range configs {
		for _, portMapping := range cfg.Container.Ports {
			port := h.parsePort(portMapping.External)
			if port > 0 && h.isPortInUse(port) {
				conflicts = append(conflicts, portConflict{
					service: cfg.Name,
					port:    port,
				})
			}
		}
	}
	return conflicts
}

func (h *ConflictsHandler) parsePort(portStr string) int {
	var port int
	_, _ = fmt.Sscanf(portStr, "%d", &port)
	return port
}

// isPortInUse returns true if the TCP port is already bound on the host.
func (h *ConflictsHandler) isPortInUse(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return true
	}
	_ = listener.Close()
	return false
}

// pairKey returns a canonical, order-independent key for a service pair.
func pairKey(a, b string) string {
	if a > b {
		a, b = b, a
	}
	return a + "/" + b
}

// ValidateArgs validates the command arguments
func (h *ConflictsHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *ConflictsHandler) GetRequiredFlags() []string {
	return []string{}
}
