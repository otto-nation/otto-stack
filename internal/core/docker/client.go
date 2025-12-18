package docker

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"slices"

	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type ResourceType string

const (
	ResourceContainer ResourceType = "container"
	ResourceVolume    ResourceType = "volume"
	ResourceNetwork   ResourceType = "network"
	ResourceImage     ResourceType = "image"
)

// InitContainerConfig holds configuration for init containers
type InitContainerConfig struct {
	Image       string
	Command     []string
	Environment map[string]string
	Volumes     []string
	WorkingDir  string
	Networks    []string
}

type Client struct {
	cli       *client.Client
	logger    *slog.Logger
	resources *ResourceManager
}

func NewClient(logger *slog.Logger) (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, pkgerrors.NewDockerError("create Docker client", "", err)
	}

	dc := &Client{
		cli:    cli,
		logger: logger,
	}
	dc.resources = NewResourceManager(dc)

	return dc, nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}

// Compose operations using docker compose CLI
func (c *Client) ComposeUp(ctx context.Context, project string, services []string, options StartOptions) error {
	args := []string{"compose", "-f", DockerComposeFilePath, "-p", project, "up", "-d"}
	if options.Build {
		args = append(args, "--build")
	}
	if options.ForceRecreate {
		args = append(args, "--force-recreate")
	}
	if options.RemoveOrphans {
		args = append(args, "--remove-orphans")
	}
	args = append(args, services...)

	return c.RunCommand(ctx, args...)
}

func (c *Client) ComposeDown(ctx context.Context, project string, options StopOptions) error {
	var args []string
	if options.Remove {
		args = []string{"compose", "-f", DockerComposeFilePath, "-p", project, "down"}
		if options.RemoveOrphans {
			args = append(args, "--remove-orphans")
		}
		if options.RemoveVolumes {
			args = append(args, "--volumes")
		}
	} else {
		args = []string{"compose", "-f", DockerComposeFilePath, "-p", project, "stop"}
		if options.Timeout > 0 {
			args = append(args, "--timeout", fmt.Sprintf("%d", options.Timeout))
		}
	}

	return c.RunCommand(ctx, args...)
}

func (c *Client) ComposeLogs(ctx context.Context, project string, services []string, options LogOptions) error {
	args := []string{"compose", "-f", DockerComposeFilePath, "-p", project, "logs"}
	if options.Follow {
		args = append(args, "--follow")
	}
	if options.Timestamps {
		args = append(args, "--timestamps")
	}
	if options.Tail != "" {
		args = append(args, "--tail", options.Tail)
	}
	args = append(args, services...)

	return c.RunCommand(ctx, args...)
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

// Container status
func (c *Client) GetDockerServiceStatus(ctx context.Context, project string, services []string) ([]DockerServiceStatus, error) {
	filter := NewProjectFilter(project)
	containers, err := c.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filter})
	if err != nil {
		return nil, err
	}

	var statuses []DockerServiceStatus
	for _, container := range containers {
		serviceName := container.Labels[ComposeServiceLabel]
		if len(services) > 0 && !contains(services, serviceName) {
			continue
		}

		status := DockerServiceStatus{
			Name:   serviceName,
			State:  DockerServiceState(container.State),
			Health: c.getContainerHealth(container),
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

func (c *Client) getContainerHealth(cont container.Summary) DockerHealthStatus {
	switch cont.State {
	case StateRunning:
		return DockerHealthStatusHealthy
	case StateStopped:
		return DockerHealthStatusUnhealthy
	case StateStarting:
		return DockerHealthStatusStarting
	default:
		return DockerHealthStatusNone
	}
}

func (c *Client) RunCommand(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, DockerCmd, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunInitContainer runs an init container and waits for completion
func (c *Client) RunInitContainer(ctx context.Context, containerName string, config InitContainerConfig) error {
	args := []string{"run", "--rm", "--name", containerName}

	// Add environment variables
	for key, value := range config.Environment {
		args = append(args, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	// Add volumes
	for _, volume := range config.Volumes {
		args = append(args, "-v", volume)
	}

	// Add working directory
	if config.WorkingDir != "" {
		args = append(args, "-w", config.WorkingDir)
	}

	// Add networks
	for _, network := range config.Networks {
		args = append(args, "--network", network)
	}

	// Add image and command
	args = append(args, config.Image)
	args = append(args, config.Command...)

	return c.RunCommand(ctx, args...)
}

func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}
