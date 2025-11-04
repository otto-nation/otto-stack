package services

import (
	"fmt"
	"strings"

	"github.com/otto-nation/otto-stack/internal/config"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"gopkg.in/yaml.v3"
)

// Manager handles all service operations
type Manager struct {
	services map[string]Service
}

// New creates a new service manager
func New() (*Manager, error) {
	manager := &Manager{
		services: make(map[string]Service),
	}

	if err := manager.loadServices(); err != nil {
		return nil, fmt.Errorf("failed to load services: %w", err)
	}

	return manager, nil
}

// GetService returns a service by name
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
	return service.Dependencies, nil
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
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		serviceName := strings.TrimSuffix(entry.Name(), ".yaml")
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
