package registry

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestProjectExists(t *testing.T) {
	t.Run("returns false when project not found", func(t *testing.T) {
		exists := projectExists("nonexistent-project-xyz")
		assert.False(t, exists)
	})

	t.Run("returns true when project exists", func(t *testing.T) {
		homeDir, err := os.UserHomeDir()
		assert.NoError(t, err)

		// Create a test project directory
		testProject := filepath.Join(homeDir, "projects", "test-project-exists", core.OttoStackDir)
		err = os.MkdirAll(testProject, 0755)
		assert.NoError(t, err)
		defer os.RemoveAll(filepath.Join(homeDir, "projects", "test-project-exists"))

		exists := projectExists("test-project-exists")
		assert.True(t, exists)
	})
}
