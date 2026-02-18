//go:build unit

package project

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestDepsHandler_LoadServices(t *testing.T) {
	handler := NewDepsHandler()

	svcs, err := handler.loadServices()
	assert.NoError(t, err)
	assert.NotEmpty(t, svcs)
}

func TestDepsHandler_FormatDependencies(t *testing.T) {
	handler := NewDepsHandler()

	result := handler.formatDependencies([]string{})
	assert.Equal(t, "none", result)

	result = handler.formatDependencies([]string{"network"})
	assert.Equal(t, "network", result)

	result = handler.formatDependencies([]string{"network", "storage"})
	assert.Contains(t, result, "network")
	assert.Contains(t, result, "storage")
}

func TestDepsHandler_BuildDependencyRows(t *testing.T) {
	handler := NewDepsHandler()

	cfg := servicetypes.ServiceConfig{
		Name: services.ServicePostgres,
	}
	cfg.Service.Dependencies.Required = []string{"network"}

	allServices := map[string]servicetypes.ServiceConfig{
		services.ServicePostgres: cfg,
	}

	rows := handler.buildDependencyRows([]string{services.ServicePostgres}, allServices)
	assert.Len(t, rows, 1)
	assert.Equal(t, services.ServicePostgres, rows[0][0])

	cfg = servicetypes.ServiceConfig{Name: services.ServicePostgres}
	allServices = map[string]servicetypes.ServiceConfig{
		services.ServicePostgres: cfg,
	}

	rows = handler.buildDependencyRows([]string{services.ServicePostgres, "nonexistent"}, allServices)
	assert.Len(t, rows, 1)

	allServices = map[string]servicetypes.ServiceConfig{}
	rows = handler.buildDependencyRows([]string{}, allServices)
	assert.Empty(t, rows)
}
