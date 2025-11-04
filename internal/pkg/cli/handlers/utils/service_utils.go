package utils

import (
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// ServiceUtils provides service operations
type ServiceUtils struct {
	loader   *ServiceLoader
	resolver *DependencyResolver
}

// NewServiceUtils creates a new service utilities instance
func NewServiceUtils() *ServiceUtils {
	return &ServiceUtils{
		loader:   NewServiceLoader(),
		resolver: NewDependencyResolver(),
	}
}

// LoadServicesByCategory loads services organized by category
func (u *ServiceUtils) LoadServicesByCategory() (map[string][]types.ServiceInfo, error) {
	return u.loader.LoadServicesByCategory()
}

// LoadServiceConfig loads a service configuration
func (u *ServiceUtils) LoadServiceConfig(serviceName string) (*types.ServiceConfig, error) {
	return u.loader.LoadServiceConfig(serviceName)
}

// ResolveServices applies composite expansion and dependency resolution
func (u *ServiceUtils) ResolveServices(serviceNames []string) ([]string, error) {
	return u.resolver.ResolveServices(serviceNames)
}

// GetServicesByCategory loads services organized by category (alias for LoadServicesByCategory)
func (u *ServiceUtils) GetServicesByCategory() (map[string][]types.ServiceInfo, error) {
	return u.LoadServicesByCategory()
}

// LoadAllServiceDependencies returns empty map (deprecated functionality)
func (u *ServiceUtils) LoadAllServiceDependencies() (map[string][]string, error) {
	return make(map[string][]string), nil
}
