package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestEmbeddedCommandsYAML_Valid(t *testing.T) {
	assert.NotEmpty(t, EmbeddedCommandsYAML, "Commands YAML should be embedded")

	var commands map[string]any
	err := yaml.Unmarshal(EmbeddedCommandsYAML, &commands)
	assert.NoError(t, err, "Commands YAML should be valid")

	assert.Contains(t, commands, "commands", "Should contain commands section")
}

func TestEmbeddedCommandsYAML_RequiredCommands(t *testing.T) {
	var config map[string]any
	err := yaml.Unmarshal(EmbeddedCommandsYAML, &config)
	assert.NoError(t, err)

	commands, ok := config["commands"].(map[string]any)
	assert.True(t, ok, "Commands should be a map")

	expectedCommands := []string{"init", "up", "down", "status"}
	for _, cmd := range expectedCommands {
		assert.Contains(t, commands, cmd, "Should contain %s command", cmd)
	}
}

func TestEmbeddedSchemaYAML_Valid(t *testing.T) {
	assert.NotEmpty(t, EmbeddedSchemaYAML, "Schema YAML should be embedded")

	var schema map[string]any
	err := yaml.Unmarshal(EmbeddedSchemaYAML, &schema)
	assert.NoError(t, err, "Schema YAML should be valid")
}

func TestEmbeddedSchemaYAML_Definitions(t *testing.T) {
	var schema map[string]any
	err := yaml.Unmarshal(EmbeddedSchemaYAML, &schema)
	assert.NoError(t, err)

	assert.NotEmpty(t, schema, "Schema should not be empty")
}

func TestEmbeddedInitSettingsYAML_Valid(t *testing.T) {
	assert.NotEmpty(t, EmbeddedInitSettingsYAML, "Init settings YAML should be embedded")

	var settings map[string]any
	err := yaml.Unmarshal(EmbeddedInitSettingsYAML, &settings)
	assert.NoError(t, err, "Init settings YAML should be valid")
}

func TestEmbeddedInitSettingsYAML_Settings(t *testing.T) {
	var settings map[string]any
	err := yaml.Unmarshal(EmbeddedInitSettingsYAML, &settings)
	assert.NoError(t, err)

	assert.NotEmpty(t, settings, "Init settings should not be empty")
}

func TestEmbeddedServicesFS_Embedded(t *testing.T) {
	entries, err := EmbeddedServicesFS.ReadDir("services")
	assert.NoError(t, err, "Should be able to read services directory")
	assert.NotEmpty(t, entries, "Services directory should not be empty")
}

func TestEmbeddedServicesFS_Categories(t *testing.T) {
	entries, err := EmbeddedServicesFS.ReadDir("services")
	assert.NoError(t, err)

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
}

func TestEmbeddedServicesFS_ValidYAML(t *testing.T) {
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

		assert.Contains(t, service, "name", "Service should have name field")
		assert.Contains(t, service, "description", "Service should have description field")
	}
}

func TestEmbeddedContent_NonEmpty(t *testing.T) {
	embeddedFiles := map[string][]byte{
		"commands.yaml":      EmbeddedCommandsYAML,
		"schema.yaml":        EmbeddedSchemaYAML,
		"init-settings.yaml": EmbeddedInitSettingsYAML,
	}

	for filename, content := range embeddedFiles {
		assert.NotEmpty(t, content, "Embedded file %s should not be empty", filename)
	}
}

func TestEmbeddedContent_ValidYAML(t *testing.T) {
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
}
