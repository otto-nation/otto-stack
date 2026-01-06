package stack

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
)

// CoreSetup contains shared setup data for core commands
type CoreSetup struct {
	Config       *config.Config
	DockerClient *docker.Client
}

// SetupCoreCommand performs common initialization for core commands
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

// ResolveServiceConfigs resolves services to ServiceConfigs using consistent logic across handlers
func ResolveServiceConfigs(args []string, setup *CoreSetup) ([]services.ServiceConfig, error) {
	if len(args) > 0 {
		// Resolve specific services from args
		serviceConfigs, err := services.ResolveUpServices(args, setup.Config)
		return serviceConfigs, err
	}
	// Use enabled services from config
	serviceConfigs, err := services.ResolveUpServices(setup.Config.Stack.Enabled, setup.Config)
	return serviceConfigs, err
}
