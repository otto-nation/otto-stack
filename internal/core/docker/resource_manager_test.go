package docker

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/otto-nation/otto-stack/test/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestResourceManager_ListContainers(t *testing.T) {
	ctx := context.Background()
	mock := &testhelpers.MockDockerClient{}

	mock.ContainerListFunc = func(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
		return []types.Container{
			{ID: "container1"},
			{ID: "container2"},
		}, nil
	}

	client := NewClientWithDependencies(mock, nil, nil)
	rm := NewResourceManager(client)

	names, err := rm.listContainers(ctx, NewProjectFilter("test"))
	assert.NoError(t, err)
	assert.Len(t, names, 2)
	assert.Contains(t, names, "container1")
	assert.Contains(t, names, "container2")
}

func TestResourceManager_ListVolumes(t *testing.T) {
	ctx := context.Background()
	mock := &testhelpers.MockDockerClient{}

	mock.VolumeListFunc = func(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error) {
		return volume.ListResponse{
			Volumes: []*volume.Volume{
				{Name: "vol1"},
				{Name: "vol2"},
			},
		}, nil
	}

	client := NewClientWithDependencies(mock, nil, nil)
	rm := NewResourceManager(client)

	names, err := rm.listVolumes(ctx, NewProjectFilter("test"))
	assert.NoError(t, err)
	assert.Len(t, names, 2)
	assert.Contains(t, names, "vol1")
	assert.Contains(t, names, "vol2")
}

func TestResourceManager_ListNetworks(t *testing.T) {
	ctx := context.Background()
	mock := &testhelpers.MockDockerClient{}

	mock.NetworkListFunc = func(ctx context.Context, options network.ListOptions) ([]network.Summary, error) {
		return []network.Summary{
			{Name: "net1"},
			{Name: "net2"},
		}, nil
	}

	client := NewClientWithDependencies(mock, nil, nil)
	rm := NewResourceManager(client)

	names, err := rm.listNetworks(ctx, NewProjectFilter("test"))
	assert.NoError(t, err)
	assert.Len(t, names, 2)
	assert.Contains(t, names, "net1")
	assert.Contains(t, names, "net2")
}

func TestResourceManager_ListImages(t *testing.T) {
	ctx := context.Background()
	mock := &testhelpers.MockDockerClient{}

	mock.ImageListFunc = func(ctx context.Context, options image.ListOptions) ([]image.Summary, error) {
		return []image.Summary{
			{RepoTags: []string{"img1:latest", "img1:v1"}},
			{RepoTags: []string{"img2:latest"}},
		}, nil
	}

	client := NewClientWithDependencies(mock, nil, nil)
	rm := NewResourceManager(client)

	names, err := rm.listImages(ctx, NewProjectFilter("test"))
	assert.NoError(t, err)
	assert.Len(t, names, 3)
	assert.Contains(t, names, "img1:latest")
	assert.Contains(t, names, "img1:v1")
	assert.Contains(t, names, "img2:latest")
}

func TestResourceManager_RemoveVolumes(t *testing.T) {
	ctx := context.Background()
	mock := &testhelpers.MockDockerClient{}

	removed := []string{}
	mock.VolumeRemoveFunc = func(ctx context.Context, volumeID string, force bool) error {
		removed = append(removed, volumeID)
		return nil
	}

	client := NewClientWithDependencies(mock, nil, nil)
	rm := NewResourceManager(client)

	err := rm.removeVolumes(ctx, []string{"vol1", "vol2"})
	assert.NoError(t, err)
	assert.Len(t, removed, 2)
	assert.Contains(t, removed, "vol1")
	assert.Contains(t, removed, "vol2")
}

func TestResourceManager_RemoveNetworks(t *testing.T) {
	ctx := context.Background()
	mock := &testhelpers.MockDockerClient{}

	removed := []string{}
	mock.NetworkRemoveFunc = func(ctx context.Context, networkID string) error {
		removed = append(removed, networkID)
		return nil
	}

	client := NewClientWithDependencies(mock, nil, nil)
	rm := NewResourceManager(client)

	err := rm.removeNetworks(ctx, []string{"net1", "net2"})
	assert.NoError(t, err)
	assert.Len(t, removed, 2)
	assert.Contains(t, removed, "net1")
	assert.Contains(t, removed, "net2")
}

