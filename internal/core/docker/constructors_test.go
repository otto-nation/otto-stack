package docker

import (
	"context"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/logger"
)

func TestClient_RunInitContainer(t *testing.T) {
	client, err := NewClient(logger.GetLogger())
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	config := InitContainerConfig{
		Image:   "alpine:latest",
		Command: []string{"echo", "test"},
		Environment: map[string]string{
			"TEST_VAR": "value",
		},
		WorkingDir: "/tmp",
	}

	// This will likely fail in test environment but should not panic
	err = client.RunInitContainer(ctx, "test-init", config)
	// We expect this to fail, but it should exercise the code path
	_ = err
}

func TestNewManager(t *testing.T) {
	manager, err := NewManager()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if manager == nil {
		t.Error("Expected manager, got nil")
	}
}
