//go:build unit

package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStateManager(t *testing.T) {
	// Reset singleton for test
	stateManagerInstance = nil

	sm1 := NewStateManager()
	sm2 := NewStateManager()

	assert.NotNil(t, sm1)
	assert.Same(t, sm1, sm2, "Should return same instance (singleton)")
}

func TestStateManager_GetConfigHash(t *testing.T) {
	sm := NewStateManager()

	cfg := &config.Config{
		Project: config.ProjectConfig{Name: "test-project"},
		Stack:   config.StackConfig{Enabled: []string{"postgres", "redis"}},
	}

	hash1, err1 := sm.GetConfigHash(cfg)
	assert.NoError(t, err1)
	hash2, err2 := sm.GetConfigHash(cfg)
	assert.NoError(t, err2)

	assert.NotEmpty(t, hash1)
	assert.Equal(t, hash1, hash2, "Same config should produce same hash")

	// Different config should produce different hash
	cfg2 := &config.Config{
		Project: config.ProjectConfig{Name: "different-project"},
		Stack:   config.StackConfig{Enabled: []string{"postgres"}},
	}
	hash3, err3 := sm.GetConfigHash(cfg2)
	assert.NoError(t, err3)
	assert.NotEqual(t, hash1, hash3, "Different config should produce different hash")
}

func TestStateManager_SaveAndLoadState(t *testing.T) {
	tempDir := t.TempDir()
	ottoDir := filepath.Join(tempDir, core.OttoStackDir)
	require.NoError(t, os.MkdirAll(ottoDir, 0755))

	// Change to temp directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	sm := NewStateManager()

	// Test data
	serviceConfigs := []types.ServiceConfig{
		{Name: "postgres", Description: "PostgreSQL database"},
		{Name: "redis", Description: "Redis cache"},
	}
	configHash := "test-hash-123"

	state := &StackState{
		ServiceConfigs: serviceConfigs,
		ConfigHash:     configHash,
	}

	t.Run("save and load state successfully", func(t *testing.T) {
		err := sm.SaveState(state)
		assert.NoError(t, err)

		// Verify state file exists
		statePath := filepath.Join(ottoDir, "state.json")
		assert.FileExists(t, statePath)

		// Load and verify state
		loadedState, err := sm.LoadState()
		assert.NoError(t, err)
		assert.Equal(t, serviceConfigs, loadedState.ServiceConfigs)
		assert.Equal(t, configHash, loadedState.ConfigHash)
	})

	t.Run("load state when file missing", func(t *testing.T) {
		// Remove state file
		statePath := filepath.Join(ottoDir, "state.json")
		os.Remove(statePath)

		state, err := sm.LoadState()
		assert.NoError(t, err)
		assert.Empty(t, state.ServiceConfigs)
		assert.Empty(t, state.ConfigHash)
	})
}

func TestStateManager_StateOperations(t *testing.T) {
	tempDir := t.TempDir()
	ottoDir := filepath.Join(tempDir, core.OttoStackDir)
	require.NoError(t, os.MkdirAll(ottoDir, 0755))

	// Change to temp directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	sm := NewStateManager()

	cfg := &config.Config{
		Project: config.ProjectConfig{Name: "test-project"},
		Stack:   config.StackConfig{Enabled: []string{"postgres"}},
	}

	t.Run("can generate config hash", func(t *testing.T) {
		hash, err := sm.GetConfigHash(cfg)
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
	})

	t.Run("empty state when no file exists", func(t *testing.T) {
		state, err := sm.LoadState()
		assert.NoError(t, err)
		assert.NotNil(t, state)
		assert.Empty(t, state.ServiceConfigs)
		assert.Empty(t, state.ConfigHash)
	})
}