func TestResourceManager_RemoveImages(t *testing.T) {
	ctx := context.Background()
	mock := &testhelpers.MockDockerClient{}

	removed := []string{}
	mock.ImageRemoveFunc = func(ctx context.Context, imageID string, options image.RemoveOptions) ([]image.DeleteResponse, error) {
		removed = append(removed, imageID)
		return nil, nil
	}

	client := NewClientWithDependencies(mock, nil, nil)
	rm := NewResourceManager(client)

	err := rm.removeImages(ctx, []string{"img1", "img2"})
	assert.NoError(t, err)
	assert.Len(t, removed, 2)
	assert.Contains(t, removed, "img1")
	assert.Contains(t, removed, "img2")
}

func TestResourceManager_List(t *testing.T) {
	ctx := context.Background()
	mock := &testhelpers.MockDockerClient{}

	mock.VolumeListFunc = func(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error) {
		return volume.ListResponse{Volumes: []*volume.Volume{{Name: "vol1"}}}, nil
	}

	client := NewClientWithDependencies(mock, nil, nil)
	rm := NewResourceManager(client)

	t.Run("list volumes", func(t *testing.T) {
		names, err := rm.List(ctx, ResourceVolume, NewProjectFilter("test"))
		assert.NoError(t, err)
		assert.Len(t, names, 1)
		assert.Contains(t, names, "vol1")
	})

	t.Run("unsupported type", func(t *testing.T) {
		_, err := rm.List(ctx, "unsupported", NewProjectFilter("test"))
		assert.Error(t, err)
	})
}

func TestResourceManager_Remove(t *testing.T) {
	ctx := context.Background()
	mock := &testhelpers.MockDockerClient{}

	removed := []string{}
	mock.VolumeRemoveFunc = func(ctx context.Context, volumeID string, force bool) error {
		removed = append(removed, volumeID)
		return nil
	}

	client := NewClientWithDependencies(mock, nil, nil)
	rm := NewResourceManager(client)

	t.Run("remove volumes", func(t *testing.T) {
		err := rm.Remove(ctx, ResourceVolume, []string{"vol1"})
		assert.NoError(t, err)
		assert.Len(t, removed, 1)
	})

	t.Run("unsupported type", func(t *testing.T) {
		err := rm.Remove(ctx, "unsupported", []string{"vol1"})
		assert.Error(t, err)
	})
}

func TestResourceManager_RemoveVolumes_WithError(t *testing.T) {
	ctx := context.Background()
	mock := &testhelpers.MockDockerClient{}

	mock.VolumeRemoveFunc = func(ctx context.Context, volumeID string, force bool) error {
		if volumeID == "vol-error" {
			return assert.AnError
		}
		return nil
	}

	client := NewClientWithDependencies(mock, nil, nil)
	rm := NewResourceManager(client)

	// Should not return error even if individual removes fail (errors are logged)
	err := rm.removeVolumes(ctx, []string{"vol1", "vol2"})
	assert.NoError(t, err)
}

func TestResourceManager_ListContainers_Error(t *testing.T) {
	ctx := context.Background()
	mock := &testhelpers.MockDockerClient{}

	mock.ContainerListFunc = func(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
		return nil, assert.AnError
	}

	client := NewClientWithDependencies(mock, nil, nil)
	rm := NewResourceManager(client)

	_, err := rm.listContainers(ctx, NewProjectFilter("test"))
	assert.Error(t, err)
}

func TestResourceManager_ListVolumes_Error(t *testing.T) {
	ctx := context.Background()
	mock := &testhelpers.MockDockerClient{}

	mock.VolumeListFunc = func(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error) {
		return volume.ListResponse{}, assert.AnError
	}

	client := NewClientWithDependencies(mock, nil, nil)
	rm := NewResourceManager(client)

	_, err := rm.listVolumes(ctx, NewProjectFilter("test"))
	assert.Error(t, err)
}

func TestResourceManager_ListNetworks_Error(t *testing.T) {
	ctx := context.Background()
	mock := &testhelpers.MockDockerClient{}

	mock.NetworkListFunc = func(ctx context.Context, options network.ListOptions) ([]network.Summary, error) {
		return nil, assert.AnError
	}

	client := NewClientWithDependencies(mock, nil, nil)
	rm := NewResourceManager(client)

	_, err := rm.listNetworks(ctx, NewProjectFilter("test"))
	assert.Error(t, err)
}

func TestResourceManager_ListImages_Error(t *testing.T) {
	ctx := context.Background()
	mock := &testhelpers.MockDockerClient{}

	mock.ImageListFunc = func(ctx context.Context, options image.ListOptions) ([]image.Summary, error) {
		return nil, assert.AnError
	}

	client := NewClientWithDependencies(mock, nil, nil)
	rm := NewResourceManager(client)

	_, err := rm.listImages(ctx, NewProjectFilter("test"))
	assert.Error(t, err)
}
