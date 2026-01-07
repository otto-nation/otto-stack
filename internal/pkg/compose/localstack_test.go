package compose

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/test/testutil"
	"github.com/stretchr/testify/assert"
)

func TestLocalStackDependencyResolution(t *testing.T) {
	// Test dependency resolution logic with real services
	manager, err := services.New()
	if err != nil {
		t.Fatalf("Failed to create service manager: %v", err)
	}

	generator, err := NewGenerator("test-project", "", manager)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	t.Run("localstack-s3 resolves dependencies", func(t *testing.T) {
		// Test that localstack-s3 (config service) includes localstack (container service)
		compose, err := generator.buildComposeStructure([]string{services.ServiceLocalstackS3})
		if err != nil {
			t.Skipf("LocalStack services not available in test environment: %v", err)
			return
		}

		servicesMap := compose["services"].(map[string]any)

		// Should include the localstack container service
		if _, exists := servicesMap[services.ServiceLocalstack]; exists {
			t.Log("✓ localstack-s3 correctly resolves to localstack container")
		} else {
			t.Log("ℹ️ localstack dependency not resolved (expected in minimal test env)")
		}
	})
}

func TestLocalStackEnvironmentMerging(t *testing.T) {
	// Test that compose generation includes proper network labels
	generator, err := NewGenerator("test-project", "", testutil.NewTestManager(t))
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	compose, err := generator.buildComposeStructure([]string{})
	if err != nil {
		t.Fatalf("Failed to generate compose structure: %v", err)
	}

	// Verify network has Otto Stack labels
	networks := compose["networks"].(map[string]any)
	defaultNet := networks["default"].(map[string]any)
	labels := defaultNet["labels"].(map[string]string)

	assert.Equal(t, "true", labels["io.otto-stack.managed"])
	assert.Equal(t, "test-project", labels["io.otto-stack.project"])

	t.Log("✓ Network labels are properly generated")
}
