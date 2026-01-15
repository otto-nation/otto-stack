package common

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// StackState represents the current state of the stack
type StackState struct {
	ServiceConfigs []types.ServiceConfig `json:"service_configs"`
	ConfigHash     string                `json:"config_hash"`
}

// StateManager handles stack state persistence
type StateManager struct{}

var stateManagerInstance *StateManager

// NewStateManager creates a new state manager (singleton)
func NewStateManager() *StateManager {
	if stateManagerInstance == nil {
		stateManagerInstance = &StateManager{}
	}
	return stateManagerInstance
}

// LoadState loads the current stack state
func (sm *StateManager) LoadState() (*StackState, error) {
	statePath := filepath.Join(core.OttoStackDir, "state.json")

	data, err := os.ReadFile(statePath)
	if os.IsNotExist(err) {
		return &StackState{}, nil // Return empty state if file doesn't exist
	}
	if err != nil {
		return nil, err
	}

	var state StackState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// SaveState saves the current stack state
func (sm *StateManager) SaveState(state *StackState) error {
	statePath := filepath.Join(core.OttoStackDir, "state.json")

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(statePath), core.PermReadWriteExec); err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(statePath, data, core.PermReadWrite)
}

// GetConfigHash generates a hash for the given configuration
func (sm *StateManager) GetConfigHash(cfg *config.Config) (string, error) {
	data, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}

	hash := md5.Sum(data)
	return fmt.Sprintf("%x", hash), nil
}
