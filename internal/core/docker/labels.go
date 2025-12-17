package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

const (
	LabelManaged     = "io.otto-stack.managed"
	LabelProject     = "io.otto-stack.project"
	LabelService     = "io.otto-stack.service"
	LabelVersion     = "io.otto-stack.version"
	LabelSharingMode = "io.otto-stack.sharing-mode"
)

// ListOttoContainers returns all containers managed by Otto Stack
func (c *Client) ListOttoContainers(ctx context.Context) ([]container.Summary, error) {
	filter := filters.NewArgs()
	filter.Add("label", fmt.Sprintf("%s=true", LabelManaged))

	return c.cli.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filter,
	})
}

// ListProjectContainers returns containers for a specific project
func (c *Client) ListProjectContainers(ctx context.Context, projectName string) ([]container.Summary, error) {
	filter := filters.NewArgs()
	filter.Add("label", fmt.Sprintf("%s=true", LabelManaged))
	filter.Add("label", fmt.Sprintf("%s=%s", LabelProject, projectName))

	return c.cli.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filter,
	})
}

// GetContainerLabels returns labels for a specific container
func (c *Client) GetContainerLabels(ctx context.Context, containerID string) (map[string]string, error) {
	inspect, err := c.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}
	return inspect.Config.Labels, nil
}

// RemoveContainer removes a container by ID
func (c *Client) RemoveContainer(ctx context.Context, containerID string, force bool) error {
	options := container.RemoveOptions{
		Force:         force,
		RemoveVolumes: false,
	}
	return c.cli.ContainerRemove(ctx, containerID, options)
}
