package project

import (
	"os"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/stretchr/testify/require"
)

// setupTestDir creates a temporary directory and changes to it
// deadcode: test helper used across multiple test files
func setupTestDir(t *testing.T) (cleanup func()) {
	tempDir, err := os.MkdirTemp("", TestTempDirPattern)
	require.NoError(t, err)

	originalDir, err := os.Getwd()
	require.NoError(t, err)

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	return func() {
		os.Chdir(originalDir)
		os.RemoveAll(tempDir)
	}
}

// createTestFile creates a file with given content
// deadcode: test helper used across multiple test files
func createTestFile(t *testing.T, path, content string) {
	err := os.WriteFile(path, []byte(content), core.PermReadWrite)
	require.NoError(t, err)
}

// createTestConfig creates a config file
// deadcode: test helper used across multiple test files
func createTestConfig(t *testing.T) {
	err := os.MkdirAll(core.OttoStackDir, core.PermReadWriteExec)
	require.NoError(t, err)
	createTestFile(t, TestConfigFilePath, TestConfigContent)
}
