//go:build unit

package project

import (
	"testing"

	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestDepsHandler_LoadServices(t *testing.T) {
	handler := NewDepsHandler()

	t.Run("loads all services", func(t *testing.T) {
		services, err := handler.loadServices()
		assert.NoError(t, err)
		assert.NotEmpty(t, services)
	})
}

func TestDepsHandler_FormatDependencies(t *testing.T) {
	handler := NewDepsHandler()

	t.Run("formats empty dependencies", func(t *testing.T) {
		result := handler.formatDependencies([]string{})
		assert.Equal(t, "none", result)
	})

	t.Run("formats single dependency", func(t *testing.T) {
		result := handler.formatDependencies([]string{"network"})
		assert.Equal(t, "network", result)
	})

	t.Run("formats multiple dependencies", func(t *testing.T) {
		result := handler.formatDependencies([]string{"network", "storage"})
		assert.Contains(t, result, "network")
		assert.Contains(t, result, "storage")
	})
}

func TestDepsHandler_BuildDependencyRows(t *testing.T) {
	handler := NewDepsHandler()

	t.Run("builds rows for existing services", func(t *testing.T) {
		allServices := map[string]servicetypes.ServiceConfig{
			"postgres": {
				Name: "postgres",
				Service: servicetypes.ServiceSpec{
					Dependencies: servicetypes.DependenciesSpec{
						Required: []string{"network"},
					},
				},
			},
		}

		rows := handler.buildDependencyRows([]string{"postgres"}, allServices)
		assert.Len(t, rows, 1)
		assert.Equal(t, "postgres", rows[0][0])
	})

	t.Run("skips non-existent services", func(t *testing.T) {
		allServices := map[string]servicetypes.ServiceConfig{
			"postgres": {Name: "postgres"},
		}

		rows := handler.buildDependencyRows([]string{"postgres", "nonexistent"}, allServices)
		assert.Len(t, rows, 1)
	})

	t.Run("handles empty service list", func(t *testing.T) {
		allServices := map[string]servicetypes.ServiceConfig{}
		rows := handler.buildDependencyRows([]string{}, allServices)
		assert.Empty(t, rows)
	})
}
