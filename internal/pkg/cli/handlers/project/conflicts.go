package project

import (
	"context"
	"fmt"
	"net"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
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
	servicePorts := h.getServicePorts()
	return h.checkAllServicePorts(servicePorts)
}

func (h *ConflictsHandler) checkAllServicePorts(servicePorts map[string][]int) []string {
	var conflicts []string
	for service, ports := range servicePorts {
		conflicts = append(conflicts, h.checkServicePorts(service, ports)...)
	}
	return conflicts
}

func (h *ConflictsHandler) checkServicePorts(service string, ports []int) []string {
	var conflicts []string
	for _, port := range ports {
		if conflict := h.checkPortConflict(service, port); conflict != "" {
			conflicts = append(conflicts, conflict)
		}
	}
	return conflicts
}

func (h *ConflictsHandler) checkPortConflict(service string, port int) string {
	if h.isPortInUse(port) {
		return fmt.Sprintf("Port %d (needed by %s) is already in use", port, service)
	}
	return ""
}

// getServicePorts extracts ports from service configurations
func (h *ConflictsHandler) getServicePorts() map[string][]int {
	manager, err := services.New()
	if err != nil {
		return make(map[string][]int)
	}

	allServices := manager.GetAllServices()
	servicePorts := make(map[string][]int)

	for serviceName, service := range allServices {
		ports := h.extractPortsFromService(&service)
		if len(ports) > 0 {
			servicePorts[serviceName] = ports
		}
	}

	return servicePorts
}

func (h *ConflictsHandler) extractPortsFromService(service *types.ServiceConfig) []int {
	var ports []int

	for _, portMapping := range service.Container.Ports {
		if port := h.parsePort(portMapping.External); port > 0 {
			ports = append(ports, port)
		}
	}

	return ports
}

func (h *ConflictsHandler) parsePort(portStr string) int {
	var port int
	_, _ = fmt.Sscanf(portStr, "%d", &port)
	return port
}

// isPortInUse checks if a port is currently in use
func (h *ConflictsHandler) isPortInUse(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return true
	}
	_ = listener.Close()
	return false
}

// ValidateArgs validates the command arguments
func (h *ConflictsHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *ConflictsHandler) GetRequiredFlags() []string {
	return []string{}
}
