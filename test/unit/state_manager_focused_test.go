//go:build unit

package unit

import (
	"encoding/json"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/project"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/stack"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

func TestStateManager_NewStateManager_Focused(t *testing.T) {
	sm := stack.NewStateManager()
	if sm == nil {
		t.Fatal("NewStateManager returned nil")
	}
}

func TestStateManager_GetConfigHash_Focused(t *testing.T) {
	sm := stack.NewStateManager()

	// Test with simple config
	cfg := &config.Config{
		Project: config.ProjectConfig{
			Name:     project.TestProjectName,
			Services: []string{services.ServicePostgres, services.ServiceRedis},
		},
	}

	hash1, err := sm.GetConfigHash(cfg)
	if err != nil {
		t.Fatalf("GetConfigHash failed: %v", err)
	}

	if hash1 == "" {
		t.Error("GetConfigHash returned empty hash")
	}

	// Test hash consistency
	hash2, err := sm.GetConfigHash(cfg)
	if err != nil {
		t.Fatalf("GetConfigHash failed on second call: %v", err)
	}

	if hash1 != hash2 {
		t.Errorf("Hash inconsistency: %q != %q", hash1, hash2)
	}

	// Test hash changes with different config
	cfg2 := &config.Config{
		Project: config.ProjectConfig{
			Name:     "different-project",
			Services: []string{services.ServiceMysql, services.ServiceRedis},
		},
	}

	hash3, err := sm.GetConfigHash(cfg2)
	if err != nil {
		t.Fatalf("GetConfigHash failed for different config: %v", err)
	}

	if hash1 == hash3 {
		t.Error("Different configs should produce different hashes")
	}
}

func TestStateManager_GetConfigHash_EmptyConfig_Focused(t *testing.T) {
	sm := stack.NewStateManager()

	cfg := &config.Config{}
	hash, err := sm.GetConfigHash(cfg)

	if err != nil {
		t.Fatalf("GetConfigHash should handle empty config: %v", err)
	}

	if hash == "" {
		t.Error("GetConfigHash should return non-empty hash even for empty config")
	}
}

func TestStackState_JSONSerialization_Focused(t *testing.T) {
	state := &stack.StackState{
		Services:   []string{services.ServicePostgres, services.ServiceRedis},
		ConfigHash: "abc123",
	}

	// Test marshaling
	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	// Test unmarshaling
	var unmarshaled stack.StackState
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	// Verify data integrity
	if len(unmarshaled.Services) != len(state.Services) {
		t.Error("Services length mismatch after JSON round-trip")
	}

	if unmarshaled.ConfigHash != state.ConfigHash {
		t.Error("ConfigHash mismatch after JSON round-trip")
	}
}

func TestStackState_EmptyState_Focused(t *testing.T) {
	// Test empty state behavior
	state := &stack.StackState{}

	// Should be able to marshal empty state
	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("JSON marshal of empty state failed: %v", err)
	}

	// Should be able to unmarshal empty state
	var unmarshaled stack.StackState
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("JSON unmarshal of empty state failed: %v", err)
	}

	// Verify empty state properties
	if len(unmarshaled.Services) != 0 {
		t.Error("Empty state should have no services")
	}

	if unmarshaled.ConfigHash != "" {
		t.Error("Empty state should have empty config hash")
	}
}
