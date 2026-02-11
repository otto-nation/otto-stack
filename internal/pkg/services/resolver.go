package services

import (
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
)

// ServiceResolver handles service dependency resolution
type ServiceResolver struct {
	manager *Manager
}

// NewServiceResolver creates a new service resolver
func NewServiceResolver(manager *Manager) *ServiceResolver {
	return &ServiceResolver{manager: manager}
}

// ResolveServices resolves service names with dependencies
func (r *ServiceResolver) ResolveServices(serviceNames []string) ([]servicetypes.ServiceConfig, error) {
	// Validate services exist
	validator := NewValidator()
	if err := validator.ValidateServiceNames(serviceNames); err != nil {
		return nil, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, messages.ErrorsStackResolveServicesFailed, err)
	}

	// Resolve all dependencies recursively
	resolvedNames := make(map[string]bool)
	var allServiceNames []string

	var resolveDependencies func(string) error
	resolveDependencies = func(serviceName string) error {
		if resolvedNames[serviceName] {
			return nil
		}

		service, err := r.manager.GetService(serviceName)
		if err != nil {
			return err
		}

		// First resolve dependencies
		for _, dep := range service.Service.Dependencies.Required {
			if err := resolveDependencies(dep); err != nil {
				// Skip missing dependencies (they might be virtual or init containers)
				continue
			}
		}

		// Add this service to output
		if !resolvedNames[serviceName] {
			resolvedNames[serviceName] = true
			allServiceNames = append(allServiceNames, serviceName)
		}

		return nil
	}

	// Resolve dependencies for all requested services
	for _, serviceName := range serviceNames {
		if err := resolveDependencies(serviceName); err != nil {
			return nil, err
		}
	}

	// Load ServiceConfigs for resolved services
	var serviceConfigs []servicetypes.ServiceConfig
	for _, serviceName := range allServiceNames {
		service, err := r.manager.GetService(serviceName)
		if err != nil {
			continue // Skip services that can't be loaded
		}
		serviceConfigs = append(serviceConfigs, *service)
	}

	return serviceConfigs, nil
}
