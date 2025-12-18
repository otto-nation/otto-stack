package compose

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLocalStackDependencyResolution(t *testing.T) {
	generator, err := NewGenerator("test-project", "", nil)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	tests := []struct {
		name             string
		inputServices    []string
		expectedServices []string
		description      string
	}{
		{
			name:             "localstack-s3 resolves dependencies",
			inputServices:    []string{"localstack-s3"},
			expectedServices: []string{"localstack"},
			description:      "localstack-s3 should resolve to localstack (init containers are auto-discovered)",
		},
		{
			name:             "multiple localstack services",
			inputServices:    []string{"localstack-s3", "localstack-sqs"},
			expectedServices: []string{"localstack"},
			description:      "Multiple LocalStack services should resolve to the same core service (init containers are auto-discovered)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			composeYAML, err := generator.GenerateYAML(tt.inputServices)
			if err != nil {
				t.Fatalf("Failed to generate YAML: %v", err)
			}

			// Parse the generated YAML
			var compose map[string]any
			if err := yaml.Unmarshal(composeYAML, &compose); err != nil {
				t.Fatalf("Failed to parse generated YAML: %v", err)
			}

			// Check services section
			services, ok := compose["services"].(map[string]any)
			if !ok {
				t.Fatal("Services section not found or invalid")
			}

			// Verify expected services are present
			for _, expectedService := range tt.expectedServices {
				if _, exists := services[expectedService]; !exists {
					t.Errorf("Expected service %s not found in generated compose", expectedService)
				}
			}

			// Verify localstack has proper configuration
			coreService, exists := services["localstack"]
			if !exists {
				return
			}

			coreMap := coreService.(map[string]any)

			// Check image
			image, ok := coreMap["image"].(string)
			if !ok || image != "localstack/localstack:latest" {
				t.Errorf("localstack should have image 'localstack/localstack:latest', got: %v", image)
			}

			// Check environment variables
			env, ok := coreMap["environment"].(map[string]any)
			if !ok {
				t.Error("localstack should have environment variables")
				return
			}

			expectedEnvVars := []string{"AWS_ACCESS_KEY_ID", "AWS_DEFAULT_REGION", "LOCALSTACK_HOST"}
			for _, envVar := range expectedEnvVars {
				if _, exists := env[envVar]; !exists {
					t.Errorf("Expected environment variable %s not found", envVar)
				}
			}

			// Check ports
			ports, ok := coreMap["ports"].([]any)
			if !ok {
				t.Error("localstack should have ports configured")
				return
			}

			portFound := false
			for _, port := range ports {
				if portStr, ok := port.(string); ok && strings.Contains(portStr, "4566:4566") {
					portFound = true
					break
				}
			}
			if !portFound {
				t.Error("localstack should expose port 4566")
			}

			// Note: localstack-init is now auto-discovered based on config files
			// It will only be present if service-configs/localstack-*.yml files exist
			// This test focuses on the core service dependency resolution

			t.Logf("✓ %s: %s", tt.name, tt.description)
		})
	}
}

func TestLocalStackEnvironmentMerging(t *testing.T) {
	generator, err := NewGenerator("test-project", "", nil)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	composeYAML, err := generator.GenerateYAML([]string{"localstack-s3"})
	if err != nil {
		t.Fatalf("Failed to generate YAML: %v", err)
	}

	// Parse the generated YAML
	var compose map[string]any
	if err := yaml.Unmarshal(composeYAML, &compose); err != nil {
		t.Fatalf("Failed to parse generated YAML: %v", err)
	}

	services := compose["services"].(map[string]any)

	// Test that localstack has environment variables
	// Note: localstack-init is auto-discovered and only present when config files exist
	for _, serviceName := range []string{"localstack"} {
		if service, exists := services[serviceName]; exists {
			serviceMap := service.(map[string]any)
			if env, ok := serviceMap["environment"].(map[string]any); ok {
				if len(env) == 0 {
					t.Errorf("%s should have environment variables", serviceName)
				}
				t.Logf("✓ %s has %d environment variables", serviceName, len(env))
			} else {
				t.Errorf("%s should have environment section", serviceName)
			}
		}
	}
}
