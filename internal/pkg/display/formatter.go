package display

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	pkgServices "github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"gopkg.in/yaml.v3"
)

const (
	// Table formatting constants
	TableWidthCompact  = 42
	TableWidthStandard = 75
	TableWidthCatalog  = 80
	TableWidthHealth   = 90

	// Column widths
	ColWidthService        = 15
	ColWidthServiceCompact = 20
	ColWidthState          = 10
	ColWidthHealth         = 10
	ColWidthUptime         = 8
	ColWidthPorts          = 12
	ColWidthUpdated        = 10
	ColWidthCategory       = 15
	ColWidthDescription    = 20
	ColWidthCheck          = 25
	ColWidthMessage        = 40

	// Duration formatting
	SecondsPerMinute = 60
	MinutesPerHour   = 60
	HoursPerDay      = 24

	// Port display limits
	MaxPortsDisplay = 12
)

// FormatHandler defines the interface for format-specific handlers
type FormatHandler interface {
	Handle(data any) error
}

// Formatter handles all output formatting
type Formatter struct {
	writer   io.Writer
	output   base.Output
	handlers map[string]FormatHandler
}

// New creates a new formatter
func New(writer io.Writer, output base.Output) *Formatter {
	f := &Formatter{
		writer: writer,
		output: output,
	}
	f.initHandlers()
	return f
}

// initHandlers initializes format handlers
func (f *Formatter) initHandlers() {
	f.handlers = map[string]FormatHandler{
		pkgServices.ServiceCatalogJSONFormat: &JSONHandler{writer: f.writer},
		pkgServices.ServiceCatalogYAMLFormat: &YAMLHandler{writer: f.writer},
	}
}

// FormatStatus formats service status information
func (f *Formatter) FormatStatus(services []ServiceStatus, options Options) error {
	if handler, exists := f.handlers[options.Format]; exists {
		return handler.Handle(map[string]any{
			"services": services,
			"summary":  f.createSummary(services),
		})
	}
	return f.formatStatusTable(services, options)
}

// FormatServiceCatalog formats service catalog information
func (f *Formatter) FormatServiceCatalog(catalog ServiceCatalog, options Options) error {
	if handler, exists := f.handlers[options.Format]; exists {
		return handler.Handle(catalog)
	}

	switch options.Format {
	case pkgServices.ServiceCatalogTableFormat:
		return f.formatCatalogTable(catalog)
	default:
		return f.formatCatalogGroup(catalog)
	}
}

// FormatValidation formats validation results
func (f *Formatter) FormatValidation(result ValidationResult, options Options) error {
	if handler, exists := f.handlers[options.Format]; exists {
		return handler.Handle(result)
	}
	return f.formatValidationTable(result)
}

// FormatVersion formats version information
func (f *Formatter) FormatVersion(info VersionInfo, options Options) error {
	if handler, exists := f.handlers[options.Format]; exists {
		return handler.Handle(info)
	}
	return f.formatVersionTable(info, options)
}

// FormatHealth formats health check results
func (f *Formatter) FormatHealth(report HealthReport, options Options) error {
	if handler, exists := f.handlers[options.Format]; exists {
		return handler.Handle(report)
	}
	return f.formatHealthTable(report, options)
}

// JSONHandler handles JSON formatting
type JSONHandler struct {
	writer io.Writer
}

