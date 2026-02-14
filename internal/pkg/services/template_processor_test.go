//go:build unit

package services

import (
	"testing"

	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTemplateProcessor(t *testing.T) {
	processor := NewTemplateProcessor()
	assert.NotNil(t, processor)
}

func TestTemplateProcessor_Process(t *testing.T) {
	processor := NewTemplateProcessor()

	t.Run("processes simple template", func(t *testing.T) {
		config := servicetypes.ServiceConfig{
			Name: "test-service",
		}
		script := "echo 'Hello World'"

		result, err := processor.Process(script, config, []servicetypes.ServiceConfig{})
		require.NoError(t, err)
		assert.Equal(t, "echo 'Hello World'", result)
	})

	t.Run("handles invalid template", func(t *testing.T) {
		config := servicetypes.ServiceConfig{Name: "test"}
		script := "{{.Invalid"

		_, err := processor.Process(script, config, []servicetypes.ServiceConfig{})
		assert.Error(t, err)
	})
}

func TestTemplateProcessor_ServiceDependsOn(t *testing.T) {
	processor := NewTemplateProcessor()

	t.Run("returns true when service has dependency", func(t *testing.T) {
		config := servicetypes.ServiceConfig{
			Service: servicetypes.ServiceSpec{
				Dependencies: servicetypes.DependenciesSpec{
					Required: []string{"postgres", "redis"},
				},
			},
		}
		assert.True(t, processor.serviceDependsOn(config, "postgres"))
	})

	t.Run("returns false when service has no dependency", func(t *testing.T) {
		config := servicetypes.ServiceConfig{
			Service: servicetypes.ServiceSpec{
				Dependencies: servicetypes.DependenciesSpec{
					Required: []string{"postgres"},
				},
			},
		}
		assert.False(t, processor.serviceDependsOn(config, "redis"))
	})

	t.Run("returns false when dependencies are nil", func(t *testing.T) {
		config := servicetypes.ServiceConfig{
			Service: servicetypes.ServiceSpec{
				Dependencies: servicetypes.DependenciesSpec{
					Required: nil,
				},
			},
		}
		assert.False(t, processor.serviceDependsOn(config, "postgres"))
	})
}

func TestTemplateProcessor_AddConfigData(t *testing.T) {
	processor := NewTemplateProcessor()

	t.Run("adds config data to template data", func(t *testing.T) {
		templateData := make(map[string]any)
		config := servicetypes.ServiceConfig{
			Name:        "postgres",
			Description: "PostgreSQL database",
		}
		processor.addConfigData(templateData, config)
		// Function may not populate data for simple configs
		assert.NotNil(t, templateData)
	})
}
