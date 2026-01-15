package docker

import (
	"context"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/logger"
)

func TestResourceManager_List(t *testing.T) {
	client, err := NewClient(logger.GetLogger())
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	filter := NewProjectFilter("test-project")

	// Test listing containers
	containers, err := client.resources.List(ctx, ResourceContainer, filter)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	_ = containers

	// Test listing volumes
	volumes, err := client.resources.List(ctx, ResourceVolume, filter)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	_ = volumes

	// Test listing networks
	networks, err := client.resources.List(ctx, ResourceNetwork, filter)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	_ = networks

	// Test listing images
	images, err := client.resources.List(ctx, ResourceImage, filter)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	_ = images
}

func TestResourceManager_Remove(t *testing.T) {
	client, err := NewClient(logger.GetLogger())
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Test removing non-existent resources (should not error)
	err = client.resources.Remove(ctx, ResourceContainer, []string{"non-existent"})
	if err != nil {
		t.Fatalf("Expected no error for non-existent container, got %v", err)
	}

	err = client.resources.Remove(ctx, ResourceVolume, []string{"non-existent"})
	if err != nil {
		t.Fatalf("Expected no error for non-existent volume, got %v", err)
	}

	err = client.resources.Remove(ctx, ResourceNetwork, []string{"non-existent"})
	if err != nil {
		t.Fatalf("Expected no error for non-existent network, got %v", err)
	}

	err = client.resources.Remove(ctx, ResourceImage, []string{"non-existent"})
	if err != nil {
		t.Fatalf("Expected no error for non-existent image, got %v", err)
	}
}
