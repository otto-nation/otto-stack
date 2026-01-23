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
