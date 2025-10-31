package display

import "github.com/otto-nation/otto-stack/internal/pkg/constants"

// CountByState counts services with the specified state
func CountByState(services []ServiceStatus, state string) int {
	count := 0
	for _, service := range services {
		if service.State == state {
			count++
		}
	}
	return count
}

// CountByHealth counts services with the specified health status
func CountByHealth(services []ServiceStatus, health string) int {
	count := 0
	for _, service := range services {
		if service.Health == health {
			count++
		}
	}
	return count
}

// FilterCatalogByCategory filters service catalog by category
func FilterCatalogByCategory(catalog ServiceCatalog, category string) ServiceCatalog {
	if category == "" {
		return catalog
	}

	if services, exists := catalog.Categories[category]; exists {
		return ServiceCatalog{
			Categories: map[string][]ServiceInfo{category: services},
			Total:      len(services),
		}
	}

	// Return empty catalog for non-existent category
	return ServiceCatalog{
		Categories: make(map[string][]ServiceInfo),
		Total:      0,
	}
}

// CreateSummary creates a summary map for service status
func CreateSummary(services []ServiceStatus) map[string]any {
	return map[string]any{
		"total":   len(services),
		"running": CountByState(services, constants.StateRunning),
		"healthy": CountByHealth(services, constants.HealthHealthy),
	}
}
