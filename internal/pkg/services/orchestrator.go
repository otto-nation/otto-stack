package services

import (
	"context"
	"log/slog"
	"path/filepath"
	"time"

	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
)

// Orchestrator provides unified service management operations
type Orchestrator struct {
	docker     *docker.Client
	logger     *slog.Logger
	projectDir string
	config     *config.Config
}

// NewOrchestrator creates a new service orchestrator
func NewOrchestrator(logger *slog.Logger, projectDir string, dockerClient *docker.Client) (*Orchestrator, error) {
	return &Orchestrator{
		docker:     dockerClient,
		logger:     logger,
		projectDir: projectDir,
	}, nil
}

func (o *Orchestrator) SetConfig(config *config.Config) {
	o.config = config
}

func (o *Orchestrator) Close() error {
	return o.docker.Close()
}

// Core operations using compose directly
func (o *Orchestrator) StartServices(ctx context.Context, services []string, options docker.StartOptions) error {
	return o.docker.ComposeUp(ctx, o.getProjectName(), services, options)
}

func (o *Orchestrator) StopServices(ctx context.Context, services []string, options docker.StopOptions) error {
	return o.docker.ComposeDown(ctx, o.getProjectName(), options)
}

func (o *Orchestrator) GetServiceStatus(ctx context.Context, services []string) ([]docker.DockerServiceStatus, error) {
	statuses, err := o.docker.GetDockerServiceStatus(ctx, o.getProjectName(), services)
	if err != nil {
		return nil, err
	}

	// Add uptime calculation
	for i := range statuses {
		if statuses[i].State.IsRunning() && statuses[i].StartedAt != nil {
			statuses[i].Uptime = time.Since(*statuses[i].StartedAt)
		}
	}

	return statuses, nil
}

func (o *Orchestrator) GetLogs(ctx context.Context, services []string, options docker.LogOptions) error {
	return o.docker.ComposeLogs(ctx, o.getProjectName(), services, options)
}

func (o *Orchestrator) ExecCommand(ctx context.Context, service string, cmd []string, options docker.ExecOptions) error {
	return docker.NewComposeBuilder().
		Project(o.getProjectName()).
		File(docker.DockerComposeFilePath).
		User(options.User).
		Workdir(options.WorkingDir).
		Exec(service, cmd...).
		Run()
}

// Resource cleanup
func (o *Orchestrator) CleanupResources(ctx context.Context, options docker.CleanupOptions) error {
	project := o.getProjectName()

	// Stop all services first
	if err := o.docker.ComposeDown(ctx, project, docker.StopOptions{
		Remove:        true,
		RemoveVolumes: options.RemoveVolumes,
	}); err != nil {
		return pkgerrors.NewServiceError("system", "stop services", err)
	}

	// Clean up additional resources if requested
	if options.RemoveVolumes {
		if err := o.docker.RemoveResources(ctx, docker.ResourceVolume, project); err != nil {
			o.logger.Error("Failed to remove volumes", "error", err)
		}
	}
	if options.RemoveImages {
		if err := o.docker.RemoveResources(ctx, docker.ResourceImage, project); err != nil {
			o.logger.Error("Failed to remove images", "error", err)
		}
	}
	if options.RemoveNetworks {
		if err := o.docker.RemoveResources(ctx, docker.ResourceNetwork, project); err != nil {
			o.logger.Error("Failed to remove networks", "error", err)
		}
	}

	return nil
}

func (o *Orchestrator) getProjectName() string {
	if o.config != nil && o.config.Project.Name != "" {
		return o.config.Project.Name
	}
	return filepath.Base(o.projectDir)
}
