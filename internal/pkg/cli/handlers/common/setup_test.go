//go:build unit

package common

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestBuildStackContext(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()
	ottoDir := filepath.Join(tempDir, core.OttoStackDir)
	require.NoError(t, os.MkdirAll(ottoDir, 0755))

	// Change to temp directory so core.OttoStackDir resolves correctly
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	t.Run("success with valid config", func(t *testing.T) {
		// Create test config
		cfg := config.Config{
			Project: config.ProjectConfig{Name: "test-project"},
		}
		configData, err := yaml.Marshal(cfg)
		require.NoError(t, err)

		configPath := filepath.Join(ottoDir, core.ConfigFileName)
		require.NoError(t, os.WriteFile(configPath, configData, 0644))

		// Test BuildStackContext
		cmd := &cobra.Command{}
		args := []string{"postgres", "redis"}

		ctx, err := BuildStackContext(cmd, args)

		assert.NoError(t, err)
		assert.Equal(t, "test-project", ctx.Project.Name)
		assert.Equal(t, args, ctx.Services.Names)
	})

	t.Run("error when config missing", func(t *testing.T) {
		// Remove config file
		configPath := filepath.Join(ottoDir, core.ConfigFileName)
		os.Remove(configPath)

		cmd := &cobra.Command{}
		args := []string{}

		_, err := BuildStackContext(cmd, args)

		assert.Error(t, err)
	})
}

func TestSetupCoreCommand(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()
	ottoDir := filepath.Join(tempDir, core.OttoStackDir)
	require.NoError(t, os.MkdirAll(ottoDir, 0755))

	// Change to temp directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	t.Run("error when not initialized", func(t *testing.T) {
		// Ensure no config exists
		configPath := filepath.Join(ottoDir, core.ConfigFileName)
		os.Remove(configPath)

		ctx := context.Background()
		base := &base.BaseCommand{}

		setup, cleanup, err := SetupCoreCommand(ctx, base)

		assert.Error(t, err)
		assert.Nil(t, setup)
		assert.Nil(t, cleanup)
	})

	t.Run("success with valid config", func(t *testing.T) {
		// Create test config
		cfg := config.Config{
			Project: config.ProjectConfig{Name: "test-project"},
		}
		configData, err := yaml.Marshal(cfg)
		require.NoError(t, err)

		configPath := filepath.Join(ottoDir, core.ConfigFileName)
		require.NoError(t, os.WriteFile(configPath, configData, 0644))

		ctx := context.Background()
		base := &base.BaseCommand{}

		setup, cleanup, err := SetupCoreCommand(ctx, base)

		if err != nil {
			// Docker might not be available in test environment
			t.Skipf("Skipping Docker test: %v", err)
		}

		assert.NotNil(t, setup)
		assert.NotNil(t, cleanup)
		assert.NotNil(t, setup.Config)
		assert.Equal(t, "test-project", setup.Config.Project.Name)

		// Test cleanup function
		cleanup()
	})
}

func TestLoadProjectConfig(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("success with base config only", func(t *testing.T) {
		cfg := config.Config{
			Project: config.ProjectConfig{Name: "test-project"},
			Stack:   config.StackConfig{Enabled: []string{"postgres"}},
		}
		configData, err := yaml.Marshal(cfg)
		require.NoError(t, err)

		configPath := filepath.Join(tempDir, "config.yml")
		require.NoError(t, os.WriteFile(configPath, configData, 0644))

		result, err := LoadProjectConfig(configPath)

		assert.NoError(t, err)
		assert.Equal(t, "test-project", result.Project.Name)
		assert.Equal(t, []string{"postgres"}, result.Stack.Enabled)
	})

	t.Run("error with invalid config", func(t *testing.T) {
		configPath := filepath.Join(tempDir, "invalid.yml")
		require.NoError(t, os.WriteFile(configPath, []byte("invalid: yaml: content"), 0644))

		_, err := LoadProjectConfig(configPath)

		assert.Error(t, err)
	})

	t.Run("error with missing file", func(t *testing.T) {
		configPath := filepath.Join(tempDir, "missing.yml")

		_, err := LoadProjectConfig(configPath)

		assert.Error(t, err)
	})
}
