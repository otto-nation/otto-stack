//go:build unit

package services

import (
	"testing"

	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/test/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateProcessor_Process_NoVariables(t *testing.T) {
	processor := NewTemplateProcessor()
	script := "echo 'hello world'"
	config := fixtures.NewServiceConfig("test").Build()

	result, err := processor.Process(script, config, []servicetypes.ServiceConfig{})
	require.NoError(t, err)
	assert.Equal(t, "echo 'hello world'", result)
}

func TestTemplateProcessor_Process_WithDependencies(t *testing.T) {
	processor := NewTemplateProcessor()
	script := "echo 'setup'"
	config := fixtures.NewServiceConfig("postgres").Build()
	appConfig := fixtures.LoadService(t, "with-deps")

	result, err := processor.Process(script, config, []servicetypes.ServiceConfig{appConfig})
	require.NoError(t, err)
	assert.Contains(t, result, "setup")
}

func TestTemplateProcessor_Process_InvalidSyntax(t *testing.T) {
	processor := NewTemplateProcessor()
	script := "{{.Invalid"
	config := fixtures.NewServiceConfig("test").Build()

	_, err := processor.Process(script, config, []servicetypes.ServiceConfig{})
	assert.Error(t, err)
}

func TestTemplateProcessor_Process_EmptyData(t *testing.T) {
	processor := NewTemplateProcessor()
	script := "{{.NonExistent}}"
	config := fixtures.NewServiceConfig("test").Build()

	result, err := processor.Process(script, config, []servicetypes.ServiceConfig{})
	require.NoError(t, err)
	assert.Equal(t, "<no value>", result)
}

func TestTemplateProcessor_serviceDependsOn_True(t *testing.T) {
	processor := NewTemplateProcessor()
	config := fixtures.LoadService(t, "with-deps")

	assert.True(t, processor.serviceDependsOn(config, "postgres"))
}

func TestTemplateProcessor_serviceDependsOn_False(t *testing.T) {
	processor := NewTemplateProcessor()
	config := fixtures.LoadService(t, "with-deps")

	assert.False(t, processor.serviceDependsOn(config, "redis"))
}

func TestTemplateProcessor_serviceDependsOn_NoDeps(t *testing.T) {
	processor := NewTemplateProcessor()
	config := fixtures.LoadService(t, "minimal")
	assert.False(t, processor.serviceDependsOn(config, "postgres"))
}

func TestTemplateProcessor_serviceDependsOn_NilDeps(t *testing.T) {
	processor := NewTemplateProcessor()
	config := fixtures.LoadService(t, "minimal")
	assert.False(t, processor.serviceDependsOn(config, "postgres"))
}

func TestTemplateProcessor_collectTemplateData_WithDeps(t *testing.T) {
	processor := NewTemplateProcessor()
	config := fixtures.NewServiceConfig("postgres").Build()
	appConfig := fixtures.LoadService(t, "with-deps")

	data := processor.collectTemplateData(config, []servicetypes.ServiceConfig{appConfig})
	assert.NotNil(t, data)
}

func TestTemplateProcessor_collectTemplateData_NoDeps(t *testing.T) {
	processor := NewTemplateProcessor()
	config := fixtures.NewServiceConfig("postgres").Build()
	redisConfig := fixtures.LoadService(t, "minimal")
	redisConfig.Name = "redis"

	data := processor.collectTemplateData(config, []servicetypes.ServiceConfig{redisConfig})
	assert.NotNil(t, data)
	assert.Empty(t, data)
}

func TestTemplateProcessor_collectTemplateData_EmptyConfigs(t *testing.T) {
	processor := NewTemplateProcessor()
	config := fixtures.NewServiceConfig("postgres").Build()

	data := processor.collectTemplateData(config, []servicetypes.ServiceConfig{})
	assert.NotNil(t, data)
	assert.Empty(t, data)
}

func TestTemplateProcessor_addConfigData(t *testing.T) {
	processor := NewTemplateProcessor()
	templateData := make(map[string]any)
	config := fixtures.NewServiceConfig("postgres").Build()

	processor.addConfigData(templateData, config)
	assert.NotNil(t, templateData)
}

func TestTemplateProcessor_addConfigData_Multiple(t *testing.T) {
	processor := NewTemplateProcessor()
	templateData := make(map[string]any)
	configs := []servicetypes.ServiceConfig{
		fixtures.NewServiceConfig("postgres").Build(),
		fixtures.NewServiceConfig("redis").Build(),
	}

	for _, cfg := range configs {
		processor.addConfigData(templateData, cfg)
	}
	assert.NotNil(t, templateData)
}
