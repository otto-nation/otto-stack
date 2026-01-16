package display

import (
	"fmt"
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
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
}

// NewValidationFormatter creates a new validation formatter
func NewValidationFormatter(writer io.Writer) *ValidationFormatter {
	return &ValidationFormatter{writer: writer}
}

// FormatTable formats validation results as a table
func (vf *ValidationFormatter) FormatTable(result ValidationResult) error {
	tw := table.NewWriter()
	tw.SetOutputMirror(vf.writer)
	tw.SetStyle(table.StyleLight)

	tw.AppendHeader(table.Row{HeaderType, HeaderField, HeaderMessage})

	for _, issue := range result.Errors {
		tw.AppendRow(table.Row{"ERROR", issue.Field, issue.Message})
	}

	for _, issue := range result.Warnings {
		tw.AppendRow(table.Row{"WARNING", issue.Field, issue.Message})
	}

	tw.Render()

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
}

// NewHealthFormatter creates a new health formatter
func NewHealthFormatter(writer io.Writer) *HealthFormatter {
	return &HealthFormatter{writer: writer}
}

// FormatTable formats health report as a table
func (hf *HealthFormatter) FormatTable(report HealthReport, options Options) error {
	tw := table.NewWriter()
	tw.SetOutputMirror(hf.writer)
	tw.SetStyle(table.StyleLight)

	tw.AppendHeader(table.Row{HeaderCheck, HeaderStatus, HeaderMessage})

	for _, check := range report.Checks {
		tw.AppendRow(table.Row{
			check.Name,
			hf.getHealthIcon(check.Status) + check.Status,
			check.Message,
		})
	}

	tw.Render()

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
}

// NewVersionFormatter creates a new version formatter
func NewVersionFormatter(writer io.Writer) *VersionFormatter {
	return &VersionFormatter{writer: writer}
}

// FormatTable formats version info as a table
func (vf *VersionFormatter) FormatTable(info VersionInfo, options Options) error {
	tw := table.NewWriter()
	tw.SetOutputMirror(vf.writer)
	tw.SetStyle(table.StyleLight)

	tw.AppendHeader(table.Row{HeaderComponent, HeaderVersion, HeaderPlatform})
	tw.AppendRow(table.Row{"Otto Stack", info.Version, info.Platform})

	tw.Render()

	if options.ShowSummary {
		_, _ = fmt.Fprintln(vf.writer)
		_, _ = fmt.Fprintf(vf.writer, "Go Version: %s\n", info.GoVersion)
		for key, value := range info.BuildInfo {
			_, _ = fmt.Fprintf(vf.writer, "%s: %s\n", key, value)
		}
	}
	return nil
}
