package registry

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManager_Save_ErrorHandling(t *testing.T) {
	// Create a file where directory should be
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "file.txt")
	err := os.WriteFile(filePath, []byte("test"), 0644)
	assert.NoError(t, err)

	// Try to save registry with path that conflicts
	m := &Manager{registryPath: filepath.Join(filePath, "registry.yaml")}
	registry := &Registry{Containers: make(map[string]*ContainerInfo)}

	err = m.Save(registry)
	assert.Error(t, err)
}

func TestManager_Register_ErrorPath(t *testing.T) {
	tempDir := t.TempDir()
	registryPath := filepath.Join(tempDir, "registry.yaml")
	m := &Manager{registryPath: registryPath}

	// Register successfully
	err := m.Register("test-service", "test-container", "test-project")
	assert.NoError(t, err)

	// Verify it was registered
	container, err := m.Get("test-service")
	assert.NoError(t, err)
	assert.NotNil(t, container)
	assert.Contains(t, container.Projects, "test-project")
}

func TestManager_Unregister_ErrorPath(t *testing.T) {
	tempDir := t.TempDir()
	registryPath := filepath.Join(tempDir, "registry.yaml")
	m := &Manager{registryPath: registryPath}

	// Register and then unregister
	err := m.Register("test-service", "test-container", "test-project")
	assert.NoError(t, err)

	err = m.Unregister("test-service", "test-project")
	assert.NoError(t, err)

	// Verify it was unregistered
	container, err := m.Get("test-service")
	assert.NoError(t, err)
	assert.Nil(t, container)
}
