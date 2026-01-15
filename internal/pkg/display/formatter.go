package display

import (
	"encoding/json"
	"io"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"gopkg.in/yaml.v3"
)

const (
	// Table formatting constants
	TableWidthCompact  = 42
	TableWidthStandard = 75
	TableWidthCatalog  = 80
	TableWidthHealth   = 90

	// Column widths
	ColWidthService        = 15
	ColWidthServiceCompact = 20
	ColWidthProvider       = 12
	ColWidthState          = 12
	ColWidthHealth         = 12
	ColWidthUptime         = 8
	ColWidthPorts          = 12
	ColWidthUpdated        = 10
	ColWidthCategory       = 15
	ColWidthDescription    = 20
	ColWidthCheck          = 25
	ColWidthMessage        = 40

	// Duration formatting
	SecondsPerMinute = 60
	MinutesPerHour   = 60
	HoursPerDay      = 24

	// Port display limits
	MaxPortsDisplay = 12
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
