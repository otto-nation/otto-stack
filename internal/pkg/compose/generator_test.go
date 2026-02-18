package compose

import (
	"os"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/test/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGenerator(t *testing.T) {
	gen, err := NewGenerator("test-project")
	require.NoError(t, err)
	assert.NotNil(t, gen)
	assert.Equal(t, "test-project", gen.projectName)

	gen, err = NewGenerator("")
	require.NoError(t, err)
	assert.NotNil(t, gen)
}

func TestGenerator_BuildComposeData(t *testing.T) {
	gen, err := NewGenerator("test-project")
	require.NoError(t, err)

	data, err := gen.BuildComposeData([]types.ServiceConfig{})
	require.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.Contains(t, string(data), "test-project")

	services := []types.ServiceConfig{
		fixtures.NewServiceConfig("redis").
			WithImage("redis:latest").
			WithPort("6379", "6379").
			Build(),
	}

	data, err = gen.BuildComposeData(services)
	require.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.Contains(t, string(data), "redis")
	assert.Contains(t, string(data), "redis:latest")

	gen, err = NewGenerator("")
	require.NoError(t, err)
	_, err = gen.BuildComposeData([]types.ServiceConfig{})
	assert.Error(t, err)
}

func TestGenerator_BuildComposeDataWithHeader(t *testing.T) {
	gen, err := NewGenerator("test-project")
	require.NoError(t, err)

	header := "# Test Header\n"
	data, err := gen.BuildComposeDataWithHeader([]types.ServiceConfig{}, header)
	require.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.Contains(t, string(data), "# Test Header")

	data, err = gen.BuildComposeDataWithHeader([]types.ServiceConfig{}, "")
	require.NoError(t, err)
	assert.NotEmpty(t, data)
}

func TestGenerator_BuildServicesFromConfigs(t *testing.T) {
	gen, err := NewGenerator("test-project")
	require.NoError(t, err)

	services := []types.ServiceConfig{
		fixtures.NewServiceConfig("postgres").
			WithImage("postgres:latest").
			WithEnv("POSTGRES_PASSWORD", "secret").
			WithEnv("POSTGRES_USER", "admin").
			Build(),
	}

	result, err := gen.buildServicesFromConfigs(services)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, result, "postgres")

	services = []types.ServiceConfig{
		fixtures.NewServiceConfig("mysql").
			WithImage("mysql:latest").
			WithVolume("./data:/var/lib/mysql").
			Build(),
	}

	result, err = gen.buildServicesFromConfigs(services)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, result, "mysql")

	services = []types.ServiceConfig{
		fixtures.NewServiceConfig("redis").
			WithImage("redis:latest").
			WithRestart("always").
			WithCommand([]string{"redis-server", "--appendonly", "yes"}).
			WithMemoryLimit("512m").
			Build(),
	}

	result, err = gen.buildServicesFromConfigs(services)
	require.NoError(t, err)
	assert.NotNil(t, result)
	redisService := result["redis"].(map[string]any)
	assert.Equal(t, "always", redisService["restart"])
	assert.Equal(t, []string{"redis-server", "--appendonly", "yes"}, redisService["command"])
	assert.Equal(t, "512m", redisService["mem_limit"])

	services = []types.ServiceConfig{
		fixtures.NewServiceConfig("custom").
			WithImage("custom:latest").
			WithEntrypoint([]string{"/bin/sh", "-c"}).
			Build(),
	}

	result, err = gen.buildServicesFromConfigs(services)
	require.NoError(t, err)
	customService := result["custom"].(map[string]any)
	assert.Equal(t, []string{"/bin/sh", "-c"}, customService["entrypoint"])

	services = []types.ServiceConfig{
		{
			Name:        "localstack-sns",
			ServiceType: types.ServiceTypeConfiguration,
		},
	}

	result, err = gen.buildServicesFromConfigs(services)
	require.NoError(t, err)
	assert.NotContains(t, result, "localstack-sns")
}

func TestGenerator_HealthCheckTiming(t *testing.T) {
	gen, err := NewGenerator("test-project")
	require.NoError(t, err)

	services := []types.ServiceConfig{
		fixtures.NewServiceConfig("postgres").
			WithImage("postgres:latest").
			WithHealthCheck([]string{"CMD", "pg_isready"}, 10, 5, 3).
			Build(),
	}

	result, err := gen.buildServicesFromConfigs(services)
	require.NoError(t, err)
	assert.Contains(t, result, "postgres")
}

func TestGenerator_ResolveEnvVar(t *testing.T) {
	gen, err := NewGenerator("test-project")
	require.NoError(t, err)

	result := gen.resolveEnvVar("plain-value")
	assert.Equal(t, "plain-value", result)

	result = gen.resolveEnvVar("${NONEXISTENT_VAR:-default}")
	assert.Equal(t, "default", result)

	os.Setenv("TEST_VAR", "actual-value")
	defer os.Unsetenv("TEST_VAR")

	result = gen.resolveEnvVar("${TEST_VAR:-default}")
	assert.Equal(t, "actual-value", result)

	result = gen.resolveEnvVar("${MALFORMED")
	assert.Equal(t, "${MALFORMED", result)
}

func TestGenerator_WriteComposeFile(t *testing.T) {
	gen, err := NewGenerator("test-project")
	require.NoError(t, err)

	tempDir := t.TempDir()
	content := []byte("version: '3.8'\nservices:\n  test:\n    image: test:latest")

	err = gen.WriteComposeFile(content, tempDir)
	require.NoError(t, err)

	filePath := tempDir + "/docker-compose.yml"
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, content, data)

	tempDir = t.TempDir()
	nestedDir := tempDir + "/nested/path"
	content = []byte("test content")

	err = gen.WriteComposeFile(content, nestedDir)
	require.NoError(t, err)

	filePath = nestedDir + "/docker-compose.yml"
	_, err = os.Stat(filePath)
	assert.NoError(t, err)
}

func TestGenerator_GenerateFromServiceConfigs(t *testing.T) {
	gen, err := NewGenerator("test-project")
	require.NoError(t, err)

	tempDir := t.TempDir()
	originalPath := tempDir + "/.otto-stack/generated"
	err = os.MkdirAll(originalPath, 0755)
	require.NoError(t, err)

	configs := []types.ServiceConfig{
		fixtures.NewServiceConfig("redis").WithImage("redis:latest").Build(),
	}

	_ = gen.GenerateFromServiceConfigs(configs, "test-project")
}
