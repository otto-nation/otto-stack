package config

import (
	"os"
	"path/filepath"

	embeddedconfig "github.com/otto-nation/otto-stack/internal/config"
	"github.com/otto-nation/otto-stack/internal/core"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
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
		if err := validateSharingPolicy(baseConfig); err != nil {
			return nil, err
		}
		warnPhantomFields(baseConfig)
		return baseConfig, nil
	}

	merged := mergeConfigs(baseConfig, localConfig)
	if err := validateSharingPolicy(merged); err != nil {
		return nil, err
	}
	warnPhantomFields(merged)
	return merged, nil
}

// warnPhantomFields logs a warning for config fields that are parsed but have no effect.
func warnPhantomFields(cfg *Config) {
	if cfg.Validation != nil && len(cfg.Validation.Options) > 0 {
		logger.Warn(messages.WarningsValidationOptionsIgnored)
	}
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

	// validation.options is not yet implemented — omit it from generated configs.

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
		return nil, pkgerrors.NewConfigError(pkgerrors.ErrCodeOperationFail, core.LocalConfigFileName, messages.ErrorsConfigParseFailed, err)
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

// validateSharingPolicy validates that shared services are marked as shareable
func validateSharingPolicy(cfg *Config) error {
	if cfg.Sharing == nil || !cfg.Sharing.Enabled || len(cfg.Sharing.Services) == 0 {
		return nil
	}

	// Load service configs directly to avoid import cycle
	for svcName := range cfg.Sharing.Services {
		svcCfg, err := loadServiceConfig(svcName)
		if err != nil {
			continue // Unknown service, skip validation
		}
		if !svcCfg.Shareable {
			return pkgerrors.NewValidationErrorf(
				pkgerrors.ErrCodeInvalid,
				"sharing.services",
				messages.ValidationServiceNotShareable,
				svcName,
			)
		}
	}
	return nil
}

// loadServiceConfig loads a single service configuration from the embedded FS.
func loadServiceConfig(serviceName string) (*types.ServiceConfig, error) {
	target := serviceName + core.ExtYAML
	const embeddedServicesDir = "services"
	categories, err := embeddedconfig.EmbeddedServicesFS.ReadDir(embeddedServicesDir)
	if err != nil {
		return nil, os.ErrNotExist
	}
	for _, category := range categories {
		if !category.IsDir() {
			continue
		}
		filePath := embeddedServicesDir + "/" + category.Name() + "/" + target
		data, err := embeddedconfig.EmbeddedServicesFS.ReadFile(filePath)
		if err != nil {
			continue
		}
		var cfg types.ServiceConfig
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
		return &cfg, nil
	}
	return nil, os.ErrNotExist
}
