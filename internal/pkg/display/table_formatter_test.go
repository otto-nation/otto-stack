//go:build unit

package display

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTableFormatter(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewTableFormatter(buf)
	assert.NotNil(t, formatter)
}

func TestTableFormatter_FormatTable(t *testing.T) {
	t.Run("formats simple table", func(t *testing.T) {
		buf := &bytes.Buffer{}
		formatter := NewTableFormatter(buf)

		headers := []string{"Name", "Status"}
		rows := [][]string{
			{"service1", "running"},
			{"service2", "stopped"},
		}

		err := formatter.FormatTable(headers, rows)
		require.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "service1")
		assert.Contains(t, output, "running")
		assert.Contains(t, output, "stopped")
	})

	t.Run("handles empty rows", func(t *testing.T) {
		buf := &bytes.Buffer{}
		formatter := NewTableFormatter(buf)

		headers := []string{"Name", "Status"}
		rows := [][]string{}

		err := formatter.FormatTable(headers, rows)
		require.NoError(t, err)
		output := buf.String()
		assert.True(t, len(output) > 0)
	})
}

func TestTableFormatter_FormatTableWithSeparators(t *testing.T) {
	t.Run("formats table with groups", func(t *testing.T) {
		buf := &bytes.Buffer{}
		formatter := NewTableFormatter(buf)

		headers := []string{"Service", "Category"}
		groups := [][][]string{
			{
				{"postgres", "database"},
				{"mysql", "database"},
			},
			{
				{"redis", "cache"},
			},
		}

		err := formatter.FormatTableWithSeparators(headers, groups)
		require.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "postgres")
		assert.Contains(t, output, "redis")
		assert.Contains(t, output, "database")
	})
}

func TestRenderTable(t *testing.T) {
	buf := &bytes.Buffer{}
	headers := []string{"Name", "Status"}
	rows := [][]string{
		{"service1", "running"},
	}

	RenderTable(buf, headers, rows)
	assert.Contains(t, buf.String(), "service1")
}

func TestRenderTableWithSeparators(t *testing.T) {
	buf := &bytes.Buffer{}
	headers := []string{"Service", "Type"}
	groups := [][][]string{
		{{"postgres", "database"}},
		{{"redis", "cache"}},
	}

	RenderTableWithSeparators(buf, headers, groups)
	output := buf.String()
	assert.Contains(t, output, "postgres")
	assert.Contains(t, output, "redis")
}
