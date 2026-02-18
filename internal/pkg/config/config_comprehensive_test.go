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

func TestConfig_WithServiceConstants(t *testing.T) {
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
}

func TestConfig_EmptyStructure(t *testing.T) {
	cfg := Config{}

	assert.Empty(t, cfg.Project.Name)
	assert.Empty(t, cfg.Stack.Enabled)
}

func TestProjectConfig_Fields(t *testing.T) {
	project := ProjectConfig{
		Name: "my-project",
		Type: "docker",
	}

	assert.Equal(t, "my-project", project.Name)
	assert.Equal(t, "docker", project.Type)
}

func TestProjectConfig_Empty(t *testing.T) {
	project := ProjectConfig{}

	assert.Empty(t, project.Name)
	assert.Empty(t, project.Type)
}

func TestProjectConfig_DefaultType(t *testing.T) {
	project := ProjectConfig{
		Name: "test-project",
		Type: DefaultProjectType,
	}

	assert.Equal(t, "docker", project.Type)
	assert.Equal(t, DefaultProjectType, project.Type)
}

func TestStackConfig_EnabledServices(t *testing.T) {
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
}

func TestStackConfig_EmptyServices(t *testing.T) {
	stack := StackConfig{}
	assert.Empty(t, stack.Enabled)
}

func TestStackConfig_ServiceConstants(t *testing.T) {
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
}

func TestConfig_SerializeToYAML(t *testing.T) {
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

	yamlStr := string(data)
	assert.Contains(t, yamlStr, testServicePostgres)
	assert.Contains(t, yamlStr, "test-project")
}

func TestConfig_DeserializeFromYAML(t *testing.T) {
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
}

func TestConfig_SaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	cfg := Config{
		Project: ProjectConfig{
			Name: "file-test-project",
		},
		Stack: StackConfig{
			Enabled: []string{testServicePostgres, testServiceRedis},
		},
	}

	configPath := filepath.Join(tempDir, core.ConfigFileName)
	data, err := yaml.Marshal(cfg)
	require.NoError(t, err)

	err = os.WriteFile(configPath, data, 0644)
	require.NoError(t, err)

	loadedData, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var loadedCfg Config
	err = yaml.Unmarshal(loadedData, &loadedCfg)
	require.NoError(t, err)

	assert.Equal(t, cfg.Project.Name, loadedCfg.Project.Name)
	assert.Equal(t, cfg.Stack.Enabled, loadedCfg.Stack.Enabled)
	assert.Contains(t, loadedCfg.Stack.Enabled, testServicePostgres)
	assert.Contains(t, loadedCfg.Stack.Enabled, testServiceRedis)
}

func TestConfig_MalformedYAML(t *testing.T) {
	malformedYAML := `
project:
  name: test
  invalid: [unclosed
`

	var cfg Config
	err := yaml.Unmarshal([]byte(malformedYAML), &cfg)
	assert.Error(t, err, "Should return error for malformed YAML")
}
