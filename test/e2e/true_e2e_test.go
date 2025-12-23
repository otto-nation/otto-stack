package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/test/e2e/framework"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_PostgresLifecycle(t *testing.T) {
	projectName := fmt.Sprintf("pg-e2e-%d", time.Now().UnixNano())
	lifecycle := framework.NewTestLifecycle(t, projectName, []string{services.ServicePostgres})
	defer lifecycle.Cleanup()

	// Initialize stack
	err := lifecycle.InitializeStack()
	require.NoError(t, err)

	// Start services
	err = lifecycle.StartServices()
	require.NoError(t, err)

	// Wait and verify services are running
	err = lifecycle.WaitForServices()
	require.NoError(t, err)

	// Check status shows running
	result := lifecycle.CLI.RunExpectSuccess(core.CommandStatus)
	assert.Contains(t, result.Stdout, services.ServicePostgres)
	assert.Contains(t, result.Stdout, "running")

	// Stop services
	err = lifecycle.StopServices()
	require.NoError(t, err)
}

func TestE2E_MultiServiceLifecycle(t *testing.T) {
	projectName := fmt.Sprintf("multi-e2e-%d", time.Now().UnixNano())
	serviceList := []string{services.ServicePostgres, services.ServiceRedis}
	lifecycle := framework.NewTestLifecycle(t, projectName, serviceList)
	defer lifecycle.Cleanup()

	// Initialize stack
	err := lifecycle.InitializeStack()
	require.NoError(t, err)

	// Start all services
	err = lifecycle.StartServices()
	require.NoError(t, err)

	// Verify all services are running
	result := lifecycle.CLI.RunExpectSuccess(core.CommandStatus)
	assert.Contains(t, result.Stdout, services.ServicePostgres)
	assert.Contains(t, result.Stdout, services.ServiceRedis)
	assert.Contains(t, result.Stdout, "running")

	// Test port allocation uniqueness
	postgresPort := lifecycle.Environment.GetServicePort(services.ServicePostgres)
	redisPort := lifecycle.Environment.GetServicePort(services.ServiceRedis)
	require.NotEqual(t, postgresPort, redisPort)
}

func TestE2E_ServiceRestart(t *testing.T) {
	projectName := fmt.Sprintf("restart-e2e-%d", time.Now().UnixNano())
	lifecycle := framework.NewTestLifecycle(t, projectName, []string{services.ServicePostgres})
	defer lifecycle.Cleanup()

	// Initialize and start
	err := lifecycle.InitializeStack()
	require.NoError(t, err)
	err = lifecycle.StartServices()
	require.NoError(t, err)

	// Stop services
	err = lifecycle.StopServices()
	require.NoError(t, err)

	// Restart services
	err = lifecycle.StartServices()
	require.NoError(t, err)

	// Verify still working
	result := lifecycle.CLI.RunExpectSuccess(core.CommandStatus)
	assert.Contains(t, result.Stdout, services.ServicePostgres)
	assert.Contains(t, result.Stdout, "running")
}

func TestE2E_LocalstackIntegration(t *testing.T) {
	projectName := fmt.Sprintf("localstack-e2e-%d", time.Now().UnixNano())
	lifecycle := framework.NewTestLifecycle(t, projectName, []string{services.ServiceLocalstack})
	defer lifecycle.Cleanup()

	err := lifecycle.InitializeStack()
	require.NoError(t, err)

	err = lifecycle.StartServices()
	require.NoError(t, err)

	result := lifecycle.CLI.RunExpectSuccess(core.CommandStatus)
	assert.Contains(t, result.Stdout, services.ServiceLocalstack)
	assert.Contains(t, result.Stdout, "running")
}

func TestE2E_AllServicesIntegration(t *testing.T) {
	projectName := fmt.Sprintf("all-services-e2e-%d", time.Now().UnixNano())
	serviceList := []string{services.ServicePostgres, services.ServiceRedis, services.ServiceLocalstack}
	lifecycle := framework.NewTestLifecycle(t, projectName, serviceList)
	defer lifecycle.Cleanup()

	err := lifecycle.InitializeStack()
	require.NoError(t, err)

	err = lifecycle.StartServices()
	require.NoError(t, err)

	// Verify all services are running
	result := lifecycle.CLI.RunExpectSuccess(core.CommandStatus)
	for _, service := range serviceList {
		assert.Contains(t, result.Stdout, service)
	}
	assert.Contains(t, result.Stdout, "running")

	// Verify unique port allocation
	ports := make(map[int]bool)
	for _, service := range serviceList {
		port := lifecycle.Environment.GetServicePort(service)
		require.Greater(t, port, 0)
		require.False(t, ports[port], "Port %d is already allocated", port)
		ports[port] = true
	}
}
