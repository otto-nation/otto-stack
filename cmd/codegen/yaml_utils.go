package codegen

import (
	"os"
	"strings"

	"github.com/otto-nation/otto-stack/internal/config"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"gopkg.in/yaml.v3"
)

// LoadYAMLConfig loads and parses a YAML file into a map
func LoadYAMLConfig(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config map[string]any
	err = yaml.Unmarshal(data, &config)
	return config, err
}

// LoadCommandConfig loads command configuration from embedded YAML
func LoadCommandConfig() (map[string]any, error) {
	var commandConfig map[string]any
	if err := yaml.Unmarshal(config.EmbeddedCommandsYAML, &commandConfig); err != nil {
		return nil, pkgerrors.NewConfigError("", "parse command config", err)
	}
	return commandConfig, nil
}

// CommandConfig represents the structure of commands.yaml
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

// LoadCommandConfigStruct loads command configuration as struct
func LoadCommandConfigStruct() (*CommandConfig, error) {
	var commandConfig CommandConfig
	if err := yaml.Unmarshal(config.EmbeddedCommandsYAML, &commandConfig); err != nil {
		return nil, pkgerrors.NewConfigError("", "parse command config", err)
	}
	return &commandConfig, nil
}

// ToPascalCase converts a string to PascalCase
func ToPascalCase(s string) string {
	if s == "" {
		return ""
	}

	parts := strings.FieldsFunc(s, func(c rune) bool {
		return c == '-' || c == '_' || c == ' '
	})

	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}

	return strings.Join(parts, "")
}
