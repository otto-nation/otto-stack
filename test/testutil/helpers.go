package testutil

import (
	"os"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/stretchr/testify/require"
)

// NewTestManager creates a services manager for testing
func NewTestManager(t *testing.T) *services.Manager {
	manager, err := services.New()
	require.NoError(t, err)
	return manager
}

// WithTempDir runs a test function in a temporary directory
func WithTempDir(t *testing.T, fn func(tempDir string)) {
	tempDir := t.TempDir()
	originalDir, err := os.Getwd()
	require.NoError(t, err)

	defer func() {
		err := os.Chdir(originalDir)
		require.NoError(t, err)
	}()

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	fn(tempDir)
}

// AssertFileExists checks that a file exists and is readable
func AssertFileExists(t *testing.T, path string) {
	info, err := os.Stat(path)
	require.NoError(t, err, "file should exist: %s", path)
	require.False(t, info.IsDir(), "path should be a file, not directory: %s", path)
}

// AssertDirExists checks that a directory exists
func AssertDirExists(t *testing.T, path string) {
	info, err := os.Stat(path)
	require.NoError(t, err, "directory should exist: %s", path)
	require.True(t, info.IsDir(), "path should be a directory: %s", path)
}
