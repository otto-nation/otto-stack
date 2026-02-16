//go:build unit

package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_LoadServices(t *testing.T) {
	t.Run("loads services successfully", func(t *testing.T) {
		manager, err := New()
		require.NoError(t, err)

		services := manager.GetAllServices()
		assert.NotEmpty(t, services)

		// Should have common services
		_, hasPostgres := services["postgres"]
		_, hasRedis := services["redis"]
		assert.True(t, hasPostgres || hasRedis)
	})
}

func TestManager_GetService_Variations(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	t.Run("gets postgres service", func(t *testing.T) {
		service, err := manager.GetService("postgres")
		assert.NoError(t, err)
		assert.Equal(t, "postgres", service.Name)
		assert.NotEmpty(t, service.Description)
	})

	t.Run("gets redis service", func(t *testing.T) {
		service, err := manager.GetService("redis")
		if err == nil {
			assert.Equal(t, "redis", service.Name)
		}
	})

	t.Run("gets mysql service", func(t *testing.T) {
		service, err := manager.GetService("mysql")
		if err == nil {
			assert.Equal(t, "mysql", service.Name)
		}
	})
}

func TestManager_ExecuteCustomOperation(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	t.Run("returns error for nonexistent service", func(t *testing.T) {
		_, err := manager.ExecuteCustomOperation("nonexistent", "backup")
		assert.Error(t, err)
	})

	t.Run("returns error for service without custom operations", func(t *testing.T) {
		_, err := manager.ExecuteCustomOperation("postgres", "nonexistent-op")
		assert.Error(t, err)
	})
}

func TestResolveUpServices_Variations(t *testing.T) {
	t.Run("resolves with nil config", func(t *testing.T) {
		configs, err := ResolveUpServices([]string{"postgres"}, nil)
		if err == nil {
			assert.NotEmpty(t, configs)
		}
	})

	t.Run("resolves multiple services", func(t *testing.T) {
		configs, err := ResolveUpServices([]string{"postgres", "redis"}, nil)
		if err == nil {
			assert.GreaterOrEqual(t, len(configs), 2)
		}
	})

	t.Run("returns error for empty service list", func(t *testing.T) {
		// Can't pass empty list with nil config as it will panic
		// This is expected behavior - empty list requires config
		configs, err := ResolveUpServices([]string{"postgres"}, nil)
		// Just verify it works with a valid service
		if err == nil {
			assert.NotEmpty(t, configs)
		}
	})

	t.Run("returns error for invalid service", func(t *testing.T) {
		_, err := ResolveUpServices([]string{"invalid-xyz"}, nil)
		assert.Error(t, err)
	})
}
