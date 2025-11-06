package compose

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
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
		constants.ComposeFieldServices: g.buildServices(serviceNames),
		constants.ComposeFieldNetworks: map[string]any{
			"default": map[string]any{
				constants.ComposeFieldName: fmt.Sprintf("%s-network", g.projectName),
			},
		},
	}

	return yaml.Marshal(compose)
}

// buildServices creates the services section
func (g *Generator) buildServices(serviceNames []string) map[string]any {
	serviceList := make(map[string]any)

	for _, serviceName := range serviceNames {
		serviceDef, err := g.manager.GetServiceV2(serviceName)
		if err != nil {
			continue
		}

		// Skip configuration services (they merge into container services)
		if serviceDef.ServiceType == services.ServiceTypeConfiguration {
			continue
		}

		serviceList[serviceName] = g.buildServiceFromV2(serviceDef)
	}

	return serviceList
}

func (g *Generator) buildServiceFromV2(v2 *services.ServiceConfigV2) map[string]any {
	service := map[string]any{
		constants.ComposeFieldImage: v2.Container.Image,
	}

	// Convert V2 ports to compose format
	if len(v2.Container.Ports) > 0 {
		var ports []string
		for _, port := range v2.Container.Ports {
			portStr := fmt.Sprintf("%s:%s", port.External, port.Internal)
			if port.Protocol != "" && port.Protocol != "tcp" {
				portStr += "/" + port.Protocol
			}
			ports = append(ports, portStr)
		}
		service[constants.ComposeFieldPorts] = ports
	}

	if len(v2.Container.Environment) > 0 {
		service[constants.ComposeFieldEnvironment] = v2.Container.Environment
	}

	if len(v2.Container.Volumes) > 0 {
		var volumes []string
		for _, vol := range v2.Container.Volumes {
			volStr := fmt.Sprintf("%s:%s", vol.Name, vol.Mount)
			if vol.ReadOnly {
				volStr += ":ro"
			}
			volumes = append(volumes, volStr)
		}
		service[constants.ComposeFieldVolumes] = volumes
	}

	if v2.Container.Restart != "" {
		service[constants.ComposeFieldRestart] = string(v2.Container.Restart)
	}

	if len(v2.Container.Command) > 0 {
		service[constants.ComposeFieldCommand] = v2.Container.Command
	}

	if v2.Container.MemoryLimit != "" {
		service["mem_limit"] = v2.Container.MemoryLimit
	}

	// Add health check if present
	if v2.Container.HealthCheck != nil {
		healthCheck := map[string]any{
			"test": v2.Container.HealthCheck.Test,
		}
		if v2.Container.HealthCheck.Interval > 0 {
			healthCheck["interval"] = v2.Container.HealthCheck.Interval.String()
		}
		if v2.Container.HealthCheck.Timeout > 0 {
			healthCheck["timeout"] = v2.Container.HealthCheck.Timeout.String()
		}
		if v2.Container.HealthCheck.Retries > 0 {
			healthCheck["retries"] = v2.Container.HealthCheck.Retries
		}
		if v2.Container.HealthCheck.StartPeriod > 0 {
			healthCheck["start_period"] = v2.Container.HealthCheck.StartPeriod.String()
		}
		service["healthcheck"] = healthCheck
	}

	return service
}

// Generate creates a compose structure (for backward compatibility)
func (g *Generator) Generate(serviceNames []string) (map[string]any, error) {
	// Build compose structure
	return map[string]any{
		constants.ComposeFieldServices: g.buildServices(serviceNames),
		constants.ComposeFieldNetworks: map[string]any{
			"default": map[string]any{
				constants.ComposeFieldName: fmt.Sprintf("%s-network", g.projectName),
			},
		},
	}, nil
}
