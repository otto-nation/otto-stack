package docker

import (
	"context"
	"log/slog"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
)

// Manager wraps the official Docker Compose SDK
type Manager struct {
	service api.Compose
}

// NewManager creates a new Compose manager using the official SDK
func NewManager() (*Manager, error) {
	dockerCli, err := command.NewDockerCli()
	if err != nil {
		return nil, pkgerrors.NewServiceError(pkgerrors.ErrCodeInternal, pkgerrors.ComponentDocker, messages.ErrorsDockerCreateCliFailed, err)
	}

	if err := dockerCli.Initialize(flags.NewClientOptions()); err != nil {
		return nil, pkgerrors.NewServiceError(pkgerrors.ErrCodeInternal, pkgerrors.ComponentDocker, messages.ErrorsDockerInitializeCliFailed, err)
	}

	service, err := compose.NewComposeService(dockerCli)
	if err != nil {
		return nil, pkgerrors.NewServiceError(pkgerrors.ErrCodeInternal, pkgerrors.ComponentDocker, messages.ErrorsDockerCreateComposeFailed, err)
	}

	return &Manager{
		service: service,
	}, nil
}

// Up starts services using the official Compose SDK
func (m *Manager) Up(ctx context.Context, project *types.Project, options api.UpOptions) error {
	slog.Debug("Starting compose up",
		"project", project.Name,
		"services", len(project.Services),
		"working_dir", project.WorkingDir,
		"recreate", options.Create.Recreate,
		"remove_orphans", options.Create.RemoveOrphans)

	// Log service details at debug level
	for name, service := range project.Services {
		slog.Debug("Service configuration", "service", name, "image", service.Image)
	}

	// If force recreate is enabled, first try to remove existing containers
	if options.Create.Recreate == api.RecreateForce {
		slog.Debug("Force recreate enabled, removing existing containers", "project", project.Name)
		downOptions := api.DownOptions{
			RemoveOrphans: true,
		}
		// Ignore errors from down - containers might not exist
		_ = m.service.Down(ctx, project.Name, downOptions)
	}

	return m.service.Up(ctx, project, options)
}

// Down stops services using the official Compose SDK
func (m *Manager) Down(ctx context.Context, project *types.Project, options api.DownOptions) error {
	slog.Debug("Starting compose down", "project", project.Name, "remove_orphans", options.RemoveOrphans)

	return m.service.Down(ctx, project.Name, options)
}

// GetService returns the underlying compose service for direct access
func (m *Manager) GetService() api.Compose {
	return m.service
}

// LoadProject loads a compose project using the official SDK method
func (m *Manager) LoadProject(ctx context.Context, composePath string, projectName string) (*types.Project, error) {
	// Use the official SDK LoadProject method as shown in documentation
	project, err := m.service.LoadProject(ctx, api.ProjectLoadOptions{
		ConfigPaths: []string{composePath},
		ProjectName: projectName,
	})
	if err != nil {
		return nil, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentDocker, messages.ErrorsDockerLoadProjectFailed, err)
	}

	return project, nil
}

// Logs retrieves logs from services
func (m *Manager) Logs(ctx context.Context, projectName string, consumer api.LogConsumer, options api.LogOptions) error {
	return m.service.Logs(ctx, projectName, consumer, options)
}

// Stop stops services using the official Compose SDK
func (m *Manager) Stop(ctx context.Context, projectName string, options api.StopOptions) error {
	slog.Debug("Starting compose stop", "project", projectName, "services", options.Services)

	return m.service.Stop(ctx, projectName, options)
}

type SimpleLogConsumer struct{}

func (s *SimpleLogConsumer) Log(containerName, message string) {
	slog.Info("Container log", "container", containerName, "message", message)
}

func (s *SimpleLogConsumer) Err(containerName, message string) {
	slog.Error("Container error", "container", containerName, "message", message)
}

func (s *SimpleLogConsumer) Status(container, msg string) {
	slog.Info("Container status", "container", container, "status", msg)
}
