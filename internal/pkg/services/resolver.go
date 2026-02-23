package services

import (
	"log/slog"

	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
)

// ServiceResolver handles service dependency resolution
type ServiceResolver struct {
	manager *Manager
	logger  *slog.Logger
}

// NewServiceResolver creates a new service resolver
func NewServiceResolver(manager *Manager) *ServiceResolver {
	return &ServiceResolver{
		manager: manager,
		logger:  logger.GetLogger(),
	}
}

// ResolveServices resolves service names with dependencies
// Validates that user-requested services exist and are accessible (not hidden)
// Then recursively includes all dependencies (including hidden ones)
func (r *ServiceResolver) ResolveServices(serviceNames []string) ([]servicetypes.ServiceConfig, error) {
	validationService := NewValidationService(r.manager)
	if err := validationService.ValidateUserServices(serviceNames); err != nil {
		return nil, err
	}

	resolved := make(map[string]bool)
	var ordered []string

	for _, name := range serviceNames {
		if err := r.resolveDependencies(name, false, resolved, &ordered); err != nil {
			return nil, err
		}
	}

	configs := make([]servicetypes.ServiceConfig, 0, len(ordered))
	for _, name := range ordered {
		service, err := r.manager.GetService(name)
		if err != nil {
			continue
		}
		configs = append(configs, *service)
	}

	return configs, nil
}

func (r *ServiceResolver) resolveDependencies(serviceName string, isDependency bool, resolved map[string]bool, ordered *[]string) error {
	if resolved[serviceName] {
		return nil
	}

	service, err := r.manager.GetService(serviceName)
	if err != nil {
		if isDependency {
			r.logger.Warn("Skipping missing dependency", "service", serviceName, "error", err)
			return nil
		}
		return err
	}

	r.logger.Debug("Resolving service", "service", serviceName, "isDependency", isDependency, "hidden", service.Hidden)

	for _, dep := range service.Service.Dependencies.Required {
		if err := r.resolveDependencies(dep, true, resolved, ordered); err != nil {
			return err
		}
	}

	resolved[serviceName] = true
	*ordered = append(*ordered, serviceName)
	return nil
}
