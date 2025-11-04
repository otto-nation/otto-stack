package stack

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/utils"
	"gopkg.in/yaml.v3"
)

// loggerAdapter interface for accessing underlying slog.Logger
type loggerAdapter interface {
	SlogLogger() *slog.Logger
}

// ProjectConfig represents the otto-stack project configuration
type ProjectConfig struct {
	Project struct {
		Name        string `yaml:"name"`
		Environment string `yaml:"environment"`
	} `yaml:"project"`
	Stack struct {
		Enabled []string `yaml:"enabled"`
	} `yaml:"stack"`
}

// LoadProjectConfig loads the otto-stack project configuration with local overrides
func LoadProjectConfig(configPath string) (*ProjectConfig, error) {
	// Load base config
	baseConfig, err := loadSingleConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf(constants.Messages[constants.MsgStack_failed_load_base_config], err)
	}

	// Try to load local config
	localPath := filepath.Join(constants.OttoStackDir, constants.LocalConfigFileName)
	localConfig, err := loadSingleConfig(localPath)
	if err != nil {
		// Local config is optional, return base config if not found
		if os.IsNotExist(err) {
			return baseConfig, nil
		}
		return nil, fmt.Errorf(constants.Messages[constants.MsgStack_failed_load_local_config], err)
	}

	// Merge configs (local overrides base)
	merged := mergeProjectConfigs(baseConfig, localConfig)
	return merged, nil
}

// loadSingleConfig loads a single config file
func loadSingleConfig(configPath string) (*ProjectConfig, error) {
	data, err := utils.ReadFileLines(configPath)
	if err != nil {
		return nil, err
	}

	var builder strings.Builder
	for _, line := range data {
		builder.WriteString(line + "\n")
	}
	content := builder.String()

	var cfg ProjectConfig
	if err := yaml.Unmarshal([]byte(content), &cfg); err != nil {
		return nil, fmt.Errorf(constants.Messages[constants.MsgStack_failed_parse_config], err)
	}

	return &cfg, nil
}

// mergeProjectConfigs merges local config into base config
func mergeProjectConfigs(base, local *ProjectConfig) *ProjectConfig {
	merged := *base // Copy base

	// Override project settings if specified in local
	if local.Project.Name != "" {
		merged.Project.Name = local.Project.Name
	}
	if local.Project.Environment != "" {
		merged.Project.Environment = local.Project.Environment
	}

	// Override stack settings if specified in local
	if len(local.Stack.Enabled) > 0 {
		merged.Stack.Enabled = local.Stack.Enabled
	}

	return &merged
}
