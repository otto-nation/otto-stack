//go:build integration

package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/test/e2e/framework"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_InitCommand(t *testing.T) {
	lifecycle := framework.NewTestLifecycle(t, "init-test", []string{services.ServicePostgres})
	defer lifecycle.Cleanup()

	t.Run("init creates project structure", func(t *testing.T) {
		err := lifecycle.InitializeStack()
		require.NoError(t, err)

		configFile := filepath.Join(lifecycle.Environment.WorkDir, core.OttoStackDir, core.ConfigFileName)
		assert.FileExists(t, configFile)

		content, err := os.ReadFile(configFile)
		require.NoError(t, err)
		assert.Contains(t, string(content), services.ServicePostgres)
	})

	t.Run("validate works after init", func(t *testing.T) {
		result := lifecycle.CLI.RunExpectSuccess(core.CommandValidate)
		assert.Contains(t, result.Stdout, "Configuration is valid")
	})
}

func TestE2E_MultiServiceConfig(t *testing.T) {
	serviceList := []string{services.ServicePostgres, services.ServiceRedis}
	lifecycle := framework.NewTestLifecycle(t, "multi-service", serviceList)
	defer lifecycle.Cleanup()

	t.Run("init with multiple services", func(t *testing.T) {
		err := lifecycle.InitializeStack()
		require.NoError(t, err)

		configFile := filepath.Join(lifecycle.Environment.WorkDir, core.OttoStackDir, core.ConfigFileName)
		content, err := os.ReadFile(configFile)
		require.NoError(t, err)

		configStr := string(content)
		assert.Contains(t, configStr, services.ServicePostgres)
		assert.Contains(t, configStr, services.ServiceRedis)
	})

	t.Run("services get unique ports", func(t *testing.T) {
		postgresPort := lifecycle.Environment.GetServicePort(services.ServicePostgres)
		redisPort := lifecycle.Environment.GetServicePort(services.ServiceRedis)

		require.Greater(t, postgresPort, 0)
		require.Greater(t, redisPort, 0)
		require.NotEqual(t, postgresPort, redisPort)
	})
}

func TestE2E_ServiceConfig(t *testing.T) {
	lifecycle := framework.NewTestLifecycle(t, "service-config", []string{services.ServiceLocalstackSqs})
	defer lifecycle.Cleanup()

	t.Run("localstack service config", func(t *testing.T) {
		err := lifecycle.InitializeStack()
		require.NoError(t, err)

		port := lifecycle.Environment.GetServicePort(services.ServiceLocalstack)
		require.Greater(t, port, 0)
		require.NotEqual(t, port, services.PortLocalstack)
	})
}
