package docker

import (
	"context"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/logger"
)

func TestClient_MapToEnvSlice(t *testing.T) {
	// Test the mapToEnvSlice utility function indirectly through RunInitContainer
	client, err := NewClient(logger.GetLogger())
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	// This will exercise mapToEnvSlice internally
	config := InitContainerConfig{
		Image:   "alpine:latest",
		Command: []string{"echo", "test"},
		Environment: map[string]string{
			"TEST_VAR": "test_value",
			"ANOTHER":  "value",
		},
	}

	// This should not panic and should handle the environment mapping
	err = client.RunInitContainer(ctx, "test-init", config)
	// We expect this to fail in test environment, but it should not panic
	_ = err
}

func TestClient_GetHealthStatus(t *testing.T) {
	client, err := NewClient(logger.GetLogger())
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	// Test getHealthStatus indirectly through GetServiceStatus
	ctx := context.Background()
	services := []string{"test-service"}
	statuses, err := client.GetServiceStatus(ctx, "test-project", services)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// This exercises getHealthStatus internally
	if len(statuses) != 1 {
		t.Errorf("Expected 1 status, got %d", len(statuses))
	}
}

func TestClient_Contains(t *testing.T) {
	client, err := NewClient(logger.GetLogger())
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	// Test contains utility indirectly through GetServiceStatus
	ctx := context.Background()
	services := []string{"web", "db", "cache"}
	statuses, err := client.GetServiceStatus(ctx, "test-project", services)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// This exercises the contains function internally
	if len(statuses) != len(services) {
		t.Errorf("Expected %d statuses, got %d", len(services), len(statuses))
	}
}
