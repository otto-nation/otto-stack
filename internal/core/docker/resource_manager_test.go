//go:build unit

package docker

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestResourceManager_List_Containers_Unit(t *testing.T) {
	mockDocker := &testhelpers.MockDockerClient{
		ContainerListFunc: func(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
			return []types.Container{{ID: "c1"}}, nil
		},
	}

	client := NewClientWithDependencies(mockDocker, nil, testhelpers.MockLogger())
	ctx := context.Background()
	filter := filters.NewArgs()

	containers, err := client.GetResources().List(ctx, ResourceContainer, filter)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(containers) != 1 {
		t.Errorf("Expected 1 container, got %d", len(containers))
	}
}

func TestResourceManager_Remove_Containers_Unit(t *testing.T) {
	containerRemoved := false

	mockDocker := &testhelpers.MockDockerClient{
		ContainerRemoveFunc: func(ctx context.Context, containerID string, options container.RemoveOptions) error {
			containerRemoved = true
			return nil
		},
	}

	client := NewClientWithDependencies(mockDocker, nil, testhelpers.MockLogger())
	ctx := context.Background()

	err := client.GetResources().Remove(ctx, ResourceContainer, []string{"c1"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !containerRemoved {
		t.Error("Expected container to be removed")
	}
}
