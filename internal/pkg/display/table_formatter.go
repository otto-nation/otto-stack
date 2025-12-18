package display

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// TableFormatter handles table-specific formatting
type TableFormatter struct {
	writer io.Writer
}

// NewTableFormatter creates a new table formatter
func NewTableFormatter(writer io.Writer) *TableFormatter {
	return &TableFormatter{writer: writer}
}

// WriteHeader writes table headers
func (tf *TableFormatter) WriteHeader(headers []string, widths []int) {
	tf.writeTableRow(headers, widths)
	tf.writeTableSeparator(calculateTotalWidth(widths))
}

// WriteRow writes a table row
func (tf *TableFormatter) WriteRow(values []string, widths []int) {
	tf.writeTableRow(values, widths)
}

// WriteSeparator writes a table separator
func (tf *TableFormatter) WriteSeparator(width int) {
	tf.writeTableSeparator(width)
}

func (tf *TableFormatter) writeTableRow(values []string, widths []int) {
	for i, value := range values {
		if i < len(widths) {
			_, _ = fmt.Fprintf(tf.writer, "%-*s ", widths[i], truncateString(value, widths[i]))
		}
	}
	_, _ = fmt.Fprintln(tf.writer)
}

func (tf *TableFormatter) writeTableSeparator(width int) {
	_, _ = fmt.Fprintln(tf.writer, strings.Repeat("-", width))
}

func calculateTotalWidth(widths []int) int {
	total := 0
	for _, w := range widths {
		total += w + 1 // +1 for space
	}
	return total
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// FormatDuration formats duration for display
func (tf *TableFormatter) FormatDuration(d time.Duration) string {
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

// FormatPorts formats port list for display
func (tf *TableFormatter) FormatPorts(ports []string) string {
	if len(ports) == 0 {
		return "-"
	}
	if len(ports) <= MaxPortsDisplay {
		return strings.Join(ports, ",")
	}
	return strings.Join(ports[:MaxPortsDisplay], ",") + "..."
}
