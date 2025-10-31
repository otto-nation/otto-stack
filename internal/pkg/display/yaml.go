package display

import (
	"io"

	"gopkg.in/yaml.v3"
)

// YAMLFormatter implements YAML output formatting
type YAMLFormatter struct {
	writer io.Writer
}

// NewYAMLFormatter creates a new YAML formatter
func NewYAMLFormatter(writer io.Writer) *YAMLFormatter {
	return &YAMLFormatter{writer: writer}
}

// FormatStatus formats service status as YAML
func (f *YAMLFormatter) FormatStatus(services []ServiceStatus, options StatusOptions) error {
	output := map[string]any{
		"services": services,
		"summary": map[string]any{
			"total":   len(services),
			"running": f.countByState(services, "running"),
			"healthy": f.countByHealth(services, "healthy"),
		},
	}

	return f.writeYAML(output)
}

// FormatServiceCatalog formats service catalog as YAML
func (f *YAMLFormatter) FormatServiceCatalog(catalog ServiceCatalog, options ServiceCatalogOptions) error {
	encoder := yaml.NewEncoder(f.writer)
	defer func() { _ = encoder.Close() }()

	// Filter by category if specified
	if options.Category != "" {
		if services, exists := catalog.Categories[options.Category]; exists {
			filteredCatalog := ServiceCatalog{
				Categories: map[string][]ServiceInfo{options.Category: services},
				Total:      len(services),
			}
			return encoder.Encode(filteredCatalog)
		}
		// Return empty catalog for non-existent category
		emptyCatalog := ServiceCatalog{Categories: make(map[string][]ServiceInfo), Total: 0}
		return encoder.Encode(emptyCatalog)
	}

	return encoder.Encode(catalog)
}

// FormatValidation formats validation results as YAML
func (f *YAMLFormatter) FormatValidation(result ValidationResult, options ValidationOptions) error {
	return f.writeYAML(result)
}

// FormatVersion formats version information as YAML
func (f *YAMLFormatter) FormatVersion(info VersionInfo, options VersionOptions) error {
	return f.writeYAML(info)
}

// FormatHealth formats health check results as YAML
func (f *YAMLFormatter) FormatHealth(report HealthReport, options HealthOptions) error {
	return f.writeYAML(report)
}

// Helper methods
func (f *YAMLFormatter) writeYAML(data any) error {
	encoder := yaml.NewEncoder(f.writer)
	defer func() {
		_ = encoder.Close()
	}()
	return encoder.Encode(data)
}

func (f *YAMLFormatter) countByState(services []ServiceStatus, state string) int {
	count := 0
	for _, service := range services {
		if service.State == state {
			count++
		}
	}
	return count
}

func (f *YAMLFormatter) countByHealth(services []ServiceStatus, health string) int {
	count := 0
	for _, service := range services {
		if service.Health == health {
			count++
		}
	}
	return count
}
