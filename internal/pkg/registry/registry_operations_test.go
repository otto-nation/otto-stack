//go:build unit

package registry

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_GetListIsShared(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	t.Run("Get returns nil for nonexistent service", func(t *testing.T) {
		info, err := manager.Get("nonexistent")
		require.NoError(t, err)
		assert.Nil(t, info)
	})

	t.Run("List returns empty map for new registry", func(t *testing.T) {
		containers, err := manager.List()
		require.NoError(t, err)
		assert.NotNil(t, containers)
		assert.Empty(t, containers)
	})

	t.Run("IsShared returns false for nonexistent service", func(t *testing.T) {
		shared, err := manager.IsShared("nonexistent")
		require.NoError(t, err)
		assert.False(t, shared)
	})

	t.Run("Get returns container after registration", func(t *testing.T) {
		registry, err := manager.Load()
		require.NoError(t, err)

		registry.Containers["postgres"] = &ContainerInfo{
			Name:     "postgres-shared",
			Projects: []ProjectRef{{Name: "project1"}},
		}

		err = manager.Save(registry)
		require.NoError(t, err)

		info, err := manager.Get("postgres")
		require.NoError(t, err)
		assert.NotNil(t, info)
		assert.Equal(t, "postgres-shared", info.Name)
	})

	t.Run("List returns all containers", func(t *testing.T) {
		registry, err := manager.Load()
		require.NoError(t, err)

		registry.Containers["redis"] = &ContainerInfo{
			Name:     "redis-shared",
			Projects: []ProjectRef{{Name: "project2"}},
		}

		err = manager.Save(registry)
		require.NoError(t, err)

		containers, err := manager.List()
		require.NoError(t, err)
		assert.Len(t, containers, 2) // postgres + redis
	})

	t.Run("IsShared returns true for registered service", func(t *testing.T) {
		shared, err := manager.IsShared("postgres")
		require.NoError(t, err)
		assert.True(t, shared)
	})
}

func TestManager_LoadSave(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	t.Run("Load creates new registry if file doesn't exist", func(t *testing.T) {
		registry, err := manager.Load()
		require.NoError(t, err)
		assert.NotNil(t, registry)
		assert.NotNil(t, registry.Containers)
		assert.Empty(t, registry.Containers)
	})

	t.Run("Save persists registry to disk", func(t *testing.T) {
		registry := NewRegistry()
		registry.Containers["test"] = &ContainerInfo{
			Name:     "test-shared",
			Projects: []ProjectRef{{Name: "proj1", ConfigDir: t.TempDir()}},
		}

		err := manager.Save(registry)
		require.NoError(t, err)

		// Verify file exists
		_, err = os.Stat(manager.registryPath)
		assert.NoError(t, err)
	})

	t.Run("Load reads saved registry", func(t *testing.T) {
		registry, err := manager.Load()
		require.NoError(t, err)
		assert.Len(t, registry.Containers, 1)
		assert.NotNil(t, registry.Containers["test"])
	})

	t.Run("Save handles empty registry", func(t *testing.T) {
		registry := NewRegistry()
		err := manager.Save(registry)
		require.NoError(t, err)
	})
}

func TestManager_RegisterUnregister(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	proj1Dir := t.TempDir()
	proj2Dir := t.TempDir()

	t.Run("Register adds new container", func(t *testing.T) {
		err := manager.Register("postgres", "postgres-shared-123", ProjectRef{Name: "project1", ConfigDir: proj1Dir})
		require.NoError(t, err)

		info, err := manager.Get("postgres")
		require.NoError(t, err)
		assert.NotNil(t, info)
		assert.Equal(t, "postgres-shared-123", info.Name)
		assert.Contains(t, info.Projects, ProjectRef{Name: "project1", ConfigDir: proj1Dir})
	})

	t.Run("Register adds project to existing container", func(t *testing.T) {
		err := manager.Register("postgres", "postgres-shared-123", ProjectRef{Name: "project2", ConfigDir: proj2Dir})
		require.NoError(t, err)

		info, err := manager.Get("postgres")
		require.NoError(t, err)
		assert.Len(t, info.Projects, 2)
		assert.Contains(t, info.Projects, ProjectRef{Name: "project1", ConfigDir: proj1Dir})
		assert.Contains(t, info.Projects, ProjectRef{Name: "project2", ConfigDir: proj2Dir})
	})

	t.Run("Register doesn't duplicate projects", func(t *testing.T) {
		err := manager.Register("postgres", "postgres-shared-123", ProjectRef{Name: "project1", ConfigDir: proj1Dir})
		require.NoError(t, err)

		info, err := manager.Get("postgres")
		require.NoError(t, err)
		assert.Len(t, info.Projects, 2) // Still 2, not 3
	})

	t.Run("Unregister removes project from container", func(t *testing.T) {
		err := manager.Unregister("postgres", "project1")
		require.NoError(t, err)

		info, err := manager.Get("postgres")
		require.NoError(t, err)
		assert.NotNil(t, info)
		assert.Len(t, info.Projects, 1)
		assert.Contains(t, info.Projects, ProjectRef{Name: "project2", ConfigDir: proj2Dir})
	})

	t.Run("Unregister removes container when no projects remain", func(t *testing.T) {
		err := manager.Unregister("postgres", "project2")
		require.NoError(t, err)

		info, err := manager.Get("postgres")
		require.NoError(t, err)
		assert.Nil(t, info)
	})

	t.Run("Unregister handles nonexistent service", func(t *testing.T) {
		err := manager.Unregister("nonexistent", "project1")
		assert.NoError(t, err) // Should not error
	})
}

