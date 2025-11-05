package services

import (
	"fmt"
	"strings"

	"github.com/otto-nation/otto-stack/internal/config"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"gopkg.in/yaml.v3"
)

// Manager handles all service operations
type Manager struct {
	services   map[string]Service
	servicesV2 map[string]types.ServiceConfigV2
}

// New creates a new service manager
func New() (*Manager, error) {
	manager := &Manager{
		services:   make(map[string]Service),
		servicesV2: make(map[string]types.ServiceConfigV2),
	}

	if err := manager.loadServices(); err != nil {
		return nil, fmt.Errorf("failed to load services: %w", err)
	}

	return manager, nil
}

// GetService returns a service by name (V1 format)
func (m *Manager) GetService(name string) (Service, error) {
	service, exists := m.services[name]
	if !exists {
		return Service{}, fmt.Errorf("service not found: %s", name)
	}
	return service, nil
}

// GetServicesByCategory returns services grouped by category
func (m *Manager) GetServicesByCategory() map[string][]Service {
	categories := make(map[string][]Service)

	for _, service := range m.services {
		categories[service.Category] = append(categories[service.Category], service)
	}

	return categories
}

// GetAllServices returns all services
func (m *Manager) GetAllServices() map[string]Service {
	return m.services
}

// ValidateServices validates a list of service names
func (m *Manager) ValidateServices(serviceNames []string) error {
	for _, name := range serviceNames {
		if _, exists := m.services[name]; !exists {
			return fmt.Errorf("unknown service: %s", name)
		}
	}
	return nil
}

// GetDependencies returns dependencies for a service
func (m *Manager) GetDependencies(serviceName string) ([]string, error) {
	service, err := m.GetService(serviceName)
	if err != nil {
		return nil, err
	}
	return service.Dependencies.Required, nil
}

// BuildConnectCommand builds a connection command for a service
func (m *Manager) BuildConnectCommand(serviceName string, options map[string]string) ([]string, error) {
	service, err := m.GetService(serviceName)
	if err != nil {
		return nil, err
	}

	if service.Connection.Client == "" {
		return nil, fmt.Errorf("no connection client configured for service: %s", serviceName)
	}

	cmd := []string{service.Connection.Client}

	// Add connection parameters
	if host, ok := options["host"]; ok && service.Connection.HostFlag != "" {
		cmd = append(cmd, service.Connection.HostFlag, host)
	}

	if port, ok := options["port"]; ok && service.Connection.PortFlag != "" {
		cmd = append(cmd, service.Connection.PortFlag, port)
	} else if service.Connection.DefaultPort > 0 && service.Connection.PortFlag != "" {
		cmd = append(cmd, service.Connection.PortFlag, fmt.Sprintf("%d", service.Connection.DefaultPort))
	}

	if user, ok := options["user"]; ok && service.Connection.UserFlag != "" {
		cmd = append(cmd, service.Connection.UserFlag, user)
	} else if service.Connection.DefaultUser != "" && service.Connection.UserFlag != "" {
		cmd = append(cmd, service.Connection.UserFlag, service.Connection.DefaultUser)
	}

	if database, ok := options["database"]; ok && service.Connection.DBFlag != "" {
		cmd = append(cmd, service.Connection.DBFlag, database)
	}

	// Add extra flags
	cmd = append(cmd, service.Connection.ExtraFlags...)

	return cmd, nil
}

