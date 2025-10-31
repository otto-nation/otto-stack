package display

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
)

// TableFormatter implements table-based output formatting
type TableFormatter struct {
	writer io.Writer
}

// NewTableFormatter creates a new table formatter
func NewTableFormatter(writer io.Writer) *TableFormatter {
	return &TableFormatter{writer: writer}
}

// FormatStatus formats service status as a table
func (f *TableFormatter) FormatStatus(services []ServiceStatus, options StatusOptions) error {
	if len(services) == 0 {
		//nolint:errcheck
		fmt.Fprintln(f.writer, "No services found")
		return nil
	}

	if options.Compact {
		return f.formatCompactStatus(services)
	}
	return f.formatDetailedStatus(services, options.Quiet)
}

// FormatServiceCatalog formats service catalog as a table
func (f *TableFormatter) FormatServiceCatalog(catalog ServiceCatalog, options ServiceCatalogOptions) error {
	filteredCatalog := FilterCatalogByCategory(catalog, options.Category)

	if filteredCatalog.Total == 0 {
		if options.Category != "" {
			_, _ = fmt.Fprintf(f.writer, constants.MsgNoServicesInCategory+"\n", options.Category)
		} else {
			_, _ = fmt.Fprintln(f.writer, "No services available")
		}
		return nil
	}

	// Table format
	_, _ = fmt.Fprintf(f.writer, "%-15s %-20s %s\n", "CATEGORY", "SERVICE", "DESCRIPTION")
	_, _ = fmt.Fprintln(f.writer, strings.Repeat("-", 80))

	for categoryName, services := range filteredCatalog.Categories {
		for _, service := range services {
			_, _ = fmt.Fprintf(f.writer, "%-15s %-20s %s\n",
				categoryName, service.Name, service.Description)
		}
	}

	return nil
}

// FormatValidation formats validation results as a table
func (f *TableFormatter) FormatValidation(result ValidationResult, options ValidationOptions) error {
	if result.Valid {
		//nolint:errcheck
		fmt.Fprintln(f.writer, "✅ Configuration is valid")
		return nil
	}

	//nolint:errcheck
	fmt.Fprintln(f.writer, "❌ Configuration validation failed")
	//nolint:errcheck
	fmt.Fprintln(f.writer)

	if len(result.Errors) > 0 {
		//nolint:errcheck
		fmt.Fprintln(f.writer, "Errors:")
		//nolint:errcheck
		fmt.Fprintf(f.writer, "%-15s %-20s %-50s\n", "SEVERITY", "FIELD", "MESSAGE")
		//nolint:errcheck
		fmt.Fprintln(f.writer, strings.Repeat("-", 85))

		for _, err := range result.Errors {
			//nolint:errcheck
			fmt.Fprintf(f.writer, "%-15s %-20s %-50s\n",
				err.Severity, err.Field, err.Message)
		}
		//nolint:errcheck
		fmt.Fprintln(f.writer)
	}

	if len(result.Warnings) > 0 {
		//nolint:errcheck
		fmt.Fprintln(f.writer, "Warnings:")
		//nolint:errcheck
		fmt.Fprintf(f.writer, "%-15s %-20s %-50s\n", "TYPE", "FIELD", "MESSAGE")
		//nolint:errcheck
		fmt.Fprintln(f.writer, strings.Repeat("-", 85))

		for _, warn := range result.Warnings {
			//nolint:errcheck
			fmt.Fprintf(f.writer, "%-15s %-20s %-50s\n",
				warn.Type, warn.Field, warn.Message)
		}
		//nolint:errcheck
		fmt.Fprintln(f.writer)
	}

	f.formatValidationSummary(result.Summary)
	return nil
}

// FormatVersion formats version information as a table
func (f *TableFormatter) FormatVersion(info VersionInfo, options VersionOptions) error {
	//nolint:errcheck
	fmt.Fprintf(f.writer, "%s version %s\n", constants.AppName, info.Version)

	if options.Full {
		//nolint:errcheck
		fmt.Fprintln(f.writer)
		//nolint:errcheck
		fmt.Fprintln(f.writer, "Build Information:")
		//nolint:errcheck
		fmt.Fprintf(f.writer, "%-15s %s\n", "Go Version:", info.GoVersion)
		//nolint:errcheck
		fmt.Fprintf(f.writer, "%-15s %s\n", "Platform:", info.Platform)

		for key, value := range info.BuildInfo {
			//nolint:errcheck
			fmt.Fprintf(f.writer, "%-15s %s\n", key+":", value)
		}
	}

	return nil
}

