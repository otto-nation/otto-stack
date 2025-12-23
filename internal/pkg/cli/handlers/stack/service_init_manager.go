package stack

import (
	"context"
	"fmt"
	"maps"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

// ServiceInitManager handles service-specific init container configuration
type ServiceInitManager struct {
	stackService *services.Service
}

// NewServiceInitManager creates a new service init manager
func NewServiceInitManager() (*ServiceInitManager, error) {
	stackService, err := NewStackService(false)
	if err != nil {
		return nil, err
	}

	return &ServiceInitManager{
		stackService: stackService,
	}, nil
}

// RunInitContainers runs init containers for the specified services
func (m *ServiceInitManager) RunInitContainers(ctx context.Context, serviceConfigs map[string]*services.ServiceConfig, projectName string) error {
	for serviceName, service := range serviceConfigs {
		if service.InitContainer != nil && service.InitContainer.Enabled {
			// Convert service config to init container config
			config := m.buildInitContainerConfig(serviceName, service, projectName)

			// Execute each script as a separate init container
			for i, script := range service.InitContainer.Scripts {
				containerName := fmt.Sprintf("%s-%s-init-%d", projectName, serviceName, i)
				scriptConfig := config
				scriptConfig.Command = []string{docker.ShellSh, docker.ShellC, script.Content}

				if err := m.stackService.DockerClient.RunInitContainer(ctx, containerName, scriptConfig); err != nil {
					return fmt.Errorf("failed to run init container for %s: %w", serviceName, err)
				}
			}
		}
	}
	return nil
}

// buildInitContainerConfig converts service config to init container config
func (m *ServiceInitManager) buildInitContainerConfig(serviceName string, service *services.ServiceConfig, projectName string) docker.InitContainerConfig {
	image := "alpine:latest" // Default image
	if service.InitContainer.Image != "" {
		image = service.InitContainer.Image
	}

	env := make(map[string]string)
	env["SERVICE_NAME"] = serviceName

	// Copy service environment
	maps.Copy(env, service.Environment)

	return docker.InitContainerConfig{
		Image:       image,
		Environment: env,
		Networks:    []string{projectName + docker.NetworkNameSuffix},
		WorkingDir:  "/",
	}
}
