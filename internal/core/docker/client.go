package docker

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type ResourceType string

const (
	ResourceContainer ResourceType = "container"
	ResourceVolume    ResourceType = "volume"
	ResourceNetwork   ResourceType = "network"
	ResourceImage     ResourceType = "image"
)

// Health status constants
const (
	HealthStatusHealthy   = "healthy"
	HealthStatusUnhealthy = "unhealthy"
	HealthStatusRunning   = "running"
	HealthStatusStopped   = "stopped"
	HealthStatusUnknown   = "unknown"

	ServiceStatusNotFound = "not found"
)

// ContainerInfo represents container information
type ContainerInfo struct {
	ID      string
	Name    string
	State   string
	Status  string
	Image   string
	Service string
}

// InitContainerConfig holds configuration for init containers
type InitContainerConfig struct {
	Image       string
	Command     []string
	Environment map[string]string
	Volumes     []string
	WorkingDir  string
	Networks    []string
}

// DockerClientInterface defines the interface for Docker operations
type DockerClientInterface interface {
	Close() error
}

// ComposeManagerInterface defines the interface for Compose operations
type ComposeManagerInterface interface {
	Up(ctx context.Context, project *types.Project, options api.UpOptions) error
	Down(ctx context.Context, project *types.Project, options api.DownOptions) error
}

type Client struct {
	cli       DockerClient
	logger    *slog.Logger
	resources *ResourceManager
	compose   *Manager
}

// NewClientWithDependencies creates a client with injected dependencies (for testing)
// deadcode: used for dependency injection in unit tests
func NewClientWithDependencies(cli DockerClient, compose *Manager, logger *slog.Logger) *Client {
	dc := &Client{
		cli:     cli,
		logger:  logger,
		compose: compose,
	}
	dc.resources = NewResourceManager(dc)
	return dc
}

func NewClient(logger *slog.Logger) (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, pkgerrors.NewDockerError("create Docker client", "", err)
	}

	composeManager, err := NewManager()
	if err != nil {
		return nil, pkgerrors.NewDockerError("create compose manager", "", err)
	}

	dc := &Client{
		cli:     NewDockerClientAdapter(cli),
		logger:  logger,
		compose: composeManager,
	}
	dc.resources = NewResourceManager(dc)

	return dc, nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}

// GetCli returns the underlying Docker client
func (c *Client) GetCli() DockerClient {
	return c.cli
}

// GetComposeManager returns the compose manager
// deadcode: used for accessing internal state in unit tests
func (c *Client) GetComposeManager() *Manager {
	return c.compose
}

// Resource management
func (c *Client) ListResources(ctx context.Context, resourceType ResourceType, project string) ([]string, error) {
	filter := NewProjectFilter(project)
	return c.resources.List(ctx, resourceType, filter)
}

func (c *Client) RemoveResources(ctx context.Context, resourceType ResourceType, project string) error {
	names, err := c.ListResources(ctx, resourceType, project)
	if err != nil {
		return err
	}
	return c.resources.Remove(ctx, resourceType, names)
}

// ListContainers lists containers for a project
func (c *Client) ListContainers(ctx context.Context, project string) ([]ContainerInfo, error) {
	filter := NewProjectFilter(project)
	containers, err := c.cli.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filter,
	})
	if err != nil {
		return nil, pkgerrors.NewServiceError(ComponentDocker, "list containers", err)
	}

	var result []ContainerInfo
	for _, cont := range containers {
		service := ""
		if serviceLabel, exists := cont.Labels[ComposeServiceLabel]; exists {
			service = serviceLabel
		}

		result = append(result, ContainerInfo{
			ID:      cont.ID,
			Name:    strings.TrimPrefix(cont.Names[0], "/"),
			State:   cont.State,
			Status:  cont.Status,
			Image:   cont.Image,
			Service: service,
		})
	}

	return result, nil
}

