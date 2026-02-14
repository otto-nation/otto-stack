//go:build unit

package project

import (
	"testing"

	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestProjectManager_filterProjectServices(t *testing.T) {
	pm := &ProjectManager{}

	t.Run("returns all services when sharing is nil", func(t *testing.T) {
		configs := []types.ServiceConfig{
			{Name: "postgres", Shareable: true},
			{Name: "redis", Shareable: true},
			{Name: "app", Shareable: false},
		}

		result := pm.filterProjectServices(configs, nil)
		assert.Len(t, result, 3)
	})

	t.Run("returns all services when sharing is disabled", func(t *testing.T) {
		configs := []types.ServiceConfig{
			{Name: "postgres", Shareable: true},
			{Name: "redis", Shareable: true},
			{Name: "app", Shareable: false},
		}
		sharing := &clicontext.SharingSpec{Enabled: false}

		result := pm.filterProjectServices(configs, sharing)
		assert.Len(t, result, 3)
	})

	t.Run("filters out shareable services when sharing is enabled", func(t *testing.T) {
		configs := []types.ServiceConfig{
			{Name: "postgres", Shareable: true},
			{Name: "redis", Shareable: true},
			{Name: "app", Shareable: false},
		}
		sharing := &clicontext.SharingSpec{Enabled: true}

		result := pm.filterProjectServices(configs, sharing)
		assert.Len(t, result, 1)
		assert.Equal(t, "app", result[0].Name)
	})

	t.Run("returns empty when all services are shareable", func(t *testing.T) {
		configs := []types.ServiceConfig{
			{Name: "postgres", Shareable: true},
			{Name: "redis", Shareable: true},
		}
		sharing := &clicontext.SharingSpec{Enabled: true}

		result := pm.filterProjectServices(configs, sharing)
		assert.Len(t, result, 0)
	})

	t.Run("returns all when no services are shareable", func(t *testing.T) {
		configs := []types.ServiceConfig{
			{Name: "app1", Shareable: false},
			{Name: "app2", Shareable: false},
		}
		sharing := &clicontext.SharingSpec{Enabled: true}

		result := pm.filterProjectServices(configs, sharing)
		assert.Len(t, result, 2)
	})
}

func TestProjectManager_formatServicesList(t *testing.T) {
	pm := &ProjectManager{}

	t.Run("formats single service", func(t *testing.T) {
		result := pm.formatServicesList([]string{"postgres"})
		assert.Equal(t, "- postgres\n", result)
	})

	t.Run("formats multiple services", func(t *testing.T) {
		result := pm.formatServicesList([]string{"postgres", "redis", "mysql"})
		assert.Contains(t, result, "- postgres\n")
		assert.Contains(t, result, "- redis\n")
		assert.Contains(t, result, "- mysql\n")
	})

	t.Run("handles empty list", func(t *testing.T) {
		result := pm.formatServicesList([]string{})
		assert.Equal(t, "", result)
	})
}
