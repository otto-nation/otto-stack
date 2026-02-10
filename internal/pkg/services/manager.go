package services

import (
	"fmt"
	"maps"

	"github.com/otto-nation/otto-stack/internal/config"
	"github.com/otto-nation/otto-stack/internal/core"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"gopkg.in/yaml.v3"
)

// Manager handles all service operations
type Manager struct {
	services map[string]servicetypes.ServiceConfig
}

// New creates a new service manager
func New() (*Manager, error) {
	manager := &Manager{
		services: make(map[string]servicetypes.ServiceConfig),
	}

	if err := manager.loadServices(); err != nil {
		return nil, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentServices, messages.ErrorsServiceLoadFailed, err)
	}

	return manager, nil
}

// GetService returns a service by name
func (m *Manager) GetService(name string) (*servicetypes.ServiceConfig, error) {
	service, exists := m.services[name]
	if !exists {
		return nil, pkgerrors.NewValidationErrorf(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ErrorsServiceNotFound, name)
	}
	return &service, nil
}

// GetAllServices returns all services
func (m *Manager) GetAllServices() map[string]servicetypes.ServiceConfig {
	return m.services
}

// ValidateServices validates a list of service names
func (m *Manager) ValidateServices(serviceNames []string) error {
	for _, name := range serviceNames {
		service, exists := m.services[name]
		if !exists {
			return pkgerrors.NewValidationErrorf(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ErrorsServiceUnknown, name)
		}
		if service.Hidden {
			return pkgerrors.NewValidationErrorf(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ErrorsServiceNotAccessible, name)
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

// loadServices loads all services from embedded filesystem
func (m *Manager) loadServices() error {
	entries, err := config.EmbeddedServicesFS.ReadDir(EmbeddedServicesDir)
	if err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentServices, messages.ErrorsServiceReadDirectoryFailed, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		category := entry.Name()
		if err := m.loadCategoryServices(category); err != nil {
			return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentServices, messages.ErrorsServiceLoadCategoryFailed, err)
		}
	}

	return nil
}

// loadCategoryServices loads services from a specific category directory
func (m *Manager) loadCategoryServices(category string) error {
	categoryPath := fmt.Sprintf("%s/%s", EmbeddedServicesDir, category)
	entries, err := config.EmbeddedServicesFS.ReadDir(categoryPath)
	if err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentServices, messages.ErrorsServiceReadCategoryFailed, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !core.IsYAMLFile(entry.Name()) {
			continue
		}

		fileName := entry.Name()
		serviceName := core.TrimYAMLExt(fileName)

		if err := m.loadService(category, serviceName); err != nil {
			return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentServices, messages.ErrorsServiceLoadFailed, err)
		}
	}

	return nil
}

// loadService loads a single service from YAML
func (m *Manager) loadService(category, serviceName string) error {
	// Try exact filename first
	exactPath := fmt.Sprintf("%s/%s/%s.yaml", EmbeddedServicesDir, category, serviceName)
	if data, err := config.EmbeddedServicesFS.ReadFile(exactPath); err == nil {
		return m.parseService(data, serviceName, category)
	}

	return pkgerrors.NewValidationErrorf(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ErrorsServiceFileNotFound, serviceName)
}

func (m *Manager) parseService(data []byte, serviceName, category string) error {
	var service servicetypes.ServiceConfig
	if err := yaml.Unmarshal(data, &service); err != nil {
		return pkgerrors.NewConfigError(pkgerrors.ErrCodeOperationFail, serviceName, messages.ErrorsConfigParseServiceFailed, err)
	}

	// Set category if not specified in YAML
	if service.Category == "" {
		service.Category = category
	}

	// Use the name from the YAML file, not the filename
	actualServiceName := service.Name
	if actualServiceName == "" {
		actualServiceName = serviceName
	}

	// Combine environment variables into AllEnvironment
	service.AllEnvironment = make(map[string]string)
	maps.Copy(service.AllEnvironment, service.Environment)
	maps.Copy(service.AllEnvironment, service.Container.Environment)

	m.services[actualServiceName] = service
	return nil
}

// ExecuteCustomOperation executes custom management operations
func (m *Manager) ExecuteCustomOperation(serviceName, operationName string) ([]string, error) {
	service, err := m.GetService(serviceName)
	if err != nil {
		return nil, err
	}

	if service.Service.Management == nil || service.Service.Management.Custom == nil {
		return nil, pkgerrors.NewConfigErrorf(pkgerrors.ErrCodeOperationFail, pkgerrors.FieldServiceName, messages.ErrorsServiceNoCustomOperations, serviceName)
	}

	operation, exists := service.Service.Management.Custom[operationName]
	if !exists {
		return nil, pkgerrors.NewValidationErrorf(pkgerrors.ErrCodeInvalid, "operation", messages.ErrorsServiceOperationNotFound, operationName)
	}

	cmd := make([]string, len(operation.Command))
	copy(cmd, operation.Command)

	if args, exists := operation.Args["default"]; exists {
		cmd = append(cmd, args...)
	}

	return cmd, nil
}
