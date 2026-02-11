package operations

import (
	"time"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// StatusConverter converts container statuses to display statuses
type StatusConverter struct{}

// NewStatusConverter creates a new status converter
func NewStatusConverter() *StatusConverter {
	return &StatusConverter{}
}

// ConvertToDisplayStatuses converts container statuses to display format
func (sc *StatusConverter) ConvertToDisplayStatuses(containerStatuses []docker.ContainerStatus, serviceConfigs []types.ServiceConfig, serviceToContainer map[string]string) []display.ServiceStatus {
	containerMap := sc.buildContainerMap(containerStatuses)
	return sc.buildDisplayStatuses(serviceConfigs, serviceToContainer, containerMap)
}

func (sc *StatusConverter) buildContainerMap(containerStatuses []docker.ContainerStatus) map[string]docker.ContainerStatus {
	containerMap := make(map[string]docker.ContainerStatus, len(containerStatuses))
	for _, status := range containerStatuses {
		containerMap[status.Name] = status
	}
	return containerMap
}

func (sc *StatusConverter) buildDisplayStatuses(serviceConfigs []types.ServiceConfig, serviceToContainer map[string]string, containerMap map[string]docker.ContainerStatus) []display.ServiceStatus {
	result := make([]display.ServiceStatus, 0, len(serviceConfigs))
	for _, config := range serviceConfigs {
		if sc.shouldSkipService(config) {
			continue
		}
		result = append(result, sc.createServiceStatus(config, serviceToContainer, containerMap))
	}
	return result
}

func (sc *StatusConverter) shouldSkipService(config types.ServiceConfig) bool {
	return config.Container.Restart == types.RestartPolicyNo || config.Hidden
}

func (sc *StatusConverter) createServiceStatus(config types.ServiceConfig, serviceToContainer map[string]string, containerMap map[string]docker.ContainerStatus) display.ServiceStatus {
	provider := serviceToContainer[config.Name]
	providerName := sc.getProviderName(config.Name, provider)

	if containerStatus, exists := containerMap[provider]; exists {
		return sc.buildFoundStatus(config.Name, providerName, containerStatus)
	}

	return sc.buildNotFoundStatus(config.Name, providerName)
}

func (sc *StatusConverter) getProviderName(serviceName, provider string) string {
	if provider == serviceName {
		return ""
	}
	return provider
}

func (sc *StatusConverter) buildFoundStatus(name, provider string, containerStatus docker.ContainerStatus) display.ServiceStatus {
	uptime := time.Duration(0)
	if !containerStatus.StartedAt.IsZero() {
		uptime = time.Since(containerStatus.StartedAt)
	}

	return display.ServiceStatus{
		Name:      name,
		Provider:  provider,
		State:     containerStatus.State,
		Health:    containerStatus.Health,
		Ports:     containerStatus.Ports,
		CreatedAt: containerStatus.CreatedAt,
		UpdatedAt: containerStatus.StartedAt,
		Uptime:    uptime,
	}
}

func (sc *StatusConverter) buildNotFoundStatus(name, provider string) display.ServiceStatus {
	return display.ServiceStatus{
		Name:     name,
		Provider: provider,
		State:    "not found",
		Health:   "unknown",
	}
}
