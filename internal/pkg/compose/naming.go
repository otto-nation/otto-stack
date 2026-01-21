package compose

import (
	"fmt"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
)

const (
	sharedNetworkSuffix = "shared"
)

// NamingStrategy handles container and resource naming
type NamingStrategy struct {
	projectName string
	sharing     *config.SharingConfig
}

// NewNamingStrategy creates a new naming strategy
func NewNamingStrategy(projectName string, sharing *config.SharingConfig) *NamingStrategy {
	return &NamingStrategy{
		projectName: projectName,
		sharing:     sharing,
	}
}

// ContainerName returns the container name for a service
func (n *NamingStrategy) ContainerName(serviceName string) string {
	if n.isShared(serviceName) {
		return fmt.Sprintf("%s-%s", core.AppName, serviceName)
	}
	return fmt.Sprintf("%s-%s", n.projectName, serviceName)
}

// VolumeName returns the volume name for a service
func (n *NamingStrategy) VolumeName(serviceName, volumeSuffix string) string {
	if n.isShared(serviceName) {
		return fmt.Sprintf("%s-%s-%s", core.AppName, serviceName, volumeSuffix)
	}
	return fmt.Sprintf("%s-%s-%s", n.projectName, serviceName, volumeSuffix)
}

// NetworkName returns the network name for a service
func (n *NamingStrategy) NetworkName(serviceName string) string {
	if n.isShared(serviceName) {
		return fmt.Sprintf("%s-%s", core.AppName, sharedNetworkSuffix)
	}
	return n.projectName + docker.NetworkNameSuffix
}

// IsShared checks if a service should be shared
func (n *NamingStrategy) IsShared(serviceName string) bool {
	return n.isShared(serviceName)
}

// isShared determines if a service is shared based on config
func (n *NamingStrategy) isShared(serviceName string) bool {
	if n.sharing == nil {
		return false
	}

	// Check per-service override first
	if override, exists := n.sharing.Services[serviceName]; exists {
		return override
	}

	// Fall back to global setting
	return n.sharing.Enabled
}
