package stack

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

// PortConflict represents a port conflict
type PortConflict struct {
	Port        string
	ServiceName string
	ProcessName string
	PID         string
}

// checkPortConflicts checks for port conflicts before starting services
func checkPortConflicts(serviceNames []string, base *base.BaseCommand) error {
	conflicts := collectPortConflicts(serviceNames)
	if len(conflicts) > 0 {
		reportConflicts(conflicts, base)
		return fmt.Errorf("port conflicts detected")
	}
	return nil
}

// collectPortConflicts collects all port conflicts for given services
func collectPortConflicts(serviceNames []string) []PortConflict {
	var conflicts []PortConflict
	projectName := getProjectNameSafe()

	for _, serviceName := range serviceNames {
		serviceConflicts := checkServicePorts(serviceName, projectName)
		conflicts = append(conflicts, serviceConflicts...)
	}

	return conflicts
}

// checkServicePorts checks ports for a specific service
func checkServicePorts(serviceName, projectName string) []PortConflict {
	var conflicts []PortConflict

	// Get ports from service configuration
	ports := getServicePorts(serviceName)
	if len(ports) == 0 {
		return conflicts
	}

	// Check each port mapping
	for _, portMapping := range ports {
		// Extract host port from mapping (format: "host:container" or "port")
		parts := strings.Split(portMapping, ":")
		hostPort := parts[0]

		if conflict := detectPortConflict(hostPort, projectName, serviceName); conflict != nil {
			conflicts = append(conflicts, *conflict)
		}
	}

	return conflicts
}

// getServicePorts returns the port mappings for a service based on its configuration
func getServicePorts(serviceName string) []string {
	// Get the service configuration
	manager, err := services.New()
	if err != nil {
		return []string{}
	}

	serviceConfig, err := manager.GetService(serviceName)
	if err != nil {
		return []string{}
	}

	var ports []string
	for _, port := range serviceConfig.Container.Ports {
		// Format: "external:internal" or just "port" if they're the same
		if port.External != port.Internal {
			ports = append(ports, port.External+":"+port.Internal)
		} else {
			ports = append(ports, port.External)
		}
	}

	return ports
}

// detectPortConflict detects if a port is in use by another process
func detectPortConflict(port, projectName, serviceName string) *PortConflict {
	pid := getPortOwnerPID(port)
	if pid == "" {
		return nil // Port is free
	}

	// Check if it's our own container
	if isOwnContainer(projectName, serviceName) {
		return nil // Not a conflict if it's our own container
	}

	processName := getProcessName(pid)
	return &PortConflict{
		Port:        port,
		ServiceName: serviceName,
		ProcessName: processName,
		PID:         pid,
	}
}

// getPortOwnerPID gets the PID of the process using the port
func getPortOwnerPID(port string) string {
	cmd := exec.Command("lsof", "-ti", fmt.Sprintf(":%s", port))
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	pid := strings.TrimSpace(string(output))
	if pid == "" {
		return ""
	}

	return pid
}

// isOwnContainer checks if the port is used by our own container
func isOwnContainer(projectName, serviceName string) bool {
	expectedName := fmt.Sprintf("%s-%s", projectName, serviceName)

	// Check if container exists and is running using StackService
	stackService, err := NewStackService(false)
	if err != nil {
		return false
	}

	containers, err := stackService.DockerClient.ListContainers(context.Background(), projectName)
	if err != nil {
		return false
	}

	// Check if any container matches our expected name
	for _, container := range containers {
		if strings.Contains(container.Name, expectedName) && container.State == "running" {
			return true
		}
	}
	return false
}

// getProcessName gets the process name for a PID
func getProcessName(pid string) string {
	cmd := exec.Command("ps", "-p", pid, "-o", "comm=")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

// getProjectNameSafe safely gets the project name
func getProjectNameSafe() string {
	// This should get the project name from config
	// For now, return a default
	return "otto-stack"
}

// reportConflicts reports port conflicts to the user
func reportConflicts(conflicts []PortConflict, base *base.BaseCommand) {
	base.Output.Error("Port conflicts detected:")
	for _, conflict := range conflicts {
		base.Output.Error("  Port %s (service: %s) is used by process %s (PID: %s)",
			conflict.Port, conflict.ServiceName, conflict.ProcessName, conflict.PID)
	}
	base.Output.Info("Please stop the conflicting processes or change the port mappings")
}
