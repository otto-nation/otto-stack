package filesystem

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureDir(t *testing.T) {
	t.Run("creates directory when it doesn't exist", func(t *testing.T) {
		tempDir := t.TempDir()
		testDir := filepath.Join(tempDir, "test", "nested", "dir")

		err := EnsureDir(testDir)
		assert.NoError(t, err)

		// Verify directory was created
		info, err := os.Stat(testDir)
		assert.NoError(t, err)
		assert.True(t, info.IsDir())
	})

	t.Run("succeeds when directory already exists", func(t *testing.T) {
		tempDir := t.TempDir()
		testDir := filepath.Join(tempDir, "existing")

		// Create directory first
		err := os.MkdirAll(testDir, 0755)
		assert.NoError(t, err)

		// Should not error when directory exists
		err = EnsureDir(testDir)
		assert.NoError(t, err)
	})
}

func TestCopyFile(t *testing.T) {
	t.Run("copies file successfully", func(t *testing.T) {
		tempDir := t.TempDir()
		srcFile := filepath.Join(tempDir, "source.txt")
		dstFile := filepath.Join(tempDir, "dest", "destination.txt")

		// Create source file
		content := []byte("test content")
		err := os.WriteFile(srcFile, content, 0644)
		assert.NoError(t, err)

		// Copy file
		err = CopyFile(srcFile, dstFile)
		assert.NoError(t, err)

		// Verify destination file exists and has correct content
		copiedContent, err := os.ReadFile(dstFile)
		assert.NoError(t, err)
		assert.Equal(t, content, copiedContent)
	})

	t.Run("creates destination directory", func(t *testing.T) {
		tempDir := t.TempDir()
		srcFile := filepath.Join(tempDir, "source.txt")
		dstFile := filepath.Join(tempDir, "nested", "deep", "destination.txt")

		// Create source file
		err := os.WriteFile(srcFile, []byte("content"), 0644)
		assert.NoError(t, err)

		// Copy to nested path
		err = CopyFile(srcFile, dstFile)
		assert.NoError(t, err)

		// Verify nested directory was created
		info, err := os.Stat(filepath.Dir(dstFile))
		assert.NoError(t, err)
		assert.True(t, info.IsDir())
	})

	t.Run("returns error for nonexistent source", func(t *testing.T) {
		tempDir := t.TempDir()
		srcFile := filepath.Join(tempDir, "nonexistent.txt")
		dstFile := filepath.Join(tempDir, "destination.txt")

		err := CopyFile(srcFile, dstFile)
		assert.Error(t, err)
	})
}

func TestWriteFile(t *testing.T) {
	t.Run("writes file with directory creation", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "nested", "file.txt")
		content := []byte("test content")

		err := WriteFile(filePath, content, 0644)
		assert.NoError(t, err)

		// Verify file was written
		readContent, err := os.ReadFile(filePath)
		assert.NoError(t, err)
		assert.Equal(t, content, readContent)

		// Verify directory was created
		info, err := os.Stat(filepath.Dir(filePath))
		assert.NoError(t, err)
		assert.True(t, info.IsDir())
	})
}

func TestExpandPath(t *testing.T) {
	t.Run("expands home directory", func(t *testing.T) {
		path := "~/test/path"
		expanded := ExpandPath(path)

		// Should not start with ~ anymore
		assert.NotEqual(t, path, expanded)
		assert.NotContains(t, expanded, "~")
	})

	t.Run("expands environment variables", func(t *testing.T) {
		os.Setenv("TEST_VAR", "test_value")
		defer os.Unsetenv("TEST_VAR")

		path := "$TEST_VAR/path"
		expanded := ExpandPath(path)

		assert.Contains(t, expanded, "test_value")
		assert.NotContains(t, expanded, "$TEST_VAR")
	})

	t.Run("handles regular paths unchanged", func(t *testing.T) {
		path := "/regular/path"
		expanded := ExpandPath(path)

		assert.Equal(t, path, expanded)
	})
}
