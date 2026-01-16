package display

import (
	"io"
	"sort"

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

	tw.AppendHeader(table.Row{HeaderService, HeaderCategory, HeaderDescription})
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
	tw := table.NewWriter()
	tw.SetOutputMirror(cf.writer)
	tw.SetStyle(table.StyleLight)

	tw.AppendHeader(table.Row{HeaderService, HeaderCategory, HeaderDescription})

	categoryNames := make([]string, 0, len(catalog.Categories))
	for category := range catalog.Categories {
		categoryNames = append(categoryNames, category)
	}
	sort.Strings(categoryNames)

	for i, category := range categoryNames {
		services := catalog.Categories[category]
		sort.Slice(services, func(a, b int) bool {
			return services[a].Name < services[b].Name
		})

		for _, service := range services {
			tw.AppendRow(table.Row{service.Name, category, service.Description})
		}

		if i < len(categoryNames)-1 {
			tw.AppendSeparator()
		}
	}

	tw.Render()
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
