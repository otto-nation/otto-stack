package stack

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// CoreSetup contains shared setup data for core commands
type CoreSetup struct {
	Config       *ProjectConfig
	DockerClient *docker.Client
}

// SetupCoreCommand performs common initialization for core commands
func SetupCoreCommand(ctx context.Context, base *types.BaseCommand) (*CoreSetup, func(), error) {
	// Check if otto-stack is initialized
	configPath := filepath.Join(constants.OttoStackDir, constants.ConfigFileName)
	if !func() bool { _, err := os.Stat(configPath); return err == nil }() {
		return nil, nil, errors.New(constants.Messages[constants.MsgErrors_not_initialized])
	}

	// Load project configuration
	cfg, err := LoadProjectConfig(configPath)
	if err != nil {
		return nil, nil, fmt.Errorf(constants.Messages[constants.MsgStack_failed_load_config], err)
	}

	// Create Docker client
	logger := base.Logger
	dockerClient, err := docker.NewClient(logger.SlogLogger())
	if err != nil {
		return nil, nil, fmt.Errorf(constants.Messages[constants.MsgStack_failed_create_docker_client], err)
	}

	cleanup := func() {
		if err := dockerClient.Close(); err != nil {
			base.Logger.Error("Failed to close Docker client", "error", err)
		}
	}

	return &CoreSetup{
		Config:       cfg,
		DockerClient: dockerClient,
	}, cleanup, nil
}
