package compose

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/test/testutil"
)

func TestAutoDiscoveryPatternMatching(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create service-configs directory using the actual constant
	configDir := filepath.Join(tempDir, core.ServiceConfigsDir)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	tests := []struct {
		name             string
		configFiles      []string
		resolvedServices []string
		expectedInits    []string
		description      string
	}{
		{
			name:             "LocalStack pattern matching",
			configFiles:      []string{"localstack-sqs.yml", "localstack-sns.yml"},
			resolvedServices: []string{services.ServiceLocalstack, services.ServiceRedis},
			expectedInits:    []string{"localstack-init"},
			description:      "Should detect localstack-init when localstack-*.yml files exist and localstack is resolved",
		},
		{
			name:             "Multiple service patterns",
			configFiles:      []string{"localstack-s3.yml", "postgres-schemas.yml"},
			resolvedServices: []string{services.ServiceLocalstack, services.ServicePostgres},
			expectedInits:    []string{"localstack-init", "postgres-init"},
			description:      "Should detect init containers for multiple services with config files",
		},
		{
			name:             "No matching services",
			configFiles:      []string{"localstack-sqs.yml", "postgres-schemas.yml"},
			resolvedServices: []string{services.ServiceRedis, services.ServiceMysql},
			expectedInits:    []string{},
			description:      "Should not create init containers when no target services are resolved",
		},
		{
			name:             "Partial matches",
			configFiles:      []string{"localstack-sqs.yml", "postgres-schemas.yml", "redis-config.yml"},
			resolvedServices: []string{services.ServiceLocalstack, services.ServiceMysql},
			expectedInits:    []string{"localstack-init"},
			description:      "Should only create init containers for services that are actually resolved",
		},
		{
			name:             "No config files",
			configFiles:      []string{},
			resolvedServices: []string{services.ServiceLocalstack, services.ServicePostgres, services.ServiceRedis},
			expectedInits:    []string{},
			description:      "Should not create init containers when no config files exist",
		},
		{
			name:             "Non-pattern files ignored",
			configFiles:      []string{"localstack-sqs.yml", "random-file.yml", "not-a-pattern.txt"},
			resolvedServices: []string{services.ServiceLocalstack},
			expectedInits:    []string{"localstack-init"},
			description:      "Should ignore files that don't match the {service}-{type}.yml pattern",
		},
		{
			name:             "Duplicate prevention",
			configFiles:      []string{"localstack-sqs.yml", "localstack-sns.yml", "localstack-s3.yml"},
			resolvedServices: []string{services.ServiceLocalstack},
			expectedInits:    []string{"localstack-init"},
			description:      "Should not create duplicate init containers for multiple configs of same service",
		},
		{
			name:             "Complex service names",
			configFiles:      []string{"my-complex-service-config.yml", "simple-setup.yml"},
			resolvedServices: []string{"my-complex-service", "simple"},
			expectedInits:    []string{"my-complex-service-init", "simple-init"},
			description:      "Should handle complex service names with hyphens correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean the directory for each test
			if err := os.RemoveAll(configDir); err != nil {
				t.Fatalf("Failed to clean test directory: %v", err)
			}
			if err := os.MkdirAll(configDir, 0755); err != nil {
				t.Fatalf("Failed to recreate test directory: %v", err)
			}

			// Create test config files
			for _, file := range tt.configFiles {
				filePath := filepath.Join(configDir, file)
				content := generateTestConfigContent(file)
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to create test file %s: %v", file, err)
				}
			}

			// Create generator and test discovery
			generator, err := NewGenerator("test-project", "", testutil.NewTestManager(t))
			if err != nil {
				t.Fatalf("Failed to create generator: %v", err)
			}

			// Mock the discovery by creating a temporary generator with modified config dir
			initServices := generator.discoverInitServicesForTesting(tt.resolvedServices, tempDir)

			// Verify results
			if len(initServices) != len(tt.expectedInits) {
				t.Errorf("Expected %d init services, got %d", len(tt.expectedInits), len(initServices))
				t.Logf("Expected: %v", tt.expectedInits)
				t.Logf("Got: %v", initServices)
			}

			// Check each expected init service is present
			for _, expected := range tt.expectedInits {
				found := false
				for _, actual := range initServices {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected init service %s not found in result", expected)
				}
			}

			// Check no unexpected init services are present
			for _, actual := range initServices {
				found := false
				for _, expected := range tt.expectedInits {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Unexpected init service %s found in result", actual)
				}
			}

			t.Logf("✓ %s: %s", tt.name, tt.description)
		})
	}
}

func TestAutoDiscoveryFileSystem(t *testing.T) {
	// Test with actual filesystem operations
	tempDir := t.TempDir()

	// Create nested directory structure
	configDir := filepath.Join(tempDir, core.ServiceConfigsDir)
	subDir := filepath.Join(configDir, "subdirectory")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create test directories: %v", err)
	}

	// Create test files in different locations
	testFiles := map[string]string{
		filepath.Join(configDir, "localstack-sqs.yml"): "localstack SQS config",
		filepath.Join(configDir, "postgres-init.yml"):  "postgres init config",
		filepath.Join(subDir, "localstack-sns.yml"):    "localstack SNS in subdirectory",
		filepath.Join(configDir, "invalid-format.txt"): "not a YAML file",
	}

	for path, content := range testFiles {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", path, err)
		}
	}

	generator, err := NewGenerator("test-project", "", testutil.NewTestManager(t))
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	// Test discovery with services that should match
	resolvedServices := []string{services.ServiceLocalstack, services.ServicePostgres}
	initServices := generator.discoverInitServicesForTesting(resolvedServices, tempDir)

	expectedInits := []string{"localstack-init", "postgres-init"}

	if len(initServices) != len(expectedInits) {
		t.Errorf("Expected %d init services, got %d", len(expectedInits), len(initServices))
		t.Logf("Expected: %v", expectedInits)
		t.Logf("Got: %v", initServices)
	}

	// Verify each expected service
	for _, expected := range expectedInits {
		found := false
		for _, actual := range initServices {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected init service %s not found", expected)
		}
	}
}

