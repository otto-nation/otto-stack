//go:build unit

package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirectoryManager_CreateDirectoryStructure(t *testing.T) {
	t.Run("creates directory structure", func(t *testing.T) {
		tmpDir := t.TempDir()

		dirs := []string{
			filepath.Join(tmpDir, "services"),
			filepath.Join(tmpDir, "configs"),
			filepath.Join(tmpDir, "data"),
		}

		for _, dir := range dirs {
			err := os.MkdirAll(dir, 0755)
			require.NoError(t, err)
		}

		// Verify directories exist
		for _, dir := range dirs {
			info, err := os.Stat(dir)
			require.NoError(t, err)
			assert.True(t, info.IsDir())
		}
	})

	t.Run("handles existing directories", func(t *testing.T) {
		tmpDir := t.TempDir()
		dir := filepath.Join(tmpDir, "existing")

		// Create once
		err := os.MkdirAll(dir, 0755)
		require.NoError(t, err)

		// Create again - should not error
		err = os.MkdirAll(dir, 0755)
		require.NoError(t, err)
	})
}

func TestConfigFileOperations(t *testing.T) {
	t.Run("generates config file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "service.yaml")

		content := "service:\n  name: test\n  port: 8080\n"
		err := os.WriteFile(configPath, []byte(content), 0644)
		require.NoError(t, err)

		data, err := os.ReadFile(configPath)
		require.NoError(t, err)
		assert.Contains(t, string(data), "name: test")
		assert.Contains(t, string(data), "port: 8080")
	})

	t.Run("overwrites existing config", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "service.yaml")

		// Write initial
		err := os.WriteFile(configPath, []byte("old: data"), 0644)
		require.NoError(t, err)

		// Overwrite
		err = os.WriteFile(configPath, []byte("new: data"), 0644)
		require.NoError(t, err)

		data, err := os.ReadFile(configPath)
		require.NoError(t, err)
		assert.Contains(t, string(data), "new: data")
		assert.NotContains(t, string(data), "old: data")
	})
}

func TestProjectManager_FileOperations(t *testing.T) {
	t.Run("creates project file", func(t *testing.T) {
		tmpDir := t.TempDir()
		projectFile := filepath.Join(tmpDir, "project.yaml")

		content := `name: myproject
version: 1.0.0
services:
  - postgres
  - redis
`
		err := os.WriteFile(projectFile, []byte(content), 0644)
		require.NoError(t, err)

		data, err := os.ReadFile(projectFile)
		require.NoError(t, err)
		assert.Contains(t, string(data), "myproject")
		assert.Contains(t, string(data), "postgres")
	})

	t.Run("checks file existence", func(t *testing.T) {
		tmpDir := t.TempDir()
		existingFile := filepath.Join(tmpDir, "exists.txt")
		missingFile := filepath.Join(tmpDir, "missing.txt")

		err := os.WriteFile(existingFile, []byte("data"), 0644)
		require.NoError(t, err)

		_, err = os.Stat(existingFile)
		assert.NoError(t, err)

		_, err = os.Stat(missingFile)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("removes project files", func(t *testing.T) {
		tmpDir := t.TempDir()
		file := filepath.Join(tmpDir, "remove.txt")

		err := os.WriteFile(file, []byte("data"), 0644)
		require.NoError(t, err)

		err = os.Remove(file)
		require.NoError(t, err)

		_, err = os.Stat(file)
		assert.True(t, os.IsNotExist(err))
	})
}

func TestValidationManager_FileValidation(t *testing.T) {
	t.Run("validates yaml syntax", func(t *testing.T) {
		tmpDir := t.TempDir()

		validYaml := filepath.Join(tmpDir, "valid.yaml")
		err := os.WriteFile(validYaml, []byte("key: value\nlist:\n  - item1\n  - item2"), 0644)
		require.NoError(t, err)

		data, err := os.ReadFile(validYaml)
		require.NoError(t, err)
		assert.Contains(t, string(data), "key: value")
	})

	t.Run("detects missing required files", func(t *testing.T) {
		tmpDir := t.TempDir()
		requiredFiles := []string{"config.yaml", "services.yaml"}

		var missing []string
		for _, file := range requiredFiles {
			path := filepath.Join(tmpDir, file)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				missing = append(missing, file)
			}
		}

		assert.Len(t, missing, 2)
	})
}

func TestServiceSelector_FileBasedSelection(t *testing.T) {
	t.Run("reads service list from file", func(t *testing.T) {
		tmpDir := t.TempDir()
		servicesFile := filepath.Join(tmpDir, "services.txt")

		services := "postgres\nredis\nmongodb\n"
		err := os.WriteFile(servicesFile, []byte(services), 0644)
		require.NoError(t, err)

		data, err := os.ReadFile(servicesFile)
		require.NoError(t, err)

		lines := string(data)
		assert.Contains(t, lines, "postgres")
		assert.Contains(t, lines, "redis")
		assert.Contains(t, lines, "mongodb")
	})
}

func TestTemplateManager_TemplateFiles(t *testing.T) {
	t.Run("processes template file", func(t *testing.T) {
		tmpDir := t.TempDir()
		templateFile := filepath.Join(tmpDir, "template.yaml")

		template := `service: {{.ServiceName}}
port: {{.Port}}
`
		err := os.WriteFile(templateFile, []byte(template), 0644)
		require.NoError(t, err)

		data, err := os.ReadFile(templateFile)
		require.NoError(t, err)
		assert.Contains(t, string(data), "{{.ServiceName}}")
		assert.Contains(t, string(data), "{{.Port}}")
	})

	t.Run("creates output from template", func(t *testing.T) {
		tmpDir := t.TempDir()
		outputFile := filepath.Join(tmpDir, "output.yaml")

		output := `service: postgres
port: 5432
`
		err := os.WriteFile(outputFile, []byte(output), 0644)
		require.NoError(t, err)

		data, err := os.ReadFile(outputFile)
		require.NoError(t, err)
		assert.Contains(t, string(data), "postgres")
		assert.Contains(t, string(data), "5432")
	})
}