// RemoveContainer removes a container by ID
func (c *Client) RemoveContainer(ctx context.Context, containerID string, force bool) error {
	return c.cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: force})
}

// RunInitContainer runs an init container
func (c *Client) RunInitContainer(ctx context.Context, name string, config InitContainerConfig) error {
	// Create container
	containerConfig := &container.Config{
		Image:      config.Image,
		Cmd:        config.Command,
		Env:        mapToEnvSlice(config.Environment),
		WorkingDir: config.WorkingDir,
	}

	hostConfig := &container.HostConfig{
		AutoRemove: true,
	}

	// Add volumes
	if len(config.Volumes) > 0 {
		hostConfig.Binds = config.Volumes
	}

	networkConfig := &network.NetworkingConfig{}
	if len(config.Networks) > 0 {
		networkConfig.EndpointsConfig = make(map[string]*network.EndpointSettings)
		for _, net := range config.Networks {
			networkConfig.EndpointsConfig[net] = &network.EndpointSettings{}
		}
	}

	resp, err := c.cli.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, nil, name)
	if err != nil {
		return pkgerrors.NewServiceError(ComponentDocker, "create init container", err)
	}

	// Start container
	if err := c.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return pkgerrors.NewServiceError(ComponentDocker, "start init container", err)
	}

	// Wait for completion
	statusCh, errCh := c.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("error waiting for init container: %w", err)
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return fmt.Errorf("init container exited with code %d", status.StatusCode)
		}
	}

	return nil
}

// ContainerStatus represents basic container status
type ContainerStatus struct {
	Name      string
	State     string
	Health    string
	Ports     []string
	CreatedAt time.Time
	StartedAt time.Time
}

// GetServiceStatus gets status of services in a project
func (c *Client) GetServiceStatus(ctx context.Context, project string, services []string) ([]ContainerStatus, error) {
	containers, err := c.ListContainers(ctx, project)
	if err != nil {
		return nil, err
	}

	statusMap := make(map[string]*ContainerStatus)

	// Initialize status for requested services
	for _, service := range services {
		statusMap[service] = &ContainerStatus{
			Name:   service,
			State:  ServiceStatusNotFound,
			Health: HealthStatusUnknown,
		}
	}

	// Update with actual container status
	for _, cont := range containers {
		if cont.Service == "" {
			continue
		}

		status, exists := statusMap[cont.Service]
		if !exists {
			continue
		}

		status.State = cont.State
		status.Health = getHealthStatus(cont.Status)

		// Get detailed container info
		details, err := c.cli.ContainerInspect(ctx, cont.ID)
		if err != nil {
			continue
		}

		status.Ports = extractPorts(details.NetworkSettings.Ports)
		if created, err := time.Parse(time.RFC3339Nano, details.Created); err == nil {
			status.CreatedAt = created
		}
		if started, err := time.Parse(time.RFC3339Nano, details.State.StartedAt); err == nil {
			status.StartedAt = started
		}
	}

	var result []ContainerStatus
	for _, status := range statusMap {
		result = append(result, *status)
	}

	return result, nil
}

// extractPorts extracts port mappings from container network settings
func extractPorts(portMap nat.PortMap) []string {
	var ports []string
	for port, bindings := range portMap {
		if len(bindings) > 0 {
			ports = append(ports, bindings[0].HostPort+":"+port.Port())
		}
	}
	return ports
}

// getHealthStatus determines health status from container state and status
func getHealthStatus(status string) string {
	if strings.Contains(status, HealthStatusHealthy) {
		return HealthStatusHealthy
	}
	if strings.Contains(status, HealthStatusUnhealthy) {
		return HealthStatusUnhealthy
	}
	return "n/a"
}

// mapToEnvSlice converts a map to environment variable slice
func mapToEnvSlice(env map[string]string) []string {
	var result []string
	for key, value := range env {
		result = append(result, fmt.Sprintf("%s=%s", key, value))
	}
	return result
}

// Container status
