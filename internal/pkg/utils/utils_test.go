package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEnsureDir(t *testing.T) {
	tmpDir := t.TempDir()
	testDir := filepath.Join(tmpDir, "test", "nested", "dir")

	err := EnsureDir(testDir)
	assert.NoError(t, err)

	info, err := os.Stat(testDir)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "source.txt")
	dstFile := filepath.Join(tmpDir, "nested", "dest.txt")

	// Create source file
	content := []byte("test content")
	err := os.WriteFile(srcFile, content, 0644)
	assert.NoError(t, err)

	// Copy file
	err = CopyFile(srcFile, dstFile)
	assert.NoError(t, err)

	// Verify copy
	copiedContent, err := os.ReadFile(dstFile)
	assert.NoError(t, err)
	assert.Equal(t, content, copiedContent)
}

func TestWriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "nested", "test.txt")
	content := []byte("test content")

	err := WriteFile(testFile, content, 0644)
	assert.NoError(t, err)

	readContent, err := os.ReadFile(testFile)
	assert.NoError(t, err)
	assert.Equal(t, content, readContent)
}

func TestExpandPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{"home expansion", "~/test", "test"},
		{"env expansion", "$HOME/test", "test"},
		{"no expansion", "/absolute/path", "/absolute/path"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpandPath(tt.input)
			assert.Contains(t, result, tt.contains)
		})
	}
}

func TestRunCommand(t *testing.T) {
	var cmd, expectedOutput string

	if runtime.GOOS == "windows" {
		cmd = "echo"
		expectedOutput = "test"
	} else {
		cmd = "echo"
		expectedOutput = "test"
	}

	output, err := RunCommand(cmd, "test")
	assert.NoError(t, err)
	assert.Contains(t, output, expectedOutput)
}

func TestIsCommandAvailable(t *testing.T) {
	// Test with a command that should exist
	var existingCmd string
	if runtime.GOOS == "windows" {
		existingCmd = "cmd"
	} else {
		existingCmd = "sh"
	}

	assert.True(t, IsCommandAvailable(existingCmd))
	assert.False(t, IsCommandAvailable("nonexistent_command_12345"))
}

func TestRetry(t *testing.T) {
	t.Run("success on first attempt", func(t *testing.T) {
		attempts := 0
		err := Retry(3, time.Millisecond, func() error {
			attempts++
			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, attempts)
	})

	t.Run("success on second attempt", func(t *testing.T) {
		attempts := 0
		err := Retry(3, time.Millisecond, func() error {
			attempts++
			if attempts < 2 {
				return assert.AnError
			}
			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 2, attempts)
	})

	t.Run("failure after all attempts", func(t *testing.T) {
		attempts := 0
		err := Retry(2, time.Millisecond, func() error {
			attempts++
			return assert.AnError
		})
		assert.Error(t, err)
		assert.Equal(t, 2, attempts)
		assert.Contains(t, err.Error(), "failed after 2 attempts")
	})
}

func TestTimeout(t *testing.T) {
	t.Run("completes within timeout", func(t *testing.T) {
		err := Timeout(100*time.Millisecond, func() error {
			time.Sleep(10 * time.Millisecond)
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("times out", func(t *testing.T) {
		err := Timeout(10*time.Millisecond, func() error {
			time.Sleep(50 * time.Millisecond)
			return nil
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "operation timed out")
	})
}

func TestIsPortInUse(t *testing.T) {
	// Test with a port that's likely not in use
	assert.False(t, IsPortInUse(65432))
}

func TestGetFreePort(t *testing.T) {
	port, err := GetFreePort(50000)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, port, 50000)
	assert.Less(t, port, 50100)
}
