package registry

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestManager_ProjectExists(t *testing.T) {
	m := &Manager{}

	t.Run("returns false when project not found", func(t *testing.T) {
		tempDir := t.TempDir()
		exists := m.projectExists("nonexistent-project", tempDir)
		assert.False(t, exists)
	})

	t.Run("returns true when project exists in projects dir", func(t *testing.T) {
		tempDir := t.TempDir()
		projectPath := filepath.Join(tempDir, "projects", "test-project", core.OttoStackDir)
		err := os.MkdirAll(projectPath, 0755)
		assert.NoError(t, err)

		exists := m.projectExists("test-project", tempDir)
		assert.True(t, exists)
	})

	t.Run("returns true when project exists in git dir", func(t *testing.T) {
		tempDir := t.TempDir()
		projectPath := filepath.Join(tempDir, "git", "my-repos", "test-project", core.OttoStackDir)
		err := os.MkdirAll(projectPath, 0755)
		assert.NoError(t, err)

		exists := m.projectExists("test-project", tempDir)
		assert.True(t, exists)
	})

	t.Run("returns true when project exists in home dir", func(t *testing.T) {
		tempDir := t.TempDir()
		projectPath := filepath.Join(tempDir, "test-project", core.OttoStackDir)
		err := os.MkdirAll(projectPath, 0755)
		assert.NoError(t, err)

		exists := m.projectExists("test-project", tempDir)
		assert.True(t, exists)
	})
}
