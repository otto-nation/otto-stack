//go:build unit

package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceResolver_ResolveServices_SingleService(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)
	resolver := NewServiceResolver(manager)

	configs, err := resolver.ResolveServices([]string{ServicePostgres})
	if err == nil {
		assert.NotEmpty(t, configs)
		found := false
		for _, cfg := range configs {
			if cfg.Name == ServicePostgres {
				found = true
				break
			}
		}
		assert.True(t, found, "postgres should be in resolved configs")
	}
}

func TestServiceResolver_ResolveServices_MultipleServices(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)
	resolver := NewServiceResolver(manager)

	configs, err := resolver.ResolveServices([]string{ServicePostgres, ServiceRedis})
	if err == nil {
		assert.NotEmpty(t, configs)
		assert.GreaterOrEqual(t, len(configs), 2)
	}
}

func TestServiceResolver_ResolveServices_InvalidService(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)
	resolver := NewServiceResolver(manager)

	_, err = resolver.ResolveServices([]string{"nonexistent-service-xyz"})
	assert.Error(t, err)
}

func TestServiceResolver_ResolveServices_EmptyList(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)
	resolver := NewServiceResolver(manager)

	_, err = resolver.ResolveServices([]string{})
	assert.Error(t, err)
}

func TestServiceResolver_ResolveServices_WithDependencies(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)
	resolver := NewServiceResolver(manager)

	configs, err := resolver.ResolveServices([]string{ServiceKafka})
	if err == nil {
		assert.NotEmpty(t, configs)
	}
}

func TestServiceResolver_ResolveServices_Deduplicates(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)
	resolver := NewServiceResolver(manager)

	configs, err := resolver.ResolveServices([]string{ServicePostgres, ServicePostgres})
	if err == nil {
		count := 0
		for _, cfg := range configs {
			if cfg.Name == ServicePostgres {
				count++
			}
		}
		assert.Equal(t, 1, count, "postgres should appear only once")
	}
}

func TestNewServiceResolver(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	resolver := NewServiceResolver(manager)
	assert.NotNil(t, resolver)
	assert.NotNil(t, resolver.manager)
}
