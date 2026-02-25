package display

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
)

// StatusFormatter handles service status formatting
type StatusFormatter struct {
	writer  io.Writer
	noColor bool
}

// NewStatusFormatter creates a new status formatter
func NewStatusFormatter(writer io.Writer) *StatusFormatter {
	return &StatusFormatter{writer: writer}
}

// RenderStatusTable converts container statuses to display format and renders a table.
// It is a silent fallback — callers should treat any error as non-fatal.
func RenderStatusTable(
	writer io.Writer,
	containerStatuses []docker.ContainerStatus,
	serviceConfigs []types.ServiceConfig,
	compact bool,
	noColor bool,
) error {
	serviceToContainer := buildServiceContainerMap(serviceConfigs)
	serviceStatuses := convertToServiceStatuses(containerStatuses, serviceConfigs, serviceToContainer)
	sf := &StatusFormatter{writer: writer, noColor: noColor}
	return sf.FormatTable(serviceStatuses, Options{Compact: compact, NoColor: noColor})
}

// buildServiceContainerMap maps each service name to its container name
func buildServiceContainerMap(serviceConfigs []types.ServiceConfig) map[string]string {
	m := make(map[string]string, len(serviceConfigs))
	for _, cfg := range serviceConfigs {
		m[cfg.Name] = resolveContainerName(cfg)
	}
	return m
}

// resolveContainerName returns the container name for a service config
func resolveContainerName(cfg types.ServiceConfig) string {
	if cfg.Hidden {
		return cfg.Name
	}
	if len(cfg.Service.Dependencies.Required) > 0 {
		return cfg.Service.Dependencies.Required[0]
	}
	return cfg.Name
}

// convertToServiceStatuses converts container statuses into display statuses
func convertToServiceStatuses(
	containerStatuses []docker.ContainerStatus,
	serviceConfigs []types.ServiceConfig,
	serviceToContainer map[string]string,
) []ServiceStatus {
	containerMap := make(map[string]docker.ContainerStatus, len(containerStatuses))
	for _, cs := range containerStatuses {
		containerMap[cs.Name] = cs
	}

	result := make([]ServiceStatus, 0, len(serviceConfigs))
	for _, cfg := range serviceConfigs {
		if cfg.Container.Restart == types.RestartPolicyNo || cfg.Hidden {
			continue
		}
		result = append(result, buildServiceStatus(cfg, serviceToContainer, containerMap))
	}
	return result
}

func buildServiceStatus(cfg types.ServiceConfig, serviceToContainer map[string]string, containerMap map[string]docker.ContainerStatus) ServiceStatus {
	provider := serviceToContainer[cfg.Name]
	providerName := ""
	if provider != cfg.Name {
		providerName = provider
	}

	scope := ScopeLocal
	if cfg.Shareable {
		scope = ScopeShared
	}

	cs, exists := containerMap[provider]
	if !exists {
		return ServiceStatus{
			Name:     cfg.Name,
			Scope:    scope,
			Provider: providerName,
			State:    StateNotFound,
			Health:   StateUnknown,
		}
	}

	uptime := time.Duration(0)
	if !cs.StartedAt.IsZero() {
		uptime = time.Since(cs.StartedAt)
	}

	return ServiceStatus{
		Name:      cfg.Name,
		Scope:     scope,
		Container: cs.Name,
		Provider:  providerName,
		State:     cs.State,
		Health:    cs.Health,
		Ports:     cs.Ports,
		CreatedAt: cs.CreatedAt,
		UpdatedAt: cs.StartedAt,
		Uptime:    uptime,
	}
}

// FormatTable formats services as a table
func (sf *StatusFormatter) FormatTable(services []ServiceStatus, options Options) error {
	sf.noColor = options.NoColor
	if options.Compact {
		return sf.formatCompact(services)
	}
	return sf.formatFull(services, options)
}

func (sf *StatusFormatter) formatCompact(services []ServiceStatus) error {
	hasProvider := sf.hasProviders(services)
	tw := table.NewWriter()
	tw.SetOutputMirror(sf.writer)
	tw.SetStyle(tableStyle)

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
	tw.SetStyle(tableStyle)

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
	headers := table.Row{HeaderService, HeaderScope}
	if hasProvider {
		headers = append(headers, HeaderProvidedBy)
	}
	headers = append(headers, HeaderState, HeaderHealth)
	return headers
}

func (sf *StatusFormatter) buildFullHeaders(hasProvider bool) table.Row {
	headers := table.Row{HeaderService, HeaderScope, HeaderContainer}
	if hasProvider {
		headers = append(headers, HeaderProvidedBy)
	}
	headers = append(headers, HeaderState, HeaderHealth, HeaderUptime, HeaderPorts, HeaderUpdated)
	return headers
}

