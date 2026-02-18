package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadServiceConfigFile_ReadError(t *testing.T) {
	_, err := loadServiceConfigFile("nonexistent-service")
	assert.Error(t, err)
}

func TestLoadServiceConfigFile_InvalidYAML(t *testing.T) {
	tempDir := t.TempDir()

	// Create .otto-stack/service-configs directory
	configDir := filepath.Join(tempDir, ".otto-stack", "service-configs")
	err := os.MkdirAll(configDir, 0755)
	assert.NoError(t, err)

	// Create invalid YAML file
	configPath := filepath.Join(configDir, "test.yml")
	err = os.WriteFile(configPath, []byte("invalid: yaml: ["), 0644)
	assert.NoError(t, err)

	// Change to temp dir
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	_, err = loadServiceConfigFile("test")
	assert.Error(t, err)
}
