package display

import (
	"encoding/json"
	"io"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"gopkg.in/yaml.v3"
)

const (
	// Duration formatting
	HoursPerDay = 24

	// Port display limits
	MaxPortsDisplay = 12

	// Table headers - Status
	HeaderService    = "SERVICE"
	HeaderProvidedBy = "PROVIDED BY"
	HeaderState      = "STATE"
	HeaderHealth     = "HEALTH"
	HeaderUptime     = "UPTIME"
	HeaderPorts      = "PORTS"
	HeaderUpdated    = "UPDATED"

	// Table headers - Catalog
	HeaderCategory    = "CATEGORY"
	HeaderDescription = "DESCRIPTION"

	// Table headers - Validation
	HeaderType    = "TYPE"
	HeaderField   = "FIELD"
	HeaderMessage = "MESSAGE"

	// Table headers - Health
	HeaderCheck  = "CHECK"
	HeaderStatus = "STATUS"

	// Table headers - Version
	HeaderComponent = "COMPONENT"
	HeaderVersion   = "VERSION"
	HeaderPlatform  = "PLATFORM"

	// Table headers - Web Interfaces
	HeaderInterface = "INTERFACE"
	HeaderURL       = "URL"

	// Table headers - Dependencies
	HeaderDependencies = "DEPENDENCIES"
)

// FormatHandler defines the interface for format-specific handlers
type FormatHandler interface {
	Handle(data any) error
}

// Formatter handles output formatting for different data types
type Formatter struct {
	writer   io.Writer
	handlers map[string]FormatHandler

	// Component formatters
	status     *StatusFormatter
	catalog    *CatalogFormatter
	validation *ValidationFormatter
	health     *HealthFormatter
	version    *VersionFormatter
}

// New creates a new formatter instance
func New(writer io.Writer, output base.Output) *Formatter {
	f := &Formatter{
		writer:     writer,
		handlers:   make(map[string]FormatHandler),
		status:     NewStatusFormatter(writer),
		catalog:    NewCatalogFormatter(writer),
		validation: NewValidationFormatter(writer),
		health:     NewHealthFormatter(writer),
		version:    NewVersionFormatter(writer),
	}
	f.initHandlers()
	return f
}

func (f *Formatter) initHandlers() {
	f.handlers["json"] = &JSONHandler{writer: f.writer}
	f.handlers["yaml"] = &YAMLHandler{writer: f.writer}
}

// FormatStatus formats service status information
func (f *Formatter) FormatStatus(services []ServiceStatus, options Options) error {
	if handler, exists := f.handlers[options.Format]; exists {
		return handler.Handle(services)
	}
	return f.status.FormatTable(services, options)
}

// FormatServiceCatalog formats service catalog information
func (f *Formatter) FormatServiceCatalog(catalog ServiceCatalog, options Options) error {
	if handler, exists := f.handlers[options.Format]; exists {
		return handler.Handle(catalog)
	}

	if options.GroupByCategory {
		return f.catalog.FormatGrouped(catalog)
	}
	return f.catalog.FormatTable(catalog)
}

// FormatValidation formats validation results
func (f *Formatter) FormatValidation(result ValidationResult, options Options) error {
	if handler, exists := f.handlers[options.Format]; exists {
		return handler.Handle(result)
	}
	return f.validation.FormatTable(result)
}

// FormatVersion formats version information
func (f *Formatter) FormatVersion(info VersionInfo, options Options) error {
	if handler, exists := f.handlers[options.Format]; exists {
		return handler.Handle(info)
	}
	return f.version.FormatTable(info, options)
}

// FormatHealth formats health report
func (f *Formatter) FormatHealth(report HealthReport, options Options) error {
	if handler, exists := f.handlers[options.Format]; exists {
		return handler.Handle(report)
	}
	return f.health.FormatTable(report, options)
}

// JSONHandler handles JSON output formatting
type JSONHandler struct {
	writer io.Writer
}

func (h *JSONHandler) Handle(data any) error {
	encoder := json.NewEncoder(h.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// YAMLHandler handles YAML output formatting
type YAMLHandler struct {
	writer io.Writer
}

func (h *YAMLHandler) Handle(data any) error {
	encoder := yaml.NewEncoder(h.writer)
	defer func() { _ = encoder.Close() }()
	return encoder.Encode(data)
}