func (sf *StatusFormatter) buildCompactRow(service ServiceStatus, hasProvider bool) table.Row {
	row := table.Row{service.Name, service.Scope}
	if hasProvider {
		provider := service.Provider
		if provider == "" {
			provider = NotApplicable
		}
		row = append(row, provider)
	}
	stateText := sf.getIcon(service.State) + service.State
	healthText := sf.getIcon(service.Health) + service.Health
	row = append(row,
		sf.colorizeState(stateText, service.State),
		sf.colorizeState(healthText, service.Health),
	)
	return row
}

func (sf *StatusFormatter) buildFullRow(service ServiceStatus, hasProvider bool) table.Row {
	container := service.Container
	if container == "" {
		container = NotApplicable
	}

	row := table.Row{service.Name, service.Scope, container}
	if hasProvider {
		provider := service.Provider
		if provider == "" {
			provider = NotApplicable
		}
		row = append(row, provider)
	}
	stateText := sf.getIcon(service.State) + service.State
	healthText := sf.getIcon(service.Health) + service.Health
	row = append(row,
		sf.colorizeState(stateText, service.State),
		sf.colorizeState(healthText, service.Health),
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
	_, _ = fmt.Fprintf(sf.writer, SummaryTotal, len(services))
	for state, count := range summary {
		if count > 0 {
			_, _ = fmt.Fprintf(sf.writer, SummaryItem, count, state)
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
	case docker.HealthStatusRunning, docker.HealthStatusHealthy:
		return ui.IconOK + " "
	case docker.StateStopped, docker.HealthStatusStopped, docker.HealthUnhealthy:
		return ui.IconFail + " "
	case docker.HealthStarting:
		return ui.IconWarn + " "
	default:
		return ui.IconUnknown + " "
	}
}

// ColorizeState applies ANSI color to a state/health cell value.
// noColor must be true when stdout is not a terminal or color is disabled.
func ColorizeState(text, state string, noColor bool) string {
	if noColor {
		return text
	}
	switch state {
	case docker.HealthStatusRunning, docker.HealthStatusHealthy:
		return ui.ColorGreen + text + ui.ColorReset
	case docker.StateStopped, docker.HealthStatusStopped, docker.HealthUnhealthy:
		return ui.ColorRed + text + ui.ColorReset
	case docker.HealthStarting:
		return ui.ColorYellow + text + ui.ColorReset
	default:
		return ui.ColorGray + text + ui.ColorReset
	}
}

// colorizeState applies ANSI color to state/health cell text based on the underlying state value
func (sf *StatusFormatter) colorizeState(text, state string) string {
	return ColorizeState(text, state, sf.noColor)
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

// RenderSharedStatusTable renders shared container statuses as a table.
// compact=true shows SERVICE, STATE, HEALTH, USED BY.
// compact=false (verbose) adds CONTAINER, UPTIME, PORTS, UPDATED.
func RenderSharedStatusTable(writer io.Writer, statuses []SharedContainerStatus, compact bool, noColor bool) error {
	sf := &StatusFormatter{writer: writer, noColor: noColor}
	if compact {
		return sf.formatSharedCompact(statuses)
	}
	return sf.formatSharedFull(statuses)
}

func (sf *StatusFormatter) formatSharedCompact(statuses []SharedContainerStatus) error {
	tw := table.NewWriter()
	tw.SetOutputMirror(sf.writer)
	tw.SetStyle(tableStyle)
	tw.AppendHeader(table.Row{HeaderService, HeaderState, HeaderHealth, HeaderUsedBy})
	for _, s := range statuses {
		stateText := sf.getIcon(s.State) + s.State
		healthText := sf.getIcon(s.Health) + s.Health
		tw.AppendRow(table.Row{
			s.Service,
			sf.colorizeState(stateText, s.State),
			sf.colorizeState(healthText, s.Health),
			formatUsedBy(s.Projects),
		})
	}
	tw.Render()
	return nil
}

func (sf *StatusFormatter) formatSharedFull(statuses []SharedContainerStatus) error {
	tw := table.NewWriter()
	tw.SetOutputMirror(sf.writer)
	tw.SetStyle(tableStyle)
	tw.AppendHeader(table.Row{HeaderService, HeaderContainer, HeaderState, HeaderHealth, HeaderUptime, HeaderPorts, HeaderUpdated, HeaderUsedBy})
	for _, s := range statuses {
		stateText := sf.getIcon(s.State) + s.State
		healthText := sf.getIcon(s.Health) + s.Health
		tw.AppendRow(table.Row{
			s.Service,
			s.Name,
			sf.colorizeState(stateText, s.State),
			sf.colorizeState(healthText, s.Health),
			sf.formatDuration(s.Uptime),
			sf.formatPorts(s.Ports),
			s.UpdatedAt.Format("15:04:05"),
			formatUsedBy(s.Projects),
		})
	}
	tw.Render()
	return nil
}

// formatUsedBy formats project names for the USED BY column.
func formatUsedBy(projects []string) string {
	if len(projects) == 0 {
		return NotApplicable
	}
	const maxShow = 3
	if len(projects) > maxShow {
		const show = 2
		return fmt.Sprintf("%s, %s, +%d more", projects[0], projects[1], len(projects)-show)
	}
	return strings.Join(projects, ", ")
}
