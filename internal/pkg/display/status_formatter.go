package display

import (
	"fmt"
	"io"
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

	headers, widths := sf.getCompactLayout(hasProvider)
	sf.table.WriteHeader(headers, widths)

	for _, service := range services {
		values := sf.getCompactValues(service, hasProvider)
		sf.table.WriteRow(values, widths)
	}
	return nil
}

func (sf *StatusFormatter) hasProviders(services []ServiceStatus) bool {
	for _, service := range services {
		if service.Provider != "" {
			return true
		}
	}
	return false
}

func (sf *StatusFormatter) getCompactLayout(hasProvider bool) ([]string, []int) {
	if hasProvider {
		return []string{"Service", "Provided By", "State", "Health"},
			[]int{ColWidthServiceCompact, ColWidthProvider, ColWidthState, ColWidthHealth}
	}
	return []string{"Service", "State", "Health"},
		[]int{ColWidthServiceCompact, ColWidthState, ColWidthHealth}
}

func (sf *StatusFormatter) getCompactValues(service ServiceStatus, hasProvider bool) []string {
	baseValues := []string{
		service.Name,
		sf.getStateIcon(service.State) + service.State,
		sf.getHealthIcon(service.Health) + service.Health,
	}

	if hasProvider {
		provider := service.Provider
		if provider == "" {
			provider = "n/a"
		}
		return []string{service.Name, provider, baseValues[1], baseValues[2]}
	}
	return baseValues
}

func (sf *StatusFormatter) formatFull(services []ServiceStatus, options Options) error {
	headers := []string{"Service", "State", "Health", "Uptime", "Ports", "Updated"}
	widths := []int{ColWidthService, ColWidthState, ColWidthHealth, ColWidthUptime, ColWidthPorts, ColWidthUpdated}

	sf.table.WriteHeader(headers, widths)

	for _, service := range services {
		values := []string{
			service.Name,
			sf.getStateIcon(service.State) + service.State,
			sf.getHealthIcon(service.Health) + service.Health,
			sf.table.FormatDuration(service.Uptime),
			sf.table.FormatPorts(service.Ports),
			service.UpdatedAt.Format("15:04:05"),
		}
		sf.table.WriteRow(values, widths)
	}

	if options.ShowSummary {
		sf.formatResourceSummary(services)
	}
	return nil
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

func (sf *StatusFormatter) getStateIcon(state string) string {
	switch state {
	case "running":
		return "✓ "
	case "stopped":
		return "✗ "
	case "starting":
		return "⟳ "
	default:
		return "? "
	}
}

func (sf *StatusFormatter) getHealthIcon(health string) string {
	switch health {
	case "healthy":
		return "✓ "
	case "unhealthy":
		return "✗ "
	case "starting":
		return "⟳ "
	default:
		return "? "
	}
}