// loadServices loads all services from embedded filesystem
func (m *Manager) loadServices() error {
	entries, err := config.EmbeddedServicesFS.ReadDir(constants.EmbeddedServicesDir)
	if err != nil {
		return fmt.Errorf("failed to read services directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		category := entry.Name()
		if err := m.loadCategoryServices(category); err != nil {
			return fmt.Errorf("failed to load category %s: %w", category, err)
		}
	}

	return nil
}

// loadCategoryServices loads services from a specific category directory
func (m *Manager) loadCategoryServices(category string) error {
	categoryPath := fmt.Sprintf("%s/%s", constants.EmbeddedServicesDir, category)
	entries, err := config.EmbeddedServicesFS.ReadDir(categoryPath)
	if err != nil {
		return fmt.Errorf("failed to read category directory %s: %w", category, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !constants.IsYAMLFile(entry.Name()) {
			continue
		}

		serviceName := constants.TrimYAMLExt(entry.Name())
		if err := m.loadService(category, serviceName); err != nil {
			return fmt.Errorf("failed to load service %s: %w", serviceName, err)
		}
	}

	return nil
}

// loadService loads a single service from YAML
func (m *Manager) loadService(category, serviceName string) error {
	servicePath := fmt.Sprintf("%s/%s/%s.yaml", constants.EmbeddedServicesDir, category, serviceName)
	data, err := config.EmbeddedServicesFS.ReadFile(servicePath)
	if err != nil {
		return fmt.Errorf("failed to read service file: %w", err)
	}

	// Detect V2 format by checking for V2-specific fields
	if isV2Format(data) {
		return m.loadV2Service(data, serviceName, category)
	}

	return m.loadV1Service(data, serviceName, category)
}

func (m *Manager) loadV2Service(data []byte, serviceName, category string) error {
	var serviceV2 types.ServiceConfigV2
	if err := yaml.Unmarshal(data, &serviceV2); err != nil {
		return fmt.Errorf("failed to parse V2 service YAML: %w", err)
	}

	m.servicesV2[serviceName] = serviceV2

	// Convert V2 to V1 for backward compatibility
	service := convertV2ToV1(serviceV2, category)
	m.services[serviceName] = service
	return nil
}

func (m *Manager) loadV1Service(data []byte, serviceName, category string) error {
	var service Service
	if err := yaml.Unmarshal(data, &service); err != nil {
		return fmt.Errorf("failed to parse service YAML: %w", err)
	}

	// Set category if not specified in YAML
	if service.Category == "" {
		service.Category = category
	}

	// Set name if not specified in YAML
	if service.Name == "" {
		service.Name = serviceName
	}

	m.services[serviceName] = service
	return nil
}

// isV2Format detects if the YAML is in V2 format
func isV2Format(data []byte) bool {
	content := string(data)
	// V2 format has both runtime and integration sections
	// and typically has structured volumes with name/mount/description
	return strings.Contains(content, "runtime:") ||
		(strings.Contains(content, "- name:") && strings.Contains(content, "mount:"))
}

// convertV2ToV1 converts V2 format to V1 for backward compatibility
func convertV2ToV1(v2 types.ServiceConfigV2, category string) Service {
	service := Service{
		Name:        v2.Name,
		Description: v2.Description,
		Category:    category,
		Type:        string(v2.Type),
		Environment: v2.Runtime.Environment,
	}

	// Convert Docker config
	service.Docker = DockerConfig{
		Image:   v2.Runtime.Image,
		Restart: string(v2.Runtime.Container.Restart),
		Command: v2.Runtime.Container.Command,
	}

	// Convert ports
	for _, port := range v2.Runtime.Ports {
		service.Docker.Ports = append(service.Docker.Ports, port.Host+":"+port.Container)
	}

	// Convert connection
	if v2.Integration.Connection != nil {
		service.Connection = ConnectionConfig{
			Client:      v2.Integration.Connection.Client,
			DefaultUser: v2.Integration.Connection.DefaultUser,
			DefaultPort: v2.Integration.Connection.DefaultPort,
			HostFlag:    v2.Integration.Connection.HostFlag,
			PortFlag:    v2.Integration.Connection.PortFlag,
			UserFlag:    v2.Integration.Connection.UserFlag,
			DBFlag:      v2.Integration.Connection.DBFlag,
		}
	}

	return service
}

// GetConnectionConfig returns connection configuration for a service
func GetConnectionConfig(serviceName string) (*ConnectionConfig, error) {
	manager, err := New()
	if err != nil {
		return nil, err
	}

	service, err := manager.GetService(serviceName)
	if err != nil {
		return nil, err
	}

	if service.Connection.Client == "" {
		return nil, fmt.Errorf("no connection client configured for service: %s", serviceName)
	}

	return &service.Connection, nil
}
