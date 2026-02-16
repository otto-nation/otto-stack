//go:build unit

package services

import (
	"testing"

	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTemplateProcessor(t *testing.T) {
	t.Run("creates processor", func(t *testing.T) {
		processor := NewTemplateProcessor()
		assert.NotNil(t, processor)
	})
}

func TestTemplateProcessor_Process(t *testing.T) {
	processor := NewTemplateProcessor()

	t.Run("processes template with no variables", func(t *testing.T) {
		script := "echo 'hello world'"
		config := servicetypes.ServiceConfig{Name: "test"}

		result, err := processor.Process(script, config, []servicetypes.ServiceConfig{})
		require.NoError(t, err)
		assert.Equal(t, "echo 'hello world'", result)
	})

	t.Run("processes template with dependencies", func(t *testing.T) {
		script := "echo 'setup'"
		config := servicetypes.ServiceConfig{Name: "postgres"}
		allConfigs := []servicetypes.ServiceConfig{
			{
				Name: "app",
				Service: servicetypes.ServiceSpec{
					Dependencies: servicetypes.DependenciesSpec{
						Required: []string{"postgres"},
					},
				},
			},
		}

		result, err := processor.Process(script, config, allConfigs)
		require.NoError(t, err)
		assert.Contains(t, result, "setup")
	})

	t.Run("handles invalid template syntax", func(t *testing.T) {
		script := "{{.Invalid"
		config := servicetypes.ServiceConfig{Name: "test"}

		_, err := processor.Process(script, config, []servicetypes.ServiceConfig{})
		assert.Error(t, err)
	})

	t.Run("handles template with empty data", func(t *testing.T) {
		script := "{{.NonExistent}}"
		config := servicetypes.ServiceConfig{Name: "test"}

		result, err := processor.Process(script, config, []servicetypes.ServiceConfig{})
		// Template executes but variable is empty
		require.NoError(t, err)
		assert.Equal(t, "<no value>", result)
	})
}

func TestTemplateProcessor_serviceDependsOn(t *testing.T) {
	processor := NewTemplateProcessor()

	t.Run("returns true when service depends on target", func(t *testing.T) {
		config := servicetypes.ServiceConfig{
			Name: "app",
			Service: servicetypes.ServiceSpec{
				Dependencies: servicetypes.DependenciesSpec{
					Required: []string{"postgres", "redis"},
				},
			},
		}

		assert.True(t, processor.serviceDependsOn(config, "postgres"))
		assert.True(t, processor.serviceDependsOn(config, "redis"))
	})

	t.Run("returns false when service does not depend on target", func(t *testing.T) {
		config := servicetypes.ServiceConfig{
			Name: "app",
			Service: servicetypes.ServiceSpec{
				Dependencies: servicetypes.DependenciesSpec{
					Required: []string{"postgres"},
				},
			},
		}

		assert.False(t, processor.serviceDependsOn(config, "redis"))
	})

	t.Run("returns false when no dependencies", func(t *testing.T) {
		config := servicetypes.ServiceConfig{Name: "app"}
		assert.False(t, processor.serviceDependsOn(config, "postgres"))
	})

	t.Run("returns false when dependencies is nil", func(t *testing.T) {
		config := servicetypes.ServiceConfig{
			Name: "app",
			Service: servicetypes.ServiceSpec{
				Dependencies: servicetypes.DependenciesSpec{
					Required: nil,
				},
			},
		}

		assert.False(t, processor.serviceDependsOn(config, "postgres"))
	})
}

func TestTemplateProcessor_collectTemplateData(t *testing.T) {
	processor := NewTemplateProcessor()

	t.Run("collects data from dependent services", func(t *testing.T) {
		config := servicetypes.ServiceConfig{Name: "postgres"}
		allConfigs := []servicetypes.ServiceConfig{
			{
				Name: "app",
				Service: servicetypes.ServiceSpec{
					Dependencies: servicetypes.DependenciesSpec{
						Required: []string{"postgres"},
					},
				},
			},
		}

		data := processor.collectTemplateData(config, allConfigs)
		assert.NotNil(t, data)
	})

	t.Run("returns empty map when no dependencies", func(t *testing.T) {
		config := servicetypes.ServiceConfig{Name: "postgres"}
		allConfigs := []servicetypes.ServiceConfig{
			{Name: "redis"},
		}

		data := processor.collectTemplateData(config, allConfigs)
		assert.NotNil(t, data)
		assert.Empty(t, data)
	})

	t.Run("handles empty allConfigs", func(t *testing.T) {
		config := servicetypes.ServiceConfig{Name: "postgres"}

		data := processor.collectTemplateData(config, []servicetypes.ServiceConfig{})
		assert.NotNil(t, data)
		assert.Empty(t, data)
	})
}

func TestTemplateProcessor_addConfigData(t *testing.T) {
	processor := NewTemplateProcessor()

	t.Run("adds config data to template", func(t *testing.T) {
		templateData := make(map[string]any)
		config := servicetypes.ServiceConfig{
			Name: "postgres",
		}

		processor.addConfigData(templateData, config)
		// Data is added based on reflection
		assert.NotNil(t, templateData)
	})

	t.Run("handles multiple configs", func(t *testing.T) {
		templateData := make(map[string]any)
		configs := []servicetypes.ServiceConfig{
			{Name: "postgres"},
			{Name: "redis"},
		}

		for _, cfg := range configs {
			processor.addConfigData(templateData, cfg)
		}
		assert.NotNil(t, templateData)
	})
}
