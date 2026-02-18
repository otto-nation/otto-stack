//go:build unit

package config

import (
	"os"
	"testing"

	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/test/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestLoadLocalConfig(t *testing.T) {
	t.Run("loads local config when it exists", func(t *testing.T) {
		tempDir := t.TempDir()

		// Save and restore working directory
		origDir, err := os.Getwd()
		require.NoError(t, err)
		defer os.Chdir(origDir)

		err = os.Chdir(tempDir)
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

		// Save and restore working directory
		origDir, err := os.Getwd()
		require.NoError(t, err)
		defer os.Chdir(origDir)

		err = os.Chdir(tempDir)
		require.NoError(t, err)

		_, err = loadLocalConfig()
		assert.Error(t, err)
	})
}

func TestMergeConfigs(t *testing.T) {
	t.Run("merges base and local configs", func(t *testing.T) {
		var base, local Config
		require.NoError(t, yaml.Unmarshal(fixtures.LoadConfigYAML(t, "minimal"), &base))
		require.NoError(t, yaml.Unmarshal(fixtures.LoadConfigYAML(t, "with-stack"), &local))

		result := mergeConfigs(&base, &local)
		assert.Equal(t, "test-project", result.Project.Name)
		assert.Contains(t, result.Stack.Enabled, "redis")
	})

	t.Run("returns base when local is nil", func(t *testing.T) {
		var base Config
		require.NoError(t, yaml.Unmarshal(fixtures.LoadConfigYAML(t, "minimal"), &base))

		result := mergeConfigs(&base, &Config{})
		assert.Equal(t, "test-project", result.Project.Name)
	})

}

func TestLoadConfig(t *testing.T) {
	t.Run("loads config from directory", func(t *testing.T) {
		tempDir := t.TempDir()

		// Save and restore working directory
		origDir, err := os.Getwd()
		require.NoError(t, err)
		defer os.Chdir(origDir)

		err = os.Chdir(tempDir)
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

		// Save and restore working directory
		origDir, err := os.Getwd()
		require.NoError(t, err)
		defer os.Chdir(origDir)

		err = os.Chdir(tempDir)
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

		// Save and restore working directory
		origDir, err := os.Getwd()
		require.NoError(t, err)
		defer os.Chdir(origDir)

		err = os.Chdir(tempDir)
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

	t.Run("returns error when config not found", func(t *testing.T) {
		tempDir := t.TempDir()

		// Save and restore working directory
		origDir, err := os.Getwd()
		require.NoError(t, err)
		defer os.Chdir(origDir)

		err = os.Chdir(tempDir)
		require.NoError(t, err)

		_, err = service.LoadConfig()
		assert.Error(t, err)
	})
}

func TestConfigService_ValidateConfig(t *testing.T) {
	service := NewConfigService()

	t.Run("validates valid config", func(t *testing.T) {
		var cfg Config
		require.NoError(t, yaml.Unmarshal(fixtures.LoadConfigYAML(t, "with-stack"), &cfg))
		err := service.ValidateConfig(&cfg)
		assert.NoError(t, err)
	})

	t.Run("rejects nil config", func(t *testing.T) {
		err := service.ValidateConfig(nil)
		assert.Error(t, err)
	})

	t.Run("rejects empty project name", func(t *testing.T) {
		cfg := &Config{Stack: StackConfig{Enabled: []string{"postgres"}}}
		err := service.ValidateConfig(cfg)
		assert.Error(t, err)
	})

	t.Run("rejects empty services", func(t *testing.T) {
		var cfg Config
		require.NoError(t, yaml.Unmarshal(fixtures.LoadConfigYAML(t, "minimal"), &cfg))
		err := service.ValidateConfig(&cfg)
		assert.Error(t, err)
	})
}

func TestConfigService_GetConfigHash(t *testing.T) {
	service := NewConfigService()

	t.Run("generates hash for config", func(t *testing.T) {
		var cfg Config
		require.NoError(t, yaml.Unmarshal(fixtures.LoadConfigYAML(t, "with-stack"), &cfg))
		hash, err := service.GetConfigHash(&cfg)
		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.Len(t, hash, 64)
	})

	t.Run("returns error for nil config", func(t *testing.T) {
		_, err := service.GetConfigHash(nil)
		assert.Error(t, err)
	})

	t.Run("same config produces same hash", func(t *testing.T) {
		var cfg Config
		require.NoError(t, yaml.Unmarshal(fixtures.LoadConfigYAML(t, "with-stack"), &cfg))
		hash1, err := service.GetConfigHash(&cfg)
		require.NoError(t, err)

		hash2, err := service.GetConfigHash(&cfg)
		require.NoError(t, err)

		assert.Equal(t, hash1, hash2)
	})
}

func TestGenerateConfig_WithSharing(t *testing.T) {
	t.Run("generates config with sharing enabled", func(t *testing.T) {
		ctx := clicontext.Context{
			Project: clicontext.ProjectSpec{Name: "test-project"},
			Services: clicontext.ServiceSpec{
				Names: []string{"postgres", "redis"},
			},
			Sharing: &clicontext.SharingSpec{
				Enabled:  true,
				Services: map[string]bool{"postgres": true},
			},
		}

		data, err := GenerateConfig(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, data)

		// Verify YAML contains sharing config
		assert.Contains(t, string(data), "sharing:")
		assert.Contains(t, string(data), "enabled: true")
	})

	t.Run("generates config without sharing", func(t *testing.T) {
		ctx := clicontext.Context{
			Project: clicontext.ProjectSpec{Name: "test-project"},
			Services: clicontext.ServiceSpec{
				Names: []string{"postgres"},
			},
		}

		data, err := GenerateConfig(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, data)
	})
}

func TestLoadBaseConfig_ErrorHandling(t *testing.T) {
	t.Run("returns error for invalid YAML", func(t *testing.T) {
		tempDir := t.TempDir()

		// Save and restore working directory
		origDir, err := os.Getwd()
		require.NoError(t, err)
		defer os.Chdir(origDir)

		err = os.Chdir(tempDir)
		require.NoError(t, err)

		err = os.MkdirAll(".otto-stack", 0755)
		require.NoError(t, err)

		// Write invalid YAML
		err = os.WriteFile(".otto-stack/config.yaml", []byte("invalid: yaml: ["), 0644)
		require.NoError(t, err)

		_, err = LoadConfig()
		assert.Error(t, err)
	})
}

func TestMergeConfigs_EdgeCases(t *testing.T) {
	t.Run("preserves base when local has empty values", func(t *testing.T) {
		base := &Config{
			Project: ProjectConfig{Name: "base-project"},
			Stack:   StackConfig{Enabled: []string{"postgres"}},
		}
		local := &Config{}

		result := mergeConfigs(base, local)
		assert.Equal(t, "base-project", result.Project.Name)
		assert.Equal(t, []string{"postgres"}, result.Stack.Enabled)
	})

	t.Run("overrides project name when local has value", func(t *testing.T) {
		base := &Config{
			Project: ProjectConfig{Name: "base-project"},
			Stack:   StackConfig{Enabled: []string{"postgres"}},
		}
		local := &Config{
			Project: ProjectConfig{Name: "local-project"},
		}

		result := mergeConfigs(base, local)
		assert.Equal(t, "local-project", result.Project.Name)
	})

	t.Run("overrides services when local has services", func(t *testing.T) {
		base := &Config{
			Project: ProjectConfig{Name: "test"},
			Stack:   StackConfig{Enabled: []string{"postgres"}},
		}
		local := &Config{
			Stack: StackConfig{Enabled: []string{"redis", "mysql"}},
		}

		result := mergeConfigs(base, local)
		assert.Equal(t, []string{"redis", "mysql"}, result.Stack.Enabled)
	})
}

func TestGenerateConfig_ErrorCases(t *testing.T) {
	t.Run("returns error for empty project name", func(t *testing.T) {
		ctx := clicontext.Context{
			Services: clicontext.ServiceSpec{
				Names: []string{"postgres"},
			},
		}

		_, err := GenerateConfig(ctx)
		assert.Error(t, err)
	})
}

func TestValidateSharingPolicy_EdgeCases(t *testing.T) {
	t.Run("allows nil sharing config", func(t *testing.T) {
		cfg := &Config{
			Project: ProjectConfig{Name: "test"},
			Stack:   StackConfig{Enabled: []string{"postgres"}},
			Sharing: nil,
		}
		err := validateSharingPolicy(cfg)
		assert.NoError(t, err)
	})

	t.Run("allows disabled sharing", func(t *testing.T) {
		cfg := &Config{
			Project: ProjectConfig{Name: "test"},
			Stack:   StackConfig{Enabled: []string{"postgres"}},
			Sharing: &SharingConfig{Enabled: false},
		}
		err := validateSharingPolicy(cfg)
		assert.NoError(t, err)
	})

	t.Run("allows empty services list", func(t *testing.T) {
		cfg := &Config{
			Project: ProjectConfig{Name: "test"},
			Stack:   StackConfig{Enabled: []string{"postgres"}},
			Sharing: &SharingConfig{
				Enabled:  true,
				Services: map[string]bool{},
			},
		}
		err := validateSharingPolicy(cfg)
		assert.NoError(t, err)
	})
}
