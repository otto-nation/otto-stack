package docker

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"slices"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

type ResourceType string

const (
	ResourceContainer ResourceType = "container"
	ResourceVolume    ResourceType = "volume"
	ResourceNetwork   ResourceType = "network"
	ResourceImage     ResourceType = "image"
)

type Client struct {
	cli       *client.Client
	logger    *slog.Logger
	resources *ResourceManager
}

func NewClient(logger *slog.Logger) (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
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
func (c *Client) ComposeUp(ctx context.Context, project string, services []string, options types.StartOptions) error {
	args := []string{"compose", "-f", constants.DockerComposeFile, "-p", project, "up", "-d"}
	if options.Build {
		args = append(args, "--build")
	}
	if options.ForceRecreate {
		args = append(args, "--force-recreate")
	}
	args = append(args, services...)

	return c.RunCommand(ctx, args...)
}

func (c *Client) ComposeDown(ctx context.Context, project string, options types.StopOptions) error {
	args := []string{"compose", "-f", constants.DockerComposeFile, "-p", project, "down"}
	if options.Remove {
		args = append(args, "--remove-orphans")
	}
	if options.RemoveVolumes {
		args = append(args, "--volumes")
	}

	return c.RunCommand(ctx, args...)
}

func (c *Client) ComposeLogs(ctx context.Context, project string, services []string, options types.LogOptions) error {
	args := []string{"compose", "-f", constants.DockerComposeFile, "-p", project, "logs"}
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
func (c *Client) GetServiceStatus(ctx context.Context, project string, services []string) ([]types.ServiceStatus, error) {
	filter := NewProjectFilter(project)
	containers, err := c.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filter})
	if err != nil {
		return nil, err
	}

	var statuses []types.ServiceStatus
	for _, container := range containers {
		serviceName := container.Labels[constants.ComposeServiceLabel]
		if len(services) > 0 && !contains(services, serviceName) {
			continue
		}

		status := types.ServiceStatus{
			Name:  serviceName,
			State: types.ServiceState(container.State),
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

func (c *Client) RunCommand(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, constants.DockerCmd, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.logger.Error("Docker command failed", "args", args, "output", string(output), "error", err)
		return fmt.Errorf("docker command failed: %w", err)
	}
	return nil
}

func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}
