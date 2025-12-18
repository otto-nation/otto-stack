package stack

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

// ServicesManagerFactory creates service managers
type ServicesManagerFactory interface {
	New() (*services.Manager, error)
}

// DefaultServicesManagerFactory implements ServicesManagerFactory
type DefaultServicesManagerFactory struct{}

func (f *DefaultServicesManagerFactory) New() (*services.Manager, error) {
	return services.New()
}

// DockerClientFactory creates Docker clients
type DockerClientFactory interface {
	NewClient(logger any) (*docker.Client, error)
}

// DefaultDockerClientFactory implements DockerClientFactory
type DefaultDockerClientFactory struct{}

func (f *DefaultDockerClientFactory) NewClient(logger any) (*docker.Client, error) {
	var loggerPtr *slog.Logger
	if logger != nil {
		if l, ok := logger.(*slog.Logger); ok {
			loggerPtr = l
		}
	}
	return docker.NewClient(loggerPtr)
}

// CoreSetup contains shared setup data for core commands
type CoreSetup struct {
	Config       *config.Config
	DockerClient *docker.Client
}

var (
	dockerFactory   DockerClientFactory    = &DefaultDockerClientFactory{}
	servicesFactory ServicesManagerFactory = &DefaultServicesManagerFactory{}
)

// SetDockerClientFactory allows injection of custom Docker client factory (for testing)
func SetDockerClientFactory(factory DockerClientFactory) {
	dockerFactory = factory
}

// SetServicesManagerFactory allows injection of custom services manager factory (for testing)
func SetServicesManagerFactory(factory ServicesManagerFactory) {
	servicesFactory = factory
}

// GetServicesManager returns a services manager using the configured factory
func GetServicesManager() (*services.Manager, error) {
	return servicesFactory.New()
}

// SetupCoreCommand performs common initialization for core commands
func SetupCoreCommand(ctx context.Context, base *base.BaseCommand) (*CoreSetup, func(), error) {
	// Check if otto-stack is initialized (redundant check for safety)
	configPath := filepath.Join(core.OttoStackDir, core.ConfigFileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil, errors.New(core.MsgErrors_not_initialized)
	}

	// Load project configuration
	cfg, err := LoadProjectConfig(configPath)
	if err != nil {
		return nil, nil, fmt.Errorf(core.MsgStack_failed_load_config, err)
	}

	// Create Docker client using factory
	dockerClient, err := dockerFactory.NewClient(nil)
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
