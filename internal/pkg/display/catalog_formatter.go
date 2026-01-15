package display

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// CatalogFormatter handles service catalog formatting
type CatalogFormatter struct {
	writer io.Writer
	table  *TableFormatter
}

// NewCatalogFormatter creates a new catalog formatter
func NewCatalogFormatter(writer io.Writer) *CatalogFormatter {
	return &CatalogFormatter{
		writer: writer,
		table:  NewTableFormatter(writer),
	}
}

// FormatTable formats catalog as a table
func (cf *CatalogFormatter) FormatTable(catalog ServiceCatalog) error {
	headers := []string{"Service", "Category", "Description"}
	widths := []int{ColWidthService, ColWidthCategory, ColWidthDescription}

	cf.table.WriteHeader(headers, widths)
	cf.writeTableRows(catalog, widths)
	return nil
}

func (cf *CatalogFormatter) writeTableRows(catalog ServiceCatalog, widths []int) {
	for category, services := range catalog.Categories {
		for _, service := range services {
			values := []string{service.Name, category, service.Description}
			cf.table.WriteRow(values, widths)
		}
	}
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
	cf.table.WriteSeparator(TableWidthCatalog)
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
