package display

import (
	"fmt"
	"io"
)

const (
	// Status symbols
	StatusSuccess  = "✓ "
	StatusError    = "✗ "
	StatusStarting = "⟳ "

	// Status strings
	StatusStringStarting = "starting"
)

// ValidationFormatter handles validation result formatting
type ValidationFormatter struct {
	writer io.Writer
	table  *TableFormatter
}

// NewValidationFormatter creates a new validation formatter
func NewValidationFormatter(writer io.Writer) *ValidationFormatter {
	return &ValidationFormatter{
		writer: writer,
		table:  NewTableFormatter(writer),
	}
}

// FormatTable formats validation results as a table
func (vf *ValidationFormatter) FormatTable(result ValidationResult) error {
	headers := []string{"Type", "Field", "Message"}
	widths := []int{ColWidthCheck, ColWidthState, ColWidthMessage}

	vf.table.WriteHeader(headers, widths)

	// Show errors
	for _, issue := range result.Errors {
		values := []string{
			"ERROR",
			issue.Field,
			issue.Message,
		}
		vf.table.WriteRow(values, widths)
	}

	// Show warnings
	for _, issue := range result.Warnings {
		values := []string{
			"WARNING",
			issue.Field,
			issue.Message,
		}
		vf.table.WriteRow(values, widths)
	}

	_, _ = fmt.Fprintln(vf.writer)
	if result.Valid {
		_, _ = fmt.Fprintln(vf.writer, StatusSuccess+"Validation passed")
	} else {
		_, _ = fmt.Fprintln(vf.writer, StatusError+"Validation failed")
	}
	return nil
}

// HealthFormatter handles health report formatting
type HealthFormatter struct {
	writer io.Writer
	table  *TableFormatter
}

// NewHealthFormatter creates a new health formatter
func NewHealthFormatter(writer io.Writer) *HealthFormatter {
	return &HealthFormatter{
		writer: writer,
		table:  NewTableFormatter(writer),
	}
}

// FormatTable formats health report as a table
func (hf *HealthFormatter) FormatTable(report HealthReport, options Options) error {
	headers := []string{"Check", "Status", "Message"}
	widths := []int{ColWidthService, ColWidthHealth, ColWidthMessage}

	hf.table.WriteHeader(headers, widths)

	for _, check := range report.Checks {
		values := []string{
			check.Name,
			hf.getHealthIcon(check.Status) + check.Status,
			check.Message,
		}
		hf.table.WriteRow(values, widths)
	}

	if options.ShowSummary {
		hf.formatHealthSummary(report)
	}
	return nil
}

func (hf *HealthFormatter) formatHealthSummary(report HealthReport) {
	_, _ = fmt.Fprintln(hf.writer)
	_, _ = fmt.Fprintf(hf.writer, "Overall Health: %s\n", report.Overall.Status)
	_, _ = fmt.Fprintf(hf.writer, "Message: %s\n", report.Overall.Message)
}

func (hf *HealthFormatter) getHealthIcon(health string) string {
	switch health {
	case "healthy":
		return StatusSuccess
	case "unhealthy":
		return StatusError
	case StatusStringStarting:
		return StatusStarting
	default:
		return "? "
	}
}

// VersionFormatter handles version info formatting
type VersionFormatter struct {
	writer io.Writer
	table  *TableFormatter
}

// NewVersionFormatter creates a new version formatter
func NewVersionFormatter(writer io.Writer) *VersionFormatter {
	return &VersionFormatter{
		writer: writer,
		table:  NewTableFormatter(writer),
	}
}

// FormatTable formats version info as a table
func (vf *VersionFormatter) FormatTable(info VersionInfo, options Options) error {
	headers := []string{"Component", "Version", "Platform"}
	widths := []int{ColWidthService, ColWidthState, ColWidthMessage}

	vf.table.WriteHeader(headers, widths)

	values := []string{
		"Otto Stack",
		info.Version,
		info.Platform,
	}
	vf.table.WriteRow(values, widths)

	if options.ShowSummary {
		_, _ = fmt.Fprintln(vf.writer)
		_, _ = fmt.Fprintf(vf.writer, "Go Version: %s\n", info.GoVersion)
		for key, value := range info.BuildInfo {
			_, _ = fmt.Fprintf(vf.writer, "%s: %s\n", key, value)
		}
	}
	return nil
}
