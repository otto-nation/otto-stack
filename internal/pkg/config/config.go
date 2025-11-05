package config

import (
	"fmt"
	"maps"
	"os"

	"github.com/otto-nation/otto-stack/internal/config"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"gopkg.in/yaml.v3"
)

// CommandConfig represents command configuration (minimal for generators)
type CommandConfig struct {
	Commands map[string]Command `yaml:"commands"`
}

// Command represents a command definition
type Command struct {
	Handler         string                `yaml:"handler"`
	Description     string                `yaml:"description"`
	LongDescription string                `yaml:"long_description"`
	Flags           map[string]FlagConfig `yaml:"flags"`
}

// FlagConfig represents a flag definition
type FlagConfig struct {
	Type        string `yaml:"type"`
	Short       string `yaml:"short"`
	Description string `yaml:"description"`
	Default     any    `yaml:"default"`
}

// LoadConfig loads otto-stack configuration with local overrides
func LoadConfig() (*types.OttoStackConfig, error) {
	// Load base config
	baseConfig, err := loadBaseConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load base config: %w", err)
	}

	// Load local overrides if they exist
	localConfig, err := loadLocalConfig()
	if err != nil {
		// Local config is optional, just use base
		return baseConfig, nil
	}

	// Merge configs (local overrides base)
	return mergeConfigs(baseConfig, localConfig), nil
}

// LoadCommandConfig loads command configuration from embedded YAML
func LoadCommandConfig() (map[string]any, error) {
	var commandConfig map[string]any
	if err := yaml.Unmarshal(config.EmbeddedCommandsYAML, &commandConfig); err != nil {
		return nil, fmt.Errorf("failed to parse commands config: %w", err)
	}
	return commandConfig, nil
}

// LoadCommandConfigStruct loads command configuration as struct (for generators)
func LoadCommandConfigStruct() (*CommandConfig, error) {
	var commandConfig CommandConfig
	if err := yaml.Unmarshal(config.EmbeddedCommandsYAML, &commandConfig); err != nil {
		return nil, fmt.Errorf("failed to parse commands config: %w", err)
	}
	return &commandConfig, nil
}

// GenerateConfig creates a new otto-stack configuration file
func GenerateConfig(projectName string, services []string) ([]byte, error) {
	config := map[string]any{
		"project": map[string]any{
			"name":        projectName,
			"environment": constants.DefaultEnvironment,
		},
		"stack": map[string]any{
			"enabled": services,
		},
	}

	return yaml.Marshal(config)
}

// loadBaseConfig loads the main configuration file
func loadBaseConfig() (*types.OttoStackConfig, error) {
	configPath := fmt.Sprintf("%s/%s", constants.OttoStackDir, constants.ConfigFileName)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	var config types.OttoStackConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// loadLocalConfig loads local configuration overrides
func loadLocalConfig() (*types.OttoStackConfig, error) {
	configPath := fmt.Sprintf("%s/%s", constants.OttoStackDir, constants.LocalConfigFileName)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err // File doesn't exist, which is fine
	}

	var config types.OttoStackConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse local config: %w", err)
	}

	return &config, nil
}

// mergeConfigs merges base config with local overrides
func mergeConfigs(base, local *types.OttoStackConfig) *types.OttoStackConfig {
	merged := *base // Copy base

	// Override project settings
	if local.Project.Name != "" {
		merged.Project.Name = local.Project.Name
	}
	if local.Project.Environment != "" {
		merged.Project.Environment = local.Project.Environment
	}

	// Override stack settings
	if len(local.Stack.Enabled) > 0 {
		merged.Stack.Enabled = local.Stack.Enabled
	}

	// Merge service configuration
	if local.ServiceConfiguration != nil {
		if merged.ServiceConfiguration == nil {
			merged.ServiceConfiguration = make(map[string]any)
		}
		maps.Copy(merged.ServiceConfiguration, local.ServiceConfiguration)
	}

	return &merged
}
