//go:build unit

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadLocalConfig(t *testing.T) {
	t.Run("loads local config when it exists", func(t *testing.T) {
		tempDir := t.TempDir()
		err := os.Chdir(tempDir)
		require.NoError(t, err)

		// Create .otto-stack directory
		err = os.MkdirAll(".otto-stack", 0755)
		require.NoError(t, err)

		// Create local config
		localConfigContent := `stack:
  enabled:
    - redis
`
		err = os.WriteFile(".otto-stack/config.local.yaml", []byte(localConfigContent), 0644)
		require.NoError(t, err)

		cfg, err := loadLocalConfig()
		require.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Contains(t, cfg.Stack.Enabled, "redis")
	})

	t.Run("returns error when local config doesn't exist", func(t *testing.T) {
		tempDir := t.TempDir()
		err := os.Chdir(tempDir)
		require.NoError(t, err)

		_, err = loadLocalConfig()
		assert.Error(t, err)
	})
}

func TestMergeConfigs(t *testing.T) {
	t.Run("merges base and local configs", func(t *testing.T) {
		base := &Config{
			Project: ProjectConfig{
				Name: "test-project",
			},
			Stack: StackConfig{
				Enabled: []string{"postgres"},
			},
		}

		local := &Config{
			Stack: StackConfig{
				Enabled: []string{"redis"},
			},
		}

		result := mergeConfigs(base, local)
		assert.Equal(t, "test-project", result.Project.Name)
		assert.Contains(t, result.Stack.Enabled, "redis")
	})

	t.Run("returns base when local is nil", func(t *testing.T) {
		base := &Config{
			Project: ProjectConfig{
				Name: "test-project",
			},
		}

		result := mergeConfigs(base, &Config{})
		assert.Equal(t, "test-project", result.Project.Name)
	})

}

func TestLoadConfig(t *testing.T) {
	t.Run("loads config from directory", func(t *testing.T) {
		tempDir := t.TempDir()
		err := os.Chdir(tempDir)
		require.NoError(t, err)

		// Create .otto-stack directory and config
		err = os.MkdirAll(".otto-stack", 0755)
		require.NoError(t, err)

		configContent := `project:
  name: test-project
  type: docker
stack:
  enabled:
    - postgres
`
		err = os.WriteFile(".otto-stack/config.yaml", []byte(configContent), 0644)
		require.NoError(t, err)

		cfg, err := LoadConfig()
		require.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "test-project", cfg.Project.Name)
		assert.Contains(t, cfg.Stack.Enabled, "postgres")
	})

	t.Run("merges with local config", func(t *testing.T) {
		tempDir := t.TempDir()
		err := os.Chdir(tempDir)
		require.NoError(t, err)

		// Create .otto-stack directory
		err = os.MkdirAll(".otto-stack", 0755)
		require.NoError(t, err)

		// Create base config
		configContent := `project:
  name: test-project
stack:
  enabled:
    - postgres
`
		err = os.WriteFile(".otto-stack/config.yaml", []byte(configContent), 0644)
		require.NoError(t, err)

		// Create local config
		localConfigContent := `stack:
  enabled:
    - redis
`
		err = os.WriteFile(".otto-stack/config.local.yaml", []byte(localConfigContent), 0644)
		require.NoError(t, err)

		cfg, err := LoadConfig()
		require.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Contains(t, cfg.Stack.Enabled, "redis")
	})
}

func TestConfigService_LoadConfig(t *testing.T) {
	service := NewConfigService()

	t.Run("loads config successfully", func(t *testing.T) {
		tempDir := t.TempDir()
		err := os.Chdir(tempDir)
		require.NoError(t, err)

		// Create config
		err = os.MkdirAll(".otto-stack", 0755)
		require.NoError(t, err)

		configContent := `project:
  name: test-project
stack:
  enabled:
    - postgres
`
		err = os.WriteFile(".otto-stack/config.yaml", []byte(configContent), 0644)
		require.NoError(t, err)

		cfg, err := service.LoadConfig()
		require.NoError(t, err)
		assert.Equal(t, "test-project", cfg.Project.Name)
	})
}

func TestConfigService_ValidateConfig(t *testing.T) {
	service := NewConfigService()

	t.Run("validates valid config", func(t *testing.T) {
		cfg := &Config{
			Project: ProjectConfig{Name: "test"},
			Stack:   StackConfig{Enabled: []string{"postgres"}},
		}
		err := service.ValidateConfig(cfg)
		assert.NoError(t, err)
	})

	t.Run("rejects nil config", func(t *testing.T) {
		err := service.ValidateConfig(nil)
		assert.Error(t, err)
	})

	t.Run("rejects empty project name", func(t *testing.T) {
		cfg := &Config{
			Stack: StackConfig{Enabled: []string{"postgres"}},
		}
		err := service.ValidateConfig(cfg)
		assert.Error(t, err)
	})

	t.Run("rejects empty services", func(t *testing.T) {
		cfg := &Config{
			Project: ProjectConfig{Name: "test"},
		}
		err := service.ValidateConfig(cfg)
		assert.Error(t, err)
	})
}

func TestConfigService_GetConfigHash(t *testing.T) {
	service := NewConfigService()

	t.Run("generates hash for config", func(t *testing.T) {
		cfg := &Config{
			Project: ProjectConfig{Name: "test"},
			Stack:   StackConfig{Enabled: []string{"postgres"}},
		}
		hash, err := service.GetConfigHash(cfg)
		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.Len(t, hash, 64) // SHA256 hex string length
	})

	t.Run("returns error for nil config", func(t *testing.T) {
		_, err := service.GetConfigHash(nil)
		assert.Error(t, err)
	})

	t.Run("same config produces same hash", func(t *testing.T) {
		cfg := &Config{
			Project: ProjectConfig{Name: "test"},
			Stack:   StackConfig{Enabled: []string{"postgres"}},
		}
		hash1, err := service.GetConfigHash(cfg)
		require.NoError(t, err)

		hash2, err := service.GetConfigHash(cfg)
		require.NoError(t, err)

		assert.Equal(t, hash1, hash2)
	})
}
