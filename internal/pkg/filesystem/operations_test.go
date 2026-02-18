package filesystem

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureDir(t *testing.T) {
	t.Run("creates directory if it doesn't exist", func(t *testing.T) {
		tempDir := t.TempDir()
		testDir := filepath.Join(tempDir, "test", "nested", "dir")

		err := EnsureDir(testDir)
		require.NoError(t, err)

		// Verify directory was created
		info, err := os.Stat(testDir)
		require.NoError(t, err)
		assert.True(t, info.IsDir())
	})

	t.Run("succeeds if directory already exists", func(t *testing.T) {
		tempDir := t.TempDir()

		// Call twice - second call should succeed
		err := EnsureDir(tempDir)
		require.NoError(t, err)

		err = EnsureDir(tempDir)
		require.NoError(t, err)
	})
}

func TestWriteFile(t *testing.T) {
	t.Run("writes file successfully", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")
		content := []byte("test content")

		err := WriteFile(testFile, content, 0644)
		require.NoError(t, err)

		// Verify file was created with correct content
		readContent, err := os.ReadFile(testFile)
		require.NoError(t, err)
		assert.Equal(t, content, readContent)
	})

	t.Run("creates nested directories", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "nested", "dirs", "test.txt")
		content := []byte("nested content")

		err := WriteFile(testFile, content, 0644)
		require.NoError(t, err)

		// Verify file exists
		readContent, err := os.ReadFile(testFile)
		require.NoError(t, err)
		assert.Equal(t, content, readContent)
	})

	t.Run("overwrites existing file", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")

		// Write initial content
		err := WriteFile(testFile, []byte("initial"), 0644)
		require.NoError(t, err)

		// Overwrite with new content
		newContent := []byte("updated")
		err = WriteFile(testFile, newContent, 0644)
		require.NoError(t, err)

		// Verify new content
		readContent, err := os.ReadFile(testFile)
		require.NoError(t, err)
		assert.Equal(t, newContent, readContent)
	})

	t.Run("fails when directory creation fails", func(t *testing.T) {
		tempDir := t.TempDir()
		existingFile := filepath.Join(tempDir, "file.txt")
		err := os.WriteFile(existingFile, []byte("test"), 0644)
		require.NoError(t, err)

		// Try to write to a path where parent is a file
		err = WriteFile(filepath.Join(existingFile, "subfile.txt"), []byte("test"), 0644)
		assert.Error(t, err)
	})
}
