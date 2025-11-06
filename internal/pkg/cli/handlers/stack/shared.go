package stack

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// CoreSetup contains shared setup data for core commands
type CoreSetup struct {
	Config       *config.Config
	DockerClient *docker.Client
}

// SetupCoreCommand performs common initialization for core commands
func SetupCoreCommand(ctx context.Context, base *types.BaseCommand) (*CoreSetup, func(), error) {
	// Check if otto-stack is initialized (redundant check for safety)
	configPath := filepath.Join(constants.OttoStackDir, constants.ConfigFileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil, errors.New(constants.MsgErrors_not_initialized)
	}

	// Load project configuration
	cfg, err := LoadProjectConfig(configPath)
	if err != nil {
		return nil, nil, fmt.Errorf(constants.MsgStack_failed_load_config, err)
	}

	// Create Docker client
	dockerClient, err := docker.NewClient(nil)
	if err != nil {
		return nil, nil, fmt.Errorf(constants.MsgStack_failed_create_docker_client, err)
	}

	cleanup := func() {
		_ = dockerClient.Close()
	}

	return &CoreSetup{
		Config:       cfg,
		DockerClient: dockerClient,
	}, cleanup, nil
}
