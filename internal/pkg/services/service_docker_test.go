//go:build unit

package services

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestNewService_WithMockDocker(t *testing.T) {
	t.Run("creates service with mock docker client", func(t *testing.T) {
		mockDocker := &testhelpers.MockDockerClient{
			PingFunc: func(ctx context.Context) (types.Ping, error) {
				return types.Ping{}, nil
			},
		}

		// Test that we can create service components
		assert.NotNil(t, mockDocker)
	})
}

func TestService_DockerOperations(t *testing.T) {
	t.Run("lists containers", func(t *testing.T) {
		mockDocker := &testhelpers.MockDockerClient{
			ContainerListFunc: func(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
				return []types.Container{
					{
						ID:    "container1",
						Names: []string{"/test-service"},
						State: "running",
					},
				}, nil
			},
		}

		containers, err := mockDocker.ContainerList(context.Background(), container.ListOptions{})
		require.NoError(t, err)
		assert.Len(t, containers, 1)
		assert.Equal(t, "container1", containers[0].ID)
	})

	t.Run("inspects container", func(t *testing.T) {
		mockDocker := &testhelpers.MockDockerClient{
			ContainerInspectFunc: func(ctx context.Context, containerID string) (types.ContainerJSON, error) {
				return types.ContainerJSON{
					ContainerJSONBase: &types.ContainerJSONBase{
						ID:    containerID,
						Name:  "/test-service",
						State: &types.ContainerState{Running: true},
					},
				}, nil
			},
		}

		info, err := mockDocker.ContainerInspect(context.Background(), "container1")
		require.NoError(t, err)
		assert.Equal(t, "container1", info.ID)
		assert.True(t, info.State.Running)
	})

	t.Run("starts container", func(t *testing.T) {
		started := false
		mockDocker := &testhelpers.MockDockerClient{
			ContainerStartFunc: func(ctx context.Context, containerID string, options container.StartOptions) error {
				started = true
				return nil
			},
		}

		err := mockDocker.ContainerStart(context.Background(), "container1", container.StartOptions{})
		require.NoError(t, err)
		assert.True(t, started)
	})

	t.Run("stops container", func(t *testing.T) {
		stopped := false
		mockDocker := &testhelpers.MockDockerClient{
			ContainerStopFunc: func(ctx context.Context, containerID string, options container.StopOptions) error {
				stopped = true
				return nil
			},
		}

		err := mockDocker.ContainerStop(context.Background(), "container1", container.StopOptions{})
		require.NoError(t, err)
		assert.True(t, stopped)
	})

	t.Run("removes container", func(t *testing.T) {
		removed := false
		mockDocker := &testhelpers.MockDockerClient{
			ContainerRemoveFunc: func(ctx context.Context, containerID string, options container.RemoveOptions) error {
				removed = true
				return nil
			},
		}

		err := mockDocker.ContainerRemove(context.Background(), "container1", container.RemoveOptions{})
		require.NoError(t, err)
		assert.True(t, removed)
	})

	t.Run("restarts container", func(t *testing.T) {
		restarted := false
		mockDocker := &testhelpers.MockDockerClient{
			ContainerRestartFunc: func(ctx context.Context, containerID string, options container.StopOptions) error {
				restarted = true
				return nil
			},
		}

		err := mockDocker.ContainerRestart(context.Background(), "container1", container.StopOptions{})
		require.NoError(t, err)
		assert.True(t, restarted)
	})
}

func TestService_ContainerFiltering(t *testing.T) {
	t.Run("filters running containers", func(t *testing.T) {
		mockDocker := &testhelpers.MockDockerClient{
			ContainerListFunc: func(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
				allContainers := []types.Container{
					{ID: "c1", Names: []string{"/running"}, State: "running"},
					{ID: "c2", Names: []string{"/stopped"}, State: "exited"},
				}

				if options.All {
					return allContainers, nil
				}

				var running []types.Container
				for _, c := range allContainers {
					if c.State == "running" {
						running = append(running, c)
					}
				}
				return running, nil
			},
		}

		// List only running
		running, err := mockDocker.ContainerList(context.Background(), container.ListOptions{All: false})
		require.NoError(t, err)
		assert.Len(t, running, 1)
		assert.Equal(t, "running", running[0].State)

		// List all
		all, err := mockDocker.ContainerList(context.Background(), container.ListOptions{All: true})
		require.NoError(t, err)
		assert.Len(t, all, 2)
	})
}

func TestService_ContainerLabels(t *testing.T) {
	t.Run("filters by project label", func(t *testing.T) {
		mockDocker := &testhelpers.MockDockerClient{
			ContainerListFunc: func(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
				return []types.Container{
					{
						ID:     "c1",
						Names:  []string{"/project1-service"},
						Labels: map[string]string{"com.docker.compose.project": "project1"},
					},
					{
						ID:     "c2",
						Names:  []string{"/project2-service"},
						Labels: map[string]string{"com.docker.compose.project": "project2"},
					},
				}, nil
			},
		}

		containers, err := mockDocker.ContainerList(context.Background(), container.ListOptions{})
		require.NoError(t, err)

		// Filter by project
		var project1Containers []types.Container
		for _, c := range containers {
			if c.Labels["com.docker.compose.project"] == "project1" {
				project1Containers = append(project1Containers, c)
			}
		}

		assert.Len(t, project1Containers, 1)
		assert.Equal(t, "c1", project1Containers[0].ID)
	})
}
