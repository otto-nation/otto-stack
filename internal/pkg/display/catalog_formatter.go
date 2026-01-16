package display

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

// CatalogFormatter handles service catalog formatting
type CatalogFormatter struct {
	writer io.Writer
}

// NewCatalogFormatter creates a new catalog formatter
func NewCatalogFormatter(writer io.Writer) *CatalogFormatter {
	return &CatalogFormatter{writer: writer}
}

// FormatTable formats catalog as a table
func (cf *CatalogFormatter) FormatTable(catalog ServiceCatalog) error {
	tw := table.NewWriter()
	tw.SetOutputMirror(cf.writer)
	tw.SetStyle(table.StyleLight)

	tw.AppendHeader(table.Row{"Service", "Category", "Description"})
	for category, services := range catalog.Categories {
		for _, service := range services {
			tw.AppendRow(table.Row{service.Name, category, service.Description})
		}
	}

	tw.Render()
	return nil
}

// FormatGrouped formats catalog grouped by category
func (cf *CatalogFormatter) FormatGrouped(catalog ServiceCatalog) error {
	categoryNames := cf.getSortedCategoryNames(catalog)

	for i, category := range categoryNames {
		cf.formatCategory(i, category, catalog.Categories[category])
	}
	return nil
}

func (cf *CatalogFormatter) getSortedCategoryNames(catalog ServiceCatalog) []string {
	categoryNames := make([]string, 0, len(catalog.Categories))
	for category := range catalog.Categories {
		categoryNames = append(categoryNames, category)
	}
	sort.Strings(categoryNames)
	return categoryNames
}

func (cf *CatalogFormatter) formatCategory(index int, category string, services []ServiceInfo) {
	if index > 0 {
		_, _ = fmt.Fprintln(cf.writer)
	}

	cf.writeCategoryHeader(category)
	cf.writeSortedServices(services)
}

func (cf *CatalogFormatter) writeCategoryHeader(category string) {
	_, _ = fmt.Fprintf(cf.writer, "%s:\n", strings.ToUpper(category))
	_, _ = fmt.Fprintln(cf.writer, strings.Repeat("-", TableWidthCatalog))
}

func (cf *CatalogFormatter) writeSortedServices(services []ServiceInfo) {
	sortedServices := cf.sortServicesByName(services)
	for _, service := range sortedServices {
		_, _ = fmt.Fprintf(cf.writer, "  %-*s %s\n", ColWidthService, service.Name, service.Description)
	}
}

func (cf *CatalogFormatter) sortServicesByName(services []ServiceInfo) []ServiceInfo {
	sorted := make([]ServiceInfo, len(services))
	copy(sorted, services)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name < sorted[j].Name
	})
	return sorted
}

// FilterCatalogByCategory filters catalog by category
func FilterCatalogByCategory(catalog ServiceCatalog, category string) ServiceCatalog {
	filtered := ServiceCatalog{Categories: make(map[string][]ServiceInfo)}
	if services, exists := catalog.Categories[category]; exists {
		filtered.Categories[category] = services
		filtered.Total = len(services)
	}
	return filtered
}
