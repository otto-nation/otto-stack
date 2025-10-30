package display

import (
	"encoding/json"
	"io"
)

// JSONFormatter implements JSON output formatting
type JSONFormatter struct {
	writer io.Writer
}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter(writer io.Writer) *JSONFormatter {
	return &JSONFormatter{writer: writer}
}

// FormatStatus formats service status as JSON
func (f *JSONFormatter) FormatStatus(services []ServiceStatus, options StatusOptions) error {
	output := map[string]any{
		"services": services,
		"summary":  CreateSummary(services),
	}

	return f.writeJSON(output)
}

// FormatServiceCatalog formats service catalog as JSON
func (f *JSONFormatter) FormatServiceCatalog(catalog ServiceCatalog, options ServiceCatalogOptions) error {
	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")

	filteredCatalog := FilterCatalogByCategory(catalog, options.Category)
	return encoder.Encode(filteredCatalog)
}

// FormatValidation formats validation results as JSON
func (f *JSONFormatter) FormatValidation(result ValidationResult, options ValidationOptions) error {
	return f.writeJSON(result)
}

// FormatVersion formats version information as JSON
func (f *JSONFormatter) FormatVersion(info VersionInfo, options VersionOptions) error {
	return f.writeJSON(info)
}

// FormatHealth formats health check results as JSON
func (f *JSONFormatter) FormatHealth(report HealthReport, options HealthOptions) error {
	return f.writeJSON(report)
}

// Helper methods
func (f *JSONFormatter) writeJSON(data any) error {
	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
