//go:build unit

package compose

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

func TestGenerator_addServiceVolumes_EmptyVolumes(t *testing.T) {
	manager := &services.Manager{}
	generator, err := NewGenerator("test-project", "/test/path", manager)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	// Test with service config with no volumes
	config := &types.ServiceConfig{
		Name: "test-service",
	}

	service := make(map[string]interface{})

	// Should handle empty volumes gracefully
	generator.addServiceVolumes(service, config)

	// Verify no volumes key is added when empty
	if _, exists := service["volumes"]; exists {
		t.Error("Expected no volumes key for service with no volumes")
	}
}

func TestGenerator_addServiceConfiguration_EmptyConfig(t *testing.T) {
	manager := &services.Manager{}
	generator, err := NewGenerator("test-project", "/test/path", manager)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	// Test with service config with no additional configuration
	config := &types.ServiceConfig{
		Name: "test-service",
	}

	service := make(map[string]interface{})

	// Should handle empty configuration gracefully
	generator.addServiceConfiguration(service, config)

	// Basic test - should not panic
	if service == nil {
		t.Error("Service should not be nil after configuration")
	}
}

func TestGenerator_addServiceHealthCheck_NoHealthCheck(t *testing.T) {
	manager := &services.Manager{}
	generator, err := NewGenerator("test-project", "/test/path", manager)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	// Test with service config with no health check
	config := &types.ServiceConfig{
		Name: "test-service",
	}

	service := make(map[string]interface{})

	// Should handle missing health check gracefully
	generator.addServiceHealthCheck(service, config)

	// Verify no healthcheck key is added when not configured
	if _, exists := service["healthcheck"]; exists {
		t.Error("Expected no healthcheck key for service with no health check")
	}
}
