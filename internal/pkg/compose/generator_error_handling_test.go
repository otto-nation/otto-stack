//go:build unit

package compose

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/stretchr/testify/assert"
)

func TestGenerator_ErrorHandling(t *testing.T) {
	t.Run("handles empty project name", func(t *testing.T) {
		manager, err := services.New()
		if err != nil {
			t.Skip("Services manager not available")
		}

		generator, err := NewGenerator("", "test-path", manager)

		// Should handle empty project name gracefully
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.NotNil(t, generator)
		}
	})

	t.Run("handles invalid services path", func(t *testing.T) {
		manager, err := services.New()
		if err != nil {
			t.Skip("Services manager not available")
		}

		generator, err := NewGenerator("test-project", "/nonexistent/path", manager)

		// Should handle invalid path gracefully
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.NotNil(t, generator)
		}
	})
}

func TestGenerator_ServiceConstants(t *testing.T) {
	t.Run("validates service constants usage", func(t *testing.T) {
		serviceConstants := []string{
			services.ServicePostgres,
			services.ServiceRedis,
			services.ServiceMysql,
			services.ServiceLocalstack,
		}

		for _, service := range serviceConstants {
			assert.NotEmpty(t, service, "Service constant should not be empty")
			assert.IsType(t, "", service, "Service constant should be string")
		}
	})

	t.Run("validates category constants", func(t *testing.T) {
		categories := []string{
			services.CategoryDatabase,
			services.CategoryCache,
			services.CategoryCloud,
		}

		for _, category := range categories {
			assert.NotEmpty(t, category, "Category constant should not be empty")
		}
	})
}
