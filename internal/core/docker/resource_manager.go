package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
)

type ResourceManager struct {
	client *Client
}

func NewResourceManager(client *Client) *ResourceManager {
	return &ResourceManager{client: client}
}

func (rm *ResourceManager) List(ctx context.Context, resourceType ResourceType, filter filters.Args) ([]string, error) {
	switch resourceType {
	case ResourceContainer:
		return rm.listContainers(ctx, filter)
	case ResourceVolume:
		return rm.listVolumes(ctx, filter)
	case ResourceNetwork:
		return rm.listNetworks(ctx, filter)
	case ResourceImage:
		return rm.listImages(ctx, filter)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

func (rm *ResourceManager) Remove(ctx context.Context, resourceType ResourceType, names []string) error {
	switch resourceType {
	case ResourceContainer:
		return rm.removeContainers(ctx, names)
	case ResourceVolume:
		return rm.removeVolumes(ctx, names)
	case ResourceNetwork:
		return rm.removeNetworks(ctx, names)
	case ResourceImage:
		return rm.removeImages(ctx, names)
	default:
		return fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

func (rm *ResourceManager) listContainers(ctx context.Context, filter filters.Args) ([]string, error) {
	containers, err := rm.client.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filter})
	if err != nil {
		return nil, err
	}
	names := make([]string, len(containers))
	for i, c := range containers {
		names[i] = c.ID
	}
	return names, nil
}

func (rm *ResourceManager) listVolumes(ctx context.Context, filter filters.Args) ([]string, error) {
	volumes, err := rm.client.cli.VolumeList(ctx, volume.ListOptions{Filters: filter})
	if err != nil {
		return nil, err
	}
	names := make([]string, len(volumes.Volumes))
	for i, v := range volumes.Volumes {
		names[i] = v.Name
	}
	return names, nil
}

func (rm *ResourceManager) listNetworks(ctx context.Context, filter filters.Args) ([]string, error) {
	networks, err := rm.client.cli.NetworkList(ctx, network.ListOptions{Filters: filter})
	if err != nil {
		return nil, err
	}
	names := make([]string, len(networks))
	for i, n := range networks {
		names[i] = n.Name
	}
	return names, nil
}

func (rm *ResourceManager) listImages(ctx context.Context, filter filters.Args) ([]string, error) {
	images, err := rm.client.cli.ImageList(ctx, image.ListOptions{Filters: filter})
	if err != nil {
		return nil, err
	}
	var names []string
	for _, img := range images {
		names = append(names, img.RepoTags...)
	}
	return names, nil
}

func (rm *ResourceManager) removeContainers(ctx context.Context, ids []string) error {
	return rm.removeResources(ids, "container", func(id string) error {
		return rm.client.cli.ContainerRemove(ctx, id, container.RemoveOptions{Force: true})
	})
}

func (rm *ResourceManager) removeVolumes(ctx context.Context, names []string) error {
	return rm.removeResources(names, "volume", func(name string) error {
		return rm.client.cli.VolumeRemove(ctx, name, false)
	})
}

func (rm *ResourceManager) removeNetworks(ctx context.Context, names []string) error {
	return rm.removeResources(names, "network", func(name string) error {
		return rm.client.cli.NetworkRemove(ctx, name)
	})
}

func (rm *ResourceManager) removeImages(ctx context.Context, names []string) error {
	return rm.removeResources(names, "image", func(name string) error {
		_, err := rm.client.cli.ImageRemove(ctx, name, image.RemoveOptions{Force: true})
		return err
	})
}

//nolint:unparam // ctx is used in closures passed as removeFn
func (rm *ResourceManager) removeResources(names []string, resourceType string, removeFn func(string) error) error {
	for _, name := range names {
		if err := removeFn(name); err != nil {
			rm.client.logger.Error("Failed to remove "+resourceType, "name", name, "error", err)
		}
	}
	return nil
}
