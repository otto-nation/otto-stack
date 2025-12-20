package project

import (
	"context"
	"fmt"
	"net"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/spf13/cobra"
)

// ConflictsHandler handles the conflicts command
type ConflictsHandler struct{}

// NewConflictsHandler creates a new conflicts handler
func NewConflictsHandler() *ConflictsHandler {
	return &ConflictsHandler{}
}

// Handle executes the conflicts command
func (h *ConflictsHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	base.Output.Header("%s", core.MsgConflicts_header)

	// Check for common port conflicts
	conflicts := h.detectPortConflicts()

	if len(conflicts) == 0 {
		base.Output.Success("No service conflicts detected")
		return nil
	}

	base.Output.Warning("Found %d potential conflicts:", len(conflicts))
	for _, conflict := range conflicts {
		base.Output.Info("  - %s", conflict)
	}

	return nil
}

// detectPortConflicts checks for actual port conflicts from services
func (h *ConflictsHandler) detectPortConflicts() []string {
	var conflicts []string

	// Get ports from actual services configuration
	servicePorts := h.getServicePorts()

	for service, ports := range servicePorts {
		for _, port := range ports {
			if h.isPortInUse(port) {
				conflicts = append(conflicts, fmt.Sprintf("Port %d (needed by %s) is already in use", port, service))
			}
		}
	}

	return conflicts
}

// getServicePorts extracts ports from service configurations
func (h *ConflictsHandler) getServicePorts() map[string][]int {
	// Simple hardcoded mapping - could be made dynamic later
	return map[string][]int{
		services.ServicePostgres: {5432},
		services.ServiceMysql:    {3306},
		services.ServiceRedis:    {6379},
	}
}

// isPortInUse checks if a port is currently in use
func (h *ConflictsHandler) isPortInUse(port int) bool {
	// Try to bind to the port to check if it's available
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return true // Port is in use
	}
	_ = listener.Close()
	return false // Port is available
}

// ValidateArgs validates the command arguments
func (h *ConflictsHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *ConflictsHandler) GetRequiredFlags() []string {
	return []string{}
}
