package display

import (
	"fmt"
	"io"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
)

// GroupFormatter implements grouped output formatting for service catalogs
type GroupFormatter struct {
	writer io.Writer
	table  *TableFormatter // Embedded for delegation
}

// NewGroupFormatter creates a new group formatter
func NewGroupFormatter(writer io.Writer) *GroupFormatter {
	return &GroupFormatter{
		writer: writer,
		table:  NewTableFormatter(writer),
	}
}

// FormatServiceCatalog formats service catalog grouped by category
// CategoryFormatter handles formatting of service categories
type CategoryFormatter struct {
	writer io.Writer
}

// formatCategoryHeader formats the category header with service count
func (cf *CategoryFormatter) formatCategoryHeader(categoryName string, serviceCount int) {
	displayInfo, exists := constants.CategoryDisplayInfo[categoryName]
	if !exists {
		displayInfo = struct{ Name, Icon string }{categoryName, "📦"}
	}

	plural := ""
	if serviceCount != 1 {
		plural = "s"
	}

	_, _ = fmt.Fprintf(cf.writer, "%s %s\n", displayInfo.Icon,
		fmt.Sprintf(constants.MsgServiceCount, displayInfo.Name, serviceCount, plural))
}

// formatServiceList formats the list of services in a category
func (cf *CategoryFormatter) formatServiceList(services []ServiceInfo) {
	for _, service := range services {
		_, _ = fmt.Fprintf(cf.writer, "  %-15s %s\n", service.Name, service.Description)
	}
	_, _ = fmt.Fprintln(cf.writer)
}

// formatCategory formats a complete category with header and services
func (cf *CategoryFormatter) formatCategory(categoryName string, services []ServiceInfo) {
	if len(services) == 0 {
		return
	}

	cf.formatCategoryHeader(categoryName, len(services))
	cf.formatServiceList(services)
}

func (f *GroupFormatter) FormatServiceCatalog(catalog ServiceCatalog, options ServiceCatalogOptions) error {
	filteredCatalog := FilterCatalogByCategory(catalog, options.Category)

	if filteredCatalog.Total == 0 {
		if options.Category != "" {
			_, _ = fmt.Fprintf(f.writer, constants.MsgNoServicesInCategory+"\n", options.Category)
		} else {
			_, _ = fmt.Fprintln(f.writer, "No services available")
		}
		return nil
	}

	_, _ = fmt.Fprintln(f.writer, constants.MsgServiceCatalogHeader)
	_, _ = fmt.Fprintln(f.writer)

	categoryFormatter := &CategoryFormatter{writer: f.writer}
	for categoryName, services := range filteredCatalog.Categories {
		categoryFormatter.formatCategory(categoryName, services)
	}

	return nil
}

// Delegate other methods to embedded table formatter
func (f *GroupFormatter) FormatStatus(services []ServiceStatus, options StatusOptions) error {
	return f.table.FormatStatus(services, options)
}

func (f *GroupFormatter) FormatValidation(result ValidationResult, options ValidationOptions) error {
	return f.table.FormatValidation(result, options)
}

func (f *GroupFormatter) FormatVersion(info VersionInfo, options VersionOptions) error {
	return f.table.FormatVersion(info, options)
}

func (f *GroupFormatter) FormatHealth(report HealthReport, options HealthOptions) error {
	return f.table.FormatHealth(report, options)
}
