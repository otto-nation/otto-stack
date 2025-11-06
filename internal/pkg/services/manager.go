package services

import (
	"fmt"

	"github.com/otto-nation/otto-stack/internal/config"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"gopkg.in/yaml.v3"
)

// Manager handles all service operations
type Manager struct {
	servicesV2 map[string]ServiceConfigV2
}

// New creates a new service manager
func New() (*Manager, error) {
	manager := &Manager{
		servicesV2: make(map[string]ServiceConfigV2),
	}

	if err := manager.loadServices(); err != nil {
		return nil, fmt.Errorf("failed to load services: %w", err)
	}

	return manager, nil
}

// GetService returns a service by name (V2 format only)
func (m *Manager) GetService(name string) (ServiceConfigV2, error) {
	service, exists := m.servicesV2[name]
	if !exists {
		return ServiceConfigV2{}, fmt.Errorf("service not found: %s", name)
	}
	return service, nil
}

// GetAllServices returns all services
func (m *Manager) GetAllServices() map[string]ServiceConfigV2 {
	return m.servicesV2
}

// ValidateServices validates a list of service names
func (m *Manager) ValidateServices(serviceNames []string) error {
	for _, name := range serviceNames {
		if _, exists := m.servicesV2[name]; !exists {
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
	return service.Service.Dependencies.Required, nil
}

// GetServiceV2 returns a service by name in V2 format
func (m *Manager) GetServiceV2(name string) (*ServiceConfigV2, error) {
	if service, exists := m.servicesV2[name]; exists {
		return &service, nil
	}
	return nil, fmt.Errorf("service not found: %s", name)
}

// BuildConnectCommand builds a connection command for a service
func (m *Manager) BuildConnectCommand(serviceName string, options map[string]string) ([]string, error) {
	v2Service, err := m.GetServiceV2(serviceName)
	if err != nil {
		return nil, err
	}

	return m.buildV2ConnectCommand(v2Service, options)
}

// buildV2ConnectCommand builds connection command from V2 management spec
func (m *Manager) buildV2ConnectCommand(v2 *ServiceConfigV2, options map[string]string) ([]string, error) {
	if v2.Service.Management == nil || v2.Service.Management.Connect == nil {
		return nil, fmt.Errorf("no connect operation configured for service: %s", v2.Name)
	}

	connect := v2.Service.Management.Connect
	if len(connect.Command) == 0 {
		return nil, fmt.Errorf("no connect command configured for service: %s", v2.Name)
	}

	cmd := make([]string, len(connect.Command))
	copy(cmd, connect.Command)

	// Use default args if no specific args provided
	if args, exists := connect.Args["default"]; exists {
		cmd = append(cmd, args...)
	}

	// Apply any provided options/overrides
	for key, value := range options {
		if key == "database" && v2.Service.Connection != nil && v2.Service.Connection.DBFlag != "" {
			cmd = append(cmd, v2.Service.Connection.DBFlag, value)
		}
		// Add more option handling as needed
	}

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

		fileName := entry.Name()
		serviceName := constants.TrimYAMLExt(fileName)

		if err := m.loadService(category, serviceName); err != nil {
			return fmt.Errorf("failed to load service %s: %w", serviceName, err)
		}
	}

	return nil
}

// loadService loads a single service from YAML (V2 only)
func (m *Manager) loadService(category, serviceName string) error {
	// Try V2 format
	v2Path := fmt.Sprintf("%s/%s/%s-v2.yaml", constants.EmbeddedServicesDir, category, serviceName)
	if data, err := config.EmbeddedServicesFS.ReadFile(v2Path); err == nil {
		return m.loadV2Service(data, serviceName, category)
	}

	// Try exact filename (for services like redis-v2)
	exactPath := fmt.Sprintf("%s/%s/%s.yaml", constants.EmbeddedServicesDir, category, serviceName)
	if data, err := config.EmbeddedServicesFS.ReadFile(exactPath); err == nil {
		return m.loadV2Service(data, serviceName, category)
	}

	return fmt.Errorf("service file not found: %s", serviceName)
}

func (m *Manager) loadV2Service(data []byte, serviceName, category string) error {
	var serviceV2 ServiceConfigV2
	if err := yaml.Unmarshal(data, &serviceV2); err != nil {
		return fmt.Errorf("failed to parse V2 service YAML: %w", err)
	}

	// Set category if not specified in YAML
	if serviceV2.Category == "" {
		serviceV2.Category = category
	}

	// Use the name from the YAML file, not the filename
	actualServiceName := serviceV2.Name
	if actualServiceName == "" {
		actualServiceName = serviceName
	}

	m.servicesV2[actualServiceName] = serviceV2
	return nil
}

// isV2Format detects if the YAML is in V2 format

// ExecuteCustomOperation executes V2 custom management operations
func (m *Manager) ExecuteCustomOperation(serviceName, operationName string) ([]string, error) {
	service, err := m.GetServiceV2(serviceName)
	if err != nil {
		return nil, err
	}

	if service.Service.Management == nil || service.Service.Management.Custom == nil {
		return nil, fmt.Errorf("no custom operations for service: %s", serviceName)
	}

	operation, exists := service.Service.Management.Custom[operationName]
	if !exists {
		return nil, fmt.Errorf("operation %s not found", operationName)
	}

	cmd := make([]string, len(operation.Command))
	copy(cmd, operation.Command)

	if args, exists := operation.Args["default"]; exists {
		cmd = append(cmd, args...)
	}

	return cmd, nil
}

// ResolveServices applies composite expansion and dependency resolution
func (m *Manager) ResolveServices(serviceNames []string) ([]string, error) {
	// Simple pass-through for now - no expansion needed
	return serviceNames, nil
}
