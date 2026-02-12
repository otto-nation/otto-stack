package services

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_GetAllServices(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	services := manager.GetAllServices()
	assert.NotEmpty(t, services, "Should return at least some services")

	// Verify services have required fields
	for name, service := range services {
		assert.NotEmpty(t, name, "Service name should not be empty")
		assert.Equal(t, name, service.Name, "Service name should match map key")
		assert.NotEmpty(t, service.Description, "Service should have description")
	}
}

func TestManager_GetDependencies(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	t.Run("service with dependencies", func(t *testing.T) {
		// localstack-sns depends on localstack
		deps, err := manager.GetDependencies("localstack-sns")
		require.NoError(t, err)
		assert.Contains(t, deps, "localstack", "localstack-sns should depend on localstack")
	})

	t.Run("service without dependencies", func(t *testing.T) {
		deps, err := manager.GetDependencies("redis")
		require.NoError(t, err)
		assert.Empty(t, deps, "redis should have no dependencies")
	})

	t.Run("nonexistent service", func(t *testing.T) {
		_, err := manager.GetDependencies("nonexistent-service")
		assert.Error(t, err, "Should error for nonexistent service")
	})
}

func TestResolveUpServices(t *testing.T) {
	t.Run("resolves single service", func(t *testing.T) {
		services, err := ResolveUpServices([]string{"redis"}, nil)
		require.NoError(t, err)
		assert.Len(t, services, 1)
		assert.Equal(t, "redis", services[0].Name)
	})

	t.Run("resolves service with dependencies", func(t *testing.T) {
		services, err := ResolveUpServices([]string{"localstack-sns"}, nil)
		require.NoError(t, err)

		// Should include both localstack-sns and its dependency localstack
		assert.GreaterOrEqual(t, len(services), 2, "Should include service and dependencies")

		serviceNames := make(map[string]bool)
		for _, svc := range services {
			serviceNames[svc.Name] = true
		}

		assert.True(t, serviceNames["localstack-sns"], "Should include localstack-sns")
		assert.True(t, serviceNames["localstack"], "Should include dependency localstack")
	})

	t.Run("resolves multiple services", func(t *testing.T) {
		services, err := ResolveUpServices([]string{"redis", "postgres"}, nil)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(services), 2)

		serviceNames := make(map[string]bool)
		for _, svc := range services {
			serviceNames[svc.Name] = true
		}

		assert.True(t, serviceNames["redis"])
		assert.True(t, serviceNames["postgres"])
	})

	t.Run("handles nonexistent service", func(t *testing.T) {
		_, err := ResolveUpServices([]string{"nonexistent"}, nil)
		assert.Error(t, err, "Should error for nonexistent service")
	})

	t.Run("handles empty service list with config", func(t *testing.T) {
		// When args is empty, it uses config.Stack.Enabled
		// Empty enabled list should error (validation requires at least one service)
		cfg := &config.Config{
			Stack: config.StackConfig{
				Enabled: []string{},
			},
		}
		_, err := ResolveUpServices([]string{}, cfg)
		assert.Error(t, err, "Should error when no services are provided")
	})
}

func TestManager_ExecuteCustomOperation(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	t.Run("nonexistent service", func(t *testing.T) {
		_, err := manager.ExecuteCustomOperation("nonexistent", "some-operation")
		assert.Error(t, err, "Should error for nonexistent service")
	})

	t.Run("service without custom operations", func(t *testing.T) {
		// Most services don't have custom operations
		_, err := manager.ExecuteCustomOperation("redis", "some-operation")
		assert.Error(t, err, "Should error when operation doesn't exist")
	})
}
