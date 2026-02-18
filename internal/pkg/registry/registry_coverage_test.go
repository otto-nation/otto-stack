package registry

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindOrphans(t *testing.T) {
	tempDir := t.TempDir()

	manager := NewManager(tempDir)

	t.Run("finds orphans with empty projects", func(t *testing.T) {
		// Create a registry with a container that has no projects
		registry := &Registry{
			Containers: map[string]*ContainerInfo{
				"test-service": {
					Name:     "otto-stack-test",
					Projects: []string{}, // Empty projects = orphan
				},
			},
		}
		err := manager.Save(registry)
		require.NoError(t, err)

		// Find orphans
		orphans, err := manager.FindOrphans()
		require.NoError(t, err)
		assert.Len(t, orphans, 1)
		assert.Equal(t, "test-service", orphans[0].Service)
	})

	t.Run("no orphans when all have projects", func(t *testing.T) {
		registry := &Registry{
			Containers: map[string]*ContainerInfo{
				"test-service": {
					Name:     "otto-stack-test",
					Projects: []string{"project1"},
				},
			},
		}
		err := manager.Save(registry)
		require.NoError(t, err)

		orphans, err := manager.FindOrphans()
		require.NoError(t, err)
		assert.Empty(t, orphans)
	})
}

func TestIsSharedContainer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"shared container", "otto-stack-redis", true},
		{"shared container with suffix", "otto-stack-postgres-1", true},
		{"non-shared container", "my-redis", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSharedContainer(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractServiceName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"basic service", "otto-stack-redis", "redis"},
		{"service with suffix", "otto-stack-postgres-1", "postgres-1"},
		{"service with multiple dashes", "otto-stack-localstack-sns", "localstack-sns"},
		{"non-shared container", "my-redis", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractServiceName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGet(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	t.Run("returns container info when exists", func(t *testing.T) {
		registry := &Registry{
			Containers: map[string]*ContainerInfo{
				"redis": {
					Name:     "otto-stack-redis",
					Projects: []string{"project1"},
				},
			},
		}
		err := manager.Save(registry)
		require.NoError(t, err)

		info, err := manager.Get("redis")
		require.NoError(t, err)
		assert.Equal(t, "otto-stack-redis", info.Name)
		assert.Contains(t, info.Projects, "project1")
	})

	t.Run("returns nil when not found", func(t *testing.T) {
		registry := &Registry{Containers: map[string]*ContainerInfo{}}
		err := manager.Save(registry)
		require.NoError(t, err)

		info, err := manager.Get("nonexistent")
		require.NoError(t, err)
		assert.Nil(t, info)
	})
}

func TestList(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	t.Run("returns all containers", func(t *testing.T) {
		registry := &Registry{
			Containers: map[string]*ContainerInfo{
				"redis":    {Name: "otto-stack-redis", Projects: []string{"p1"}},
				"postgres": {Name: "otto-stack-postgres", Projects: []string{"p2"}},
			},
		}
		err := manager.Save(registry)
		require.NoError(t, err)

		containers, err := manager.List()
		require.NoError(t, err)
		assert.Len(t, containers, 2)
	})

	t.Run("returns empty list when no containers", func(t *testing.T) {
		registry := &Registry{Containers: map[string]*ContainerInfo{}}
		err := manager.Save(registry)
		require.NoError(t, err)

		containers, err := manager.List()
		require.NoError(t, err)
		assert.Empty(t, containers)
	})
}

func TestIsShared(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	t.Run("returns true when service is shared", func(t *testing.T) {
		registry := &Registry{
			Containers: map[string]*ContainerInfo{
				"redis": {Name: "otto-stack-redis", Projects: []string{"p1"}},
			},
		}
		err := manager.Save(registry)
		require.NoError(t, err)

		shared, err := manager.IsShared("redis")
		require.NoError(t, err)
		assert.True(t, shared)
	})

	t.Run("returns false when service not shared", func(t *testing.T) {
		registry := &Registry{Containers: map[string]*ContainerInfo{}}
		err := manager.Save(registry)
		require.NoError(t, err)

		shared, err := manager.IsShared("nonexistent")
		require.NoError(t, err)
		assert.False(t, shared)
	})
}

func TestLoad_ErrorHandling(t *testing.T) {
	t.Run("handles invalid YAML", func(t *testing.T) {
		tempDir := t.TempDir()
		manager := NewManager(tempDir)

		// Write invalid YAML
		err := os.WriteFile(manager.registryPath, []byte("invalid: yaml: content: ["), 0644)
		require.NoError(t, err)

		_, err = manager.Load()
		assert.Error(t, err)
	})

	t.Run("returns empty registry for empty file", func(t *testing.T) {
		tempDir := t.TempDir()
		manager := NewManager(tempDir)

		// Create empty file
		err := os.WriteFile(manager.registryPath, []byte(""), 0644)
		require.NoError(t, err)

		registry, err := manager.Load()
		require.NoError(t, err)
		assert.NotNil(t, registry)
		assert.Empty(t, registry.Containers)
	})
}

func TestSave_ErrorHandling(t *testing.T) {
	t.Run("creates directory if not exists", func(t *testing.T) {
		tempDir := t.TempDir()
		nestedPath := filepath.Join(tempDir, "nested", "path")
		manager := NewManager(nestedPath)

		registry := &Registry{
			Containers: map[string]*ContainerInfo{
				"test": {Name: "otto-stack-test", Projects: []string{"p1"}},
			},
		}

		err := manager.Save(registry)
		require.NoError(t, err)

		// Verify file was created
		_, err = os.Stat(manager.registryPath)
		assert.NoError(t, err)
	})
}

func TestLoad_EdgeCases(t *testing.T) {
	t.Run("loads registry with nil containers", func(t *testing.T) {
		tempDir := t.TempDir()
		manager := NewManager(tempDir)

		// Write YAML without containers field
		err := os.WriteFile(manager.registryPath, []byte("version: 1\n"), 0644)
		require.NoError(t, err)

		registry, err := manager.Load()
		require.NoError(t, err)
		assert.NotNil(t, registry.Containers)
	})

	t.Run("creates new registry if file doesn't exist", func(t *testing.T) {
		tempDir := t.TempDir()
		nonExistentPath := filepath.Join(tempDir, "nonexistent")
		manager := NewManager(nonExistentPath)

		registry, err := manager.Load()
		require.NoError(t, err)
		assert.NotNil(t, registry)
		assert.NotNil(t, registry.Containers)
	})
}

func TestSave_EdgeCases(t *testing.T) {
	t.Run("saves empty registry", func(t *testing.T) {
		tempDir := t.TempDir()
		manager := NewManager(tempDir)

		registry := NewRegistry()
		err := manager.Save(registry)
		require.NoError(t, err)

		loaded, err := manager.Load()
		require.NoError(t, err)
		assert.Empty(t, loaded.Containers)
	})

	t.Run("overwrites existing registry", func(t *testing.T) {
		tempDir := t.TempDir()
		manager := NewManager(tempDir)

		registry1 := &Registry{
			Containers: map[string]*ContainerInfo{
				"test1": {Name: "otto-stack-test1", Projects: []string{"p1"}},
			},
		}
		err := manager.Save(registry1)
		require.NoError(t, err)

		registry2 := &Registry{
			Containers: map[string]*ContainerInfo{
				"test2": {Name: "otto-stack-test2", Projects: []string{"p2"}},
			},
		}
		err = manager.Save(registry2)
		require.NoError(t, err)

		loaded, err := manager.Load()
		require.NoError(t, err)
		assert.Len(t, loaded.Containers, 1)
		assert.Contains(t, loaded.Containers, "test2")
		assert.NotContains(t, loaded.Containers, "test1")
	})
}

func TestRegister(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	t.Run("registers new service", func(t *testing.T) {
		err := manager.Register("redis", "otto-stack-redis", "project1")
		require.NoError(t, err)

		info, err := manager.Get("redis")
		require.NoError(t, err)
		assert.Equal(t, "otto-stack-redis", info.Name)
		assert.Contains(t, info.Projects, "project1")
	})

	t.Run("adds project to existing service", func(t *testing.T) {
		err := manager.Register("redis", "otto-stack-redis", "project2")
		require.NoError(t, err)

		info, err := manager.Get("redis")
		require.NoError(t, err)
		assert.Contains(t, info.Projects, "project1")
		assert.Contains(t, info.Projects, "project2")
	})

	t.Run("does not duplicate projects", func(t *testing.T) {
		err := manager.Register("redis", "otto-stack-redis", "project1")
		require.NoError(t, err)

		info, err := manager.Get("redis")
		require.NoError(t, err)
		count := 0
		for _, p := range info.Projects {
			if p == "project1" {
				count++
			}
		}
		assert.Equal(t, 1, count)
	})
}

func TestUnregister(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	t.Run("removes project from service", func(t *testing.T) {
		err := manager.Register("redis", "otto-stack-redis", "project1")
		require.NoError(t, err)
		err = manager.Register("redis", "otto-stack-redis", "project2")
		require.NoError(t, err)

		err = manager.Unregister("redis", "project1")
		require.NoError(t, err)

		info, err := manager.Get("redis")
		require.NoError(t, err)
		assert.NotContains(t, info.Projects, "project1")
		assert.Contains(t, info.Projects, "project2")
	})

	t.Run("removes service when no projects left", func(t *testing.T) {
		err := manager.Unregister("redis", "project2")
		require.NoError(t, err)

		info, err := manager.Get("redis")
		require.NoError(t, err)
		assert.Nil(t, info)
	})

	t.Run("handles unregistering non-existent service", func(t *testing.T) {
		err := manager.Unregister("nonexistent", "project1")
		assert.NoError(t, err)
	})
}
