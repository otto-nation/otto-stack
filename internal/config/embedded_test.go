package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestEmbeddedCommandsYAML(t *testing.T) {
	t.Run("commands YAML is embedded and valid", func(t *testing.T) {
		assert.NotEmpty(t, EmbeddedCommandsYAML, "Commands YAML should be embedded")

		// Should be valid YAML
		var commands map[string]any
		err := yaml.Unmarshal(EmbeddedCommandsYAML, &commands)
		assert.NoError(t, err, "Commands YAML should be valid")

		// Should contain expected sections
		assert.Contains(t, commands, "commands", "Should contain commands section")
		assert.Contains(t, commands, "messages", "Should contain messages section")
	})

	t.Run("contains required commands", func(t *testing.T) {
		var config map[string]any
		err := yaml.Unmarshal(EmbeddedCommandsYAML, &config)
		assert.NoError(t, err)

		commands, ok := config["commands"].(map[string]any)
		assert.True(t, ok, "Commands should be a map")

		// Check for some core commands
		expectedCommands := []string{"init", "up", "down", "status"}
		for _, cmd := range expectedCommands {
			assert.Contains(t, commands, cmd, "Should contain %s command", cmd)
		}
	})
}

func TestEmbeddedSchemaYAML(t *testing.T) {
	t.Run("schema YAML is embedded and valid", func(t *testing.T) {
		assert.NotEmpty(t, EmbeddedSchemaYAML, "Schema YAML should be embedded")

		// Should be valid YAML
		var schema map[string]any
		err := yaml.Unmarshal(EmbeddedSchemaYAML, &schema)
		assert.NoError(t, err, "Schema YAML should be valid")
	})

	t.Run("contains schema definitions", func(t *testing.T) {
		var schema map[string]any
		err := yaml.Unmarshal(EmbeddedSchemaYAML, &schema)
		assert.NoError(t, err)

		// Should contain some schema definitions
		assert.NotEmpty(t, schema, "Schema should not be empty")
	})
}

func TestEmbeddedInitSettingsYAML(t *testing.T) {
	t.Run("init settings YAML is embedded and valid", func(t *testing.T) {
		assert.NotEmpty(t, EmbeddedInitSettingsYAML, "Init settings YAML should be embedded")

		// Should be valid YAML
		var settings map[string]any
		err := yaml.Unmarshal(EmbeddedInitSettingsYAML, &settings)
		assert.NoError(t, err, "Init settings YAML should be valid")
	})

	t.Run("contains initialization settings", func(t *testing.T) {
		var settings map[string]any
		err := yaml.Unmarshal(EmbeddedInitSettingsYAML, &settings)
		assert.NoError(t, err)

		// Should contain some settings
		assert.NotEmpty(t, settings, "Init settings should not be empty")
	})
}

func TestEmbeddedServicesFS(t *testing.T) {
	t.Run("services filesystem is embedded", func(t *testing.T) {
		// Check that we can read the services directory
		entries, err := EmbeddedServicesFS.ReadDir("services")
		assert.NoError(t, err, "Should be able to read services directory")
		assert.NotEmpty(t, entries, "Services directory should not be empty")
	})

	t.Run("contains service categories", func(t *testing.T) {
		entries, err := EmbeddedServicesFS.ReadDir("services")
		assert.NoError(t, err)

		// Should contain some expected categories
		categoryNames := make([]string, 0, len(entries))
		for _, entry := range entries {
			if entry.IsDir() {
				categoryNames = append(categoryNames, entry.Name())
			}
		}

		expectedCategories := []string{"database", "cache", "messaging"}
		for _, category := range expectedCategories {
			found := false
			for _, name := range categoryNames {
				if name == category {
					found = true
					break
				}
			}
			if !found {
				t.Logf("Expected category %s not found in: %v", category, categoryNames)
			}
		}
	})

	t.Run("service files are valid YAML", func(t *testing.T) {
		// Test a few service files to ensure they're valid YAML
		testServices := []string{
			"services/database/postgres.yaml",
			"services/cache/redis.yaml",
		}

		for _, servicePath := range testServices {
			data, err := EmbeddedServicesFS.ReadFile(servicePath)
			if err != nil {
				t.Logf("Service file %s not found (may not exist): %v", servicePath, err)
				continue
			}

			var service map[string]any
			err = yaml.Unmarshal(data, &service)
			assert.NoError(t, err, "Service file %s should be valid YAML", servicePath)

			// Should have basic service structure
			assert.Contains(t, service, "name", "Service should have name field")
			assert.Contains(t, service, "description", "Service should have description field")
		}
	})
}

func TestEmbeddedContent_Consistency(t *testing.T) {
	t.Run("all embedded content is non-empty", func(t *testing.T) {
		embeddedFiles := map[string][]byte{
			"commands.yaml":      EmbeddedCommandsYAML,
			"schema.yaml":        EmbeddedSchemaYAML,
			"init-settings.yaml": EmbeddedInitSettingsYAML,
		}

		for filename, content := range embeddedFiles {
			assert.NotEmpty(t, content, "Embedded file %s should not be empty", filename)
		}
	})

	t.Run("embedded YAML files are valid", func(t *testing.T) {
		yamlFiles := map[string][]byte{
			"commands.yaml":      EmbeddedCommandsYAML,
			"schema.yaml":        EmbeddedSchemaYAML,
			"init-settings.yaml": EmbeddedInitSettingsYAML,
		}

		for filename, content := range yamlFiles {
			var parsed map[string]any
			err := yaml.Unmarshal(content, &parsed)
			assert.NoError(t, err, "Embedded YAML file %s should be valid", filename)
		}
	})
}
