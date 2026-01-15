package docker

import (
	"context"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/logger"
)

func TestClient_ListOttoContainers(t *testing.T) {
	client, err := NewClient(logger.GetLogger())
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	containers, err := client.ListOttoContainers(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return empty list in test environment
	_ = containers
}

func TestClient_ListProjectContainers(t *testing.T) {
	client, err := NewClient(logger.GetLogger())
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	containers, err := client.ListProjectContainers(ctx, "test-project")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return empty list in test environment
	_ = containers
}

func TestClient_GetContainerLabels(t *testing.T) {
	client, err := NewClient(logger.GetLogger())
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	labels, err := client.GetContainerLabels(ctx, "non-existent-container")
	// Expect error for non-existent container
	if err == nil {
		t.Error("Expected error for non-existent container")
	}
	_ = labels
}

func TestNewDefaultProjectLoader(t *testing.T) {
	loader, err := NewDefaultProjectLoader()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if loader == nil {
		t.Error("Expected project loader, got nil")
	}
}

func TestProjectLoader_Load(t *testing.T) {
	loader, err := NewDefaultProjectLoader()
	if err != nil {
		t.Fatalf("Expected no error creating loader, got %v", err)
	}

	project, err := loader.Load("test-project")
	// Expect error for non-existent compose file
	if err == nil {
		t.Error("Expected error for non-existent compose file")
	}
	_ = project
}
