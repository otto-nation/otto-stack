package display

import (
	"fmt"
	"io"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
)

// GroupFormatter implements grouped output formatting for service catalogs
type GroupFormatter struct {
	writer io.Writer
}

// NewGroupFormatter creates a new group formatter
func NewGroupFormatter(writer io.Writer) *GroupFormatter {
	return &GroupFormatter{writer: writer}
}

// FormatServiceCatalog formats service catalog grouped by category
func (f *GroupFormatter) FormatServiceCatalog(catalog ServiceCatalog, options ServiceCatalogOptions) error {
	if catalog.Total == 0 {
		_, _ = fmt.Fprintln(f.writer, "No services available")
		return nil
	}

	// Filter by category if specified
	categories := catalog.Categories
	if options.Category != "" {
		if services, exists := catalog.Categories[options.Category]; exists {
			categories = map[string][]ServiceInfo{options.Category: services}
		} else {
			_, _ = fmt.Fprintf(f.writer, constants.MsgNoServicesInCategory+"\n", options.Category)
			return nil
		}
	}

	_, _ = fmt.Fprintln(f.writer, constants.MsgServiceCatalogHeader)
	_, _ = fmt.Fprintln(f.writer)

	for categoryName, services := range categories {
		if len(services) == 0 {
			continue
		}

		// Get category display info
		displayInfo, exists := constants.CategoryDisplayInfo[categoryName]
		if !exists {
			displayInfo = struct{ Name, Icon string }{categoryName, "📦"}
		}

		// Category header with count
		serviceCount := len(services)
		plural := ""
		if serviceCount != 1 {
			plural = "s"
		}
		_, _ = fmt.Fprintf(f.writer, "%s %s\n", displayInfo.Icon,
			fmt.Sprintf(constants.MsgServiceCount, displayInfo.Name, serviceCount, plural))

		// List services in category
		for _, service := range services {
			_, _ = fmt.Fprintf(f.writer, "  %-15s %s\n", service.Name, service.Description)
		}
		_, _ = fmt.Fprintln(f.writer)
	}

	return nil
}

// Implement other required methods (delegating to table formatter for now)
func (f *GroupFormatter) FormatStatus(services []ServiceStatus, options StatusOptions) error {
	tableFormatter := NewTableFormatter(f.writer)
	return tableFormatter.FormatStatus(services, options)
}

func (f *GroupFormatter) FormatValidation(result ValidationResult, options ValidationOptions) error {
	tableFormatter := NewTableFormatter(f.writer)
	return tableFormatter.FormatValidation(result, options)
}

func (f *GroupFormatter) FormatVersion(info VersionInfo, options VersionOptions) error {
	tableFormatter := NewTableFormatter(f.writer)
	return tableFormatter.FormatVersion(info, options)
}

func (f *GroupFormatter) FormatHealth(report HealthReport, options HealthOptions) error {
	tableFormatter := NewTableFormatter(f.writer)
	return tableFormatter.FormatHealth(report, options)
}
