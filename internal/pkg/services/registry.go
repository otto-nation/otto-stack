package services

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/otto-nation/otto-stack/internal/config"
)

// ServiceRegistry manages service definitions and validation
type ServiceRegistry struct {
	services map[string]ServiceDefinition
}

// ServiceOption represents a configurable option for a service
type ServiceOption struct {
	Name        string   `yaml:"name"`
	Type        string   `yaml:"type"`
	Description string   `yaml:"description"`
	Default     string   `yaml:"default,omitempty"`
	Example     string   `yaml:"example,omitempty"`
	Required    bool     `yaml:"required,omitempty"`
	Values      []string `yaml:"values,omitempty"` // for enum types
}

// ServiceDefinition represents a complete service definition from services.yaml
type ServiceDefinition struct {
	Name                 string          `yaml:"name"`
	Description          string          `yaml:"description"`
	Category             string          `yaml:"category"`
	Type                 string          `yaml:"type,omitempty"`       // "container", "configuration", or "composite"
	Visibility           string          `yaml:"visibility,omitempty"` // "hidden" to hide from interactive selection
	Components           []string        `yaml:"components,omitempty"` // List of component services for composite type
	ServiceConfiguration []ServiceOption `yaml:"service_configuration,omitempty"`
	Documentation        Documentation   `yaml:"documentation,omitempty"`

	// Configuration service fields
	TargetService        string            `yaml:"target_service,omitempty"`
	EnvironmentAdditions map[string]string `yaml:"environment_additions,omitempty"`

	// Dependencies
	Dependencies ServiceDependencies `yaml:"dependencies,omitempty"`

	// Docker configuration
	Docker DockerConfig `yaml:"docker,omitempty"`

	// Extended properties
	DefaultPort int               `yaml:"default_port,omitempty"`
	HealthCheck HealthCheckConfig `yaml:"health_check,omitempty"`
	Tags        []string          `yaml:"tags,omitempty"`
}

// Documentation represents structured documentation for a service
type Documentation struct {
	Examples      []string          `yaml:"examples,omitempty"`
	UsageNotes    string            `yaml:"usage_notes,omitempty"`
	Links         []string          `yaml:"links,omitempty"`
	UseCases      []string          `yaml:"use_cases,omitempty"`
	WebInterfaces []WebInterface    `yaml:"web_interfaces,omitempty"`
	Docs          []DocLink         `yaml:"docs,omitempty"`
	CLICommands   map[string]string `yaml:"cli_commands,omitempty"`
	SpringConfig  SpringConfig      `yaml:"spring_config,omitempty"`
}

// WebInterface represents a web interface for a service
type WebInterface struct {
	Name        string `yaml:"name"`
	URL         string `yaml:"url"`
	Description string `yaml:"description"`
}

// DocLink represents a documentation link
type DocLink struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// SpringConfig represents Spring Boot configuration
type SpringConfig struct {
	Properties []string `yaml:"properties,omitempty"`
}

// ServiceDependencies represents service dependency configuration
type ServiceDependencies struct {
	Required  []string `yaml:"required,omitempty"`
	Soft      []string `yaml:"soft,omitempty"`
	Conflicts []string `yaml:"conflicts,omitempty"`
	Provides  []string `yaml:"provides,omitempty"`
}

// DockerConfig represents Docker Compose specific configuration
type DockerConfig struct {
	Image         string            `yaml:"image,omitempty"`
	Restart       string            `yaml:"restart,omitempty"`
	Networks      []string          `yaml:"networks,omitempty"`
	MemoryLimit   string            `yaml:"memory_limit,omitempty"`
	Environment   []string          `yaml:"environment,omitempty"`
	Command       []string          `yaml:"command,omitempty"`
	HealthCheck   DockerHealthCheck `yaml:"health_check,omitempty"`
	Ports         []string          `yaml:"ports,omitempty"`
	Volumes       []VolumeConfig    `yaml:"volumes,omitempty"`        // Complex volume configurations
	SimpleVolumes []string          `yaml:"simple_volumes,omitempty"` // Simple string volumes (host:container)
	DependsOn     []string          `yaml:"depends_on,omitempty"`
}

// DockerHealthCheck represents Docker health check configuration
type DockerHealthCheck struct {
	Test        []string `yaml:"test,omitempty"`
	Interval    string   `yaml:"interval,omitempty"`
	Timeout     string   `yaml:"timeout,omitempty"`
	Retries     int      `yaml:"retries,omitempty"`
	StartPeriod string   `yaml:"start_period,omitempty"`
}

// VolumeConfig represents volume configuration
type VolumeConfig struct {
	Name        string `yaml:"name"`
	Mount       string `yaml:"mount"`
	Description string `yaml:"description,omitempty"`
}

