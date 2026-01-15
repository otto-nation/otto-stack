//go:build unit

package compose_test

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/compose"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestNewGenerator(t *testing.T) {
	t.Run("creates generator successfully", func(t *testing.T) {
		generator, err := compose.NewGenerator("test-project")
		require.NoError(t, err)
		assert.NotNil(t, generator)
	})

	t.Run("handles empty project name", func(t *testing.T) {
		generator, err := compose.NewGenerator("")
		require.NoError(t, err)
		assert.NotNil(t, generator)
	})
}

func TestGenerator_BuildComposeData(t *testing.T) {
	generator, err := compose.NewGenerator("test-project")
	require.NoError(t, err)

	t.Run("generates valid YAML from empty configs", func(t *testing.T) {
		data, err := generator.BuildComposeData([]types.ServiceConfig{})
		require.NoError(t, err)
		assert.NotEmpty(t, data)

		// Verify it's valid YAML
		var result map[string]any
		err = yaml.Unmarshal(data, &result)
		require.NoError(t, err)
		assert.Contains(t, result, "services")
		assert.Contains(t, result, "networks")
	})

	t.Run("generates YAML with service configs", func(t *testing.T) {
		configs := []types.ServiceConfig{
			{
				Name: "test-service",
				Container: types.ContainerSpec{
					Image: "test:latest",
				},
			},
		}

		data, err := generator.BuildComposeData(configs)
		require.NoError(t, err)
		assert.NotEmpty(t, data)

		var result map[string]any
		err = yaml.Unmarshal(data, &result)
		require.NoError(t, err)

		services := result["services"].(map[string]any)
		assert.Contains(t, services, "test-service")
	})
}
