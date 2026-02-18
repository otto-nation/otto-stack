package docker

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/system"
	"github.com/docker/docker/api/types/volume"
	"github.com/otto-nation/otto-stack/test/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetDockerClient(t *testing.T) {
	mock := &testhelpers.MockDockerClient{}
	client := &Client{cli: mock}

	result := client.GetDockerClient()
	assert.Equal(t, mock, result)
}

func TestClient_RemoveResources(t *testing.T) {
	ctx := context.Background()
	mock := &testhelpers.MockDockerClient{}

	mock.ContainerListFunc = func(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
		return []types.Container{{ID: "test-id"}}, nil
	}
	mock.ContainerRemoveFunc = func(ctx context.Context, containerID string, options container.RemoveOptions) error {
		return nil
	}

	client := NewClientWithDependencies(mock, nil, nil)

	err := client.RemoveResources(ctx, "container", "test-project")
	assert.NoError(t, err)
}

func TestDockerClientAdapter_Methods(t *testing.T) {
	ctx := context.Background()
	mock := &testhelpers.MockDockerClient{}

	t.Run("ContainerList", func(t *testing.T) {
		called := false
		mock.ContainerListFunc = func(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
			called = true
			return []types.Container{}, nil
		}
		_, err := mock.ContainerList(ctx, container.ListOptions{})
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("ContainerRemove", func(t *testing.T) {
		called := false
		mock.ContainerRemoveFunc = func(ctx context.Context, containerID string, options container.RemoveOptions) error {
			called = true
			return nil
		}
		err := mock.ContainerRemove(ctx, "test", container.RemoveOptions{})
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("ContainerInspect", func(t *testing.T) {
		called := false
		mock.ContainerInspectFunc = func(ctx context.Context, containerID string) (types.ContainerJSON, error) {
			called = true
			return types.ContainerJSON{}, nil
		}
		_, err := mock.ContainerInspect(ctx, "test")
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("ContainerStop", func(t *testing.T) {
		called := false
		mock.ContainerStopFunc = func(ctx context.Context, containerID string, options container.StopOptions) error {
			called = true
			return nil
		}
		err := mock.ContainerStop(ctx, "test", container.StopOptions{})
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("ContainerRestart", func(t *testing.T) {
		called := false
		mock.ContainerRestartFunc = func(ctx context.Context, containerID string, options container.StopOptions) error {
			called = true
			return nil
		}
		err := mock.ContainerRestart(ctx, "test", container.StopOptions{})
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("ContainerLogs", func(t *testing.T) {
		called := false
		mock.ContainerLogsFunc = func(ctx context.Context, container string, options container.LogsOptions) (io.ReadCloser, error) {
			called = true
			return io.NopCloser(strings.NewReader("")), nil
		}
		_, err := mock.ContainerLogs(ctx, "test", container.LogsOptions{})
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("VolumeList", func(t *testing.T) {
		called := false
		mock.VolumeListFunc = func(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error) {
			called = true
			return volume.ListResponse{}, nil
		}
		_, err := mock.VolumeList(ctx, volume.ListOptions{})
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("VolumeRemove", func(t *testing.T) {
		called := false
		mock.VolumeRemoveFunc = func(ctx context.Context, volumeID string, force bool) error {
			called = true
			return nil
		}
		err := mock.VolumeRemove(ctx, "test", false)
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("NetworkList", func(t *testing.T) {
		called := false
		mock.NetworkListFunc = func(ctx context.Context, options network.ListOptions) ([]network.Summary, error) {
			called = true
			return []network.Summary{}, nil
		}
		_, err := mock.NetworkList(ctx, network.ListOptions{})
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("NetworkRemove", func(t *testing.T) {
		called := false
		mock.NetworkRemoveFunc = func(ctx context.Context, networkID string) error {
			called = true
			return nil
		}
		err := mock.NetworkRemove(ctx, "test")
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("ImageList", func(t *testing.T) {
		called := false
		mock.ImageListFunc = func(ctx context.Context, options image.ListOptions) ([]image.Summary, error) {
			called = true
			return []image.Summary{}, nil
		}
		_, err := mock.ImageList(ctx, image.ListOptions{})
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("ImageRemove", func(t *testing.T) {
		called := false
		mock.ImageRemoveFunc = func(ctx context.Context, imageID string, options image.RemoveOptions) ([]image.DeleteResponse, error) {
			called = true
			return []image.DeleteResponse{}, nil
		}
		_, err := mock.ImageRemove(ctx, "test", image.RemoveOptions{})
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("Info", func(t *testing.T) {
		called := false
		mock.InfoFunc = func(ctx context.Context) (system.Info, error) {
			called = true
			return system.Info{}, nil
		}
		_, err := mock.Info(ctx)
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("Ping", func(t *testing.T) {
		called := false
		mock.PingFunc = func(ctx context.Context) (types.Ping, error) {
			called = true
			return types.Ping{}, nil
		}
		_, err := mock.Ping(ctx)
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("Close", func(t *testing.T) {
		called := false
		mock.CloseFunc = func() error {
			called = true
			return nil
		}
		err := mock.Close()
		assert.NoError(t, err)
		assert.True(t, called)
	})
}
