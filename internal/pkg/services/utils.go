package services

import (
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
)

// ExtractServiceNames extracts service names from ServiceConfigs
func ExtractServiceNames(serviceConfigs []servicetypes.ServiceConfig) []string {
	if len(serviceConfigs) == 0 {
		return nil
	}
	serviceNames := make([]string, len(serviceConfigs))
	for i, config := range serviceConfigs {
		serviceNames[i] = config.Name
	}
	return serviceNames
}

// ServiceUtils provides service operations
type ServiceUtils struct {
	manager *Manager
}

// NewServiceUtils creates a new service utilities instance
func NewServiceUtils() *ServiceUtils {
	manager, _ := New()
	return &ServiceUtils{manager: manager}
}

// LoadServicesByCategory loads services organized by category
func (u *ServiceUtils) LoadServicesByCategory() (map[string][]servicetypes.ServiceConfig, error) {
	allServices := u.manager.GetAllServices()
	categories := make(map[string][]servicetypes.ServiceConfig)

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
func (u *ServiceUtils) LoadServiceConfig(serviceName string) (*servicetypes.ServiceConfig, error) {
	service, err := u.manager.GetService(serviceName)
	if err != nil {
		return nil, err
	}
	if service.Hidden {
		return nil, pkgerrors.NewValidationErrorf(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, "service not accessible: %s", serviceName)
	}
	return service, nil
}

// GetServicesByCategory loads services organized by category (alias)
func (u *ServiceUtils) GetServicesByCategory() (map[string][]servicetypes.ServiceConfig, error) {
	return u.LoadServicesByCategory()
}
