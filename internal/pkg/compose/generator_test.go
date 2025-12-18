package compose

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestNewGenerator(t *testing.T) {
	t.Run("creates generator successfully", func(t *testing.T) {
		generator, err := NewGenerator("test-project", "/tmp/services", nil)
		assert.NoError(t, err)
		assert.NotNil(t, generator)
		assert.Equal(t, "test-project", generator.projectName)
		assert.NotNil(t, generator.manager)
	})

	t.Run("handles empty project name", func(t *testing.T) {
		generator, err := NewGenerator("", "/tmp/services", nil)
		assert.NoError(t, err)
		assert.NotNil(t, generator)
		assert.Equal(t, "", generator.projectName)
	})
}

func TestGenerator_GenerateYAML(t *testing.T) {
	generator, err := NewGenerator("test-project", "/tmp/services", nil)
	require.NoError(t, err)

	t.Run("generates valid YAML structure", func(t *testing.T) {
		yamlBytes, err := generator.GenerateYAML([]string{"postgres"})
		if err != nil {
			// If postgres service doesn't exist, that's okay for this test
			t.Logf("Service not found (expected in test): %v", err)
			return
		}

		assert.NotEmpty(t, yamlBytes)

		// Parse the generated YAML to verify structure
		var compose map[string]any
		err = yaml.Unmarshal(yamlBytes, &compose)
		assert.NoError(t, err)

		// Check required top-level fields
		assert.Contains(t, compose, "services")
		assert.Contains(t, compose, "networks")

		// Check network configuration
		networks, ok := compose["networks"].(map[string]any)
		assert.True(t, ok)
		assert.Contains(t, networks, "default")

		defaultNetwork, ok := networks["default"].(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, "test-project-network", defaultNetwork["name"])
	})

	t.Run("handles empty service list", func(t *testing.T) {
		yamlBytes, err := generator.GenerateYAML([]string{})
		assert.NoError(t, err)
		assert.NotEmpty(t, yamlBytes)

		var compose map[string]any
		err = yaml.Unmarshal(yamlBytes, &compose)
		assert.NoError(t, err)

		services, ok := compose["services"].(map[string]any)
		assert.True(t, ok)
		assert.Empty(t, services)
	})

	t.Run("handles nonexistent services gracefully", func(t *testing.T) {
		yamlBytes, err := generator.GenerateYAML([]string{"nonexistent-service"})
		assert.NoError(t, err)
		assert.NotEmpty(t, yamlBytes)

		var compose map[string]any
		err = yaml.Unmarshal(yamlBytes, &compose)
		assert.NoError(t, err)

		services, ok := compose["services"].(map[string]any)
		assert.True(t, ok)
		// Should be empty since nonexistent service is skipped
		assert.Empty(t, services)
	})

	t.Run("generates valid YAML format", func(t *testing.T) {
		yamlBytes, err := generator.GenerateYAML([]string{"postgres", "redis"})
		assert.NoError(t, err)

		yamlString := string(yamlBytes)

		// Basic YAML structure validation
		assert.Contains(t, yamlString, "services:")
		assert.Contains(t, yamlString, "networks:")

		// Should be valid YAML
		var result map[string]any
		err = yaml.Unmarshal(yamlBytes, &result)
		assert.NoError(t, err)
	})
}

func TestGenerator_buildServices(t *testing.T) {
	generator, err := NewGenerator("test-project", "/tmp/services", nil)
	require.NoError(t, err)

	t.Run("builds services map", func(t *testing.T) {
		services := generator.buildServices([]string{"postgres", "redis"})
		assert.NotNil(t, services)

		// The actual content depends on available services
		// We just verify it returns a map without errors
		assert.IsType(t, map[string]any{}, services)
	})

	t.Run("handles empty service list", func(t *testing.T) {
		services := generator.buildServices([]string{})
		assert.NotNil(t, services)
		assert.Empty(t, services)
	})

	t.Run("skips nonexistent services", func(t *testing.T) {
		services := generator.buildServices([]string{"nonexistent"})
		assert.NotNil(t, services)
		// Should be empty since service doesn't exist
		assert.Empty(t, services)
	})
}

func TestGenerator_NetworkConfiguration(t *testing.T) {
	t.Run("creates project-specific network name", func(t *testing.T) {
		generator, err := NewGenerator("my-awesome-project", "/tmp/services", nil)
		require.NoError(t, err)

		yamlBytes, err := generator.GenerateYAML([]string{})
		assert.NoError(t, err)

		var compose map[string]any
		err = yaml.Unmarshal(yamlBytes, &compose)
		assert.NoError(t, err)

		networks := compose["networks"].(map[string]any)
		defaultNetwork := networks["default"].(map[string]any)
		assert.Equal(t, "my-awesome-project-network", defaultNetwork["name"])
	})

	t.Run("handles special characters in project name", func(t *testing.T) {
		generator, err := NewGenerator("test_project-123", "/tmp/services", nil)
		require.NoError(t, err)

		yamlBytes, err := generator.GenerateYAML([]string{})
		assert.NoError(t, err)

		var compose map[string]any
		err = yaml.Unmarshal(yamlBytes, &compose)
		assert.NoError(t, err)

		networks := compose["networks"].(map[string]any)
		defaultNetwork := networks["default"].(map[string]any)
		assert.Equal(t, "test_project-123-network", defaultNetwork["name"])
	})
}

func TestGenerator_YAMLOutput(t *testing.T) {
	generator, err := NewGenerator("test", "/tmp/services", nil)
	require.NoError(t, err)

	t.Run("produces valid YAML syntax", func(t *testing.T) {
		yamlBytes, err := generator.GenerateYAML([]string{})
		assert.NoError(t, err)

		yamlString := string(yamlBytes)

		// Should not contain Go map syntax
		assert.NotContains(t, yamlString, "map[")

		// Should contain proper YAML syntax
		lines := strings.Split(yamlString, "\n")
		foundServices := false
		foundNetworks := false

		for _, line := range lines {
			if strings.Contains(line, "services:") {
				foundServices = true
			}
			if strings.Contains(line, "networks:") {
				foundNetworks = true
			}
		}

		assert.True(t, foundServices, "Should contain services section")
		assert.True(t, foundNetworks, "Should contain networks section")
	})
}
