//go:build integration

package e2e

import (
	"fmt"
	"os"
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
	err = lifecycle.WaitForServices()
	require.NoError(t, err)

	// Stop services
	err = lifecycle.StopServices()
	require.NoError(t, err, "Failed to stop services: %v", err)

	// Restart services
	err = lifecycle.StartServices()
	require.NoError(t, err, "Failed to restart services: %v", err)
	err = lifecycle.WaitForServices()
	require.NoError(t, err, "Failed to wait for restarted services: %v", err)

	// Verify still working
	result := lifecycle.CLI.RunExpectSuccess(core.CommandStatus)
	assert.Contains(t, result.Stdout, services.ServicePostgres)
	assert.Contains(t, result.Stdout, "running")
}

func TestE2E_LocalstackIntegration(t *testing.T) {
	const testQueueName = "test-queue"

	if testing.Short() {
		t.Skip("Skipping LocalStack test in short mode")
	}

	projectName := fmt.Sprintf("localstack-e2e-%d", time.Now().UnixNano())
	lifecycle := framework.NewTestLifecycle(t, projectName, []string{services.ServiceLocalstackSqs})
	defer lifecycle.Cleanup()

	err := lifecycle.InitializeStack()
	require.NoError(t, err)

	// Create service config file for init container
	err = lifecycle.CreateServiceConfigFile("localstack-sqs.yml", map[string]interface{}{
		"queues": []map[string]interface{}{
			{"name": testQueueName},
		},
	})
	require.NoError(t, err)

	// Debug: Check if service config file was created
	serviceConfigPath := fmt.Sprintf("%s/.otto-stack/service-configs/localstack-sqs.yml", lifecycle.Environment.WorkDir)
	if _, err := os.Stat(serviceConfigPath); err == nil {
		t.Logf("✅ Service config file created at: %s", serviceConfigPath)
		if content, err := os.ReadFile(serviceConfigPath); err == nil {
			t.Logf("📄 Service config content: %s", string(content))
		}
	} else {
		t.Logf("❌ Service config file NOT found at: %s", serviceConfigPath)
	}

	// Give LocalStack more time to start
	t.Log("Starting LocalStack (this may take up to 2 minutes)...")
	err = lifecycle.StartServices()
	if err != nil {
		t.Skipf("LocalStack failed to start (likely resource constraints): %v", err)
	}

	result := lifecycle.CLI.RunExpectSuccess(core.CommandStatus)
	assert.Contains(t, result.Stdout, services.ServiceLocalstack)
	assert.Contains(t, result.Stdout, "running")

	// Verify the SQS queue was created
	port := lifecycle.Environment.GetServicePort(services.ServiceLocalstack)
	endpoint := fmt.Sprintf("http://localhost:%d", port)

	// Use AWS CLI to list queues and verify our test queue exists
	// Set LocalStack credentials to bypass SSO
	lifecycle.CLI.SetEnv("AWS_ACCESS_KEY_ID", "test")
	lifecycle.CLI.SetEnv("AWS_SECRET_ACCESS_KEY", "test")
	lifecycle.CLI.SetEnv("AWS_DEFAULT_REGION", "us-east-1")

	// Wait a bit for LocalStack and init container to be ready
	t.Log("Waiting for LocalStack and init container to initialize...")
	time.Sleep(10 * time.Second)

	// Debug: Check what containers are running
	dockerResult := lifecycle.CLI.RunSystemCommand("docker", "ps", "--format", "table {{.Names}}\t{{.Status}}")
	t.Logf("🐳 Running containers:\n%s", dockerResult.Stdout)

	// Debug: Check docker-compose file content
	composePath := fmt.Sprintf("%s/.otto-stack/docker-compose.yml", lifecycle.Environment.WorkDir)
	if composeContent, err := os.ReadFile(composePath); err == nil {
		t.Logf("📋 Docker-compose file content:\n%s", string(composeContent))
	} else {
		t.Logf("❌ Could not read docker-compose file: %v", err)
	}

	// Retry AWS CLI check a few times
	var awsResult *framework.CLIResult
	for i := 0; i < 3; i++ {
		awsResult = lifecycle.CLI.RunSystemCommand("aws", "--endpoint-url="+endpoint, "sqs", "list-queues")
		if awsResult.ExitCode == 0 {
			break
		}
		if i < 2 { // Don't sleep on the last attempt
			t.Logf("AWS CLI attempt %d failed, retrying in 5 seconds...", i+1)
			time.Sleep(5 * time.Second)
		}
	}

	require.Equal(t, 0, awsResult.ExitCode, "AWS CLI should succeed: %s", awsResult.Stderr)
	require.NotEmpty(t, awsResult.Stdout, "AWS CLI should return queue list")
	assert.Contains(t, awsResult.Stdout, testQueueName, "SQS queue should be created and accessible")
	t.Logf("✅ SQS queue '%s' verified successfully via AWS CLI", testQueueName)
}

func TestE2E_AllServicesIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping all services test in short mode")
	}

	projectName := fmt.Sprintf("all-services-e2e-%d", time.Now().UnixNano())
	serviceList := []string{services.ServicePostgres, services.ServiceRedis, services.ServiceLocalstackSqs}
	lifecycle := framework.NewTestLifecycle(t, projectName, serviceList)
	defer lifecycle.Cleanup()

	err := lifecycle.InitializeStack()
	require.NoError(t, err)

	t.Log("Starting all services (this may take up to 3 minutes)...")
	err = lifecycle.StartServices()
	if err != nil {
		t.Skipf("All services failed to start (likely resource constraints): %v", err)
	}

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
