package utils

import (
	"fmt"
	"maps"

	"github.com/otto-nation/otto-stack/internal/config"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"gopkg.in/yaml.v3"
)

// ServiceLoader handles loading service configurations
type ServiceLoader struct{}

// NewServiceLoader creates a new service loader
func NewServiceLoader() *ServiceLoader {
	return &ServiceLoader{}
}

// LoadServicesByCategory loads all services organized by category
func (sl *ServiceLoader) LoadServicesByCategory() (map[string][]types.ServiceInfo, error) {
	categories, err := sl.getCategories()
	if err != nil {
		return nil, err
	}

	result := make(map[string][]types.ServiceInfo)
	for _, category := range categories {
		services := sl.loadServicesInCategory(category)
		if len(services) > 0 {
			result[category] = services
		}
	}

	return result, nil
}

// LoadServiceConfig loads a specific service configuration
func (sl *ServiceLoader) LoadServiceConfig(serviceName string) (*types.ServiceConfig, error) {
	categories, err := sl.getCategories()
	if err != nil {
		return nil, err
	}

	for _, category := range categories {
		if config := sl.loadServiceFromCategory(category, serviceName); config != nil {
			return config, nil
		}
	}

	return nil, fmt.Errorf("service not found: %s", serviceName)
}

// LoadAllServices loads all service configurations
func (sl *ServiceLoader) LoadAllServices() (map[string]*types.ServiceConfig, error) {
	categories, err := sl.getCategories()
	if err != nil {
		return nil, err
	}

	result := make(map[string]*types.ServiceConfig)
	for _, category := range categories {
		services := sl.loadAllInCategory(category)
		maps.Copy(result, services)
	}

	return result, nil
}

func (sl *ServiceLoader) getCategories() ([]string, error) {
	entries, err := config.EmbeddedServicesFS.ReadDir("services")
	if err != nil {
		return nil, fmt.Errorf("failed to read services directory: %w", err)
	}

	var categories []string
	for _, entry := range entries {
		if entry.IsDir() {
			categories = append(categories, entry.Name())
		}
	}
	return categories, nil
}

func (sl *ServiceLoader) loadServicesInCategory(category string) []types.ServiceInfo {
	categoryPath := fmt.Sprintf("services/%s", category)
	entries, err := config.EmbeddedServicesFS.ReadDir(categoryPath)
	if err != nil {
		return nil
	}

	var services []types.ServiceInfo
	for _, entry := range entries {
		if !constants.IsYAMLFile(entry.Name()) {
			continue
		}

		serviceName := constants.TrimYAMLExt(entry.Name())
		if info := sl.parseServiceInfo(categoryPath, entry.Name(), serviceName, category); info != nil {
			services = append(services, *info)
		}
	}

	return services
}

func (sl *ServiceLoader) loadServiceFromCategory(category, serviceName string) *types.ServiceConfig {
	servicePath := fmt.Sprintf("services/%s/%s%s", category, serviceName, constants.ServiceConfigExtension)
	data, err := config.EmbeddedServicesFS.ReadFile(servicePath)
	if err != nil {
		return nil
	}

	var config types.ServiceConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil
	}

	return &config
}

func (sl *ServiceLoader) loadAllInCategory(category string) map[string]*types.ServiceConfig {
	categoryPath := fmt.Sprintf("services/%s", category)
	entries, err := config.EmbeddedServicesFS.ReadDir(categoryPath)
	if err != nil {
		return nil
	}

	result := make(map[string]*types.ServiceConfig)
	for _, entry := range entries {
		if !constants.IsYAMLFile(entry.Name()) {
			continue
		}

		serviceName := constants.TrimYAMLExt(entry.Name())
		if config := sl.loadServiceFromCategory(category, serviceName); config != nil {
			result[serviceName] = config
		}
	}

	return result
}

func (sl *ServiceLoader) parseServiceInfo(categoryPath, fileName, serviceName, category string) *types.ServiceInfo {
	serviceFile := fmt.Sprintf("%s/%s", categoryPath, fileName)
	data, err := config.EmbeddedServicesFS.ReadFile(serviceFile)
	if err != nil {
		return nil
	}

	var serviceData map[string]any
	if err := yaml.Unmarshal(data, &serviceData); err != nil {
		return nil
	}

	// Skip hidden services
	if visibility, _ := serviceData["visibility"].(string); visibility == "hidden" {
		return nil
	}

	return &types.ServiceInfo{
		Name:        serviceName,
		Category:    category,
		Description: getStringValue(serviceData, "description"),
		Type:        getStringValue(serviceData, "type"),
		Visibility:  getStringValue(serviceData, "visibility"),
	}
}

func getStringValue(data map[string]any, key string) string {
	if val, exists := data[key]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
