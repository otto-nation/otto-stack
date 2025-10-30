package stack

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/types"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	utilsPkg "github.com/otto-nation/otto-stack/internal/pkg/utils"
)

// CoreSetup contains shared setup data for core commands
type CoreSetup struct {
	Config       *ProjectConfig
	DockerClient *docker.Client
}

// SetupCoreCommand performs common initialization for core commands
func SetupCoreCommand(ctx context.Context, base *types.BaseCommand) (*CoreSetup, func(), error) {
	// Check if otto-stack is initialized
	configPath := filepath.Join(constants.DevStackDir, constants.ConfigFileName)
	if !utilsPkg.FileExists(configPath) {
		return nil, nil, errors.New(constants.ErrNotInitialized)
	}

	// Load project configuration
	cfg, err := LoadProjectConfig(configPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create Docker client
	logger := base.Logger.(loggerAdapter)
	dockerClient, err := docker.NewClient(logger.SlogLogger())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Docker client: %w", err)
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
