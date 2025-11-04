package stack

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadProjectConfigWithLocalOverride(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "otto-config-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create base config
	baseConfig := `
project:
  name: "test-project"
  environment: "development"
stack:
  enabled:
    - postgres
    - redis
`
	baseConfigPath := filepath.Join(tmpDir, "otto-stack-config.yaml")
	err = os.WriteFile(baseConfigPath, []byte(baseConfig), constants.FilePermReadWrite)
	require.NoError(t, err)

	// Test loading base config only
	cfg, err := LoadProjectConfig(baseConfigPath)
	require.NoError(t, err)
	assert.Equal(t, "test-project", cfg.Project.Name)
	assert.Equal(t, "development", cfg.Project.Environment)
	assert.Equal(t, []string{"postgres", "redis"}, cfg.Stack.Enabled)

	// Create local override config
	localConfig := `
project:
  name: "local-project"
  environment: "local"
stack:
  enabled:
    - postgres
    - mongodb
`
	localConfigPath := filepath.Join(tmpDir, "otto-stack-config.local.yaml")
	err = os.WriteFile(localConfigPath, []byte(localConfig), constants.FilePermReadWrite)
	require.NoError(t, err)

	// Test loading with local override
	cfg, err = LoadProjectConfig(baseConfigPath)
	require.NoError(t, err)
	assert.Equal(t, "local-project", cfg.Project.Name)
	assert.Equal(t, "local", cfg.Project.Environment)
	assert.Equal(t, []string{"postgres", "mongodb"}, cfg.Stack.Enabled)
}

func TestLoadProjectConfigPartialOverride(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "otto-config-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create base config
	baseConfig := `
project:
  name: "test-project"
  environment: "development"
stack:
  enabled:
    - postgres
    - redis
`
	baseConfigPath := filepath.Join(tmpDir, "otto-stack-config.yaml")
	err = os.WriteFile(baseConfigPath, []byte(baseConfig), constants.FilePermReadWrite)
	require.NoError(t, err)

	// Create partial local override (only environment)
	localConfig := `
project:
  environment: "local"
`
	localConfigPath := filepath.Join(tmpDir, "otto-stack-config.local.yaml")
	err = os.WriteFile(localConfigPath, []byte(localConfig), constants.FilePermReadWrite)
	require.NoError(t, err)

	// Test loading with partial override
	cfg, err := LoadProjectConfig(baseConfigPath)
	require.NoError(t, err)
	assert.Equal(t, "test-project", cfg.Project.Name)                 // unchanged
	assert.Equal(t, "local", cfg.Project.Environment)                 // overridden
	assert.Equal(t, []string{"postgres", "redis"}, cfg.Stack.Enabled) // unchanged
}
