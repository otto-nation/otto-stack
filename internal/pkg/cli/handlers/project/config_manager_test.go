//go:build unit

package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigManager_GenerateServiceConfigContent(t *testing.T) {
	cm := NewConfigManager()

	t.Run("generates config content", func(t *testing.T) {
		content := cm.generateServiceConfigContent("postgres")
		assert.NotEmpty(t, content)
		assert.Contains(t, content, "postgres")
	})
}

func TestConfigManager_GenerateServiceConfig(t *testing.T) {
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	os.Chdir(tempDir)
	configDir := filepath.Join(core.OttoStackDir, core.ServiceConfigsDir)
	require.NoError(t, os.MkdirAll(configDir, 0755))

	cm := NewConfigManager()

	t.Run("generates service config file", func(t *testing.T) {
		err := cm.generateServiceConfig("postgres")
		require.NoError(t, err)

		configPath := filepath.Join(configDir, "postgres"+core.YMLFileExtension)
		_, err = os.Stat(configPath)
		assert.NoError(t, err)
	})
}

func TestConfigManager_GenerateServiceConfigs(t *testing.T) {
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	os.Chdir(tempDir)
	configDir := filepath.Join(core.OttoStackDir, core.ServiceConfigsDir)
	require.NoError(t, os.MkdirAll(configDir, 0755))

	cm := NewConfigManager()
	baseCmd := &base.BaseCommand{Output: ui.NewOutput()}

	t.Run("generates configs for non-hidden services", func(t *testing.T) {
		configs := []types.ServiceConfig{
			{Name: "postgres", Hidden: false},
			{Name: "redis", Hidden: true},
		}
		cm.GenerateServiceConfigs(configs, false, baseCmd)

		postgresPath := filepath.Join(configDir, "postgres"+core.YMLFileExtension)
		redisPath := filepath.Join(configDir, "redis"+core.YMLFileExtension)

		_, err := os.Stat(postgresPath)
		assert.NoError(t, err)

		_, err = os.Stat(redisPath)
		assert.Error(t, err)
	})

	t.Run("skips shared services when sharing enabled", func(t *testing.T) {
		configs := []types.ServiceConfig{
			{Name: "mysql", Shareable: true},
		}
		cm.GenerateServiceConfigs(configs, true, baseCmd)

		mysqlPath := filepath.Join(configDir, "mysql"+core.YMLFileExtension)
		_, err := os.Stat(mysqlPath)
		assert.Error(t, err)
	})
}
