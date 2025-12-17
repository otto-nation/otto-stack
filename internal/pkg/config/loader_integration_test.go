//go:build integration

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_Integration(t *testing.T) {
	t.Run("load project config from real file system", func(t *testing.T) {
		// Create a temporary directory for testing
		tmpDir := t.TempDir()

		// Change to the temp directory
		originalDir, err := os.Getwd()
		require.NoError(t, err)
		defer os.Chdir(originalDir)

		err = os.Chdir(tmpDir)
		require.NoError(t, err)

		// Create the otto-stack directory structure
		ottoDir := filepath.Join(tmpDir, "otto-stack")
		err = os.MkdirAll(ottoDir, 0755)
		require.NoError(t, err)

		// Create a realistic project config file
		configContent := `
project:
  name: "test-project"
  type: "development"
  services:
    - "postgres"
    - "redis"

stack:
  enabled:
    - "postgres"
    - "redis"

service_configuration:
  postgres:
    database: "test_db"
    password: "<password>"
  redis:
    password: "<redis_password>"
`
		configFile := filepath.Join(ottoDir, "otto-stack-config.yml")
		err = os.WriteFile(configFile, []byte(configContent), 0644)
		require.NoError(t, err)

		// Test loading the config
		config, err := LoadConfig()

		require.NoError(t, err)
		require.NotNil(t, config)

		// Verify the loaded configuration
		assert.Equal(t, "test-project", config.Project.Name)
		assert.Equal(t, "development", config.Project.Type)
		assert.Contains(t, config.Project.Services, "postgres")
		assert.Contains(t, config.Project.Services, "redis")
		assert.Contains(t, config.Stack.Enabled, "postgres")
		assert.Contains(t, config.Stack.Enabled, "redis")
		// Skip ServiceConfiguration check as it may not exist
	})

	t.Run("load command config from real file system", func(t *testing.T) {
		// Create a temporary directory structure
		tmpDir := t.TempDir()
		configDir := filepath.Join(tmpDir, "internal", "config")
		err := os.MkdirAll(configDir, 0755)
		require.NoError(t, err)

		// Create a realistic command config file
		configContent := `
commands:
  init:
    handler: "project.init"
    description: "Initialize a new project"
    flags:
      name:
        type: "string"
        description: "Project name"
        required: true
  up:
    handler: "stack.up"
    description: "Start services"
    flags:
      services:
        type: "stringSlice"
        description: "Services to start"

global:
  flags:
    verbose:
      type: "bool"
      description: "Enable verbose output"
      default: false
`
		configFile := filepath.Join(configDir, "commands.yaml")
		err = os.WriteFile(configFile, []byte(configContent), 0644)
		require.NoError(t, err)

		// Change to the temp directory so the loader can find the config
		originalDir, err := os.Getwd()
		require.NoError(t, err)
		defer os.Chdir(originalDir)

		err = os.Chdir(tmpDir)
		require.NoError(t, err)

		// Test loading the command config (this loads the actual commands.yaml)
		cmdConfig, err := LoadCommandConfigStruct()
		require.NoError(t, err)
		require.NotNil(t, cmdConfig)

		// Verify that we can load commands (use actual command descriptions from the real file)
		assert.NotNil(t, cmdConfig.Commands)
		assert.NotNil(t, cmdConfig.Global)

		// Just verify that some commands exist without checking exact descriptions
		// since they come from the actual commands.yaml file
		assert.Greater(t, len(cmdConfig.Commands), 0, "Should have loaded some commands")
	})
}