func TestManager_PurgeNonShareable(t *testing.T) {
	t.Run("RemovesBadEntries", func(t *testing.T) {
		tmpDir := t.TempDir()
		manager := NewManager(tmpDir)

		err := manager.Register("localstack-sqs", "otto-stack-localstack-sqs", ProjectRef{Name: "global", ConfigDir: t.TempDir()})
		require.NoError(t, err)

		shareableMap := map[string]bool{
			"redis":    true,
			"postgres": true,
		}
		err = manager.PurgeNonShareable(shareableMap)
		require.NoError(t, err)

		info, err := manager.Get("localstack-sqs")
		require.NoError(t, err)
		assert.Nil(t, info, "non-shareable entry should be removed")
	})

	t.Run("KeepsValidEntries", func(t *testing.T) {
		tmpDir := t.TempDir()
		manager := NewManager(tmpDir)

		err := manager.Register("redis", "otto-stack-redis", ProjectRef{Name: "project1", ConfigDir: t.TempDir()})
		require.NoError(t, err)

		shareableMap := map[string]bool{
			"redis": true,
		}
		err = manager.PurgeNonShareable(shareableMap)
		require.NoError(t, err)

		info, err := manager.Get("redis")
		require.NoError(t, err)
		assert.NotNil(t, info, "shareable entry should be kept")
	})

	t.Run("NoChanges", func(t *testing.T) {
		tmpDir := t.TempDir()
		manager := NewManager(tmpDir)

		err := manager.Register("redis", "otto-stack-redis", ProjectRef{Name: "project1", ConfigDir: t.TempDir()})
		require.NoError(t, err)

		shareableMap := map[string]bool{
			"redis":    true,
			"postgres": true,
		}
		// No error and no write needed when nothing changes
		err = manager.PurgeNonShareable(shareableMap)
		require.NoError(t, err)

		info, err := manager.Get("redis")
		require.NoError(t, err)
		assert.NotNil(t, info)
	})

	t.Run("EmptyRegistry", func(t *testing.T) {
		tmpDir := t.TempDir()
		manager := NewManager(tmpDir)

		shareableMap := map[string]bool{"redis": true}
		err := manager.PurgeNonShareable(shareableMap)
		require.NoError(t, err)
	})

	t.Run("EmptyShareableMap", func(t *testing.T) {
		tmpDir := t.TempDir()
		manager := NewManager(tmpDir)

		err := manager.Register("redis", "otto-stack-redis", ProjectRef{Name: "project1", ConfigDir: t.TempDir()})
		require.NoError(t, err)
		err = manager.Register("postgres", "otto-stack-postgres", ProjectRef{Name: "project1", ConfigDir: t.TempDir()})
		require.NoError(t, err)

		err = manager.PurgeNonShareable(map[string]bool{})
		require.NoError(t, err)

		containers, err := manager.List()
		require.NoError(t, err)
		assert.Empty(t, containers, "all entries should be removed when shareable map is empty")
	})
}

func TestManager_ValidateAgainstDocker(t *testing.T) {
	t.Skip("Skipping ValidateAgainstDocker tests that require a running Docker daemon")
}

func TestManager_LoadWithCorruptedFile(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	t.Run("heals corrupted registry file via Docker rebuild", func(t *testing.T) {
		// Create a corrupted file
		err := os.WriteFile(manager.registryPath, []byte("invalid yaml {{{"), 0644)
		require.NoError(t, err)

		// Load should no longer error — it calls rebuildFromDocker instead.
		// If Docker is unavailable (common in CI), rebuildFromDocker returns NewRegistry().
		reg, err := manager.Load()
		require.NoError(t, err)
		assert.NotNil(t, reg, "should return a valid registry even from corrupt file")
		assert.NotNil(t, reg.Containers)
	})
}

func TestManager_rebuildFromDocker(t *testing.T) {
	t.Run("SkippedDockerDependentCases", func(t *testing.T) {
		t.Skip("Skipping rebuildFromDocker tests that require a running Docker daemon")
	})
}
