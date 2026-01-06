package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestLoadCommandConfig(t *testing.T) {
	t.Run("loads command config successfully", func(t *testing.T) {
		config, err := LoadCommandConfig()
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.IsType(t, map[string]any{}, config)
	})

	t.Run("returns valid structure", func(t *testing.T) {
		config, err := LoadCommandConfig()
		require.NoError(t, err)

		// Should contain expected top-level keys
		expectedKeys := []string{"commands", "global", "messages"}
		for _, key := range expectedKeys {
			assert.Contains(t, config, key, "Config should contain %s section", key)
		}
	})
}

func TestLoadCommandConfigStruct(t *testing.T) {
	t.Run("loads command config as struct", func(t *testing.T) {
		config, err := LoadCommandConfigStruct()
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.IsType(t, &CommandConfig{}, config)
	})

	t.Run("populates struct fields", func(t *testing.T) {
		config, err := LoadCommandConfigStruct()
		require.NoError(t, err)

		assert.NotNil(t, config.Commands)
		assert.NotNil(t, config.Global)
		assert.NotEmpty(t, config.Commands, "Should have at least some commands")
	})

	t.Run("validates command structure", func(t *testing.T) {
		config, err := LoadCommandConfigStruct()
		require.NoError(t, err)

		// Check that commands have required fields
		for cmdName, cmd := range config.Commands {
			assert.NotEmpty(t, cmd.Description, "Command %s should have description", cmdName)
			// Handler and flags are optional
		}
	})
}

func TestGenerateConfig(t *testing.T) {
	t.Run("generates valid config YAML", func(t *testing.T) {
		projectName := "test-project"
		serviceNames := []string{"postgres", "redis"}

		configBytes, err := GenerateConfig(projectName, serviceNames, nil)
		assert.NoError(t, err)
		assert.NotEmpty(t, configBytes)

		// Should be valid YAML
		var config Config
		err = yaml.Unmarshal(configBytes, &config)
		assert.NoError(t, err)

		// Verify content
		assert.Equal(t, projectName, config.Project.Name)
		assert.Equal(t, []string{"postgres", "redis"}, config.Stack.Enabled)
	})

	t.Run("handles empty project name", func(t *testing.T) {
		_, err := GenerateConfig("", []string{"postgres"}, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "project name cannot be empty")
	})

	t.Run("handles empty services list", func(t *testing.T) {
		configBytes, err := GenerateConfig("test", []string{}, nil)
		assert.NoError(t, err)

		var config Config
		err = yaml.Unmarshal(configBytes, &config)
		assert.NoError(t, err)
		assert.Empty(t, config.Stack.Enabled)
	})

	t.Run("sets project type", func(t *testing.T) {
		configBytes, err := GenerateConfig("test", []string{"postgres"}, nil)
		require.NoError(t, err)

		var config Config
		err = yaml.Unmarshal(configBytes, &config)
		require.NoError(t, err)
		assert.NotEmpty(t, config.Project.Type)
	})
}

func TestConfig_Structure(t *testing.T) {
	t.Run("config struct has required fields", func(t *testing.T) {
		config := Config{
			Project: ProjectConfig{
				Name: "test",
				Type: "application",
			},
			Stack: StackConfig{
				Enabled: []string{"postgres"}, // Using string literal that matches ServicePostgres constant
			},
		}

		// Should marshal to YAML without errors
		yamlBytes, err := yaml.Marshal(config)
		assert.NoError(t, err)
		assert.NotEmpty(t, yamlBytes)

		// Should unmarshal back correctly
		var unmarshaled Config
		err = yaml.Unmarshal(yamlBytes, &unmarshaled)
		assert.NoError(t, err)
		assert.Equal(t, config.Project.Name, unmarshaled.Project.Name)
		assert.Equal(t, config.Stack.Enabled, unmarshaled.Stack.Enabled)
	})
}

func TestProjectConfig_Timestamps(t *testing.T) {
	t.Run("handles timestamp fields", func(t *testing.T) {
		now := time.Now()
		project := ProjectConfig{
			Name:      "test",
			Type:      "app",
			CreatedAt: now,
			UpdatedAt: now,
		}

		yamlBytes, err := yaml.Marshal(project)
		assert.NoError(t, err)

		var unmarshaled ProjectConfig
		err = yaml.Unmarshal(yamlBytes, &unmarshaled)
		assert.NoError(t, err)

		// Times should be preserved (within reasonable precision)
		assert.WithinDuration(t, now, unmarshaled.CreatedAt, time.Second)
		assert.WithinDuration(t, now, unmarshaled.UpdatedAt, time.Second)
	})
}

func TestFlagConfig_Types(t *testing.T) {
	t.Run("supports different flag types", func(t *testing.T) {
		flags := map[string]FlagConfig{
			"verbose": {
				Type:        "bool",
				Short:       "v",
				Description: "Enable verbose output",
				Default:     false,
			},
			"count": {
				Type:        "int",
				Description: "Number of items",
				Default:     10,
			},
			"name": {
				Type:        "string",
				Description: "Project name",
				Default:     "default",
			},
		}

		yamlBytes, err := yaml.Marshal(flags)
		assert.NoError(t, err)

		var unmarshaled map[string]FlagConfig
		err = yaml.Unmarshal(yamlBytes, &unmarshaled)
		assert.NoError(t, err)

		assert.Equal(t, "bool", unmarshaled["verbose"].Type)
		assert.Equal(t, "v", unmarshaled["verbose"].Short)
		assert.Equal(t, false, unmarshaled["verbose"].Default)
		assert.Equal(t, 10, unmarshaled["count"].Default)
		assert.Equal(t, "default", unmarshaled["name"].Default)
	})
}