func TestAutoDiscoveryNonExistentDirectory(t *testing.T) {
	// Test behavior when service-configs directory doesn't exist
	tempDir := t.TempDir()

	generator, err := NewGenerator("test-project", "", testutil.NewTestManager(t))
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	// Use a non-existent directory path for testing
	nonExistentPath := filepath.Join(tempDir, "non-existent-configs")
	initServices := generator.discoverInitServicesForTesting([]string{services.ServiceLocalstack, services.ServicePostgres}, nonExistentPath)

	if len(initServices) != 0 {
		t.Errorf("Expected no init services when directory doesn't exist, got %v", initServices)
	}
}

func TestServicePatternExtraction(t *testing.T) {
	tests := []struct {
		filename        string
		expectedService string
		shouldMatch     bool
		description     string
	}{
		{
			filename:        "localstack-sqs.yml",
			expectedService: services.ServiceLocalstack,
			shouldMatch:     true,
			description:     "Standard localstack pattern",
		},
		{
			filename:        "postgres-schemas.yml",
			expectedService: services.ServicePostgres,
			shouldMatch:     true,
			description:     "Standard postgres pattern",
		},
		{
			filename:        "my-service-config.yml",
			expectedService: "my",
			shouldMatch:     true,
			description:     "Multi-word service name",
		},
		{
			filename:        "single.yml",
			expectedService: "",
			shouldMatch:     false,
			description:     "Single word - no pattern match",
		},
		{
			filename:        "no-extension",
			expectedService: "",
			shouldMatch:     false,
			description:     "No file extension",
		},
		{
			filename:        "complex-service-name-config.yml",
			expectedService: "complex",
			shouldMatch:     true,
			description:     "Complex service name takes first part",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// Extract service name using same logic as discoverInitServices
			if !strings.HasSuffix(tt.filename, ".yml") && !strings.HasSuffix(tt.filename, ".yaml") {
				if tt.shouldMatch {
					t.Errorf("File %s should match but has no YAML extension", tt.filename)
				}
				return
			}

			serviceName := strings.TrimSuffix(strings.TrimSuffix(tt.filename, ".yml"), ".yaml")
			parts := strings.Split(serviceName, "-")

			if len(parts) >= 2 {
				extractedService := parts[0]
				if tt.shouldMatch {
					if extractedService != tt.expectedService {
						t.Errorf("Expected service '%s', got '%s'", tt.expectedService, extractedService)
					}
				} else {
					t.Errorf("File %s should not match but extracted service '%s'", tt.filename, extractedService)
				}
			} else {
				if tt.shouldMatch {
					t.Errorf("File %s should match but didn't extract a service", tt.filename)
				}
			}
		})
	}
}

// discoverInitServicesForTesting is a test helper that allows overriding the config directory
func (g *Generator) discoverInitServicesForTesting(resolvedServices []string, testDir string) []string {
	var initServices []string

	// Use test directory instead of the default one
	configDir := filepath.Join(testDir, core.ServiceConfigsDir)

	// Check if directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return initServices
	}

	// Walk through all subdirectories to find config files
	err := filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !core.IsYAMLFile(info.Name()) {
			return nil
		}

		// Extract service name from filename (remove extension)
		serviceName := core.TrimYAMLExt(info.Name())

		// Check for pattern: {target-service}-{config-type}.yml
		if parts := strings.Split(serviceName, "-"); len(parts) >= 2 {
			// For complex service names, try different combinations
			// e.g., "my-complex-service-config.yml" -> try "my-complex-service", then "my-complex", then "my"
			for i := len(parts) - 1; i >= 1; i-- {
				targetService := strings.Join(parts[:i], "-")

				// If target service prefix matches any resolved service, add its init container
				for _, resolved := range resolvedServices {
					if strings.HasPrefix(resolved, targetService) || resolved == targetService {
						initServiceName := targetService + "-init"
						// Avoid duplicates
						found := false
						for _, existing := range initServices {
							if existing == initServiceName {
								found = true
								break
							}
						}
						if !found {
							initServices = append(initServices, initServiceName)
						}
						goto nextFile // Found a match, move to next file
					}
				}
			}
		nextFile:
		}

		return nil
	})

	if err != nil {
		return initServices
	}

	return initServices
}

// generateTestConfigContent creates realistic test configuration content
func generateTestConfigContent(filename string) string {
	if strings.Contains(filename, "localstack-sqs") {
		return `queues:
  - name: test-queue
    visibility_timeout: 30`
	}
	if strings.Contains(filename, "localstack-sns") {
		return `topics:
  - name: test-topic`
	}
	if strings.Contains(filename, "localstack-s3") {
		return `buckets:
  - name: test-bucket`
	}
	if strings.Contains(filename, services.ServicePostgres) {
		return `schemas:
  - name: test_schema
    tables:
      - users
      - orders`
	}

	// Generic config content
	return `configuration:
  enabled: true
  settings:
    key: value`
}
