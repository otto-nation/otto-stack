package docker

import (
	"context"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/logger"
)

func TestClient_GetComposeManager(t *testing.T) {
	client, err := NewClient(logger.GetLogger())
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	manager := client.GetComposeManager()
	if manager == nil {
		t.Error("Expected compose manager, got nil")
	}
}

func TestClient_ListResources(t *testing.T) {
	client, err := NewClient(logger.GetLogger())
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	resources, err := client.ListResources(ctx, ResourceContainer, "test-project")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resources == nil {
		t.Error("Expected resources list, got nil")
	}
}

func TestClient_RemoveResources(t *testing.T) {
	client, err := NewClient(logger.GetLogger())
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	err = client.RemoveResources(ctx, ResourceContainer, "test-project")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestClient_ListContainers(t *testing.T) {
	client, err := NewClient(logger.GetLogger())
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	containers, err := client.ListContainers(ctx, "test-project")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test passes if no error - empty slice is valid
	_ = containers
}

func TestClient_RemoveContainer(t *testing.T) {
	client, err := NewClient(logger.GetLogger())
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	err = client.RemoveContainer(ctx, "non-existent-container", false)
	// Expect error for non-existent container, but function should not panic
	if err == nil {
		t.Error("Expected error for non-existent container")
	}
}

func TestClient_GetServiceStatus(t *testing.T) {
	client, err := NewClient(logger.GetLogger())
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	services := []string{"web", "db"}
	statuses, err := client.GetServiceStatus(ctx, "test-project", services)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(statuses) != len(services) {
		t.Errorf("Expected %d statuses, got %d", len(services), len(statuses))
	}
}

func TestClient_GetDockerServiceStatus(t *testing.T) {
	client, err := NewClient(logger.GetLogger())
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	services := []string{"test-service"}
	statuses, err := client.GetDockerServiceStatus(ctx, "test-project", services)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Status should be valid even for non-existent service
	_ = statuses
}
