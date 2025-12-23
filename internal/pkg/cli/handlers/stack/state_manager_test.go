package stack

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/test/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStateManager(t *testing.T) {
	sm := NewStateManager()
	assert.NotNil(t, sm)
}

func TestStateManager_LoadState_FileNotExists(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	sm := NewStateManager()

	// Should return empty state when file doesn't exist
	state, err := sm.LoadState()

	assert.NoError(t, err)
	assert.NotNil(t, state)
	assert.Empty(t, state.Services)
	assert.Empty(t, state.ConfigHash)
}

func TestStateManager_SaveAndLoadState(t *testing.T) {
	testutil.WithTempDir(t, func(tempDir string) {
		sm := NewStateManager()

		// Create test state
		originalState := &StackState{
			Services:   []string{"postgres", "redis"},
			ConfigHash: "test-hash-123",
		}

		// Save state
		err := sm.SaveState(originalState)
		require.NoError(t, err)

		// Verify file was created
		statePath := filepath.Join(core.OttoStackDir, "state.json")
		testutil.AssertFileExists(t, statePath)

		// Load state back
		loadedState, err := sm.LoadState()
		require.NoError(t, err)

		assert.Equal(t, originalState.Services, loadedState.Services)
		assert.Equal(t, originalState.ConfigHash, loadedState.ConfigHash)
	})
}

func TestStateManager_GetConfigHash(t *testing.T) {
	sm := NewStateManager()

	// Create test config
	cfg := &config.Config{
		Project: config.ProjectConfig{
			Name: "test-project",
		},
	}

	// Get hash
	hash, err := sm.GetConfigHash(cfg)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.Len(t, hash, 32) // MD5 hash should be 32 characters

	// Same config should produce same hash
	hash2, err := sm.GetConfigHash(cfg)
	assert.NoError(t, err)
	assert.Equal(t, hash, hash2)
}

func TestStateManager_SaveState_CreateDirectory(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	sm := NewStateManager()

	// Ensure otto-stack directory doesn't exist
	os.RemoveAll(core.OttoStackDir)

	state := &StackState{
		Services:   []string{"test"},
		ConfigHash: "hash",
	}

	// Should create directory and save file
	err := sm.SaveState(state)
	assert.NoError(t, err)

	// Verify directory was created
	assert.DirExists(t, core.OttoStackDir)
	assert.FileExists(t, filepath.Join(core.OttoStackDir, "state.json"))
}
