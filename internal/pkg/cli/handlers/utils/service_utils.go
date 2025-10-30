package utils

import (
	"fmt"
	"strings"

	"github.com/otto-nation/otto-stack/internal/config"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/types"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"gopkg.in/yaml.v3"
)

// ServiceUtils provides shared utilities for service operations
type ServiceUtils struct{}

// NewServiceUtils creates a new service utilities instance
func NewServiceUtils() *ServiceUtils {
	return &ServiceUtils{}
}

// GetServicesByCategory loads services organized by category
func (u *ServiceUtils) GetServicesByCategory() (map[string][]types.ServiceInfo, error) {
	categories, err := u.getCategories()
	if err != nil {
		return nil, err
	}

	result := make(map[string][]types.ServiceInfo)
	for _, category := range categories {
		services, err := u.getServicesInCategory(category)
		if err != nil {
			continue
		}
		if len(services) > 0 {
			result[category] = services
		}
	}

	return result, nil
}

// LoadServicesByCategory is an alias for backward compatibility
func (u *ServiceUtils) LoadServicesByCategory() (map[string][]types.ServiceInfo, error) {
	return u.GetServicesByCategory()
}

// LoadServiceConfig loads a service configuration
func (u *ServiceUtils) LoadServiceConfig(serviceName string) (*types.ServiceConfig, error) {
	categories, err := u.getCategories()
	if err != nil {
		return nil, err
	}

	for _, category := range categories {
		config, err := u.loadServiceFromCategory(category, serviceName)
		if err == nil {
			return config, nil
		}
	}

	return nil, fmt.Errorf("service %s not found", serviceName)
}

// LoadAllServiceDependencies loads dependencies for all services
func (u *ServiceUtils) LoadAllServiceDependencies() (map[string][]string, error) {
	categories, err := u.getCategories()
	if err != nil {
		return nil, err
	}

	result := make(map[string][]string)
	for _, category := range categories {
		deps, err := u.getDependenciesInCategory(category)
		if err != nil {
			continue
		}
		for service, serviceDeps := range deps {
			result[service] = serviceDeps
		}
	}

	return result, nil
}

// ResolveDependencies resolves service dependencies and returns ordered list
func (u *ServiceUtils) ResolveDependencies(selectedServices []string) ([]string, error) {
	serviceMap, err := u.LoadAllServiceDependencies()
	if err != nil {
		return selectedServices, err
	}

	// Pre-expand all composite dependencies
	expandedServiceMap := make(map[string][]string)
	for service, deps := range serviceMap {
		expandedDeps, err := u.ExpandCompositeServices(deps)
		if err != nil {
			return selectedServices, err
		}
		expandedServiceMap[service] = expandedDeps
	}

	visited := make(map[string]bool)
	visiting := make(map[string]bool)
	var result []string

	var visit func(string) error
	visit = func(serviceName string) error {
		if visiting[serviceName] {
			return fmt.Errorf("circular dependency detected: %s", serviceName)
		}
		if visited[serviceName] {
			return nil
		}

		visiting[serviceName] = true
		for _, dep := range expandedServiceMap[serviceName] {
			if err := visit(dep); err != nil {
				return err
			}
		}
		visiting[serviceName] = false
		visited[serviceName] = true

		// Only add container services to result
		serviceConfig, err := u.LoadServiceConfig(serviceName)
		if err == nil && (serviceConfig.Type == "configuration" || serviceConfig.Type == "composite") {
			return nil
		}

		result = append(result, serviceName)
		return nil
	}

	for _, service := range selectedServices {
		if err := visit(service); err != nil {
			return selectedServices, err
		}
	}

	return result, nil
}

// ExpandCompositeServices expands composite services to their component services
func (u *ServiceUtils) ExpandCompositeServices(selectedServices []string) ([]string, error) {
	var expandedServices []string

	for _, serviceName := range selectedServices {
		serviceConfig, err := u.LoadServiceConfig(serviceName)
		if err != nil {
			// If service not found, keep as-is (might be a direct service name)
			expandedServices = append(expandedServices, serviceName)
			continue
		}

		if serviceConfig.Type == "composite" && len(serviceConfig.Components) > 0 {
			// Expand composite service to its components
			expandedServices = append(expandedServices, serviceConfig.Components...)
		} else {
			// Regular service, keep as-is
			expandedServices = append(expandedServices, serviceName)
		}
	}

	return expandedServices, nil
}

