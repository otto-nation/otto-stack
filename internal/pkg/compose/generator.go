package compose

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/scripts"
)

// QuotedStringSlice ensures strings are quoted in YAML output
type QuotedStringSlice []string

func (q QuotedStringSlice) MarshalYAML() (interface{}, error) {
	if len(q) == 0 {
		return nil, nil
	}

	// Create a slice of yaml.Node with quoted strings
	var nodes []*yaml.Node
	for _, s := range q {
		node := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: s,
			Style: yaml.DoubleQuotedStyle,
		}
		nodes = append(nodes, node)
	}

	return &yaml.Node{
		Kind:    yaml.SequenceNode,
		Content: nodes,
	}, nil
}

// ComposeService represents a Docker Compose service configuration
type ComposeService struct {
	Image       string            `yaml:"image"`
	Ports       []string          `yaml:"ports,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	DependsOn   []string          `yaml:"depends_on,omitempty"`
	Command     QuotedStringSlice `yaml:"command,omitempty"`
	HealthCheck *HealthCheck      `yaml:"healthcheck,omitempty"`
	Restart     string            `yaml:"restart,omitempty"`
}

// HealthCheck represents Docker health check configuration
type HealthCheck struct {
	Test     QuotedStringSlice `yaml:"test"`
	Interval string            `yaml:"interval,omitempty"`
	Timeout  string            `yaml:"timeout,omitempty"`
	Retries  int               `yaml:"retries,omitempty"`
}

// ComposeFile represents a complete Docker Compose file
type ComposeFile struct {
	Services map[string]ComposeService `yaml:"services"`
	Volumes  map[string]any            `yaml:"volumes,omitempty"`
	Networks map[string]any            `yaml:"networks,omitempty"`
}

// Generator handles docker-compose file generation
type Generator struct {
	projectName string
	registry    *services.ServiceRegistry
}

// NewGenerator creates a new compose generator
func NewGenerator(projectName string, servicesPath string) (*Generator, error) {
	registry, err := services.NewServiceRegistry(servicesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create service registry: %w", err)
	}

	return &Generator{
		projectName: projectName,
		registry:    registry,
	}, nil
}

// Generate creates a docker-compose file for the specified services
func (g *Generator) Generate(serviceNames []string) (*ComposeFile, error) {
	compose := &ComposeFile{
		Services: make(map[string]ComposeService),
		Volumes:  make(map[string]any),
		Networks: make(map[string]any),
	}

	// Add default network
	compose.Networks["default"] = map[string]any{
		"name": fmt.Sprintf("%s-network", g.projectName),
	}

	// Use service utils for proper dependency resolution and composite expansion
	serviceUtils := utils.NewServiceUtils()

	// First expand composite services, then resolve dependencies
	expandedServices, err := serviceUtils.ExpandCompositeServices(serviceNames)
	if err != nil {
		return nil, fmt.Errorf("failed to expand composite services: %w", err)
	}

	resolvedServices, err := serviceUtils.ResolveDependencies(expandedServices)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	// Use the resolved services list (which excludes configuration and composite services)
	allServices := resolvedServices

	// Write scripts only if needed
	if err := g.writeRequiredScripts(allServices); err != nil {
		return nil, err
	}

	// Separate container vs configuration services
	containerServices := []string{}
	configServices := make(map[string][]services.ServiceDefinition) // target -> configs

	for _, serviceName := range allServices {
		serviceDef, exists := g.registry.GetService(serviceName)
		if !exists {
			return nil, fmt.Errorf("service %s not found", serviceName)
		}

		if serviceDef.Type == constants.ServiceTypeConfiguration {
			target := serviceDef.TargetService
			configServices[target] = append(configServices[target], serviceDef)
		} else {
			containerServices = append(containerServices, serviceName)
		}
	}

	// Add container services with merged configurations
	for _, serviceName := range containerServices {
		if err := g.addService(compose, serviceName); err != nil {
			return nil, fmt.Errorf("failed to add service %s: %w", serviceName, err)
		}

		// Apply configuration services
		if configs, exists := configServices[serviceName]; exists {
			if err := g.mergeConfigurations(compose, serviceName, configs); err != nil {
				return nil, fmt.Errorf("failed to merge configurations for %s: %w", serviceName, err)
			}
		}
	}

	return compose, nil
}

// GetRequiredPorts extracts all ports from the generated compose file
func (g *Generator) GetRequiredPorts(serviceNames []string) (map[string][]string, error) {
	compose, err := g.Generate(serviceNames)
	if err != nil {
		return nil, err
	}

	ports := make(map[string][]string)
	for serviceName, service := range compose.Services {
		if len(service.Ports) > 0 {
			var servicePorts []string
			for _, port := range service.Ports {
				// Extract host port from "host:container" format
				if strings.Contains(port, ":") {
					hostPort := strings.Split(port, ":")[0]
					servicePorts = append(servicePorts, hostPort)
				} else {
					servicePorts = append(servicePorts, port)
				}
			}
			if len(servicePorts) > 0 {
				ports[serviceName] = servicePorts
			}
		}
	}

	return ports, nil
}

// mergeConfigurations applies configuration services to a container service
func (g *Generator) mergeConfigurations(compose *ComposeFile, serviceName string, configs []services.ServiceDefinition) error {
	service := compose.Services[serviceName]

	for _, config := range configs {
		// Merge environment additions
		if len(config.EnvironmentAdditions) > 0 {
			if service.Environment == nil {
				service.Environment = make(map[string]string)
			}
			for key, value := range config.EnvironmentAdditions {
				// For SERVICES key, append instead of replace
				if key == "SERVICES" {
					if existing, exists := service.Environment[key]; exists && existing != "" {
						service.Environment[key] = existing + "," + value
					} else {
						service.Environment[key] = value
					}
				} else {
					service.Environment[key] = value
				}
			}
		}
	}

	compose.Services[serviceName] = service
	return nil
}

// resolveDependencies resolves all dependencies recursively and returns a unique list
//
//nolint:unused // Used internally for dependency resolution
func (g *Generator) resolveDependencies(serviceNames []string) ([]string, error) {
	resolved := make(map[string]bool)
	var result []string

	var resolve func(serviceName string) error
	resolve = func(serviceName string) error {
		if resolved[serviceName] {
			return nil
		}

		// Get service definition
		serviceDef, exists := g.registry.GetService(serviceName)
		if !exists {
			return fmt.Errorf("service %s not found", serviceName)
		}

		// Resolve dependencies first
		for _, dep := range serviceDef.Dependencies.Required {
			if err := resolve(dep); err != nil {
				return err
			}
		}

		// Add this service
		if !resolved[serviceName] {
			resolved[serviceName] = true
			result = append(result, serviceName)
		}

		return nil
	}

	// Resolve all requested services
	for _, serviceName := range serviceNames {
		if err := resolve(serviceName); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// expandCompositeServices expands composite services to their component services
//
//nolint:unused // Used for recursive composite service expansion
func (g *Generator) expandCompositeServices(serviceNames []string) ([]string, error) {
	var expandedServices []string

	for _, serviceName := range serviceNames {
		if g.registry != nil {
			if serviceDef, exists := g.registry.GetService(serviceName); exists {
				if serviceDef.Type == constants.ServiceTypeComposite && len(serviceDef.Components) > 0 {
					// Recursively expand composite services (in case components are also composite)
					componentServices, err := g.expandCompositeServices(serviceDef.Components)
					if err != nil {
						return nil, fmt.Errorf("failed to expand components of %s: %w", serviceName, err)
					}
					expandedServices = append(expandedServices, componentServices...)
				} else {
					// Regular service, keep as-is
					expandedServices = append(expandedServices, serviceName)
				}
			} else {
				// Service not found in registry, keep as-is
				expandedServices = append(expandedServices, serviceName)
			}
		} else {
			// No registry, keep as-is
			expandedServices = append(expandedServices, serviceName)
		}
	}

	return expandedServices, nil
}

// GenerateYAML generates the docker-compose YAML content
func (g *Generator) GenerateYAML(serviceNames []string) ([]byte, error) {
	compose, err := g.Generate(serviceNames)
	if err != nil {
		return nil, err
	}

	return yaml.Marshal(compose)
}

// addService adds a specific service to the compose file based on service definitions
func (g *Generator) addService(compose *ComposeFile, serviceName string) error {
	// Try to get service from registry first
	if g.registry != nil {
		if serviceDef, exists := g.registry.GetService(serviceName); exists {
			return g.addServiceFromDefinition(compose, serviceName, serviceDef)
		}
	}
	return nil
}

// addServiceFromDefinition creates a compose service from a service definition
func (g *Generator) addServiceFromDefinition(compose *ComposeFile, serviceName string, def services.ServiceDefinition) error {
	// Skip composite services - they don't have their own containers
	if def.Type == constants.ServiceTypeComposite {
		return nil
	}

	// Hidden services can still be processed if they're needed as dependencies
	// Only skip them if they're truly not needed (this filtering happens at selection time)

	service := ComposeService{
		Restart: "unless-stopped",
	}

	// Set image from docker.image or docker_image option in service_configuration
	if def.Docker.Image != "" {
		service.Image = def.Docker.Image
	} else {
		// Check service_configuration for docker_image option
		for _, option := range def.ServiceConfiguration {
			if option.Name == "docker_image" && option.Default != "" {
				service.Image = option.Default
				break
			}
		}
	}

	// Configure environment variables from Docker section only
	if len(def.Docker.Environment) > 0 {
		if service.Environment == nil {
			service.Environment = make(map[string]string)
		}
		for _, env := range def.Docker.Environment {
			// Parse KEY=VALUE format
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				service.Environment[parts[0]] = parts[1]
			}
		}
	}

	// Configure ports from Docker section only
	if len(def.Docker.Ports) > 0 {
		service.Ports = def.Docker.Ports
	}

	// Configure volumes - handle both simple string volumes and VolumeConfig
	if len(def.Docker.SimpleVolumes) > 0 {
		for _, vol := range def.Docker.SimpleVolumes {
			// Resolve relative paths to prevent nested otto-stack directories
			resolvedVol := strings.Replace(vol, "./"+constants.DevStackDir+"/", "./", 1)
			service.Volumes = append(service.Volumes, resolvedVol)
		}
	}

	// Configure complex volumes - process VolumeConfig to create volume mounts
	for _, vol := range def.Docker.Volumes {
		volumeName := fmt.Sprintf("%s-%s", g.projectName, vol.Name)
		compose.Volumes[volumeName] = map[string]interface{}{}

		// Add volume mount to service
		volumeMount := fmt.Sprintf("%s:%s", volumeName, vol.Mount)
		service.Volumes = append(service.Volumes, volumeMount)
	}

	// Configure health check
	if len(def.Docker.HealthCheck.Test) > 0 {
		service.HealthCheck = &HealthCheck{
			Test:     QuotedStringSlice(def.Docker.HealthCheck.Test),
			Interval: def.Docker.HealthCheck.Interval,
			Timeout:  def.Docker.HealthCheck.Timeout,
			Retries:  def.Docker.HealthCheck.Retries,
		}
	}

	// Configure restart policy - handle init containers generically
	if def.Docker.Restart != "" {
		service.Restart = def.Docker.Restart
	}

	// Services with restart: "no" are treated as init containers
	// They should run once and exit, not appear as persistent services

	// Configure command
	if len(def.Docker.Command) > 0 {
		service.Command = QuotedStringSlice(def.Docker.Command)
	}

	// Configure dependencies
	if len(def.Dependencies.Required) > 0 {
		service.DependsOn = def.Dependencies.Required
	}
	if len(def.Docker.DependsOn) > 0 {
		service.DependsOn = append(service.DependsOn, def.Docker.DependsOn...)
	}

	compose.Services[serviceName] = service
	return nil
}

// writeRequiredScripts writes only the scripts needed for enabled services
func (g *Generator) writeRequiredScripts(services []string) error {
	scriptMap := map[string]string{
		constants.ServiceKafkaTopics:    constants.KafkaTopicsInitScript,
		constants.ServiceLocalstackInit: constants.LocalstackInitScript,
	}

	scriptContent := map[string]string{
		constants.KafkaTopicsInitScript: scripts.KafkaTopicsInitScript,
		constants.LocalstackInitScript:  scripts.LocalstackInitScript,
	}

	var scriptsToWrite []string
	for _, service := range services {
		if scriptFile, exists := scriptMap[service]; exists {
			scriptsToWrite = append(scriptsToWrite, scriptFile)
		}
	}

	if len(scriptsToWrite) == 0 {
		return nil
	}

	scriptsDir := filepath.Join(constants.DevStackDir, constants.ScriptsDir)
	if err := os.MkdirAll(scriptsDir, 0755); err != nil {
		return fmt.Errorf("failed to create scripts directory: %w", err)
	}

	for _, scriptFile := range scriptsToWrite {
		scriptPath := filepath.Join(scriptsDir, scriptFile)
		content := scriptContent[scriptFile]
		if err := os.WriteFile(scriptPath, []byte(content), 0755); err != nil {
			return fmt.Errorf("failed to write script %s: %w", scriptFile, err)
		}
	}

	return nil
}

// isServiceEnabled checks if a service is in the list of enabled services
//
//nolint:unused // Used by writeRequiredScripts method
func (g *Generator) isServiceEnabled(services []string, serviceName string) bool {
	for _, service := range services {
		if service == serviceName {
			return true
		}
	}
	return false
}
