//go:build unit

package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestCheckInitialization(t *testing.T) {
	originalDir, _ := os.Getwd()

	t.Run("returns error when not initialized", func(t *testing.T) {
		tempDir := t.TempDir()
		defer os.Chdir(originalDir)

		// Change to temp directory where no config exists
		os.Chdir(tempDir)

		err := CheckInitialization()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("returns nil when initialized", func(t *testing.T) {
		tempDir := t.TempDir()
		defer os.Chdir(originalDir)

		// Change to temp directory
		os.Chdir(tempDir)

		// Create otto-stack directory and config file
		ottoDir := filepath.Join(tempDir, core.OttoStackDir)
		err := os.MkdirAll(ottoDir, 0755)
		assert.NoError(t, err)

		configPath := filepath.Join(ottoDir, core.ConfigFileName)
		err = os.WriteFile(configPath, []byte("project:\n  name: test"), 0644)
		assert.NoError(t, err)

		// Now check should pass
		err = CheckInitialization()
		assert.NoError(t, err)
	})

	t.Run("handles directory exists but no config file", func(t *testing.T) {
		tempDir := t.TempDir()
		defer os.Chdir(originalDir)

		// Change to temp directory
		os.Chdir(tempDir)

		// Create otto-stack directory but no config file
		ottoDir := filepath.Join(tempDir, core.OttoStackDir)
		err := os.MkdirAll(ottoDir, 0755)
		assert.NoError(t, err)

		// Should still return error since config file doesn't exist
		err = CheckInitialization()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("handles permission issues", func(t *testing.T) {
		tempDir := t.TempDir()
		defer os.Chdir(originalDir)

		// Change to temp directory
		os.Chdir(tempDir)

		// Create directory with restricted permissions (if possible)
		restrictedDir := filepath.Join(tempDir, "restricted")
		err := os.MkdirAll(restrictedDir, 0000)
		if err == nil {
			defer os.Chmod(restrictedDir, 0755) // Cleanup

			os.Chdir(restrictedDir)

			// Should handle permission errors gracefully
			err = CheckInitialization()
			assert.Error(t, err) // Should error due to permissions or missing file
		}
	})
}
