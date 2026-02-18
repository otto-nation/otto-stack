package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadSingleConfig_InvalidYAML(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Write invalid YAML
	err := os.WriteFile(configPath, []byte("invalid: yaml: content: ["), 0644)
	assert.NoError(t, err)

	_, err = loadSingleConfig(configPath)
	assert.Error(t, err)
}

func TestLoadProjectConfig_InvalidYAML(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Write invalid YAML
	err := os.WriteFile(configPath, []byte("invalid: yaml: content: ["), 0644)
	assert.NoError(t, err)

	_, err = LoadProjectConfig(configPath)
	assert.Error(t, err)
}

func TestLoadProjectConfig_LocalConfigInvalidYAML(t *testing.T) {
	tempDir := t.TempDir()

	// Create .otto-stack directory
	ottoDir := filepath.Join(tempDir, ".otto-stack")
	err := os.MkdirAll(ottoDir, 0755)
	assert.NoError(t, err)

	// Create valid base config
	baseConfigPath := filepath.Join(ottoDir, "config.yaml")
	err = os.WriteFile(baseConfigPath, []byte("project:\n  name: test\nstack:\n  services:\n    - postgres"), 0644)
	assert.NoError(t, err)

	// Create invalid local config
	localConfigPath := filepath.Join(ottoDir, "config.local.yaml")
	err = os.WriteFile(localConfigPath, []byte("invalid: yaml: ["), 0644)
	assert.NoError(t, err)

	// Change to temp dir so paths work
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	_, err = LoadProjectConfig(filepath.Join(".otto-stack", "config.yaml"))
	assert.Error(t, err)
}
