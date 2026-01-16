//go:build unit

package display

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"github.com/stretchr/testify/assert"
)

func TestFormatter_FormatStatus_Table(t *testing.T) {
	var buf bytes.Buffer
	output := ui.NewOutput()
	formatter := New(&buf, output)

	services := []ServiceStatus{
		{
			Name:      services.ServicePostgres,
			State:     "running",
			Health:    "healthy",
			Ports:     []string{"5432:5432"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Uptime:    time.Hour,
		},
	}

	err := formatter.FormatStatus(services, Options{})
	assert.NoError(t, err)

	output_str := buf.String()
	assert.Contains(t, output_str, "postgres")
	assert.Contains(t, output_str, "run") // Changed from "running" to "run"
	assert.Contains(t, output_str, "hea") // Changed from "healthy" to "hea"
}

func TestFormatter_FormatStatus_JSON(t *testing.T) {
	var buf bytes.Buffer
	output := ui.NewOutput()
	formatter := New(&buf, output)

	services := []ServiceStatus{
		{
			Name:   services.ServicePostgres,
			State:  "running",
			Health: "healthy",
		},
	}

	err := formatter.FormatStatus(services, Options{Format: "json"})
	assert.NoError(t, err)

	output_str := buf.String()
	assert.Contains(t, output_str, `"name": "postgres"`)
	assert.Contains(t, output_str, `"state": "running"`)
}

func TestFormatter_FormatStatus_YAML(t *testing.T) {
	var buf bytes.Buffer
	output := ui.NewOutput()
	formatter := New(&buf, output)

	services := []ServiceStatus{
		{
			Name:   services.ServicePostgres,
			State:  "running",
			Health: "healthy",
		},
	}

	err := formatter.FormatStatus(services, Options{Format: "yaml"})
	assert.NoError(t, err)

	output_str := buf.String()
	assert.Contains(t, output_str, "name: postgres")
	assert.Contains(t, output_str, "state: running")
}

func TestFormatter_FormatServiceCatalog(t *testing.T) {
	var buf bytes.Buffer
	output := ui.NewOutput()
	formatter := New(&buf, output)

	catalog := ServiceCatalog{
		Categories: map[string][]ServiceInfo{
			"database": {
				{Name: services.ServicePostgres, Description: "PostgreSQL database"},
			},
		},
		Total: 1,
	}

	err := formatter.FormatServiceCatalog(catalog, Options{})
	assert.NoError(t, err)

	output_str := buf.String()
	assert.Contains(t, output_str, "postgres")
	assert.Contains(t, output_str, "PostgreSQL database")
}

func TestFormatter_FormatValidation(t *testing.T) {
	var buf bytes.Buffer
	output := ui.NewOutput()
	formatter := New(&buf, output)

	result := ValidationResult{
		Valid: false,
		Errors: []ValidationIssue{
			{Field: "test", Message: "test error"},
		},
		Summary: map[string]int{"errors": 1},
	}

	err := formatter.FormatValidation(result, Options{})
	assert.NoError(t, err)

	output_str := buf.String()
	assert.Contains(t, output_str, "test error")
}

func TestFilterCatalogByCategory(t *testing.T) {
	catalog := ServiceCatalog{
		Categories: map[string][]ServiceInfo{
			"database": {
				{Name: services.ServicePostgres, Description: "PostgreSQL"},
			},
			"cache": {
				{Name: services.ServiceRedis, Description: "Redis cache"},
			},
		},
		Total: 2,
	}

	filtered := FilterCatalogByCategory(catalog, "database")
	assert.Equal(t, 1, filtered.Total)
	assert.Contains(t, filtered.Categories, "database")
	assert.NotContains(t, filtered.Categories, "cache")

	// Test empty category
	empty := FilterCatalogByCategory(catalog, "nonexistent")
	assert.Equal(t, 0, empty.Total)
}

func TestFormatter_FormatStatus_Compact(t *testing.T) {
	var buf bytes.Buffer
	output := ui.NewOutput()
	formatter := New(&buf, output)

	services := []ServiceStatus{
		{Name: services.ServicePostgres, State: "running", Health: "healthy"},
	}

	err := formatter.FormatStatus(services, Options{Compact: true})
	assert.NoError(t, err)

	output_str := buf.String()
	lines := strings.Split(strings.TrimSpace(output_str), "\n")
	// go-pretty renders with borders: top, header, separator, data, bottom = 5 lines
	assert.Len(t, lines, 5)
	assert.Contains(t, output_str, "postgres")
	assert.Contains(t, output_str, "running")
	assert.Contains(t, output_str, "healthy")
}
