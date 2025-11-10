package compose

import (
	"fmt"

	"gopkg.in/yaml.v3"

	dockerConstants "github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

// Generator handles docker-compose file generation
type Generator struct {
	projectName string
	manager     *services.Manager
}

// NewGenerator creates a new compose generator
func NewGenerator(projectName string, servicesPath string) (*Generator, error) {
	manager, err := services.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create service manager: %w", err)
	}

	return &Generator{
		projectName: projectName,
		manager:     manager,
	}, nil
}

// GenerateYAML creates a docker-compose YAML for the specified services
func (g *Generator) GenerateYAML(serviceNames []string) ([]byte, error) {
	// Build compose structure with project-specific default network
	compose := map[string]any{
		dockerConstants.ComposeFieldServices: g.buildServices(serviceNames),
		dockerConstants.ComposeFieldNetworks: map[string]any{
			"default": map[string]any{
				dockerConstants.ComposeFieldName: fmt.Sprintf("%s-network", g.projectName),
			},
		},
	}

	return yaml.Marshal(compose)
}

// buildServices creates the services section
func (g *Generator) buildServices(serviceNames []string) map[string]any {
	serviceList := make(map[string]any)

	for _, serviceName := range serviceNames {
		serviceDef, err := g.manager.GetService(serviceName)
		if err != nil {
			continue
		}

		// Skip configuration services (they merge into container services)
		if serviceDef.ServiceType == services.ServiceTypeConfiguration {
			continue
		}

		serviceList[serviceName] = g.buildService(serviceDef)
	}

	return serviceList
}

func (g *Generator) buildService(config *services.ServiceConfig) map[string]any {
	service := map[string]any{
		dockerConstants.ComposeFieldImage: config.Container.Image,
	}

	// Convert ports to compose format
	if len(config.Container.Ports) > 0 {
		var ports []string
		for _, port := range config.Container.Ports {
			portStr := fmt.Sprintf("%s:%s", port.External, port.Internal)
			if port.Protocol != "" && port.Protocol != "tcp" {
				portStr += "/" + port.Protocol
			}
			ports = append(ports, portStr)
		}
		service[dockerConstants.ComposeFieldPorts] = ports
	}

	if len(config.Container.Environment) > 0 {
		service[dockerConstants.ComposeFieldEnvironment] = config.Container.Environment
	}

	if len(config.Container.Volumes) > 0 {
		var volumes []string
		for _, vol := range config.Container.Volumes {
			volStr := fmt.Sprintf("%s:%s", vol.Name, vol.Mount)
			if vol.ReadOnly {
				volStr += ":ro"
			}
			volumes = append(volumes, volStr)
		}
		service[dockerConstants.ComposeFieldVolumes] = volumes
	}

	if config.Container.Restart != "" {
		service[dockerConstants.ComposeFieldRestart] = string(config.Container.Restart)
	}

	if len(config.Container.Command) > 0 {
		service[dockerConstants.ComposeFieldCommand] = config.Container.Command
	}

	if config.Container.MemoryLimit != "" {
		service["mem_limit"] = config.Container.MemoryLimit
	}

	// Add health check if present
	if config.Container.HealthCheck != nil {
		healthCheck := map[string]any{
			"test": config.Container.HealthCheck.Test,
		}
		if config.Container.HealthCheck.Interval > 0 {
			healthCheck["interval"] = config.Container.HealthCheck.Interval.String()
		}
		if config.Container.HealthCheck.Timeout > 0 {
			healthCheck["timeout"] = config.Container.HealthCheck.Timeout.String()
		}
		if config.Container.HealthCheck.Retries > 0 {
			healthCheck["retries"] = config.Container.HealthCheck.Retries
		}
		if config.Container.HealthCheck.StartPeriod > 0 {
			healthCheck["start_period"] = config.Container.HealthCheck.StartPeriod.String()
		}
		service["healthcheck"] = healthCheck
	}

	return service
}

// Generate creates a compose structure (for backward compatibility)
func (g *Generator) Generate(serviceNames []string) (map[string]any, error) {
	// Build compose structure
	return map[string]any{
		dockerConstants.ComposeFieldServices: g.buildServices(serviceNames),
		dockerConstants.ComposeFieldNetworks: map[string]any{
			"default": map[string]any{
				dockerConstants.ComposeFieldName: fmt.Sprintf("%s-network", g.projectName),
			},
		},
	}, nil
}
