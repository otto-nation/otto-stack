//go:build unit

package docker

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestClient_ListProjectContainers_Unit(t *testing.T) {
	mockDocker := &testhelpers.MockDockerClient{
		ContainerListFunc: func(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
			return []types.Container{
				{
					ID:    "c1",
					Names: []string{"/test-project-web-1"},
					Labels: map[string]string{
						"com.docker.compose.project": "test-project",
					},
				},
			}, nil
		},
	}

	client := NewClientWithDependencies(mockDocker, nil, testhelpers.MockLogger())

	ctx := context.Background()
	containers, err := client.ListProjectContainers(ctx, "test-project")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(containers) != 1 {
		t.Errorf("Expected 1 container, got %d", len(containers))
	}
}

func TestNewDefaultProjectLoader(t *testing.T) {
	loader, err := NewDefaultProjectLoader()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if loader == nil {
		t.Error("Expected non-nil loader")
	}
}
