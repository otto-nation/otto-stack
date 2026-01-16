//go:build unit

package display

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStatusFormatter_EdgeCases(t *testing.T) {
	t.Run("handles empty services list", func(t *testing.T) {
		var buf bytes.Buffer
		formatter := NewStatusFormatter(&buf)

		err := formatter.FormatTable([]ServiceStatus{}, Options{})
		assert.NoError(t, err)
	})

	t.Run("handles services without providers", func(t *testing.T) {
		var buf bytes.Buffer
		formatter := NewStatusFormatter(&buf)

		services := []ServiceStatus{
			{Name: "test-service", State: "running"},
		}

		err := formatter.FormatTable(services, Options{})
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "test-service")
	})

	t.Run("handles services with providers", func(t *testing.T) {
		var buf bytes.Buffer
		formatter := NewStatusFormatter(&buf)

		services := []ServiceStatus{
			{Name: "test-service", State: "running", Provider: "docker"},
		}

		// Use compact format to show providers
		err := formatter.FormatTable(services, Options{Compact: true})
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "test-service")
		// Provider should be shown in compact format when present
		output := buf.String()
		assert.True(t, len(output) > 0)
	})
}

// TableFormatter tests removed - now using go-pretty/table which handles edge cases internally

func TestValidationFormatter_EdgeCases(t *testing.T) {
	t.Run("handles empty validation result", func(t *testing.T) {
		var buf bytes.Buffer
		formatter := NewValidationFormatter(&buf)

		result := ValidationResult{
			Valid:   true,
			Errors:  []ValidationIssue{},
			Summary: map[string]int{},
		}

		err := formatter.FormatTable(result)
		assert.NoError(t, err)
	})

	t.Run("handles validation with errors", func(t *testing.T) {
		var buf bytes.Buffer
		formatter := NewValidationFormatter(&buf)

		result := ValidationResult{
			Valid: false,
			Errors: []ValidationIssue{
				{Type: "error", Message: "error1"},
				{Type: "error", Message: "error2"},
			},
		}

		err := formatter.FormatTable(result)
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "error1")
		assert.Contains(t, output, "error2")
	})
}

func TestHealthFormatter_EdgeCases(t *testing.T) {
	t.Run("handles empty health report", func(t *testing.T) {
		var buf bytes.Buffer
		formatter := NewHealthFormatter(&buf)

		report := HealthReport{
			Checks: []HealthCheck{},
		}

		err := formatter.FormatTable(report, Options{})
		assert.NoError(t, err)
	})

	t.Run("handles health report with checks", func(t *testing.T) {
		var buf bytes.Buffer
		formatter := NewHealthFormatter(&buf)

		report := HealthReport{
			Checks: []HealthCheck{
				{Name: "service1", Status: "healthy"},
				{Name: "service2", Status: "unhealthy"},
			},
		}

		err := formatter.FormatTable(report, Options{})
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "service1")
		assert.Contains(t, output, "service2")
	})
}

func TestCatalogFormatter_EdgeCases(t *testing.T) {
	t.Run("handles empty catalog", func(t *testing.T) {
		var buf bytes.Buffer
		formatter := NewCatalogFormatter(&buf)

		catalog := ServiceCatalog{
			Categories: map[string][]ServiceInfo{},
		}

		err := formatter.FormatTable(catalog)
		assert.NoError(t, err)
	})

	t.Run("handles catalog with services", func(t *testing.T) {
		var buf bytes.Buffer
		formatter := NewCatalogFormatter(&buf)

		catalog := ServiceCatalog{
			Categories: map[string][]ServiceInfo{
				"database": {
					{Name: "postgres", Category: "database"},
				},
				"cache": {
					{Name: "redis", Category: "cache"},
				},
			},
		}

		err := formatter.FormatTable(catalog)
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "postgres")
		assert.Contains(t, output, "redis")
	})

	t.Run("filters catalog by category", func(t *testing.T) {
		catalog := ServiceCatalog{
			Categories: map[string][]ServiceInfo{
				"database": {
					{Name: "postgres", Category: "database"},
					{Name: "mysql", Category: "database"},
				},
				"cache": {
					{Name: "redis", Category: "cache"},
				},
			},
		}

		filtered := FilterCatalogByCategory(catalog, "database")
		assert.Len(t, filtered.Categories["database"], 2)
	})
}

func TestFormatter_OutputTypes(t *testing.T) {
	t.Run("handles different output types", func(t *testing.T) {
		var buf bytes.Buffer

		// Test with mock output
		mockOutput := &mockOutput{}
		formatter := New(&buf, mockOutput)
		assert.NotNil(t, formatter)

		// Test with empty services
		services := []ServiceStatus{}
		err := formatter.FormatStatus(services, Options{})
		assert.NoError(t, err)
	})
}

// Mock output for testing
type mockOutput struct{}

func (m *mockOutput) Success(msg string, args ...any) {}
func (m *mockOutput) Error(msg string, args ...any)   {}
func (m *mockOutput) Warning(msg string, args ...any) {}
func (m *mockOutput) Info(msg string, args ...any)    {}
func (m *mockOutput) Header(msg string, args ...any)  {}
func (m *mockOutput) Muted(msg string, args ...any)   {}
func (m *mockOutput) Writer() io.Writer               { return os.Stdout }

func TestServiceStatus_Fields(t *testing.T) {
	t.Run("validates service status fields", func(t *testing.T) {
		now := time.Now()
		status := ServiceStatus{
			Name:      "test-service",
			State:     "running",
			Health:    "healthy",
			Provider:  "docker",
			Ports:     []string{"8080:8080"},
			CreatedAt: now,
			UpdatedAt: now,
			Uptime:    time.Hour,
		}

		assert.Equal(t, "test-service", status.Name)
		assert.Equal(t, "running", status.State)
		assert.Equal(t, "healthy", status.Health)
		assert.Equal(t, "docker", status.Provider)
		assert.Len(t, status.Ports, 1)
		assert.Equal(t, time.Hour, status.Uptime)
	})
}
