package services

import (
	"github.com/otto-nation/otto-stack/internal/core"
)

// ServiceUtils provides service operations
type ServiceUtils struct {
	manager *Manager
}

// NewServiceUtils creates a new service utilities instance
func NewServiceUtils() *ServiceUtils {
	manager, _ := New()
	return &ServiceUtils{manager: manager}
}

// ResolveServices applies composite expansion and dependency resolution
func (u *ServiceUtils) ResolveServices(serviceNames []string) ([]string, error) {
	return u.manager.ResolveServices(serviceNames)
}

// LoadServicesByCategory loads services organized by category
func (u *ServiceUtils) LoadServicesByCategory() (map[string][]ServiceConfig, error) {
	allServices := u.manager.GetAllServices()
	categories := make(map[string][]ServiceConfig)

	for _, service := range allServices {
		// Skip hidden services
		if service.Hidden {
			continue
		}
		categories[service.Category] = append(categories[service.Category], service)
	}

	return categories, nil
}

// LoadServiceConfig loads a specific service configuration
func (u *ServiceUtils) LoadServiceConfig(serviceName string) (*ServiceConfig, error) {
	return u.manager.GetService(serviceName)
}

// GetServicesByCategory loads services organized by category (alias)
func (u *ServiceUtils) GetServicesByCategory() (map[string][]ServiceConfig, error) {
	return u.LoadServicesByCategory()
}

// IsYAMLFile checks if filename is a YAML file (alias to constants function)
func IsYAMLFile(filename string) bool {
	return core.IsYAMLFile(filename)
}

// TrimYAMLExt removes YAML extension from filename (alias to constants function)
func TrimYAMLExt(filename string) string {
	return core.TrimYAMLExt(filename)
}
