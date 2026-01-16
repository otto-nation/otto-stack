package display

import (
	"fmt"
	"io"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
)

// StatusFormatter handles service status formatting
type StatusFormatter struct {
	writer io.Writer
	table  *TableFormatter
}

// NewStatusFormatter creates a new status formatter
func NewStatusFormatter(writer io.Writer) *StatusFormatter {
	return &StatusFormatter{
		writer: writer,
		table:  NewTableFormatter(writer),
	}
}

// FormatTable formats services as a table
func (sf *StatusFormatter) FormatTable(services []ServiceStatus, options Options) error {
	if options.Compact {
		return sf.formatCompact(services)
	}
	return sf.formatFull(services, options)
}

func (sf *StatusFormatter) formatCompact(services []ServiceStatus) error {
	hasProvider := sf.hasProviders(services)
	headers, widths := sf.buildCompactLayout(hasProvider)

	sf.table.WriteHeader(headers, widths)
	for _, service := range services {
		values := sf.buildCompactValues(service, hasProvider)
		sf.table.WriteRow(values, widths)
	}
	return nil
}

func (sf *StatusFormatter) formatFull(services []ServiceStatus, options Options) error {
	hasProvider := sf.hasProviders(services)
	headers, widths := sf.buildFullLayout(hasProvider)

	sf.table.WriteHeader(headers, widths)
	for _, service := range services {
		values := sf.buildFullValues(service, hasProvider)
		sf.table.WriteRow(values, widths)
	}

	if options.ShowSummary {
		sf.formatResourceSummary(services)
	}
	return nil
}

func (sf *StatusFormatter) buildCompactLayout(hasProvider bool) ([]string, []int) {
	headers := []string{ui.StatusHeaderService}
	widths := []int{ColWidthServiceCompact}

	if hasProvider {
		headers = append(headers, ui.StatusHeaderProvidedBy)
		widths = append(widths, ColWidthProvider)
	}

	headers = append(headers, ui.StatusHeaderState, ui.StatusHeaderHealth)
	widths = append(widths, ColWidthState, ColWidthHealth)

	return headers, widths
}

func (sf *StatusFormatter) buildFullLayout(hasProvider bool) ([]string, []int) {
	headers := []string{ui.StatusHeaderService}
	widths := []int{ColWidthService}

	if hasProvider {
		headers = append(headers, ui.StatusHeaderProvidedBy)
		widths = append(widths, ColWidthProvider)
	}

	headers = append(headers, ui.StatusHeaderState, ui.StatusHeaderHealth, ui.StatusHeaderUptime, ui.StatusHeaderPorts, ui.StatusHeaderUpdated)
	widths = append(widths, ColWidthState, ColWidthHealth, ColWidthUptime, ColWidthPorts, ColWidthUpdated)

	return headers, widths
}

func (sf *StatusFormatter) buildCompactValues(service ServiceStatus, hasProvider bool) []string {
	values := []string{service.Name}

	if hasProvider {
		provider := service.Provider
		if provider == "" {
			provider = "n/a"
		}
		values = append(values, provider)
	}

	values = append(values,
		sf.getIcon(service.State)+service.State,
		sf.getIcon(service.Health)+service.Health,
	)

	return values
}

func (sf *StatusFormatter) buildFullValues(service ServiceStatus, hasProvider bool) []string {
	values := []string{service.Name}

	if hasProvider {
		provider := service.Provider
		if provider == "" {
			provider = "n/a"
		}
		values = append(values, provider)
	}

	values = append(values,
		sf.getIcon(service.State)+service.State,
		sf.getIcon(service.Health)+service.Health,
		sf.table.FormatDuration(service.Uptime),
		sf.table.FormatPorts(service.Ports),
		service.UpdatedAt.Format("15:04:05"),
	)

	return values
}

func (sf *StatusFormatter) hasProviders(services []ServiceStatus) bool {
	for _, service := range services {
		if service.Provider != "" {
			return true
		}
	}
	return false
}

func (sf *StatusFormatter) formatResourceSummary(services []ServiceStatus) {
	summary := sf.createSummary(services)
	_, _ = fmt.Fprintln(sf.writer)
	_, _ = fmt.Fprintf(sf.writer, "Summary: %d total", len(services))
	for state, count := range summary {
		if count > 0 {
			_, _ = fmt.Fprintf(sf.writer, ", %d %s", count, state)
		}
	}
	_, _ = fmt.Fprintln(sf.writer)
}

func (sf *StatusFormatter) createSummary(services []ServiceStatus) map[string]int {
	summary := make(map[string]int)
	for _, service := range services {
		summary[service.State]++
	}
	return summary
}

func (sf *StatusFormatter) getIcon(state string) string {
	switch state {
	case docker.HealthStatusRunning:
		return ui.IconSuccess + " "
	case docker.HealthStatusHealthy:
		return ui.IconHealthy + " "
	case docker.HealthStatusStopped:
		return ui.IconError + " "
	case docker.HealthUnhealthy:
		return ui.IconUnhealthy + " "
	case docker.HealthStarting:
		return ui.IconStarting + " "
	default:
		return ui.IconUnknown + " "
	}
}
