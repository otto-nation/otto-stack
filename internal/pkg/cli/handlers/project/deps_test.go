//go:build unit

package project

import (
	"testing"

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
