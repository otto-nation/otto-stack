package common

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"gopkg.in/yaml.v3"
)

// CoreSetup contains common setup data for handlers
type CoreSetup struct {
	Config       *config.Config
	DockerClient *docker.Client
}

// SetupCoreCommand provides common setup for handlers that need Docker and config
func SetupCoreCommand(ctx context.Context, base *base.BaseCommand) (*CoreSetup, func(), error) {
	// Check if otto-stack is initialized
	configPath := filepath.Join(core.OttoStackDir, core.ConfigFileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil, errors.New(core.MsgErrors_not_initialized)
	}

	// Load project configuration
	cfg, err := LoadProjectConfig(configPath)
	if err != nil {
		return nil, nil, fmt.Errorf(core.MsgStack_failed_load_config, err)
	}

	// Create Docker client
	dockerClient, err := docker.NewClient(nil)
	if err != nil {
		return nil, nil, fmt.Errorf(core.MsgStack_failed_create_docker_client, err)
	}

	cleanup := func() {
		_ = dockerClient.Close()
	}

	return &CoreSetup{
		Config:       cfg,
		DockerClient: dockerClient,
	}, cleanup, nil
}

// LoadProjectConfig loads the project configuration from the given path
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

// ResolveServiceConfigs resolves service configurations from args or enabled services
func ResolveServiceConfigs(args []string, setup *CoreSetup) ([]types.ServiceConfig, error) {
	var serviceConfigs []types.ServiceConfig
	var err error

	if len(args) > 0 {
		serviceConfigs, err = services.ResolveUpServices(args, setup.Config)
	} else {
		serviceConfigs, err = services.ResolveUpServices(setup.Config.Stack.Enabled, setup.Config)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to resolve services: %w", err)
	}

	return serviceConfigs, nil
}
