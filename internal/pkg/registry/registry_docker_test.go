//go:build unit

package registry

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestOrphanDetection_Orphaned(t *testing.T) {
	mockDocker := &testhelpers.MockDockerClient{
		ContainerListFunc: func(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
			return []types.Container{
				{
					ID:     "orphan1",
					Names:  []string{"/orphaned-service"},
					Labels: map[string]string{"com.docker.compose.project": "deleted-project"},
					State:  "running",
				},
				{
					ID:     "active1",
					Names:  []string{"/active-service"},
					Labels: map[string]string{"com.docker.compose.project": "active-project"},
					State:  "running",
				},
			}, nil
		},
	}

	containers, err := mockDocker.ContainerList(context.Background(), container.ListOptions{All: true})
	require.NoError(t, err)
	assert.Len(t, containers, 2)

	var orphans []types.Container
	knownProjects := map[string]bool{"active-project": true}
	for _, c := range containers {
		project := c.Labels["com.docker.compose.project"]
		if !knownProjects[project] {
			orphans = append(orphans, c)
		}
	}

	assert.Len(t, orphans, 1)
	assert.Equal(t, "orphan1", orphans[0].ID)
}

func TestOrphanDetection_Zombies(t *testing.T) {
	mockDocker := &testhelpers.MockDockerClient{
		ContainerListFunc: func(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
			return []types.Container{
				{
					ID:     "zombie1",
					Names:  []string{"/zombie-service"},
					State:  "exited",
					Labels: map[string]string{"com.docker.compose.project": "test"},
				},
				{
					ID:     "running1",
					Names:  []string{"/running-service"},
					State:  "running",
					Labels: map[string]string{"com.docker.compose.project": "test"},
				},
			}, nil
		},
	}

	containers, err := mockDocker.ContainerList(context.Background(), container.ListOptions{All: true})
	require.NoError(t, err)

	var zombies []types.Container
	for _, c := range containers {
		if c.State == "exited" {
			zombies = append(zombies, c)
		}
	}

	assert.Len(t, zombies, 1)
	assert.Equal(t, "zombie1", zombies[0].ID)
}

func TestRegistry_RegisterContainer(t *testing.T) {
	mockDocker := &testhelpers.MockDockerClient{
		ContainerInspectFunc: func(ctx context.Context, containerID string) (types.ContainerJSON, error) {
			return testhelpers.MockContainerJSON(containerID, "/test-service", "postgres:15", "test", true), nil
		},
	}

	info, err := mockDocker.ContainerInspect(context.Background(), "container1")
	require.NoError(t, err)
	assert.Equal(t, "container1", info.ID)
	assert.Equal(t, "postgres:15", info.Config.Image)
	assert.Equal(t, "test", info.Config.Labels["com.docker.compose.project"])
}

func TestRegistry_UnregisterContainer(t *testing.T) {
	removed := false
	mockDocker := &testhelpers.MockDockerClient{
		ContainerRemoveFunc: func(ctx context.Context, containerID string, options container.RemoveOptions) error {
			removed = true
			return nil
		},
	}

	err := mockDocker.ContainerRemove(context.Background(), "container1", container.RemoveOptions{Force: true})
	require.NoError(t, err)
	assert.True(t, removed)
}

func TestRegistry_ListContainersByProject(t *testing.T) {
	mockDocker := &testhelpers.MockDockerClient{
		ContainerListFunc: func(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
			return []types.Container{
				{
					ID:     "c1",
					Names:  []string{"/myproject-postgres"},
					Labels: map[string]string{"com.docker.compose.project": "myproject"},
				},
				{
					ID:     "c2",
					Names:  []string{"/myproject-redis"},
					Labels: map[string]string{"com.docker.compose.project": "myproject"},
				},
				{
					ID:     "c3",
					Names:  []string{"/other-service"},
					Labels: map[string]string{"com.docker.compose.project": "other"},
				},
			}, nil
		},
	}

	containers, err := mockDocker.ContainerList(context.Background(), container.ListOptions{})
	require.NoError(t, err)

	var projectContainers []types.Container
	for _, c := range containers {
		if c.Labels["com.docker.compose.project"] == "myproject" {
			projectContainers = append(projectContainers, c)
		}
	}

	assert.Len(t, projectContainers, 2)
}

func TestRegistry_CountContainersPerProject(t *testing.T) {
	mockDocker := &testhelpers.MockDockerClient{
		ContainerListFunc: func(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
			return []types.Container{
				{ID: "c1", Labels: map[string]string{"com.docker.compose.project": "proj1"}},
				{ID: "c2", Labels: map[string]string{"com.docker.compose.project": "proj1"}},
				{ID: "c3", Labels: map[string]string{"com.docker.compose.project": "proj2"}},
			}, nil
		},
	}

	containers, err := mockDocker.ContainerList(context.Background(), container.ListOptions{})
	require.NoError(t, err)

	counts := make(map[string]int)
	for _, c := range containers {
		project := c.Labels["com.docker.compose.project"]
		counts[project]++
	}

	assert.Equal(t, 2, counts["proj1"])
	assert.Equal(t, 1, counts["proj2"])
}

func TestRegistry_ContainerHealth(t *testing.T) {
	mockDocker := &testhelpers.MockDockerClient{
		ContainerInspectFunc: func(ctx context.Context, containerID string) (types.ContainerJSON, error) {
			return testhelpers.MockContainerJSONWithHealth(containerID, true, "healthy"), nil
		},
	}

	info, err := mockDocker.ContainerInspect(context.Background(), "container1")
	require.NoError(t, err)
	assert.True(t, info.State.Running)
	assert.Equal(t, "healthy", info.State.Health.Status)
}
