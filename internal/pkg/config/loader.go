package config

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"

	"github.com/otto-nation/otto-stack/internal/config"
	"gopkg.in/yaml.v3"
)

// Loader handles loading and parsing of command configurations
type Loader struct {
	configPath string
	cache      *CommandConfig
}

// NewLoader creates a new configuration loader
func NewLoader(configPath string) *Loader {
	return &Loader{
		configPath: configPath,
	}
}

// Load loads and parses the command configuration from YAML
func (l *Loader) Load() (*CommandConfig, error) {
	if l.cache != nil {
		return l.cache, nil
	}

	var data []byte
	var err error

	// Try to resolve config path
	configPath, pathErr := l.resolveConfigPath()
	if pathErr != nil {
		// If no file found, use embedded config
		data = config.EmbeddedCommandsYAML
	} else {
		// Read the YAML file
		data, err = os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
		}
	}

	// Parse YAML
	var config CommandConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config YAML: %w", err)
	}

	// Validate configuration
	if err := l.validateConfig(&config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Post-process configuration
	l.postProcessConfig(&config)

	l.cache = &config
	return &config, nil
}

// LoadFromPath loads configuration from a specific path
func (l *Loader) LoadFromPath(path string) (*CommandConfig, error) {
	oldPath := l.configPath
	l.configPath = path
	l.cache = nil // Clear cache

	config, err := l.Load()
	l.configPath = oldPath // Restore original path
	return config, err
}

// Reload clears cache and reloads configuration
func (l *Loader) Reload() (*CommandConfig, error) {
	l.cache = nil
	return l.Load()
}

// GetConfigPath returns the resolved configuration path
func (l *Loader) GetConfigPath() (string, error) {
	return l.resolveConfigPath()
}

// resolveConfigPath resolves the configuration file path
func (l *Loader) resolveConfigPath() (string, error) {
	if l.configPath == "" {
		// Try default locations
		candidates := []string{
			"internal/config/commands.yaml",
			"config/commands.yaml",
			".otto-stack/commands.yaml",
		}

		for _, candidate := range candidates {
			if _, err := os.Stat(candidate); err == nil {
				return filepath.Abs(candidate)
			}
		}

		return "", fmt.Errorf("no commands.yaml found in default locations: %v", candidates)
	}

	// Use provided path
	if !filepath.IsAbs(l.configPath) {
		return filepath.Abs(l.configPath)
	}

	return l.configPath, nil
}

// validateConfig performs basic structural validation
func (l *Loader) validateConfig(config *CommandConfig) error {
	if config.Metadata.Version == "" {
		return fmt.Errorf("metadata.version is required")
	}

	if len(config.Commands) == 0 {
		return fmt.Errorf("no commands defined")
	}

	// Validate each command has required fields
	for name, cmd := range config.Commands {
		if cmd.Description == "" {
			return fmt.Errorf("command %s: description is required", name)
		}
		if cmd.Usage == "" {
			return fmt.Errorf("command %s: usage is required", name)
		}
	}

	return nil
}

// postProcessConfig performs post-processing tasks
func (l *Loader) postProcessConfig(config *CommandConfig) {
	l.addCommandsToCategories(config)
	l.setDefaultFlagTypes(config)
}

// addCommandsToCategories ensures all commands are added to their categories
func (l *Loader) addCommandsToCategories(config *CommandConfig) {
	for cmdName, cmd := range config.Commands {
		if cmd.Category == "" {
			continue
		}

		category, exists := config.Categories[cmd.Category]
		if !exists {
			continue
		}

		if !l.categoryContainsCommand(category, cmdName) {
			category.Commands = append(category.Commands, cmdName)
			config.Categories[cmd.Category] = category
		}
	}
}

// categoryContainsCommand checks if a category already contains a command
func (l *Loader) categoryContainsCommand(category Category, cmdName string) bool {
	return slices.Contains(category.Commands, cmdName)
}

// setDefaultFlagTypes sets default types for flags that don't have them
func (l *Loader) setDefaultFlagTypes(config *CommandConfig) {
	for cmdName, cmd := range config.Commands {
		for flagName, flag := range cmd.Flags {
			if flag.Type == "" {
				flag.Type = "string" // Default type
				cmd.Flags[flagName] = flag
			}
		}
		config.Commands[cmdName] = cmd
	}
}

// LoadDefault loads the default commands configuration
func LoadDefault() (*CommandConfig, error) {
	loader := NewLoader("")
	return loader.Load()
}

// LoadFromBytes loads configuration from byte data
func LoadFromBytes(data []byte) (*CommandConfig, error) {
	var config CommandConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config YAML: %w", err)
	}

	loader := &Loader{}
	if err := loader.validateConfig(&config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	loader.postProcessConfig(&config)

	return &config, nil
}

// SaveConfig saves a configuration to a YAML file
func SaveConfig(config *CommandConfig, path string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, constants.DirPermReadWriteExec); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write file
	if err := os.WriteFile(path, data, constants.FilePermReadWrite); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", path, err)
	}

	return nil
}

// MergeConfigs merges multiple configurations, with later configs taking precedence
func MergeConfigs(configs ...*CommandConfig) (*CommandConfig, error) {
	if len(configs) == 0 {
		return nil, fmt.Errorf("no configurations to merge")
	}

	result := &CommandConfig{
		Metadata:   configs[0].Metadata,
		Global:     GlobalConfig{Flags: make(map[string]Flag)},
		Categories: make(map[string]Category),
		Commands:   make(map[string]Command),
		Help:       make(map[string]string),
	}

	for _, config := range configs {
		// Merge metadata (last wins)
		result.Metadata = config.Metadata

		// Merge global flags
		maps.Copy(result.Global.Flags, config.Global.Flags)

		// Merge categories
		maps.Copy(result.Categories, config.Categories)

		// Merge commands
		maps.Copy(result.Commands, config.Commands)

		// Merge help
		maps.Copy(result.Help, config.Help)
	}

	return result, nil
}
