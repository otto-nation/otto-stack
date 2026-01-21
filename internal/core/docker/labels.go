package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

const (
	LabelOttoManaged     = "io.otto-stack.managed"
	LabelOttoProject     = "io.otto-stack.project"
	LabelOttoService     = "io.otto-stack.service"
	LabelOttoVersion     = "io.otto-stack.version"
	LabelOttoSharingMode = "io.otto-stack.sharing-mode"
	LabelOttoShared      = "io.otto-stack.shared"
)

// ListOttoContainers returns all containers managed by Otto Stack
func (c *Client) ListOttoContainers(ctx context.Context) ([]container.Summary, error) {
	filter := filters.NewArgs()
	filter.Add("label", fmt.Sprintf("%s=true", LabelOttoManaged))

	return c.cli.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filter,
	})
}

// ListProjectContainers returns containers for a specific project
func (c *Client) ListProjectContainers(ctx context.Context, projectName string) ([]container.Summary, error) {
	filter := filters.NewArgs()
	filter.Add("label", fmt.Sprintf("%s=true", LabelOttoManaged))
	filter.Add("label", fmt.Sprintf("%s=%s", LabelOttoProject, projectName))

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
