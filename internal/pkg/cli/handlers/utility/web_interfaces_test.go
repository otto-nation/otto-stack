//go:build unit

package utility

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestWebInterfacesHandler_extractWebInterfaces(t *testing.T) {
	handler := NewWebInterfacesHandler()

	t.Run("extracts interfaces from all services when showAll is true", func(t *testing.T) {
		configs := []types.ServiceConfig{
			{
				Name: "postgres",
				Documentation: types.DocumentationSpec{
					WebInterfaces: []types.WebInterface{
						{Name: "Admin", URL: "http://localhost:5432", Description: "DB Admin"},
					},
				},
			},
			{
				Name: "redis",
				Documentation: types.DocumentationSpec{
					WebInterfaces: []types.WebInterface{
						{Name: "Commander", URL: "http://localhost:6379", Description: "Redis UI"},
					},
				},
			},
		}

		result := handler.extractWebInterfaces(configs, nil, true)
		assert.Len(t, result, 2)
		assert.Equal(t, "postgres", result[0].Service)
		assert.Equal(t, "Admin", result[0].Name)
		assert.Equal(t, "redis", result[1].Service)
	})

	t.Run("filters by running services when showAll is false", func(t *testing.T) {
		configs := []types.ServiceConfig{
			{
				Name: "postgres",
				Documentation: types.DocumentationSpec{
					WebInterfaces: []types.WebInterface{
						{Name: "Admin", URL: "http://localhost:5432"},
					},
				},
			},
			{
				Name: "redis",
				Documentation: types.DocumentationSpec{
					WebInterfaces: []types.WebInterface{
						{Name: "Commander", URL: "http://localhost:6379"},
					},
				},
			},
		}
		running := map[string]bool{"postgres": true, "redis": false}

		result := handler.extractWebInterfaces(configs, running, false)
		assert.Len(t, result, 1)
		assert.Equal(t, "postgres", result[0].Service)
	})

	t.Run("handles services with no web interfaces", func(t *testing.T) {
		configs := []types.ServiceConfig{
			{Name: "postgres", Documentation: types.DocumentationSpec{WebInterfaces: []types.WebInterface{}}},
		}

		result := handler.extractWebInterfaces(configs, nil, true)
		assert.Len(t, result, 0)
	})

	t.Run("handles multiple interfaces per service", func(t *testing.T) {
		configs := []types.ServiceConfig{
			{
				Name: "grafana",
				Documentation: types.DocumentationSpec{
					WebInterfaces: []types.WebInterface{
						{Name: "Dashboard", URL: "http://localhost:3000"},
						{Name: "API", URL: "http://localhost:3000/api"},
					},
				},
			},
		}

		result := handler.extractWebInterfaces(configs, nil, true)
		assert.Len(t, result, 2)
		assert.Equal(t, "Dashboard", result[0].Name)
		assert.Equal(t, "API", result[1].Name)
	})
}

func TestWebInterfacesHandler_createWebInterfaces(t *testing.T) {
	handler := NewWebInterfacesHandler()

	t.Run("creates web interfaces from config", func(t *testing.T) {
		webIfaces := []types.WebInterface{
			{Name: "Admin", URL: "http://localhost:8080", Description: "Admin Panel"},
			{Name: "API", URL: "http://localhost:8080/api", Description: "REST API"},
		}

		result := handler.createWebInterfaces("myservice", webIfaces)
		assert.Len(t, result, 2)
		assert.Equal(t, "myservice", result[0].Service)
		assert.Equal(t, "Admin", result[0].Name)
		assert.Equal(t, "http://localhost:8080", result[0].URL)
		assert.Equal(t, "Admin Panel", result[0].Description)
		assert.Equal(t, "API", result[1].Name)
	})

	t.Run("handles empty web interfaces", func(t *testing.T) {
		result := handler.createWebInterfaces("myservice", []types.WebInterface{})
		assert.Len(t, result, 0)
	})
}

func TestWebInterfacesHandler_shouldIncludeService(t *testing.T) {
	handler := NewWebInterfacesHandler()

	t.Run("includes all services when showAll is true", func(t *testing.T) {
		assert.True(t, handler.shouldIncludeService("postgres", nil, true))
		assert.True(t, handler.shouldIncludeService("redis", map[string]bool{"redis": false}, true))
	})

	t.Run("includes all services when runningServices is nil", func(t *testing.T) {
		assert.True(t, handler.shouldIncludeService("postgres", nil, false))
	})

	t.Run("filters by running status when showAll is false", func(t *testing.T) {
		running := map[string]bool{"postgres": true, "redis": false}
		assert.True(t, handler.shouldIncludeService("postgres", running, false))
		assert.False(t, handler.shouldIncludeService("redis", running, false))
	})

	t.Run("excludes services not in running map", func(t *testing.T) {
		running := map[string]bool{"postgres": true}
		assert.False(t, handler.shouldIncludeService("redis", running, false))
	})
}

func TestWebInterfacesHandler_formatStatus(t *testing.T) {
	handler := NewWebInterfacesHandler()

	t.Run("formats available status", func(t *testing.T) {
		result := handler.formatStatus(true)
		assert.Contains(t, result, "Available")
	})

	t.Run("formats unavailable status", func(t *testing.T) {
		result := handler.formatStatus(false)
		assert.Contains(t, result, "Not Available")
	})
}

func TestWebInterfacesHandler_formatStatusFromResponse(t *testing.T) {
	handler := NewWebInterfacesHandler()

	t.Run("formats status from 200 response", func(t *testing.T) {
		result := handler.formatStatusFromResponse(200)
		assert.Contains(t, result, "Available")
	})

	t.Run("formats status from 404 response", func(t *testing.T) {
		result := handler.formatStatusFromResponse(404)
		assert.Contains(t, result, "Not Available")
	})

	t.Run("formats status from 500 response", func(t *testing.T) {
		result := handler.formatStatusFromResponse(500)
		assert.Contains(t, result, "Not Available")
	})
}

func TestWebInterfacesHandler_ValidateArgs(t *testing.T) {
	handler := NewWebInterfacesHandler()

	t.Run("validates empty args", func(t *testing.T) {
		err := handler.ValidateArgs([]string{})
		assert.NoError(t, err)
	})

	t.Run("validates with args", func(t *testing.T) {
		err := handler.ValidateArgs([]string{"postgres", "redis"})
		assert.NoError(t, err)
	})
}
