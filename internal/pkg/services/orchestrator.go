package services

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// Orchestrator provides unified service management operations
type Orchestrator struct {
	docker     *docker.Client
	logger     *slog.Logger
	projectDir string
	config     *config.Config
}

// NewOrchestrator creates a new service orchestrator
func NewOrchestrator(logger *slog.Logger, projectDir string) (*Orchestrator, error) {
	dockerClient, err := docker.NewClient(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

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
func (o *Orchestrator) StartServices(ctx context.Context, services []string, options types.StartOptions) error {
	return o.docker.ComposeUp(ctx, o.getProjectName(), services, options)
}

func (o *Orchestrator) StopServices(ctx context.Context, services []string, options types.StopOptions) error {
	return o.docker.ComposeDown(ctx, o.getProjectName(), options)
}

func (o *Orchestrator) GetServiceStatus(ctx context.Context, services []string) ([]types.ServiceStatus, error) {
	statuses, err := o.docker.GetServiceStatus(ctx, o.getProjectName(), services)
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

func (o *Orchestrator) GetLogs(ctx context.Context, services []string, options types.LogOptions) error {
	return o.docker.ComposeLogs(ctx, o.getProjectName(), services, options)
}

func (o *Orchestrator) ExecCommand(ctx context.Context, service string, cmd []string, options types.ExecOptions) error {
	// Use docker compose exec for better integration
	args := []string{"compose", "-f", constants.DockerComposeFile, "-p", o.getProjectName(), "exec"}
	if options.User != "" {
		args = append(args, "--user", options.User)
	}
	if options.WorkingDir != "" {
		args = append(args, "--workdir", options.WorkingDir)
	}
	args = append(args, service)
	args = append(args, cmd...)

	return o.docker.RunCommand(ctx, args...)
}

// Resource cleanup
func (o *Orchestrator) CleanupResources(ctx context.Context, options types.CleanupOptions) error {
	project := o.getProjectName()

	// Stop all services first
	if err := o.docker.ComposeDown(ctx, project, types.StopOptions{
		Remove:        true,
		RemoveVolumes: options.RemoveVolumes,
	}); err != nil {
		return fmt.Errorf("failed to stop services: %w", err)
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
