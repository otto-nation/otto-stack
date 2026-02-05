package config

import (
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"gopkg.in/yaml.v3"
)

const (
	// Default configuration values
	DefaultProjectType = "docker"

	// Error message templates
	ErrConfigNotFound     = "config file not found: %s"
	ErrConfigParse        = "failed to parse config: %w"
	ErrLocalConfigParse   = "failed to parse local config: %w"
	ErrServiceConfigParse = "failed to parse %s config: %w"
	ErrServiceNotFound    = "%s config not found for: %s"
	ErrCommandConfigParse = "failed to parse commands config: %w"
	ErrLoadBaseConfig     = "failed to load base config: %w"

	// Config types for error messages
	ConfigTypeService      = "service"
	ConfigTypeLocalService = "local service"
)

// CommandConfig represents command configuration (minimal for generators)
type CommandConfig struct {
	Commands map[string]Command `yaml:"commands"`
	Global   GlobalConfig       `yaml:"global"`
}

// GlobalConfig represents global configuration
type GlobalConfig struct {
	Flags map[string]FlagConfig `yaml:"flags"`
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
func LoadConfig() (*Config, error) {
	baseConfig, err := loadBaseConfig()
	if err != nil {
		return nil, pkgerrors.NewConfigError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentConfig, messages.ErrorsConfigLoadFailed, err)
	}

	localConfig, err := loadLocalConfig()
	if err != nil {
		// Local config is optional, just use base
		return baseConfig, nil
	}

	return mergeConfigs(baseConfig, localConfig), nil
}

// GenerateConfig creates a new otto-stack configuration file
func GenerateConfig(ctx clicontext.Context) ([]byte, error) {
	if ctx.Project.Name == "" {
		return nil, pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldProjectName, messages.ValidationProjectNameEmpty, nil)
	}

	config := Config{
		Project: ProjectConfig{
			Name: ctx.Project.Name,
			Type: DefaultProjectType,
		},
		Stack: StackConfig{
			Enabled: ctx.Services.Names,
		},
	}

	if ctx.Sharing != nil {
		config.Sharing = &SharingConfig{
			Enabled:  ctx.Sharing.Enabled,
			Services: ctx.Sharing.Services,
		}
	}

	if ctx.Options.Validation != nil {
		config.Validation = &ValidationConfig{
			Options: ctx.Options.Validation,
		}
	}

	return yaml.Marshal(config)
}

// getConfigPath returns the path to the main config file
func getConfigPath() string {
	return filepath.Join(core.OttoStackDir, core.ConfigFileName)
}

// getLocalConfigPath returns the path to the local config file
func getLocalConfigPath() string {
	return filepath.Join(core.OttoStackDir, core.LocalConfigFileName)
}

// loadBaseConfig loads the main configuration file
func loadBaseConfig() (*Config, error) {
	configPath := getConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, pkgerrors.NewConfigErrorf(pkgerrors.ErrCodeNotFound, configPath, messages.ErrorsConfigNotFound, configPath)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, pkgerrors.NewConfigError(pkgerrors.ErrCodeOperationFail, configPath, messages.ErrorsConfigParseFailed, err)
	}

	return &config, nil
}

// loadLocalConfig loads local configuration overrides
func loadLocalConfig() (*Config, error) {
	configPath := getLocalConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err // File doesn't exist, which is fine
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, pkgerrors.NewConfigError(pkgerrors.ErrCodeOperationFail, "config.local.yaml", messages.ErrorsConfigParseFailed, err)
	}

	return &config, nil
}

// mergeConfigs merges base config with local overrides
func mergeConfigs(base, local *Config) *Config {
	merged := *base // Copy base

	// Override project settings
	if local.Project.Name != "" {
		merged.Project.Name = local.Project.Name
	}

	// Override stack settings
	if len(local.Stack.Enabled) > 0 {
		merged.Stack.Enabled = local.Stack.Enabled
	}

	return &merged
}

// loadServiceConfigFile loads service configuration (base or local)
