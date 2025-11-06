package display

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"gopkg.in/yaml.v3"
)

// Formatter handles all output formatting
type Formatter struct {
	writer io.Writer
	output types.Output
}

// New creates a new formatter
func New(writer io.Writer, output types.Output) *Formatter {
	return &Formatter{
		writer: writer,
		output: output,
	}
}

// FormatStatus formats service status information
func (f *Formatter) FormatStatus(services []ServiceStatus, options Options) error {
	switch options.Format {
	case constants.ServiceCatalogJSONFormat:
		return f.writeJSON(map[string]any{
			"services": services,
			"summary":  f.createSummary(services),
		})
	case constants.ServiceCatalogYAMLFormat:
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
	case constants.ServiceCatalogJSONFormat:
		return f.writeJSON(catalog)
	case constants.ServiceCatalogYAMLFormat:
		return f.writeYAML(catalog)
	case constants.ServiceCatalogTableFormat:
		return f.formatCatalogTable(catalog)
	default:
		return f.formatCatalogGroup(catalog)
	}
}

// FormatValidation formats validation results
func (f *Formatter) FormatValidation(result ValidationResult, options Options) error {
	switch options.Format {
	case constants.ServiceCatalogJSONFormat:
		return f.writeJSON(result)
	case constants.ServiceCatalogYAMLFormat:
		return f.writeYAML(result)
	default:
		return f.formatValidationTable(result)
	}
}

// FormatVersion formats version information
func (f *Formatter) FormatVersion(info VersionInfo, options Options) error {
	switch options.Format {
	case constants.ServiceCatalogJSONFormat:
		return f.writeJSON(info)
	case constants.ServiceCatalogYAMLFormat:
		return f.writeYAML(info)
	default:
		return f.formatVersionTable(info, options)
	}
}

// FormatHealth formats health check results
func (f *Formatter) FormatHealth(report HealthReport, options Options) error {
	switch options.Format {
	case constants.ServiceCatalogJSONFormat:
		return f.writeJSON(report)
	case constants.ServiceCatalogYAMLFormat:
		return f.writeYAML(report)
	default:
		return f.formatHealthTable(report, options)
	}
}

// Table formatting methods
func (f *Formatter) formatStatusTable(services []ServiceStatus, options Options) error {
	if len(services) == 0 {
		f.output.Info(constants.MsgServices_no_description)
		return nil
	}

	if options.Compact {
		_, _ = fmt.Fprintf(f.writer, "%-20s %-10s %-12s\n",
			constants.StatusHeaderService, constants.StatusHeaderState, constants.StatusHeaderHealth)
		_, _ = fmt.Fprintln(f.writer, strings.Repeat(constants.StatusSeparator, constants.TableWidth42))
		for _, service := range services {
			_, _ = fmt.Fprintf(f.writer, "%-20s %-10s %-12s\n",
				service.Name, f.getStateIcon(service.State)+" "+service.State,
				f.getHealthIcon(service.Health)+" "+service.Health)
		}
		return nil
	}

	_, _ = fmt.Fprintf(f.writer, "%-15s %-10s %-10s %-8s %-12s %-10s\n",
		constants.StatusHeaderService, constants.StatusHeaderState, constants.StatusHeaderHealth,
		"UPTIME", "PORTS", "UPDATED")
	_, _ = fmt.Fprintln(f.writer, strings.Repeat(constants.StatusSeparator, constants.TableWidth75))

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
		f.output.Info(constants.MsgServices_no_description)
		return nil
	}

	_, _ = fmt.Fprintf(f.writer, "%-15s %-20s %s\n", "CATEGORY", constants.StatusHeaderService, "DESCRIPTION")
	_, _ = fmt.Fprintln(f.writer, strings.Repeat(constants.StatusSeparator, constants.TableWidth80))

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
		f.output.Info(constants.MsgServices_no_description)
		return nil
	}

	f.output.Header(constants.MsgServiceCatalogHeader)

	for categoryName, serviceList := range catalog.Categories {
		if len(serviceList) == 0 {
			continue
		}

		icon := "📦"
		if displayInfo, exists := services.CategoryDisplayInfo[categoryName]; exists {
			icon = displayInfo.Icon
		}

		plural := ""
		if len(serviceList) != 1 {
			plural = "s"
		}
		_, _ = fmt.Fprintf(f.writer, "%s %s\n", icon,
			fmt.Sprintf(constants.MsgServiceCount, categoryName, len(serviceList), plural))

		for _, service := range serviceList {
			description := service.Description
			if description == "" {
				description = constants.MsgServices_no_description
			}
			_, _ = fmt.Fprintf(f.writer, "  %-15s %s\n", service.Name, description)
		}
		_, _ = fmt.Fprintln(f.writer)
	}
	return nil
}

