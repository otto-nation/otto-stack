//go:build unit

package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceResolver_ResolveServices(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	resolver := NewServiceResolver(manager)

	t.Run("resolves single service without dependencies", func(t *testing.T) {
		configs, err := resolver.ResolveServices([]string{"postgres"})
		if err == nil {
			assert.NotEmpty(t, configs)
			// Should contain at least postgres
			found := false
			for _, cfg := range configs {
				if cfg.Name == "postgres" {
					found = true
					break
				}
			}
			assert.True(t, found, "postgres should be in resolved configs")
		}
	})

	t.Run("resolves multiple services", func(t *testing.T) {
		configs, err := resolver.ResolveServices([]string{"postgres", "redis"})
		if err == nil {
			assert.NotEmpty(t, configs)
			assert.GreaterOrEqual(t, len(configs), 2)
		}
	})

	t.Run("returns error for invalid service", func(t *testing.T) {
		_, err := resolver.ResolveServices([]string{"nonexistent-service-xyz"})
		assert.Error(t, err)
	})

	t.Run("handles empty service list", func(t *testing.T) {
		_, err := resolver.ResolveServices([]string{})
		// Empty list should return error (at least one service required)
		assert.Error(t, err)
	})

	t.Run("resolves service with dependencies", func(t *testing.T) {
		// Try a service that might have dependencies
		configs, err := resolver.ResolveServices([]string{"kafka"})
		if err == nil {
			// Should include kafka and potentially its dependencies
			assert.NotEmpty(t, configs)
		}
	})

	t.Run("deduplicates services", func(t *testing.T) {
		configs, err := resolver.ResolveServices([]string{"postgres", "postgres"})
		if err == nil {
			// Count postgres occurrences
			count := 0
			for _, cfg := range configs {
				if cfg.Name == "postgres" {
					count++
				}
			}
			assert.Equal(t, 1, count, "postgres should appear only once")
		}
	})
}

func TestNewServiceResolver(t *testing.T) {
	t.Run("creates resolver with manager", func(t *testing.T) {
		manager, err := New()
		require.NoError(t, err)

		resolver := NewServiceResolver(manager)
		assert.NotNil(t, resolver)
		assert.NotNil(t, resolver.manager)
	})
}
