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

// Formatter handles all output formatting
type Formatter struct {
	writer io.Writer
	output base.Output
}

// New creates a new formatter
func New(writer io.Writer, output base.Output) *Formatter {
	return &Formatter{
		writer: writer,
		output: output,
	}
}

// FormatStatus formats service status information
func (f *Formatter) FormatStatus(services []ServiceStatus, options Options) error {
	switch options.Format {
	case pkgServices.ServiceCatalogJSONFormat:
		return f.writeJSON(map[string]any{
			"services": services,
			"summary":  f.createSummary(services),
		})
	case pkgServices.ServiceCatalogYAMLFormat:
		return f.writeYAML(map[string]any{
			"services": services,
			"summary":  f.createSummary(services),
		})
	default:
		return f.formatStatusTable(services, options)
	}
}

// FormatServiceCatalog formats service catalog information
func (f *Formatter) FormatServiceCatalog(catalog ServiceCatalog, options Options) error {
	switch options.Format {
	case pkgServices.ServiceCatalogJSONFormat:
		return f.writeJSON(catalog)
	case pkgServices.ServiceCatalogYAMLFormat:
		return f.writeYAML(catalog)
	case pkgServices.ServiceCatalogTableFormat:
		return f.formatCatalogTable(catalog)
	default:
		return f.formatCatalogGroup(catalog)
	}
}

// FormatValidation formats validation results
func (f *Formatter) FormatValidation(result ValidationResult, options Options) error {
	switch options.Format {
	case pkgServices.ServiceCatalogJSONFormat:
		return f.writeJSON(result)
	case pkgServices.ServiceCatalogYAMLFormat:
		return f.writeYAML(result)
	default:
		return f.formatValidationTable(result)
	}
}

// FormatVersion formats version information
func (f *Formatter) FormatVersion(info VersionInfo, options Options) error {
	switch options.Format {
	case pkgServices.ServiceCatalogJSONFormat:
		return f.writeJSON(info)
	case pkgServices.ServiceCatalogYAMLFormat:
		return f.writeYAML(info)
	default:
		return f.formatVersionTable(info, options)
	}
}

// FormatHealth formats health check results
func (f *Formatter) FormatHealth(report HealthReport, options Options) error {
	switch options.Format {
	case pkgServices.ServiceCatalogJSONFormat:
		return f.writeJSON(report)
	case pkgServices.ServiceCatalogYAMLFormat:
		return f.writeYAML(report)
	default:
		return f.formatHealthTable(report, options)
	}
}

// Table formatting methods
func (f *Formatter) formatStatusTable(services []ServiceStatus, options Options) error {
	if len(services) == 0 {
		f.output.Info(core.MsgServices_no_description)
		return nil
	}

	if options.Compact {
		_, _ = fmt.Fprintf(f.writer, "%-20s %-10s %-12s\n",
			ui.StatusHeaderService, ui.StatusHeaderState, ui.StatusHeaderHealth)
		_, _ = fmt.Fprintln(f.writer, strings.Repeat(ui.StatusSeparator, ui.TableWidth42))
		for _, service := range services {
			_, _ = fmt.Fprintf(f.writer, "%-20s %-10s %-12s\n",
				service.Name, f.getStateIcon(service.State)+" "+service.State,
				f.getHealthIcon(service.Health)+" "+service.Health)
		}
		return nil
	}

	_, _ = fmt.Fprintf(f.writer, "%-15s %-10s %-10s %-8s %-12s %-10s\n",
		ui.StatusHeaderService, ui.StatusHeaderState, ui.StatusHeaderHealth,
		"UPTIME", "PORTS", "UPDATED")
	_, _ = fmt.Fprintln(f.writer, strings.Repeat(ui.StatusSeparator, ui.TableWidth75))

	for _, service := range services {
		_, _ = fmt.Fprintf(f.writer, "%-15s %-10s %-10s %-8s %-12s %-10s\n",
			service.Name,
			f.getStateIcon(service.State)+" "+service.State,
			f.getHealthIcon(service.Health)+" "+service.Health,
			f.formatDuration(service.Uptime),
			f.formatPorts(service.Ports),
			service.UpdatedAt.Format("15:04:05"))
	}

	if !options.Quiet {
		f.formatResourceSummary(services)
	}
	return nil
}

func (f *Formatter) formatCatalogTable(catalog ServiceCatalog) error {
	if catalog.Total == 0 {
		f.output.Info(core.MsgServices_no_description)
		return nil
	}

	_, _ = fmt.Fprintf(f.writer, "%-15s %-20s %s\n", "CATEGORY", ui.StatusHeaderService, "DESCRIPTION")
	_, _ = fmt.Fprintln(f.writer, strings.Repeat(ui.StatusSeparator, ui.TableWidth80))

	for categoryName, services := range catalog.Categories {
		for _, service := range services {
			_, _ = fmt.Fprintf(f.writer, "%-15s %-20s %s\n",
				categoryName, service.Name, service.Description)
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

		icon := "📦"
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
	_, _ = fmt.Fprintf(f.writer, "%-25s %-10s %-15s %-40s\n",
		"CHECK", "STATUS", "CATEGORY", "MESSAGE")
	_, _ = fmt.Fprintln(f.writer, strings.Repeat(ui.StatusSeparator, ui.TableWidth90))

	for _, check := range report.Checks {
		icon := f.getHealthIcon(check.Status)
		_, _ = fmt.Fprintf(f.writer, "%-25s %-10s %-15s %-40s\n",
			check.Name, icon+" "+check.Status, check.Category, check.Message)

		if options.Verbose && check.Suggestion != "" {
			_, _ = fmt.Fprintf(f.writer, "   💡 %s\n", check.Suggestion)
		}
	}
	return nil
}

// Helper methods
func (f *Formatter) writeJSON(data any) error {
	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (f *Formatter) writeYAML(data any) error {
	encoder := yaml.NewEncoder(f.writer)
	defer func() { _ = encoder.Close() }()
	return encoder.Encode(data)
}

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
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < core.HoursPerDay*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/core.HoursPerDay))
}

func (f *Formatter) formatPorts(ports []string) string {
	joined := strings.Join(ports, ",")
	if len(joined) > core.MaxCategoryCommands {
		return joined[:core.MaxCategoryCommands] + "..."
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