func (f *Formatter) formatValidationTable(result ValidationResult) error {
	if result.Valid {
		f.output.Success(constants.MsgSuccess_config_valid)
		return nil
	}

	f.output.Error(constants.MsgValidation_config_failed, len(result.Errors))

	if len(result.Errors) > 0 {
		_, _ = fmt.Fprintln(f.writer, "Errors:")
		for _, err := range result.Errors {
			_, _ = fmt.Fprintf(f.writer, "  %s: %s\n", err.Field, err.Message)
		}
		_, _ = fmt.Fprintln(f.writer)
	}

	if len(result.Warnings) > 0 {
		f.output.Warning(constants.MsgValidation_warnings, len(result.Warnings))
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
	f.output.Info(constants.MsgVersion_version_label, info.Version)

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
	if report.Overall.Status == constants.HealthHealthy {
		f.output.Success(constants.MsgSuccess_all_checks_passed, "system")
	} else {
		f.output.Error(constants.MsgDoctor_some_issues)
	}

	if report.Overall.Message != "" {
		f.output.Info("   %s", report.Overall.Message)
	}

	if len(report.Checks) == 0 {
		return nil
	}

	f.output.Header(constants.MsgDoctor_health_check_header, "System")
	_, _ = fmt.Fprintf(f.writer, "%-25s %-10s %-15s %-40s\n",
		"CHECK", "STATUS", "CATEGORY", "MESSAGE")
	_, _ = fmt.Fprintln(f.writer, strings.Repeat(constants.StatusSeparator, constants.TableWidth90))

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
	summary := map[string]int{constants.SummaryTotal: len(services)}
	for _, service := range services {
		if service.State == constants.StateRunning {
			summary[constants.SummaryRunning]++
		}
		if service.Health == constants.HealthHealthy {
			summary[constants.SummaryHealthy]++
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
	if d < constants.HoursPerDay*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/constants.HoursPerDay))
}

func (f *Formatter) formatPorts(ports []string) string {
	joined := strings.Join(ports, ",")
	if len(joined) > constants.MaxCategoryCommands {
		return joined[:constants.MaxCategoryCommands] + "..."
	}
	return joined
}

func (f *Formatter) formatResourceSummary(services []ServiceStatus) {
	summary := f.createSummary(services)
	_, _ = fmt.Fprintln(f.writer)
	_, _ = fmt.Fprintln(f.writer, "Resource Summary:")
	_, _ = fmt.Fprintf(f.writer, "  Total Services: %d\n", summary[constants.SummaryTotal])
	_, _ = fmt.Fprintf(f.writer, "  Running: %d\n", summary[constants.SummaryRunning])
	_, _ = fmt.Fprintf(f.writer, "  Healthy: %d\n", summary[constants.SummaryHealthy])
}

func (f *Formatter) getStateIcon(state string) string {
	switch state {
	case constants.StateRunning:
		return constants.Icons[constants.IconState_running]
	case constants.StateStopped:
		return constants.Icons[constants.IconState_stopped]
	case constants.StateCreated:
		return constants.Icons[constants.IconState_created]
	case constants.StateStarting:
		return constants.Icons[constants.IconState_starting]
	case constants.StatePaused:
		return constants.Icons[constants.IconState_paused]
	default:
		return constants.Icons[constants.IconState_default]
	}
}

func (f *Formatter) getHealthIcon(health string) string {
	switch health {
	case constants.HealthHealthy:
		return constants.Icons[constants.IconHealth_healthy]
	case constants.HealthUnhealthy:
		return constants.Icons[constants.IconHealth_unhealthy]
	case constants.HealthStarting:
		return constants.Icons[constants.IconHealth_starting]
	case constants.HealthNone:
		return constants.Icons[constants.IconHealth_none]
	default:
		return constants.Icons[constants.IconHealth_default]
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
