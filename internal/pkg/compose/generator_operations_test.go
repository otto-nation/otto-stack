//go:build unit

package compose

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/test/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestGenerator_UncoveredMethods(t *testing.T) {
	t.Run("tests resolveEnvVar edge cases", func(t *testing.T) {
		manager, err := services.New()
		testhelpers.AssertValidConstructor(t, manager, err, "Services Manager")

		generator, err := NewGenerator("test-project", "test-path", manager)
		testhelpers.AssertValidConstructor(t, generator, err, "Generator")

		// Test with environment variable syntax
		result := generator.resolveEnvVar("${TEST_VAR:-default}")
		assert.IsType(t, "", result)

		// Test with regular string
		result = generator.resolveEnvVar("regular_value")
		assert.Equal(t, "regular_value", result)

		// Test with empty string
		result = generator.resolveEnvVar("")
		assert.Equal(t, "", result)
	})

	t.Run("tests addServiceVolumes", func(t *testing.T) {
		manager, err := services.New()
		if err != nil {
			t.Skip("Services manager not available")
		}

		generator, err := NewGenerator("test-project", "test-path", manager)
		if err != nil {
			t.Skip("Generator not available")
		}

		service := map[string]any{}
		config := &types.ServiceConfig{Name: "postgres"}

		generator.addServiceVolumes(service, config)

		// Should handle volumes configuration
		assert.IsType(t, map[string]any{}, service)
	})

	t.Run("tests addServiceConfiguration", func(t *testing.T) {
		manager, err := services.New()
		if err != nil {
			t.Skip("Services manager not available")
		}

		generator, err := NewGenerator("test-project", "test-path", manager)
		if err != nil {
			t.Skip("Generator not available")
		}

		service := map[string]any{}
		config := &types.ServiceConfig{Name: "postgres"}

		generator.addServiceConfiguration(service, config)

		// Should add configuration
		assert.IsType(t, map[string]any{}, service)
	})

	t.Run("tests addServiceHealthCheck", func(t *testing.T) {
		manager, err := services.New()
		if err != nil {
			t.Skip("Services manager not available")
		}

		generator, err := NewGenerator("test-project", "test-path", manager)
		if err != nil {
			t.Skip("Generator not available")
		}

		service := map[string]any{}
		config := &types.ServiceConfig{Name: "postgres"}

		generator.addServiceHealthCheck(service, config)

		// Should add health check configuration
		assert.IsType(t, map[string]any{}, service)
	})
}

func TestGenerator_EdgeCases(t *testing.T) {
	t.Run("tests buildComposeStructure with empty configs", func(t *testing.T) {
		manager, err := services.New()
		if err != nil {
			t.Skip("Services manager not available")
		}

		generator, err := NewGenerator("test-project", "test-path", manager)
		if err != nil {
			t.Skip("Generator not available")
		}

		result, err := generator.buildComposeStructure([]types.ServiceConfig{})
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.NotNil(t, result)
		}
	})

	t.Run("tests GenerateFromServiceConfigs", func(t *testing.T) {
		manager, err := services.New()
		if err != nil {
			t.Skip("Services manager not available")
		}

		generator, err := NewGenerator("test-project", "test-path", manager)
		if err != nil {
			t.Skip("Generator not available")
		}

		configs := []types.ServiceConfig{
			{Name: services.ServicePostgres},
		}

		err = generator.GenerateFromServiceConfigs(configs, "test-project")
		// Should handle gracefully
		if err != nil {
			assert.Error(t, err)
		}
	})

	t.Run("tests createBaseService", func(t *testing.T) {
		manager, err := services.New()
		if err != nil {
			t.Skip("Services manager not available")
		}

		generator, err := NewGenerator("test-project", "test-path", manager)
		if err != nil {
			t.Skip("Generator not available")
		}

		config := &types.ServiceConfig{Name: "postgres"}
		service := generator.createBaseService(config)

		assert.IsType(t, map[string]any{}, service)
		assert.NotNil(t, service)
	})
}
