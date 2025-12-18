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

	// Flatten categories into services list
	for category, services := range catalog.Categories {
		for _, service := range services {
			values := []string{
				service.Name,
				category,
				service.Description,
			}
			cf.table.WriteRow(values, widths)
		}
	}
	return nil
}

// FormatGrouped formats catalog grouped by category
func (cf *CatalogFormatter) FormatGrouped(catalog ServiceCatalog) error {
	// Sort categories
	var categoryNames []string
	for category := range catalog.Categories {
		categoryNames = append(categoryNames, category)
	}
	sort.Strings(categoryNames)

	for i, category := range categoryNames {
		if i > 0 {
			_, _ = fmt.Fprintln(cf.writer)
		}

		_, _ = fmt.Fprintf(cf.writer, "%s:\n", strings.ToUpper(category))
		cf.table.WriteSeparator(TableWidthCatalog)

		// Sort services within category
		services := catalog.Categories[category]
		sort.Slice(services, func(i, j int) bool {
			return services[i].Name < services[j].Name
		})

		for _, service := range services {
			_, _ = fmt.Fprintf(cf.writer, "  %-*s %s\n", ColWidthService, service.Name, service.Description)
		}
	}
	return nil
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
