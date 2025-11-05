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

// Manager provides unified service management operations
type Manager struct {
	docker     *docker.Client
	logger     *slog.Logger
	projectDir string
	config     *config.Config
}

// NewManager creates a new service manager
func NewManager(logger *slog.Logger, projectDir string) (*Manager, error) {
	dockerClient, err := docker.NewClient(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	return &Manager{
		docker:     dockerClient,
		logger:     logger,
		projectDir: projectDir,
	}, nil
}

func (m *Manager) SetConfig(config *config.Config) {
	m.config = config
}

func (m *Manager) Close() error {
	return m.docker.Close()
}

// Core operations using compose directly
func (m *Manager) StartServices(ctx context.Context, services []string, options types.StartOptions) error {
	return m.docker.ComposeUp(ctx, m.getProjectName(), services, options)
}

func (m *Manager) StopServices(ctx context.Context, services []string, options types.StopOptions) error {
	return m.docker.ComposeDown(ctx, m.getProjectName(), options)
}

func (m *Manager) GetServiceStatus(ctx context.Context, services []string) ([]types.ServiceStatus, error) {
	statuses, err := m.docker.GetServiceStatus(ctx, m.getProjectName(), services)
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

func (m *Manager) GetLogs(ctx context.Context, services []string, options types.LogOptions) error {
	return m.docker.ComposeLogs(ctx, m.getProjectName(), services, options)
}

func (m *Manager) ExecCommand(ctx context.Context, service string, cmd []string, options types.ExecOptions) error {
	// Use docker compose exec for better integration
	args := []string{"compose", "-f", constants.DockerComposeFile, "-p", m.getProjectName(), "exec"}
	if options.User != "" {
		args = append(args, "--user", options.User)
	}
	if options.WorkingDir != "" {
		args = append(args, "--workdir", options.WorkingDir)
	}
	args = append(args, service)
	args = append(args, cmd...)

	return m.docker.RunCommand(ctx, args...)
}

// Resource cleanup
func (m *Manager) CleanupResources(ctx context.Context, options types.CleanupOptions) error {
	project := m.getProjectName()

	// Stop all services first
	if err := m.docker.ComposeDown(ctx, project, types.StopOptions{
		Remove:        true,
		RemoveVolumes: options.RemoveVolumes,
	}); err != nil {
		return fmt.Errorf("failed to stop services: %w", err)
	}

	// Clean up additional resources if requested
	if options.RemoveVolumes {
		if err := m.docker.RemoveResources(ctx, docker.ResourceVolume, project); err != nil {
			m.logger.Error("Failed to remove volumes", "error", err)
		}
	}
	if options.RemoveImages {
		if err := m.docker.RemoveResources(ctx, docker.ResourceImage, project); err != nil {
			m.logger.Error("Failed to remove images", "error", err)
		}
	}
	if options.RemoveNetworks {
		if err := m.docker.RemoveResources(ctx, docker.ResourceNetwork, project); err != nil {
			m.logger.Error("Failed to remove networks", "error", err)
		}
	}

	return nil
}

func (m *Manager) getProjectName() string {
	if m.config != nil && m.config.Project.Name != "" {
		return m.config.Project.Name
	}
	return filepath.Base(m.projectDir)
}