// ResolveServices applies composite expansion and dependency resolution
func (u *ServiceUtils) ResolveServices(serviceNames []string) ([]string, error) {
	expandedServices, err := u.ExpandCompositeServices(serviceNames)
	if err != nil {
		return serviceNames, fmt.Errorf("failed to expand composite services: %w", err)
	}

	resolvedServices, err := u.ResolveDependencies(expandedServices)
	if err != nil {
		return serviceNames, fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	return resolvedServices, nil
}

// Helper methods
func (u *ServiceUtils) getCategories() ([]string, error) {
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

func (u *ServiceUtils) getServicesInCategory(category string) ([]types.ServiceInfo, error) {
	categoryPath := fmt.Sprintf("services/%s", category)
	entries, err := config.EmbeddedServicesFS.ReadDir(categoryPath)
	if err != nil {
		return nil, err
	}

	var services []types.ServiceInfo
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), constants.ServiceConfigExtension) {
			continue
		}

		serviceName := strings.TrimSuffix(entry.Name(), constants.ServiceConfigExtension)
		serviceInfo, err := u.parseServiceInfo(categoryPath, entry.Name(), serviceName, category)
		if err != nil {
			continue
		}

		// Filter out hidden services from interactive selection
		if serviceInfo.Visibility == "hidden" {
			continue
		}

		services = append(services, serviceInfo)
	}

	return services, nil
}

func (u *ServiceUtils) parseServiceInfo(categoryPath, fileName, serviceName, category string) (types.ServiceInfo, error) {
	serviceFile := fmt.Sprintf("%s/%s", categoryPath, fileName)
	data, err := config.EmbeddedServicesFS.ReadFile(serviceFile)
	if err != nil {
		return types.ServiceInfo{}, err
	}

	var serviceData map[string]any
	if err := yaml.Unmarshal(data, &serviceData); err != nil {
		return types.ServiceInfo{}, err
	}

	return types.ServiceInfo{
		Name:                 serviceName,
		Category:             category,
		Description:          getString(serviceData, "description"),
		Type:                 getString(serviceData, "type"),
		Visibility:           getString(serviceData, "visibility"),
		Components:           getStringSlice(serviceData["components"]),
		Dependencies:         getDependencies(serviceData),
		ServiceConfiguration: convertServiceConfiguration(serviceData),
		Documentation:        parseDocumentation(serviceData),
	}, nil
}

func (u *ServiceUtils) loadServiceFromCategory(category, serviceName string) (*types.ServiceConfig, error) {
	servicePath := fmt.Sprintf("services/%s/%s%s", category, serviceName, constants.ServiceConfigExtension)
	data, err := config.EmbeddedServicesFS.ReadFile(servicePath)
	if err != nil {
		return nil, err
	}

	var serviceConfig types.ServiceConfig
	if err := yaml.Unmarshal(data, &serviceConfig); err != nil {
		return nil, fmt.Errorf("failed to parse service config for %s: %w", serviceName, err)
	}

	return &serviceConfig, nil
}

func (u *ServiceUtils) getDependenciesInCategory(category string) (map[string][]string, error) {
	categoryPath := fmt.Sprintf("services/%s", category)
	entries, err := config.EmbeddedServicesFS.ReadDir(categoryPath)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]string)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), constants.ServiceConfigExtension) {
			continue
		}

		serviceName := strings.TrimSuffix(entry.Name(), constants.ServiceConfigExtension)
		deps, err := u.parseServiceDependencies(categoryPath, entry.Name())
		if err != nil {
			continue
		}
		result[serviceName] = deps
	}

	return result, nil
}

func (u *ServiceUtils) parseServiceDependencies(categoryPath, fileName string) ([]string, error) {
	serviceFile := fmt.Sprintf("%s/%s", categoryPath, fileName)
	data, err := config.EmbeddedServicesFS.ReadFile(serviceFile)
	if err != nil {
		return nil, err
	}

	var serviceData map[string]any
	if err := yaml.Unmarshal(data, &serviceData); err != nil {
		return nil, err
	}

	return getDependencies(serviceData), nil
}

