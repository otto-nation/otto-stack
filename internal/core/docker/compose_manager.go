package docker

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
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
func (m *Manager) Up(ctx context.Context, project *types.Project, options UpOptions) error {
	sdkOptions := options.ToSDK()

	slog.Debug("Starting compose up",
		"project", project.Name,
		"services", len(project.Services),
		"working_dir", project.WorkingDir,
		"recreate", sdkOptions.Create.Recreate,
		"remove_orphans", sdkOptions.Create.RemoveOrphans,
		"pull_latest", options.PullLatestImages)

	// Pull latest images before starting when requested.
	if options.PullLatestImages {
		if err := m.service.Pull(ctx, project, api.PullOptions{}); err != nil {
			slog.Warn("Failed to pull latest images", "error", err)
		}
	}

	// Log service details at debug level
	for name, service := range project.Services {
		slog.Debug("Service configuration", "service", name, "image", service.Image)
	}

	// If force recreate is enabled, first try to remove existing containers.
	// When cleanup_on_recreate is set, also purge volumes for a full reset.
	if sdkOptions.Create.Recreate == api.RecreateForce {
		slog.Debug("Force recreate enabled, removing existing containers", "project", project.Name)
		downOptions := api.DownOptions{
			RemoveOrphans: true,
			Volumes:       options.CleanupOnRecreate,
		}
		// Ignore errors from down — containers might not exist yet
		_ = m.service.Down(ctx, project.Name, downOptions)
	}

	return m.service.Up(ctx, project, sdkOptions)
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

// LoadProject loads a compose project using the official SDK method.
// configPaths is a list of compose files to load; later entries override earlier ones.
func (m *Manager) LoadProject(ctx context.Context, configPaths []string, projectName string) (*types.Project, error) {
	project, err := m.service.LoadProject(ctx, api.ProjectLoadOptions{
		ConfigPaths: configPaths,
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

// Exec executes a command in a service container
func (m *Manager) Exec(ctx context.Context, projectName string, options api.RunOptions) (int, error) {
	slog.Debug("Starting compose exec", "project", projectName, "service", options.Service, "command", options.Command)

	return m.service.Exec(ctx, projectName, options)
}

// logColorPalette is the ordered set of ANSI colors assigned to service name
// prefixes when tailing logs from multiple services. Red is intentionally
// omitted — it is reserved for error lines.
var logColorPalette = []string{
	ui.ColorCyan,
	ui.ColorBlue,
	ui.ColorYellow,
	ui.ColorMagenta,
	ui.ColorGreen,
}

// ServiceLogConsumer writes container log lines to an io.Writer, optionally
// prefixing each line with a color-coded service name when multiple services
// are being tailed.
type ServiceLogConsumer struct {
	writer       io.Writer
	noColor      bool
	multiService bool
	mu           sync.Mutex
	serviceColor map[string]string
	nextColorIdx int
}

// NewServiceLogConsumer creates a ServiceLogConsumer. When serviceCount is 1
// no per-service prefix is printed; when serviceCount > 1 each service gets
// a distinct color and its name is prepended to every line.
func NewServiceLogConsumer(writer io.Writer, noColor bool, serviceCount int) *ServiceLogConsumer {
	return &ServiceLogConsumer{
		writer:       writer,
		noColor:      noColor,
		multiService: serviceCount > 1,
		serviceColor: make(map[string]string),
	}
}

func (c *ServiceLogConsumer) Log(containerName, message string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, _ = fmt.Fprintln(c.writer, c.formatLine(containerName, message, false))
}

func (c *ServiceLogConsumer) Err(containerName, message string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, _ = fmt.Fprintln(c.writer, c.formatLine(containerName, message, true))
}

// Status carries container lifecycle events (e.g. "container started").
// These are debug-level noise and are not shown to the user.
func (c *ServiceLogConsumer) Status(container, msg string) {
	slog.Debug("Container status", "container", container, "status", msg)
}

// formatLine builds the output line for a single log entry. When isErr is
// true the message itself is colored red to distinguish stderr from stdout.
func (c *ServiceLogConsumer) formatLine(containerName, message string, isErr bool) string {
	if isErr && !c.noColor {
		message = ui.ColorRed + message + ui.ColorReset
	}

	if !c.multiService {
		return message
	}

	prefix := c.prefixFor(containerName)
	return prefix + " | " + message
}

// prefixFor returns a color-coded service name prefix, assigning a new color
// from the palette on first encounter.
func (c *ServiceLogConsumer) prefixFor(containerName string) string {
	color, ok := c.serviceColor[containerName]
	if !ok {
		color = logColorPalette[c.nextColorIdx%len(logColorPalette)]
		c.serviceColor[containerName] = color
		c.nextColorIdx++
	}

	if c.noColor {
		return containerName
	}
	return color + containerName + ui.ColorReset
}
