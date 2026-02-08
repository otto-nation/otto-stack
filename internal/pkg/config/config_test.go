//go:build unit

package config

import (
	"testing"
	"time"

	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestGenerateConfig(t *testing.T) {
	t.Run("generates valid config YAML", func(t *testing.T) {
		ctx := clicontext.NewBuilder().
			WithProject("test-project", "").
			WithServices([]string{"postgres", "redis"}, nil).
			Build()

		configBytes, err := GenerateConfig(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, configBytes)

		// Should be valid YAML
		var config Config
		err = yaml.Unmarshal(configBytes, &config)
		assert.NoError(t, err)

		// Verify content
		assert.Equal(t, "test-project", config.Project.Name)
		assert.Equal(t, []string{"postgres", "redis"}, config.Stack.Enabled)
	})

	t.Run("handles empty project name", func(t *testing.T) {
		ctx := clicontext.NewBuilder().
			WithProject("", "").
			WithServices([]string{"postgres"}, nil).
			Build()

		_, err := GenerateConfig(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), messages.ValidationProjectNameEmpty)
	})

	t.Run("handles empty services list", func(t *testing.T) {
		ctx := clicontext.NewBuilder().
			WithProject("test", "").
			WithServices([]string{}, nil).
			Build()

		configBytes, err := GenerateConfig(ctx)
		assert.NoError(t, err)

		var config Config
		err = yaml.Unmarshal(configBytes, &config)
		assert.NoError(t, err)
		assert.Empty(t, config.Stack.Enabled)
	})

	t.Run("sets project type", func(t *testing.T) {
		ctx := clicontext.NewBuilder().
			WithProject("test", "").
			WithServices([]string{"postgres"}, nil).
			Build()

		configBytes, err := GenerateConfig(ctx)
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

func TestValidateSharingPolicy(t *testing.T) {
	t.Run("allows shareable service", func(t *testing.T) {
		cfg := &Config{
			Sharing: &SharingConfig{
				Enabled:  true,
				Services: map[string]bool{"redis": true},
			},
		}
		err := validateSharingPolicy(cfg)
		assert.NoError(t, err)
	})

	t.Run("rejects non-shareable service", func(t *testing.T) {
		t.Skip("Skipping - requires running from project root")
		cfg := &Config{
			Sharing: &SharingConfig{
				Enabled:  true,
				Services: map[string]bool{"kafka": true},
			},
		}
		err := validateSharingPolicy(cfg)
		assert.Error(t, err)
		if err != nil {
			assert.Contains(t, err.Error(), "kafka")
		}
	})

	t.Run("skips validation when sharing disabled", func(t *testing.T) {
		cfg := &Config{
			Sharing: &SharingConfig{
				Enabled:  false,
				Services: map[string]bool{"kafka": true},
			},
		}
		err := validateSharingPolicy(cfg)
		assert.NoError(t, err)
	})

	t.Run("skips unknown services", func(t *testing.T) {
		cfg := &Config{
			Sharing: &SharingConfig{
				Enabled:  true,
				Services: map[string]bool{"unknown-service": true},
			},
		}
		err := validateSharingPolicy(cfg)
		assert.NoError(t, err)
	})
}
