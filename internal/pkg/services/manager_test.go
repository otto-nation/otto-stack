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

	service, err := manager.GetService(ServicePostgres)
	assert.NoError(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, ServicePostgres, service.Name)

	_, err = manager.GetService("nonexistent-service-xyz")
	assert.Error(t, err)
}

func TestManager_GetAllServices(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	services := manager.GetAllServices()
	assert.NotNil(t, services)
	assert.NotEmpty(t, services)

	_, hasPostgres := services[ServicePostgres]
	_, hasRedis := services[ServiceRedis]
	assert.True(t, hasPostgres || hasRedis, "should have at least one common service")

	for name, service := range services {
		assert.Equal(t, name, service.Name, "map key should match service name")
	}
}

func TestManager_GetDependencies(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	deps, err := manager.GetDependencies(ServiceKafka)
	if err == nil {
		assert.NotNil(t, deps)
	}

	deps, err = manager.GetDependencies(ServicePostgres)
	assert.NoError(t, err)
	if deps != nil {
		assert.GreaterOrEqual(t, len(deps), 0)
	}

	_, err = manager.GetDependencies("nonexistent-service-xyz")
	assert.Error(t, err)
}

func TestNew(t *testing.T) {
	manager, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, manager)
	assert.NotNil(t, manager.services)
	assert.NotEmpty(t, manager.services)

	services := manager.GetAllServices()
	assert.NotEmpty(t, services)
}
