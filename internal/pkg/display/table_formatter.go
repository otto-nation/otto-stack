package display

import (
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
)

// TableFormatter handles generic table formatting
type TableFormatter struct {
	writer io.Writer
}

// NewTableFormatter creates a new table formatter
func NewTableFormatter(writer io.Writer) *TableFormatter {
	return &TableFormatter{writer: writer}
}

// FormatTable formats data as a table with given headers and rows
func (tf *TableFormatter) FormatTable(headers []string, rows [][]string) error {
	tw := table.NewWriter()
	tw.SetOutputMirror(tf.writer)
	tw.SetStyle(table.StyleLight)

	// Convert headers to table.Row
	headerRow := make(table.Row, len(headers))
	for i, h := range headers {
		headerRow[i] = h
	}
	tw.AppendHeader(headerRow)

	// Convert rows to table.Row
	for _, row := range rows {
		tableRow := make(table.Row, len(row))
		for i, cell := range row {
			tableRow[i] = cell
		}
		tw.AppendRow(tableRow)
	}

	tw.Render()
	return nil
}

// FormatTableWithSeparators formats data as a table with separators between groups
func (tf *TableFormatter) FormatTableWithSeparators(headers []string, groups [][][]string) error {
	tw := table.NewWriter()
	tw.SetOutputMirror(tf.writer)
	tw.SetStyle(table.StyleLight)

	// Convert headers to table.Row
	headerRow := make(table.Row, len(headers))
	for i, h := range headers {
		headerRow[i] = h
	}
	tw.AppendHeader(headerRow)

	// Add rows with separators between groups
	for i, group := range groups {
		for _, row := range group {
			tableRow := make(table.Row, len(row))
			for j, cell := range row {
				tableRow[j] = cell
			}
			tw.AppendRow(tableRow)
		}
		if i < len(groups)-1 {
			tw.AppendSeparator()
		}
	}

	tw.Render()
	return nil
}

// RenderTable is a convenience function to render a simple table
func RenderTable(writer io.Writer, headers []string, rows [][]string) {
	tf := NewTableFormatter(writer)
	_ = tf.FormatTable(headers, rows)
}

// RenderTableWithSeparators is a convenience function to render a table with group separators
func RenderTableWithSeparators(writer io.Writer, headers []string, groups [][][]string) {
	tf := NewTableFormatter(writer)
	_ = tf.FormatTableWithSeparators(headers, groups)
}
