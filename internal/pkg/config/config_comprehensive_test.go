//go:build unit

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// Service constants to avoid import cycle
const (
	testServicePostgres   = "postgres"
	testServiceRedis      = "redis"
	testServiceMysql      = "mysql"
	testServiceLocalstack = "localstack"
)

func TestConfig_Validation(t *testing.T) {
	t.Run("validates config with service constants", func(t *testing.T) {
		cfg := Config{
			Project: ProjectConfig{
				Name: "test-project",
			},
			Stack: StackConfig{
				Enabled: []string{testServicePostgres, testServiceRedis},
			},
		}

		assert.Equal(t, "test-project", cfg.Project.Name)
		assert.Contains(t, cfg.Stack.Enabled, testServicePostgres)
		assert.Contains(t, cfg.Stack.Enabled, testServiceRedis)
	})

	t.Run("validates empty config structure", func(t *testing.T) {
		cfg := Config{}

		assert.Empty(t, cfg.Project.Name)
		assert.Empty(t, cfg.Stack.Enabled)
	})
}

func TestProjectConfig_Validation(t *testing.T) {
	t.Run("validates project config fields", func(t *testing.T) {
		project := ProjectConfig{
			Name: "my-project",
			Type: "docker",
		}

		assert.Equal(t, "my-project", project.Name)
		assert.Equal(t, "docker", project.Type)
	})

	t.Run("handles empty project config", func(t *testing.T) {
		project := ProjectConfig{}

		assert.Empty(t, project.Name)
		assert.Empty(t, project.Type)
	})

	t.Run("validates default project type constant", func(t *testing.T) {
		project := ProjectConfig{
			Name: "test-project",
			Type: DefaultProjectType,
		}

		assert.Equal(t, "docker", project.Type)
		assert.Equal(t, DefaultProjectType, project.Type)
	})
}

func TestStackConfig_ServiceManagement(t *testing.T) {
	t.Run("manages enabled services using constants", func(t *testing.T) {
		stack := StackConfig{
			Enabled: []string{
				testServicePostgres,
				testServiceRedis,
				testServiceMysql,
			},
		}

		assert.Len(t, stack.Enabled, 3)
		assert.Contains(t, stack.Enabled, testServicePostgres)
		assert.Contains(t, stack.Enabled, testServiceRedis)
		assert.Contains(t, stack.Enabled, testServiceMysql)
	})

	t.Run("handles empty enabled services", func(t *testing.T) {
		stack := StackConfig{}

		assert.Empty(t, stack.Enabled)
	})

	t.Run("validates service constants in config", func(t *testing.T) {
		serviceConstants := []string{
			testServicePostgres,
			testServiceRedis,
			testServiceMysql,
			testServiceLocalstack,
		}

		for _, service := range serviceConstants {
			stack := StackConfig{
				Enabled: []string{service},
			}

			assert.Contains(t, stack.Enabled, service)
			assert.NotEmpty(t, service, "Service constant should not be empty")
		}
	})
}

func TestConfig_YAMLSerialization(t *testing.T) {
	t.Run("serializes config to YAML using constants", func(t *testing.T) {
		cfg := Config{
			Project: ProjectConfig{
				Name: "test-project",
			},
			Stack: StackConfig{
				Enabled: []string{testServicePostgres},
			},
		}

		data, err := yaml.Marshal(cfg)
		assert.NoError(t, err)
		assert.NotEmpty(t, data)

		// Should contain service constant
		yamlStr := string(data)
		assert.Contains(t, yamlStr, testServicePostgres)
		assert.Contains(t, yamlStr, "test-project")
	})

	t.Run("deserializes YAML to config using constants", func(t *testing.T) {
		yamlData := `
project:
  name: test-project
stack:
  enabled:
    - postgres
    - redis
`

		var cfg Config
		err := yaml.Unmarshal([]byte(yamlData), &cfg)
		assert.NoError(t, err)

		assert.Equal(t, "test-project", cfg.Project.Name)
		assert.Contains(t, cfg.Stack.Enabled, testServicePostgres)
		assert.Contains(t, cfg.Stack.Enabled, testServiceRedis)
	})
}

func TestConfig_FileOperations(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("saves and loads config file using constants", func(t *testing.T) {
		cfg := Config{
			Project: ProjectConfig{
				Name: "file-test-project",
			},
			Stack: StackConfig{
				Enabled: []string{testServicePostgres, testServiceRedis},
			},
		}

		// Save config
		configPath := filepath.Join(tempDir, core.ConfigFileName)
		data, err := yaml.Marshal(cfg)
		require.NoError(t, err)

		err = os.WriteFile(configPath, data, 0644)
		require.NoError(t, err)

		// Load config
		loadedData, err := os.ReadFile(configPath)
		require.NoError(t, err)

		var loadedCfg Config
		err = yaml.Unmarshal(loadedData, &loadedCfg)
		require.NoError(t, err)

		// Verify using constants
		assert.Equal(t, cfg.Project.Name, loadedCfg.Project.Name)
		assert.Equal(t, cfg.Stack.Enabled, loadedCfg.Stack.Enabled)
		assert.Contains(t, loadedCfg.Stack.Enabled, testServicePostgres)
		assert.Contains(t, loadedCfg.Stack.Enabled, testServiceRedis)
	})

	t.Run("handles malformed YAML gracefully", func(t *testing.T) {
		malformedYAML := `
project:
  name: test
  invalid: [unclosed
`

		var cfg Config
		err := yaml.Unmarshal([]byte(malformedYAML), &cfg)
		assert.Error(t, err, "Should return error for malformed YAML")
	})
}
