package compose

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
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
	// Use utils for service resolution (no duplication)
	serviceUtils := utils.NewServiceUtils()
	resolvedServices, err := serviceUtils.ResolveServices(serviceNames)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve services: %w", err)
	}

	// Build compose structure
	compose := map[string]any{
		constants.ComposeFieldServices: g.buildServices(resolvedServices),
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
	services := make(map[string]any)

	for _, serviceName := range serviceNames {
		serviceDef, err := g.manager.GetService(serviceName)
		if err != nil {
			continue
		}

		// Skip configuration services (they merge into container services)
		if serviceDef.Type == constants.ServiceTypeConfiguration {
			continue
		}

		services[serviceName] = g.buildService(serviceDef)
	}

	return services
}

// buildService creates a single service definition
func (g *Generator) buildService(def services.Service) map[string]any {
	service := map[string]any{
		constants.ComposeFieldImage: def.Docker.Image,
	}

	if len(def.Docker.Ports) > 0 {
		service[constants.ComposeFieldPorts] = def.Docker.Ports
	}

	if len(def.Docker.Environment) > 0 {
		service[constants.ComposeFieldEnvironment] = def.Docker.Environment
	}

	if len(def.Docker.SimpleVolumes) > 0 {
		service[constants.ComposeFieldVolumes] = def.Docker.SimpleVolumes
	}

	if def.Docker.Restart != "" {
		service[constants.ComposeFieldRestart] = def.Docker.Restart
	}

	if len(def.Docker.Command) > 0 {
		service[constants.ComposeFieldCommand] = def.Docker.Command
	}

	if len(def.Docker.DependsOn) > 0 {
		service[constants.ComposeFieldDependsOn] = def.Docker.DependsOn
	}

	return service
}

// Generate creates a compose structure (for backward compatibility)
func (g *Generator) Generate(serviceNames []string) (map[string]any, error) {
	// Use utils for service resolution
	serviceUtils := utils.NewServiceUtils()
	resolvedServices, err := serviceUtils.ResolveServices(serviceNames)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve services: %w", err)
	}

	// Build compose structure
	return map[string]any{
		constants.ComposeFieldServices: g.buildServices(resolvedServices),
		constants.ComposeFieldNetworks: map[string]any{
			"default": map[string]any{
				constants.ComposeFieldName: fmt.Sprintf("%s-network", g.projectName),
			},
		},
	}, nil
}