// parseDocumentation parses documentation section from service data
func parseDocumentation(serviceData map[string]any) types.ServiceDocumentation {
	doc := types.ServiceDocumentation{}

	// Check for new documentation structure
	if docData, exists := serviceData["documentation"]; exists {
		if docMap, ok := docData.(map[string]any); ok {
			doc.Examples = getStringSlice(docMap["examples"])
			doc.UsageNotes = getString(docMap, "usage_notes")
			doc.Links = getStringSlice(docMap["links"])
			doc.UseCases = getStringSlice(docMap["use_cases"])
			doc.WebInterfaces = parseWebInterfaces(docMap["web_interfaces"])
			return doc
		}
	}

	// Fallback to old structure for backward compatibility
	doc.Examples = getStringSlice(serviceData["examples"])
	doc.UsageNotes = getString(serviceData, "usage_notes")
	doc.Links = getStringSlice(serviceData["links"])
	doc.UseCases = getStringSlice(serviceData["use_cases"])
	doc.WebInterfaces = parseWebInterfaces(serviceData["web_interfaces"])

	return doc
}

// parseWebInterfaces parses web interfaces from service data
func parseWebInterfaces(data any) []types.WebInterface {
	if data == nil {
		return nil
	}

	if interfaces, ok := data.([]any); ok {
		var result []types.WebInterface
		for _, item := range interfaces {
			if interfaceMap, ok := item.(map[string]any); ok {
				result = append(result, types.WebInterface{
					Name:        getString(interfaceMap, "name"),
					URL:         getString(interfaceMap, "url"),
					Description: getString(interfaceMap, "description"),
				})
			}
		}
		return result
	}

	return nil
}

// Helper functions
func getString(data map[string]any, key string) string {
	if val, exists := data[key]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getStringSlice(val any) []string {
	if val == nil {
		return nil
	}

	if slice, ok := val.([]any); ok {
		var result []string
		for _, item := range slice {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
		return result
	}

	return nil
}

func getDependencies(serviceData map[string]any) []string {
	deps, exists := serviceData["dependencies"]
	if !exists {
		return nil
	}

	depsMap, ok := deps.(map[string]any)
	if !ok {
		return nil
	}

	required, exists := depsMap["required"]
	if !exists {
		return nil
	}

	return getStringSlice(required)
}

// convertServiceConfiguration converts service configuration from YAML to CLI ServiceOption format
func convertServiceConfiguration(serviceData map[string]any) []types.ServiceOption {
	// Check for new service_configuration key first
	if data := serviceData["service_configuration"]; data != nil {
		return convertOptionsData(data)
	}

	// Fallback to old options key for backward compatibility
	if data := serviceData["options"]; data != nil {
		return convertOptionsData(data)
	}

	return nil
}

// convertOptionsData converts options from YAML to CLI ServiceOption format
func convertOptionsData(data any) []types.ServiceOption {
	if data == nil {
		return nil
	}

	// Handle both old string slice format and new structured format
	switch v := data.(type) {
	case []any:
		var options []types.ServiceOption
		for _, item := range v {
			if optMap, ok := item.(map[any]any); ok {
				// New structured format
				option := types.ServiceOption{
					Name:        getStringFromInterface(optMap, "name"),
					Type:        getStringFromInterface(optMap, "type"),
					Description: getStringFromInterface(optMap, "description"),
					Default:     getStringFromInterface(optMap, "default"),
					Example:     getStringFromInterface(optMap, "example"),
					Required:    getBool(optMap, "required"),
					Values:      getStringSlice(optMap["values"]),
				}
				options = append(options, option)
			} else if str, ok := item.(string); ok {
				// Old string format - convert to structured
				option := types.ServiceOption{
					Name:        str,
					Type:        "string",
					Description: fmt.Sprintf("Configuration option: %s", str),
					Required:    false,
				}
				options = append(options, option)
			}
		}
		return options
	default:
		return nil
	}
}

// getStringFromInterface safely extracts a string value from any map
func getStringFromInterface(data map[any]any, key string) string {
	if val, exists := data[key]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// getBool safely extracts a boolean value from a map
func getBool(data map[any]any, key string) bool {
	if val, exists := data[key]; exists {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}
