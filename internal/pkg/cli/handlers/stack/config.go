package stack

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/config"

	"gopkg.in/yaml.v3"
)

// LoadProjectConfig loads the otto-stack project configuration with local overrides
func LoadProjectConfig(configPath string) (*config.Config, error) {
	// Load base config
	baseConfig, err := loadSingleConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf(core.MsgStack_failed_load_base_config, err)
	}

	// Try to load local config
	localPath := filepath.Join(core.OttoStackDir, core.LocalConfigFileName)
	localConfig, err := loadSingleConfig(localPath)
	if err != nil {
		// Local config is optional, return base config if not found
		if os.IsNotExist(err) {
			return baseConfig, nil
		}
		return nil, fmt.Errorf(core.MsgStack_failed_load_local_config, err)
	}

	// Merge configs (local overrides base)
	merged := mergeProjectConfigs(baseConfig, localConfig)
	return merged, nil
}

// loadSingleConfig loads a single config file
func loadSingleConfig(configPath string) (*config.Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf(core.MsgStack_failed_parse_config, err)
	}

	return &cfg, nil
}

// mergeProjectConfigs merges local config into base config
func mergeProjectConfigs(base, local *config.Config) *config.Config {
	merged := *base // Copy base

	// Override project settings if specified in local
	if local.Project.Name != "" {
		merged.Project.Name = local.Project.Name
	}

	// Override stack settings if specified in local
	if len(local.Stack.Enabled) > 0 {
		merged.Stack.Enabled = local.Stack.Enabled
	}

	return &merged
}
