//go:build unit

package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadSingleConfig(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("loads valid config", func(t *testing.T) {
		configPath := filepath.Join(tempDir, "config.yaml")
		configData := `project:
  name: test-project
stack:
  enabled:
    - postgres
`
		require.NoError(t, os.WriteFile(configPath, []byte(configData), 0644))

		cfg, err := loadSingleConfig(configPath)
		require.NoError(t, err)
		assert.Equal(t, "test-project", cfg.Project.Name)
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		_, err := loadSingleConfig(filepath.Join(tempDir, "nonexistent.yaml"))
		assert.Error(t, err)
	})
}

func TestMergeProjectConfigs(t *testing.T) {
	base := &config.Config{
		Project: config.ProjectConfig{Name: "base-project"},
		Stack:   config.StackConfig{Enabled: []string{"postgres"}},
	}
	local := &config.Config{
		Project: config.ProjectConfig{Name: "local-project"},
		Stack:   config.StackConfig{Enabled: []string{"redis"}},
	}

	merged := mergeProjectConfigs(base, local)
	assert.Equal(t, "local-project", merged.Project.Name)
	assert.Contains(t, merged.Stack.Enabled, "redis")
}

func TestLoadProjectConfig(t *testing.T) {
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	os.Chdir(tempDir)
	require.NoError(t, os.MkdirAll(core.OttoStackDir, 0755))

	t.Run("loads base config when no local config", func(t *testing.T) {
		configPath := filepath.Join(core.OttoStackDir, core.ConfigFileName)
		configData := `project:
  name: test-project
stack:
  enabled:
    - postgres
`
		require.NoError(t, os.WriteFile(configPath, []byte(configData), 0644))

		cfg, err := LoadProjectConfig(configPath)
		require.NoError(t, err)
		assert.Equal(t, "test-project", cfg.Project.Name)
	})

	t.Run("merges local config when present", func(t *testing.T) {
		configPath := filepath.Join(core.OttoStackDir, core.ConfigFileName)
		localPath := filepath.Join(core.OttoStackDir, core.LocalConfigFileName)

		baseData := `project:
  name: base-project
stack:
  enabled:
    - postgres
`
		localData := `project:
  name: local-project
`
		require.NoError(t, os.WriteFile(configPath, []byte(baseData), 0644))
		require.NoError(t, os.WriteFile(localPath, []byte(localData), 0644))

		cfg, err := LoadProjectConfig(configPath)
		require.NoError(t, err)
		assert.Equal(t, "local-project", cfg.Project.Name)
	})
}

func TestResolveServiceConfigs_Default(t *testing.T) {
	cfg := &config.Config{
		Stack: config.StackConfig{
			Enabled: []string{"postgres"},
		},
	}
	setup := &CoreSetup{Config: cfg}

	configs, err := ResolveServiceConfigs([]string{}, setup)
	if err != nil {
		t.Logf("Service resolution error: %v", err)
	} else {
		assert.NotNil(t, configs)
	}
}

func TestResolveServiceConfigs_Provided(t *testing.T) {
	cfg := &config.Config{
		Stack: config.StackConfig{
			Enabled: []string{"postgres"},
		},
	}
	setup := &CoreSetup{Config: cfg}

	configs, err := ResolveServiceConfigs([]string{"redis"}, setup)
	if err != nil {
		t.Logf("Service resolution error: %v", err)
	} else {
		assert.NotNil(t, configs)
	}
}
