package compose

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGenerator(t *testing.T) {
	t.Run("creates generator successfully", func(t *testing.T) {
		gen, err := NewGenerator("test-project")
		require.NoError(t, err)
		assert.NotNil(t, gen)
		assert.Equal(t, "test-project", gen.projectName)
	})

	t.Run("creates generator with empty name", func(t *testing.T) {
		gen, err := NewGenerator("")
		require.NoError(t, err)
		assert.NotNil(t, gen)
	})
}

func TestGenerator_BuildComposeData(t *testing.T) {
	t.Run("builds compose with empty services", func(t *testing.T) {
		gen, err := NewGenerator("test-project")
		require.NoError(t, err)

		data, err := gen.BuildComposeData([]types.ServiceConfig{})
		require.NoError(t, err)
		assert.NotEmpty(t, data)
		assert.Contains(t, string(data), "test-project")
	})

	t.Run("builds compose with single service", func(t *testing.T) {
		gen, err := NewGenerator("test-project")
		require.NoError(t, err)

		services := []types.ServiceConfig{
			{
				Name: "redis",
				Container: types.ContainerSpec{
					Image: "redis:latest",
					Ports: []types.PortSpec{
						{External: "6379", Internal: "6379"},
					},
				},
			},
		}

		data, err := gen.BuildComposeData(services)
		require.NoError(t, err)
		assert.NotEmpty(t, data)
		assert.Contains(t, string(data), "redis")
		assert.Contains(t, string(data), "redis:latest")
	})

	t.Run("fails with empty project name", func(t *testing.T) {
		gen, err := NewGenerator("")
		require.NoError(t, err)

		_, err = gen.BuildComposeData([]types.ServiceConfig{})
		assert.Error(t, err, "Should error with empty project name")
	})
}

func TestGenerator_BuildComposeDataWithHeader(t *testing.T) {
	t.Run("builds compose with header", func(t *testing.T) {
		gen, err := NewGenerator("test-project")
		require.NoError(t, err)

		header := "# Test Header\n"
		data, err := gen.BuildComposeDataWithHeader([]types.ServiceConfig{}, header)
		require.NoError(t, err)
		assert.NotEmpty(t, data)
		assert.Contains(t, string(data), "# Test Header")
	})

	t.Run("builds compose without header", func(t *testing.T) {
		gen, err := NewGenerator("test-project")
		require.NoError(t, err)

		data, err := gen.BuildComposeDataWithHeader([]types.ServiceConfig{}, "")
		require.NoError(t, err)
		assert.NotEmpty(t, data)
	})
}

func TestGenerator_BuildServicesFromConfigs(t *testing.T) {
	gen, err := NewGenerator("test-project")
	require.NoError(t, err)

	t.Run("handles service with environment variables", func(t *testing.T) {
		services := []types.ServiceConfig{
			{
				Name: "postgres",
				Container: types.ContainerSpec{
					Image: "postgres:latest",
					Environment: map[string]string{
						"POSTGRES_PASSWORD": "secret",
						"POSTGRES_USER":     "admin",
					},
				},
			},
		}

		result, err := gen.buildServicesFromConfigs(services)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Contains(t, result, "postgres")
	})

	t.Run("handles service with volumes", func(t *testing.T) {
		services := []types.ServiceConfig{
			{
				Name: "mysql",
				Container: types.ContainerSpec{
					Image: "mysql:latest",
					Volumes: []types.VolumeSpec{
						{Mount: "./data:/var/lib/mysql"},
					},
				},
			},
		}

		result, err := gen.buildServicesFromConfigs(services)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Contains(t, result, "mysql")
	})

	t.Run("handles configuration service type", func(t *testing.T) {
		services := []types.ServiceConfig{
			{
				Name:        "localstack-sns",
				ServiceType: types.ServiceTypeConfiguration,
			},
		}

		result, err := gen.buildServicesFromConfigs(services)
		require.NoError(t, err)
		// Configuration services should be skipped
		assert.NotContains(t, result, "localstack-sns")
	})
}
