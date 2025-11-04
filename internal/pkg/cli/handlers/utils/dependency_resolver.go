package utils

import (
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// DependencyResolver handles service dependency resolution
type DependencyResolver struct {
	loader *ServiceLoader
}

// NewDependencyResolver creates a new dependency resolver
func NewDependencyResolver() *DependencyResolver {
	return &DependencyResolver{
		loader: NewServiceLoader(),
	}
}

// ResolveServices expands composites and resolves dependencies
func (dr *DependencyResolver) ResolveServices(serviceNames []string) ([]string, error) {
	// Load all services once
	allServices, err := dr.loader.LoadAllServices()
	if err != nil {
		return serviceNames, err
	}

	// Expand composite services
	expanded := dr.expandComposites(serviceNames, allServices)

	// For dependency resolution, we need ServiceInfo which has Dependencies
	// For now, just return expanded services since ServiceConfig doesn't have Dependencies
	return dr.removeDuplicates(expanded), nil
}

func (dr *DependencyResolver) expandComposites(services []string, allServices map[string]*types.ServiceConfig) []string {
	var result []string

	for _, serviceName := range services {
		config, exists := allServices[serviceName]
		if !exists {
			result = append(result, serviceName)
			continue
		}

		if config.Type == "composite" && len(config.Components) > 0 {
			result = append(result, config.Components...)
		} else {
			result = append(result, serviceName)
		}
	}

	return result
}

func (dr *DependencyResolver) removeDuplicates(services []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, service := range services {
		if !seen[service] {
			seen[service] = true
			result = append(result, service)
		}
	}

	return result
}
