//go:build unit

package docker

import (
	"context"
	"errors"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestClient_ListResources_Unit(t *testing.T) {
	mockDocker := &testhelpers.MockDockerClient{
		ContainerListFunc: func(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
			return []types.Container{
				{ID: "container1", Names: []string{"/test-container"}},
			}, nil
		},
	}

	client := NewClientWithDependencies(mockDocker, nil, testhelpers.MockLogger())

	ctx := context.Background()
	resources, err := client.ListResources(ctx, ResourceContainer, "test-project")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(resources) != 1 {
		t.Errorf("Expected 1 resource, got %d", len(resources))
	}
}

func TestClient_RemoveContainer_Unit(t *testing.T) {
	removeCalled := false
	mockDocker := &testhelpers.MockDockerClient{
		ContainerRemoveFunc: func(ctx context.Context, containerID string, options container.RemoveOptions) error {
			removeCalled = true
			if containerID == "non-existent" {
				return errors.New("container not found")
			}
			return nil
		},
	}

	client := NewClientWithDependencies(mockDocker, nil, testhelpers.MockLogger())

	ctx := context.Background()
	err := client.RemoveContainer(ctx, "non-existent", false)

	if err == nil {
		t.Error("Expected error for non-existent container")
	}

	if !removeCalled {
		t.Error("Expected ContainerRemove to be called")
	}
}

func TestClient_GetServiceStatus_Unit(t *testing.T) {
	mockDocker := &testhelpers.MockDockerClient{
		ContainerListFunc: func(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
			return []types.Container{
				{
					ID:    "web1",
					Names: []string{"/test-project-web-1"},
					State: "running",
					Labels: map[string]string{
						"com.docker.compose.project": "test-project",
						"com.docker.compose.service": "web",
					},
				},
			}, nil
		},
		ContainerInspectFunc: func(ctx context.Context, containerID string) (types.ContainerJSON, error) {
			return types.ContainerJSON{
				ContainerJSONBase: &types.ContainerJSONBase{
					ID:      containerID,
					Created: "2024-01-01T00:00:00Z",
					State: &types.ContainerState{
						StartedAt: "2024-01-01T00:00:01Z",
					},
				},
				NetworkSettings: &types.NetworkSettings{},
			}, nil
		},
	}

	client := NewClientWithDependencies(mockDocker, nil, testhelpers.MockLogger())

	ctx := context.Background()
	services := []string{"web", "db"}
	statuses, err := client.GetServiceStatus(ctx, "test-project", services)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(statuses) != len(services) {
		t.Errorf("Expected %d statuses, got %d", len(services), len(statuses))
	}
}
