package stack

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"

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
stack:
  enabled:
    - postgres
    - redis
`
	baseConfigPath := filepath.Join(tmpDir, "otto-stack-config.yaml")
	err = os.WriteFile(baseConfigPath, []byte(baseConfig), core.PermReadWrite)
	require.NoError(t, err)

	// Test loading base config only
	cfg, err := LoadProjectConfig(baseConfigPath)
	require.NoError(t, err)
	assert.Equal(t, "test-project", cfg.Project.Name)
	assert.Equal(t, []string{"postgres", "redis"}, cfg.Stack.Enabled)

	// Create local override config in otto-stack directory
	ottoStackDir := filepath.Join(tmpDir, core.OttoStackDir)
	err = os.MkdirAll(ottoStackDir, core.PermReadWriteExec)
	require.NoError(t, err)

	localConfig := `
project:
  name: "local-project"
stack:
  enabled:
    - postgres
    - mongodb
`
	localConfigPath := filepath.Join(ottoStackDir, core.LocalConfigFileName)
	err = os.WriteFile(localConfigPath, []byte(localConfig), core.PermReadWrite)
	require.NoError(t, err)

	// Change to temp directory so LoadProjectConfig can find the local config
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Test loading with local override
	cfg, err = LoadProjectConfig(baseConfigPath)
	require.NoError(t, err)
	assert.Equal(t, "local-project", cfg.Project.Name)
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
stack:
  enabled:
    - postgres
    - redis
`
	baseConfigPath := filepath.Join(tmpDir, "otto-stack-config.yaml")
	err = os.WriteFile(baseConfigPath, []byte(baseConfig), core.PermReadWrite)
	require.NoError(t, err)

	// Create partial local override (only environment) in otto-stack directory
	ottoStackDir := filepath.Join(tmpDir, core.OttoStackDir)
	err = os.MkdirAll(ottoStackDir, core.PermReadWriteExec)
	require.NoError(t, err)

	localConfig := `
project:
`
	localConfigPath := filepath.Join(ottoStackDir, core.LocalConfigFileName)
	err = os.WriteFile(localConfigPath, []byte(localConfig), core.PermReadWrite)
	require.NoError(t, err)

	// Change to temp directory so LoadProjectConfig can find the local config
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Test loading with partial override
	cfg, err := LoadProjectConfig(baseConfigPath)
	require.NoError(t, err)
	assert.Equal(t, "test-project", cfg.Project.Name)                 // unchanged
	assert.Equal(t, []string{"postgres", "redis"}, cfg.Stack.Enabled) // unchanged
}
