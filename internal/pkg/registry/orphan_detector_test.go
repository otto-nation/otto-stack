package registry

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/stretchr/testify/assert"
)

func TestOrphanDetector_buildContainerMap(t *testing.T) {
	detector := &OrphanDetector{}

	containers := []docker.ContainerInfo{
		{Name: "container1"},
		{Name: "container2"},
		{Name: "container3"},
	}

	containerMap := detector.buildContainerMap(containers)

	assert.Len(t, containerMap, 3)
	assert.True(t, containerMap["container1"])
	assert.True(t, containerMap["container2"])
	assert.True(t, containerMap["container3"])
	assert.False(t, containerMap["nonexistent"])
}

func TestOrphanDetector_buildContainerMap_Empty(t *testing.T) {
	detector := &OrphanDetector{}

	containerMap := detector.buildContainerMap([]docker.ContainerInfo{})

	assert.Empty(t, containerMap)
}

func TestOrphanDetector_FindOrphans_LoadError(t *testing.T) {
	tempDir := t.TempDir()
	registryPath := filepath.Join(tempDir, "registry.yaml")

	// Use tab character in YAML which is invalid per YAML spec
	err := os.WriteFile(registryPath, []byte("shared_containers:\n\tinvalid"), 0644)
	assert.NoError(t, err)

	manager := NewManager(registryPath)
	detector := NewOrphanDetector(manager)

	_, err = detector.FindOrphans()
	assert.Error(t, err)
}
