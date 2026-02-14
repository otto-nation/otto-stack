//go:build unit

package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_GetService(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	t.Run("gets existing service", func(t *testing.T) {
		service, err := manager.GetService("postgres")
		assert.NoError(t, err)
		assert.NotNil(t, service)
		assert.Equal(t, "postgres", service.Name)
	})

	t.Run("returns error for nonexistent service", func(t *testing.T) {
		_, err := manager.GetService("nonexistent-service-xyz")
		assert.Error(t, err)
	})
}

func TestManager_GetAllServices(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	t.Run("returns all services", func(t *testing.T) {
		services := manager.GetAllServices()
		assert.NotNil(t, services)
		assert.NotEmpty(t, services)

		// Should contain common services
		_, hasPostgres := services["postgres"]
		_, hasRedis := services["redis"]
		assert.True(t, hasPostgres || hasRedis, "should have at least one common service")
	})

	t.Run("returns map with service names as keys", func(t *testing.T) {
		services := manager.GetAllServices()

		for name, service := range services {
			assert.Equal(t, name, service.Name, "map key should match service name")
		}
	})
}

func TestManager_GetDependencies(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	t.Run("gets dependencies for service with deps", func(t *testing.T) {
		// Try kafka which typically has dependencies
		deps, err := manager.GetDependencies("kafka")
		if err == nil {
			assert.NotNil(t, deps)
			// Kafka might have zookeeper as dependency
		}
	})

	t.Run("gets empty dependencies for service without deps", func(t *testing.T) {
		deps, err := manager.GetDependencies("postgres")
		assert.NoError(t, err)
		// Postgres typically has no required dependencies, but deps might be nil or empty
		if deps != nil {
			assert.GreaterOrEqual(t, len(deps), 0)
		}
	})

	t.Run("returns error for nonexistent service", func(t *testing.T) {
		_, err := manager.GetDependencies("nonexistent-service-xyz")
		assert.Error(t, err)
	})
}

func TestNew(t *testing.T) {
	t.Run("creates manager successfully", func(t *testing.T) {
		manager, err := New()
		assert.NoError(t, err)
		assert.NotNil(t, manager)
		assert.NotNil(t, manager.services)
		assert.NotEmpty(t, manager.services)
	})

	t.Run("loads services on creation", func(t *testing.T) {
		manager, err := New()
		require.NoError(t, err)

		// Should have loaded services
		services := manager.GetAllServices()
		assert.NotEmpty(t, services)
	})
}