// HealthCheckConfig represents health check configuration
type HealthCheckConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint,omitempty"`
	Interval string `yaml:"interval,omitempty"`
	Timeout  string `yaml:"timeout,omitempty"`
	Retries  int    `yaml:"retries,omitempty"`
}

// ServiceManifest represents the structure of services.yaml
type ServiceManifest map[string]ServiceDefinition

// NewServiceRegistry creates a new service registry
func NewServiceRegistry(configPath string) (*ServiceRegistry, error) {
	registry := &ServiceRegistry{
		services: make(map[string]ServiceDefinition),
	}

	if err := registry.LoadFromEmbedded(); err != nil {
		return nil, fmt.Errorf("failed to load service registry: %w", err)
	}

	return registry, nil
}

// LoadFromEmbedded loads services from the embedded file system
func (r *ServiceRegistry) LoadFromEmbedded() error {
	// Walk through all service directories in the embedded FS
	categories := []string{"database", "cache", "messaging", "cloud", "observability"}

	for _, category := range categories {
		categoryPath := fmt.Sprintf("services/%s", category)
		entries, err := config.EmbeddedServicesFS.ReadDir(categoryPath)
		if err != nil {
			continue // Skip if category doesn't exist
		}

		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
				servicePath := fmt.Sprintf("%s/%s", categoryPath, entry.Name())
				if err := r.loadServiceFromEmbedded(servicePath); err != nil {
					return fmt.Errorf("failed to load service from %s: %w", servicePath, err)
				}
			}
		}
	}

	return nil
}

// loadServiceFromEmbedded loads a single service from the embedded FS
func (r *ServiceRegistry) loadServiceFromEmbedded(servicePath string) error {
	data, err := config.EmbeddedServicesFS.ReadFile(servicePath)
	if err != nil {
		return fmt.Errorf("failed to read service file: %w", err)
	}

	var serviceDef ServiceDefinition
	if err := yaml.Unmarshal(data, &serviceDef); err != nil {
		return fmt.Errorf("failed to parse service YAML: %w", err)
	}

	// Extract filename without extension
	filename := filepath.Base(servicePath)
	expectedName := strings.TrimSuffix(filename, ".yaml")

	// Validate that filename matches service name
	if serviceDef.Name != expectedName {
		return fmt.Errorf("service name '%s' does not match filename '%s' (expected '%s')",
			serviceDef.Name, filename, expectedName)
	}

	// Validate and store the service
	if err := r.validateServiceDefinition(serviceDef.Name, serviceDef); err != nil {
		return fmt.Errorf("invalid service definition for %s: %w", serviceDef.Name, err)
	}

	r.services[serviceDef.Name] = serviceDef
	return nil
}

// GetService returns a service definition by name
func (r *ServiceRegistry) GetService(name string) (ServiceDefinition, bool) {
	service, exists := r.services[name]
	return service, exists
}

// GetAllServices returns all service definitions
func (r *ServiceRegistry) GetAllServices() map[string]ServiceDefinition {
	// Return a copy to prevent modification
	result := make(map[string]ServiceDefinition)
	for name, service := range r.services {
		result[name] = service
	}
	return result
}

// GetServiceNames returns all service names
func (r *ServiceRegistry) GetServiceNames() []string {
	names := make([]string, 0, len(r.services))
	for name := range r.services {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GetServicesByCategory returns services in a specific category
func (r *ServiceRegistry) GetServicesByCategory(category string) []string {
	var services []string
	for name, service := range r.services {
		if service.Category == category {
			services = append(services, name)
		}
	}
	sort.Strings(services)
	return services
}

// GetServicesByTag returns services with a specific tag
func (r *ServiceRegistry) GetServicesByTag(tag string) []string {
	var services []string
	for name, service := range r.services {
		for _, serviceTag := range service.Tags {
			if serviceTag == tag {
				services = append(services, name)
				break
			}
		}
	}
	sort.Strings(services)
	return services
}

// ValidateService validates that a service exists
func (r *ServiceRegistry) ValidateService(name string) error {
	if _, exists := r.services[name]; !exists {
		available := r.GetServiceNames()
		return fmt.Errorf("unknown service '%s'. Available services: %v", name, available)
	}
	return nil
}

// ValidateServices validates multiple service names
func (r *ServiceRegistry) ValidateServices(names []string) error {
	for _, name := range names {
		if err := r.ValidateService(name); err != nil {
			return err
		}
	}
	return nil
}

// GetServiceDependencies returns the dependencies for a service
func (r *ServiceRegistry) GetServiceDependencies(name string) ([]string, error) {
	service, exists := r.GetService(name)
	if !exists {
		return nil, fmt.Errorf("service %s not found", name)
	}
	return service.Dependencies.Required, nil
}

// GetAllCategories returns all unique service categories
func (r *ServiceRegistry) GetAllCategories() []string {
	categories := make(map[string]bool)
	for _, service := range r.services {
		if service.Category != "" {
			categories[service.Category] = true
		}
	}

	result := make([]string, 0, len(categories))
	for category := range categories {
		result = append(result, category)
	}
	sort.Strings(result)
	return result
}

// GetAllTags returns all unique service tags
func (r *ServiceRegistry) GetAllTags() []string {
	tags := make(map[string]bool)
	for _, service := range r.services {
		for _, tag := range service.Tags {
			tags[tag] = true
		}
	}

	result := make([]string, 0, len(tags))
	for tag := range tags {
		result = append(result, tag)
	}
	sort.Strings(result)
	return result
}

// SearchServices searches for services by name, description, or tags
func (r *ServiceRegistry) SearchServices(query string) []string {
	var matches []string
	query = strings.ToLower(query)

	for name, service := range r.services {
		// Check name
		if strings.Contains(strings.ToLower(name), query) {
			matches = append(matches, name)
			continue
		}

		// Check description
		if strings.Contains(strings.ToLower(service.Description), query) {
			matches = append(matches, name)
			continue
		}

		// Check tags
		for _, tag := range service.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				matches = append(matches, name)
				break
			}
		}
	}

	sort.Strings(matches)
	return matches
}

