//go:build unit

package services

import (
	"testing"

	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/test/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_LoadServices(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	services := manager.GetAllServices()
	assert.NotEmpty(t, services)

	_, hasPostgres := services[ServicePostgres]
	_, hasRedis := services[ServiceRedis]
	assert.True(t, hasPostgres || hasRedis)
}

func TestManager_GetService_Postgres(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	service, err := manager.GetService(ServicePostgres)
	assert.NoError(t, err)
	assert.Equal(t, ServicePostgres, service.Name)
	assert.NotEmpty(t, service.Description)
}

func TestManager_GetService_Redis(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	service, err := manager.GetService(ServiceRedis)
	if err == nil {
		assert.Equal(t, ServiceRedis, service.Name)
	}
}

func TestManager_GetService_Mysql(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	service, err := manager.GetService(ServiceMysql)
	if err == nil {
		assert.Equal(t, ServiceMysql, service.Name)
	}
}

func TestManager_ExecuteCustomOperation_NonexistentService(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)
	_, err = manager.ExecuteCustomOperation("nonexistent", "backup")
	assert.Error(t, err)
}

func TestManager_ExecuteCustomOperation_NoCustomOps(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)
	_, err = manager.ExecuteCustomOperation(ServicePostgres, "nonexistent-op")
	assert.Error(t, err)
}

func TestManager_ExecuteCustomOperation_Success(t *testing.T) {
	config := fixtures.LoadService(t, "custom-ops")
	manager := &Manager{
		services: map[string]servicetypes.ServiceConfig{"testservice": config},
	}

	cmd, err := manager.ExecuteCustomOperation("testservice", "backup")
	require.NoError(t, err)
	assert.Equal(t, []string{"pg_dump", "-U", "postgres"}, cmd)
}

func TestManager_ExecuteCustomOperation_NoDefaultArgs(t *testing.T) {
	config := fixtures.LoadService(t, "custom-ops")
	manager := &Manager{
		services: map[string]servicetypes.ServiceConfig{"testservice": config},
	}

	cmd, err := manager.ExecuteCustomOperation("testservice", "restore")
	require.NoError(t, err)
	assert.Equal(t, []string{"pg_restore"}, cmd)
}

func TestManager_ExecuteCustomOperation_NonexistentOp(t *testing.T) {
	config := fixtures.LoadService(t, "custom-ops")
	manager := &Manager{
		services: map[string]servicetypes.ServiceConfig{"testservice": config},
	}

	_, err := manager.ExecuteCustomOperation("testservice", "nonexistent")
	assert.Error(t, err)
}

func TestResolveUpServices_NilConfig(t *testing.T) {
	configs, err := ResolveUpServices([]string{"postgres"}, nil)
	if err == nil {
		assert.NotEmpty(t, configs)
	}
}

func TestResolveUpServices_MultipleServices(t *testing.T) {
	configs, err := ResolveUpServices([]string{"postgres", "redis"}, nil)
	if err == nil {
		assert.GreaterOrEqual(t, len(configs), 2)
	}
}

func TestResolveUpServices_InvalidService(t *testing.T) {
	_, err := ResolveUpServices([]string{"invalid-xyz"}, nil)
	assert.Error(t, err)
}
