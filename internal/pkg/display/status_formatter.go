package display

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
)

// StatusFormatter handles service status formatting
type StatusFormatter struct {
	writer io.Writer
}

// NewStatusFormatter creates a new status formatter
func NewStatusFormatter(writer io.Writer) *StatusFormatter {
	return &StatusFormatter{writer: writer}
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
	tw := table.NewWriter()
	tw.SetOutputMirror(sf.writer)
	tw.SetStyle(table.StyleLight)

	headers := sf.buildCompactHeaders(hasProvider)
	tw.AppendHeader(headers)

	for _, service := range services {
		tw.AppendRow(sf.buildCompactRow(service, hasProvider))
	}

	tw.Render()
	return nil
}

func (sf *StatusFormatter) formatFull(services []ServiceStatus, options Options) error {
	hasProvider := sf.hasProviders(services)
	tw := table.NewWriter()
	tw.SetOutputMirror(sf.writer)
	tw.SetStyle(table.StyleLight)

	headers := sf.buildFullHeaders(hasProvider)
	tw.AppendHeader(headers)

	for _, service := range services {
		tw.AppendRow(sf.buildFullRow(service, hasProvider))
	}

	tw.Render()

	if options.ShowSummary {
		sf.formatResourceSummary(services)
	}
	return nil
}

func (sf *StatusFormatter) buildCompactHeaders(hasProvider bool) table.Row {
	headers := table.Row{HeaderService}
	if hasProvider {
		headers = append(headers, HeaderProvidedBy)
	}
	headers = append(headers, HeaderState, HeaderHealth)
	return headers
}

func (sf *StatusFormatter) buildFullHeaders(hasProvider bool) table.Row {
	headers := table.Row{HeaderService}
	if hasProvider {
		headers = append(headers, HeaderProvidedBy)
	}
	headers = append(headers, HeaderState, HeaderHealth, HeaderUptime, HeaderPorts, HeaderUpdated)
	return headers
}

func (sf *StatusFormatter) buildCompactRow(service ServiceStatus, hasProvider bool) table.Row {
	row := table.Row{service.Name}
	if hasProvider {
		provider := service.Provider
		if provider == "" {
			provider = "n/a"
		}
		row = append(row, provider)
	}
	row = append(row,
		sf.getIcon(service.State)+service.State,
		sf.getIcon(service.Health)+service.Health,
	)
	return row
}

func (sf *StatusFormatter) buildFullRow(service ServiceStatus, hasProvider bool) table.Row {
	row := table.Row{service.Name}
	if hasProvider {
		provider := service.Provider
		if provider == "" {
			provider = "n/a"
		}
		row = append(row, provider)
	}
	row = append(row,
		sf.getIcon(service.State)+service.State,
		sf.getIcon(service.Health)+service.Health,
		sf.formatDuration(service.Uptime),
		sf.formatPorts(service.Ports),
		service.UpdatedAt.Format("15:04:05"),
	)
	return row
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

func (sf *StatusFormatter) formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < HoursPerDay*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours())/HoursPerDay)
}

func (sf *StatusFormatter) formatPorts(ports []string) string {
	if len(ports) == 0 {
		return "-"
	}
	if len(ports) <= MaxPortsDisplay {
		return strings.Join(ports, ",")
	}
	return strings.Join(ports[:MaxPortsDisplay], ",") + "..."
}