// GetServiceInfo returns formatted service information
func (r *ServiceRegistry) GetServiceInfo(name string) (string, error) {
	service, exists := r.GetService(name)
	if !exists {
		return "", fmt.Errorf("service %s not found", name)
	}

	var info strings.Builder
	info.WriteString(fmt.Sprintf("Service: %s\n", name))
	info.WriteString(fmt.Sprintf("Description: %s\n", service.Description))

	if service.Category != "" {
		info.WriteString(fmt.Sprintf("Category: %s\n", service.Category))
	}

	if service.DefaultPort > 0 {
		info.WriteString(fmt.Sprintf("Default Port: %d\n", service.DefaultPort))
	}

	if len(service.Dependencies.Required) > 0 {
		info.WriteString(fmt.Sprintf("Dependencies: %s\n", strings.Join(service.Dependencies.Required, ", ")))
	}

	if len(service.Tags) > 0 {
		info.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(service.Tags, ", ")))
	}

	if len(service.ServiceConfiguration) > 0 {
		info.WriteString("\nConfiguration Options:\n")
		for _, option := range service.ServiceConfiguration {
			info.WriteString(fmt.Sprintf("  • %s (%s)\n", option.Name, option.Type))
			info.WriteString(fmt.Sprintf("    %s\n", option.Description))
			if option.Default != "" {
				info.WriteString(fmt.Sprintf("    Default: %s\n", option.Default))
			}
			if option.Example != "" {
				info.WriteString(fmt.Sprintf("    Example: %s\n", option.Example))
			}
			if option.Required {
				info.WriteString("    Required: Yes\n")
			}
			if len(option.Values) > 0 {
				info.WriteString(fmt.Sprintf("    Allowed values: %s\n", strings.Join(option.Values, ", ")))
			}
			info.WriteString("\n")
		}
	}

	if len(service.Documentation.Examples) > 0 {
		info.WriteString("\nExamples:\n")
		for _, example := range service.Documentation.Examples {
			info.WriteString(fmt.Sprintf("  %s\n", example))
		}
	}

	if service.Documentation.UsageNotes != "" {
		info.WriteString(fmt.Sprintf("\nUsage Notes:\n%s\n", service.Documentation.UsageNotes))
	}

	if len(service.Documentation.Links) > 0 {
		info.WriteString("\nLinks:\n")
		for _, link := range service.Documentation.Links {
			info.WriteString(fmt.Sprintf("  - %s\n", link))
		}
	}

	return info.String(), nil
}

// validateServiceDefinition validates a service definition
func (r *ServiceRegistry) validateServiceDefinition(name string, definition ServiceDefinition) error {
	if definition.Description == "" {
		return fmt.Errorf("description is required")
	}

	// Validate port if specified
	if definition.DefaultPort < 0 || definition.DefaultPort > 65535 {
		return fmt.Errorf("invalid default port: %d", definition.DefaultPort)
	}

	// Validate health check configuration
	if definition.HealthCheck.Enabled {
		if definition.HealthCheck.Endpoint == "" {
			return fmt.Errorf("health check endpoint is required when health check is enabled")
		}
	}

	return nil
}

// LoadDefaultServiceRegistry loads the service registry from embedded services
func LoadDefaultServiceRegistry() (*ServiceRegistry, error) {
	return NewServiceRegistry("")
}

// LoadServiceRegistryFromPath loads the service registry from embedded services (ignores path)
func LoadServiceRegistryFromPath(path string) (*ServiceRegistry, error) {
	return NewServiceRegistry("")
}