// FormatHealth formats health check results as a table
func (f *TableFormatter) FormatHealth(report HealthReport, options HealthOptions) error {
	// Display overall status
	statusIcon := "✅"
	if report.Overall.Status != "healthy" {
		statusIcon = "❌"
	}

	//nolint:errcheck
	fmt.Fprintf(f.writer, "%s Overall Status: %s\n", statusIcon, report.Overall.Status)
	if report.Overall.Message != "" {
		//nolint:errcheck
		fmt.Fprintf(f.writer, "   %s\n", report.Overall.Message)
	}
	//nolint:errcheck
	fmt.Fprintln(f.writer)

	if len(report.Checks) == 0 {
		return nil
	}

	// Display detailed checks
	//nolint:errcheck
	fmt.Fprintln(f.writer, "Health Checks:")
	//nolint:errcheck
	fmt.Fprintf(f.writer, "%-25s %-10s %-15s %-40s\n",
		"CHECK", "STATUS", "CATEGORY", "MESSAGE")
	//nolint:errcheck
	fmt.Fprintln(f.writer, strings.Repeat("-", 90))

	for _, check := range report.Checks {
		icon := GetHealthIcon(check.Status)
		_, _ = fmt.Fprintf(f.writer, "%-25s %-10s %-15s %-40s\n",
			check.Name, icon+" "+check.Status, check.Category, check.Message)

		if options.Verbose && check.Suggestion != "" {
			_, _ = fmt.Fprintf(f.writer, "   💡 %s\n", check.Suggestion)
		}
	}

	return nil
}

// Helper methods
func (f *TableFormatter) formatCompactStatus(services []ServiceStatus) error {
	_, _ = fmt.Fprintf(f.writer, "%-20s %-10s %-12s\n", "SERVICE", "STATE", "HEALTH")
	_, _ = fmt.Fprintln(f.writer, strings.Repeat("-", 42))

	for _, service := range services {
		stateIcon := GetStateIcon(service.State)
		healthIcon := GetHealthIcon(service.Health)

		_, _ = fmt.Fprintf(f.writer, "%-20s %-10s %-12s\n",
			service.Name, stateIcon+" "+service.State, healthIcon+" "+service.Health)
	}

	return nil
}

// ServiceRowData holds formatted service data for display
type ServiceRowData struct {
	Name    string
	State   string
	Health  string
	Uptime  string
	Ports   string
	Updated string
}

// formatServiceRow converts ServiceStatus to display format
func (f *TableFormatter) formatServiceRow(service ServiceStatus) ServiceRowData {
	stateIcon := GetStateIcon(service.State)
	healthIcon := GetHealthIcon(service.Health)
	uptime := f.formatDuration(service.Uptime)
	ports := f.formatPorts(service.Ports)
	updated := service.UpdatedAt.Format("15:04:05")

	return ServiceRowData{
		Name:    service.Name,
		State:   stateIcon + " " + service.State,
		Health:  healthIcon + " " + service.Health,
		Uptime:  uptime,
		Ports:   ports,
		Updated: updated,
	}
}

// formatPorts truncates port list if too long
func (f *TableFormatter) formatPorts(ports []string) string {
	joined := strings.Join(ports, ",")
	if len(joined) > 10 {
		return joined[:10] + "..."
	}
	return joined
}

func (f *TableFormatter) formatDetailedStatus(services []ServiceStatus, quiet bool) error {
	_, _ = fmt.Fprintf(f.writer, "%-15s %-10s %-10s %-8s %-12s %-10s\n",
		"SERVICE", "STATE", "HEALTH", "UPTIME", "PORTS", "UPDATED")
	_, _ = fmt.Fprintln(f.writer, strings.Repeat("-", 75))

	for _, service := range services {
		row := f.formatServiceRow(service)
		_, _ = fmt.Fprintf(f.writer, "%-15s %-10s %-10s %-8s %-12s %-10s\n",
			row.Name, row.State, row.Health, row.Uptime, row.Ports, row.Updated)
	}

	if !quiet {
		f.formatResourceSummary(services)
	}

	return nil
}

func (f *TableFormatter) formatValidationSummary(summary ValidationSummary) {
	//nolint:errcheck
	fmt.Fprintln(f.writer, "Summary:")
	//nolint:errcheck
	fmt.Fprintf(f.writer, "  Total Commands: %d\n", summary.TotalCommands)
	//nolint:errcheck
	fmt.Fprintf(f.writer, "  Valid Commands: %d\n", summary.ValidCommands)
	//nolint:errcheck
	fmt.Fprintf(f.writer, "  Errors: %d\n", summary.ErrorCount)
	//nolint:errcheck
	fmt.Fprintf(f.writer, "  Warnings: %d\n", summary.WarningCount)
	//nolint:errcheck
	fmt.Fprintf(f.writer, "  Coverage: %d%%\n", summary.CoveragePercentage)
}

func (f *TableFormatter) formatResourceSummary(services []ServiceStatus) {
	running := CountByState(services, constants.StateRunning)
	healthy := CountByHealth(services, constants.HealthHealthy)

	_, _ = fmt.Fprintln(f.writer)
	_, _ = fmt.Fprintln(f.writer, "Resource Summary:")
	_, _ = fmt.Fprintf(f.writer, "  Total Services: %d\n", len(services))
	_, _ = fmt.Fprintf(f.writer, "  Running: %d\n", running)
	_, _ = fmt.Fprintf(f.writer, "  Healthy: %d\n", healthy)
}

func (f *TableFormatter) formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}
