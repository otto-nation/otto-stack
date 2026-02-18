//go:build unit

package context

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDetector(t *testing.T) {
	detector, err := NewDetector()
	require.NoError(t, err)
	assert.NotNil(t, detector)
	assert.NotEmpty(t, detector.homeDir)
}

func TestDetector_DetectContext(t *testing.T) {
	detector, err := NewDetector()
	require.NoError(t, err)

	mode, err := detector.DetectContext()
	require.NoError(t, err)
	assert.NotNil(t, mode)
}

func TestDetector_Detect(t *testing.T) {
	detector, err := NewDetector()
	require.NoError(t, err)

	ctx, err := detector.Detect()
	require.NoError(t, err)
	assert.NotNil(t, ctx)
	assert.NotNil(t, ctx.SharedContainers)
}

func TestDetector_FindProjectRoot(t *testing.T) {
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	t.Run("finds project root when exists", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "project")
		configDir := filepath.Join(projectDir, core.OttoStackDir)
		require.NoError(t, os.MkdirAll(configDir, 0755))

		configFile := filepath.Join(configDir, core.ConfigFileName)
		require.NoError(t, os.WriteFile(configFile, []byte("{}"), 0644))

		os.Chdir(projectDir)

		detector, _ := NewDetector()
		project, err := detector.findProjectRoot()
		require.NoError(t, err)
		assert.NotNil(t, project)
	})

	t.Run("returns nil when no project found", func(t *testing.T) {
		emptyDir := filepath.Join(tempDir, "empty")
		require.NoError(t, os.MkdirAll(emptyDir, 0755))
		os.Chdir(emptyDir)

		detector, _ := NewDetector()
		project, err := detector.findProjectRoot()
		require.NoError(t, err)
		assert.Nil(t, project)
	})
}

func TestExecutionMode_SharedMode(t *testing.T) {
	mode := &SharedMode{Shared: &SharedInfo{Root: "/test"}}
	assert.Equal(t, "/test", mode.SharedRoot())
}

func TestExecutionMode_ProjectMode(t *testing.T) {
	mode := &ProjectMode{
		Project: &ProjectInfo{Root: "/project"},
		Shared:  &SharedInfo{Root: "/shared"},
	}
	assert.Equal(t, "/shared", mode.SharedRoot())
}