func (h *JSONHandler) Handle(data any) error {
	encoder := json.NewEncoder(h.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// YAMLHandler handles YAML formatting
type YAMLHandler struct {
	writer io.Writer
}

func (h *YAMLHandler) Handle(data any) error {
	encoder := yaml.NewEncoder(h.writer)
	defer func() { _ = encoder.Close() }()
	return encoder.Encode(data)
}

// Table formatting methods
func (f *Formatter) formatStatusTable(services []ServiceStatus, options Options) error {
	if len(services) == 0 {
		f.output.Info(core.MsgServices_no_description)
		return nil
	}

	if options.Compact {
		return f.formatCompactStatusTable(services)
	}
	return f.formatFullStatusTable(services, options)
}

func (f *Formatter) formatCompactStatusTable(services []ServiceStatus) error {
	f.writeTableHeader([]string{ui.StatusHeaderService, ui.StatusHeaderState, ui.StatusHeaderHealth},
		[]int{ColWidthServiceCompact, ColWidthState, ColWidthHealth})
	f.writeTableSeparator(TableWidthCompact)

	for _, service := range services {
		f.writeTableRow([]string{
			service.Name,
			f.getStateIcon(service.State) + " " + service.State,
			f.getHealthIcon(service.Health) + " " + service.Health,
		}, []int{ColWidthServiceCompact, ColWidthState, ColWidthHealth})
	}
	return nil
}

func (f *Formatter) formatFullStatusTable(services []ServiceStatus, options Options) error {
	headers := []string{ui.StatusHeaderService, ui.StatusHeaderState, ui.StatusHeaderHealth, "UPTIME", "PORTS", "UPDATED"}
	widths := []int{ColWidthService, ColWidthState, ColWidthHealth, ColWidthUptime, ColWidthPorts, ColWidthUpdated}

	f.writeTableHeader(headers, widths)
	f.writeTableSeparator(TableWidthStandard)

	for _, service := range services {
		f.writeTableRow([]string{
			service.Name,
			f.getStateIcon(service.State) + " " + service.State,
			f.getHealthIcon(service.Health) + " " + service.Health,
			f.formatDuration(service.Uptime),
			f.formatPorts(service.Ports),
			service.UpdatedAt.Format("15:04:05"),
		}, widths)
	}

	if !options.Quiet {
		f.formatResourceSummary(services)
	}
	return nil
}

func (f *Formatter) writeTableHeader(headers []string, widths []int) {
	for i, header := range headers {
		_, _ = fmt.Fprintf(f.writer, "%-*s ", widths[i], header)
	}
	_, _ = fmt.Fprintln(f.writer)
}

func (f *Formatter) writeTableRow(values []string, widths []int) {
	for i, value := range values {
		_, _ = fmt.Fprintf(f.writer, "%-*s ", widths[i], value)
	}
	_, _ = fmt.Fprintln(f.writer)
}

func (f *Formatter) writeTableSeparator(width int) {
	_, _ = fmt.Fprintln(f.writer, strings.Repeat(ui.StatusSeparator, width))
}

func (f *Formatter) formatCatalogTable(catalog ServiceCatalog) error {
	if catalog.Total == 0 {
		f.output.Info(core.MsgServices_no_description)
		return nil
	}

	headers := []string{"CATEGORY", ui.StatusHeaderService, "DESCRIPTION"}
	widths := []int{ColWidthCategory, ColWidthDescription, 0} // Last column flexible

	f.writeTableHeader(headers, widths)
	f.writeTableSeparator(TableWidthCatalog)

	for categoryName, services := range catalog.Categories {
		for _, service := range services {
			f.writeTableRow([]string{categoryName, service.Name, service.Description}, widths)
		}
	}
	return nil
}

func (f *Formatter) formatCatalogGroup(catalog ServiceCatalog) error {
	if catalog.Total == 0 {
		f.output.Info(core.MsgServices_no_description)
		return nil
	}

	f.output.Header(pkgServices.MsgServiceCatalogHeader)

	for categoryName, serviceList := range catalog.Categories {
		if len(serviceList) == 0 {
			continue
		}

		icon := ui.IconBox
		if displayInfo, exists := pkgServices.CategoryDisplayInfo[categoryName]; exists {
			icon = displayInfo.Icon
		}

		plural := ""
		if len(serviceList) != 1 {
			plural = "s"
		}
		_, _ = fmt.Fprintf(f.writer, "%s %s\n", icon,
			fmt.Sprintf(pkgServices.MsgCategoryServiceCount, categoryName, len(serviceList), plural))

		for _, service := range serviceList {
			description := service.Description
			if description == "" {
				description = core.MsgServices_no_description
			}
			_, _ = fmt.Fprintf(f.writer, "  %-15s %s\n", service.Name, description)
		}
		_, _ = fmt.Fprintln(f.writer)
	}
	return nil
}

func (f *Formatter) formatValidationTable(result ValidationResult) error {
	if result.Valid {
		f.output.Success(core.MsgSuccess_config_valid)
		return nil
	}

	f.output.Error(core.MsgValidation_config_failed, len(result.Errors))

	if len(result.Errors) > 0 {
		_, _ = fmt.Fprintln(f.writer, "Errors:")
		for _, err := range result.Errors {
			_, _ = fmt.Fprintf(f.writer, "  %s: %s\n", err.Field, err.Message)
		}
		_, _ = fmt.Fprintln(f.writer)
	}

	if len(result.Warnings) > 0 {
		f.output.Warning(core.MsgValidation_warnings, len(result.Warnings))
		for _, warn := range result.Warnings {
			_, _ = fmt.Fprintf(f.writer, "  %s: %s\n", warn.Field, warn.Message)
		}
		_, _ = fmt.Fprintln(f.writer)
	}

	_, _ = fmt.Fprintln(f.writer, "Summary:")
	for key, value := range result.Summary {
		_, _ = fmt.Fprintf(f.writer, "  %s: %d\n", key, value)
	}
	return nil
}

func (f *Formatter) formatVersionTable(info VersionInfo, options Options) error {
	f.output.Info(core.MsgVersion_version_label, info.Version)

	if options.Full {
		_, _ = fmt.Fprintln(f.writer)
		_, _ = fmt.Fprintln(f.writer, "Build Information:")
		_, _ = fmt.Fprintf(f.writer, "%-15s %s\n", "Go Version:", info.GoVersion)
		_, _ = fmt.Fprintf(f.writer, "%-15s %s\n", "Platform:", info.Platform)

		for key, value := range info.BuildInfo {
			_, _ = fmt.Fprintf(f.writer, "%-15s %s\n", key+":", value)
		}
	}
	return nil
}

func (f *Formatter) formatHealthTable(report HealthReport, options Options) error {
	if report.Overall.Status == docker.HealthHealthy {
		f.output.Success(core.MsgSuccess_all_checks_passed, "system")
	} else {
		f.output.Error(core.MsgDoctor_some_issues)
	}

	if report.Overall.Message != "" {
		f.output.Info("   %s", report.Overall.Message)
	}

	if len(report.Checks) == 0 {
		return nil
	}

	f.output.Header(core.MsgDoctor_health_check_header, "System")
	headers := []string{"CHECK", "STATUS", "CATEGORY", "MESSAGE"}
	widths := []int{ColWidthCheck, ColWidthState, ColWidthCategory, ColWidthMessage}

	f.writeTableHeader(headers, widths)
	f.writeTableSeparator(TableWidthHealth)

	for _, check := range report.Checks {
		icon := f.getHealthIcon(check.Status)
		f.writeTableRow([]string{
			check.Name,
			icon + " " + check.Status,
			check.Category,
			check.Message,
		}, widths)

		if options.Verbose && check.Suggestion != "" {
			_, _ = fmt.Fprintf(f.writer, "   💡 %s\n", check.Suggestion)
		}
	}
	return nil
}

// Helper methods
func (f *Formatter) createSummary(services []ServiceStatus) map[string]int {
	summary := map[string]int{pkgServices.SummaryTotal: len(services)}
	for _, service := range services {
		if service.State == docker.StateRunning {
			summary[pkgServices.SummaryRunning]++
		}
		if service.Health == docker.HealthHealthy {
			summary[pkgServices.SummaryHealthy]++
		}
	}
	return summary
}

func (f *Formatter) formatDuration(d time.Duration) string {
	seconds := int(d.Seconds())
	if seconds < SecondsPerMinute {
		return fmt.Sprintf("%ds", seconds)
	}

	minutes := seconds / SecondsPerMinute
	if minutes < MinutesPerHour {
		return fmt.Sprintf("%dm", minutes)
	}

	hours := minutes / MinutesPerHour
	if hours < HoursPerDay {
		return fmt.Sprintf("%dh", hours)
	}

	return fmt.Sprintf("%dd", hours/HoursPerDay)
}

func (f *Formatter) formatPorts(ports []string) string {
	joined := strings.Join(ports, ",")
	if len(joined) > MaxPortsDisplay {
		return joined[:MaxPortsDisplay] + "..."
	}
	return joined
}

func (f *Formatter) formatResourceSummary(services []ServiceStatus) {
	summary := f.createSummary(services)
	_, _ = fmt.Fprintln(f.writer)
	_, _ = fmt.Fprintln(f.writer, "Resource Summary:")
	_, _ = fmt.Fprintf(f.writer, "  Total Services: %d\n", summary[pkgServices.SummaryTotal])
	_, _ = fmt.Fprintf(f.writer, "  Running: %d\n", summary[pkgServices.SummaryRunning])
	_, _ = fmt.Fprintf(f.writer, "  Healthy: %d\n", summary[pkgServices.SummaryHealthy])
}

func (f *Formatter) getStateIcon(state string) string {
	switch state {
	case docker.StateRunning:
		return core.Icons[core.IconState_running]
	case docker.StateStopped:
		return core.Icons[core.IconState_stopped]
	case docker.StateCreated:
		return core.Icons[core.IconState_created]
	case docker.StateStarting:
		return core.Icons[core.IconState_starting]
	case docker.StatePaused:
		return core.Icons[core.IconState_paused]
	default:
		return core.Icons[core.IconState_default]
	}
}

func (f *Formatter) getHealthIcon(health string) string {
	switch health {
	case docker.HealthHealthy:
		return core.Icons[core.IconHealth_healthy]
	case docker.HealthUnhealthy:
		return core.Icons[core.IconHealth_unhealthy]
	case docker.HealthStarting:
		return core.Icons[core.IconHealth_starting]
	case docker.HealthNone:
		return core.Icons[core.IconHealth_none]
	default:
		return core.Icons[core.IconHealth_default]
	}
}

// FilterCatalogByCategory filters service catalog by category
func FilterCatalogByCategory(catalog ServiceCatalog, category string) ServiceCatalog {
	if category == "" {
		return catalog
	}

	if services, exists := catalog.Categories[category]; exists {
		return ServiceCatalog{
			Categories: map[string][]ServiceInfo{category: services},
			Total:      len(services),
		}
	}

	return ServiceCatalog{
		Categories: make(map[string][]ServiceInfo),
		Total:      0,
	}
}
