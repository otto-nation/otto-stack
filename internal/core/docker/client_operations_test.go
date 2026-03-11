//go:build unit

package docker

import (
	"context"
	"errors"
	"testing"

	composetypes "github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/docker/api/types/container"
	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestClient_ListResources_Unit(t *testing.T) {
	mockDocker := &testhelpers.MockDockerClient{
		ContainerListFunc: func(ctx context.Context, options container.ListOptions) ([]container.Summary, error) {
			return []container.Summary{
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

func TestClient_ResolveServicesToStart_AllRunning(t *testing.T) {
	mockDocker := &testhelpers.MockDockerClient{
		ContainerInspectFunc: func(ctx context.Context, containerID string) (container.InspectResponse, error) {
			return testhelpers.MockContainerJSON(containerID, containerID, "redis:7-alpine", "shared", true), nil
		},
	}
	client := NewClientWithDependencies(mockDocker, nil, testhelpers.MockLogger())

	proj := &composetypes.Project{
		Name:     "shared",
		Services: composetypes.Services{"redis": {ContainerName: "otto-stack-redis"}},
	}

	toCreate, err := client.ResolveServicesToStart(context.Background(), proj)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(toCreate) != 0 {
		t.Errorf("Expected no services to create, got %v", toCreate)
	}
}

func TestClient_ResolveServicesToStart_NoneExist(t *testing.T) {
	mockDocker := &testhelpers.MockDockerClient{
		ContainerInspectFunc: func(ctx context.Context, containerID string) (container.InspectResponse, error) {
			return container.InspectResponse{}, errors.New("No such container")
		},
	}
	client := NewClientWithDependencies(mockDocker, nil, testhelpers.MockLogger())

	proj := &composetypes.Project{
		Name:     "shared",
		Services: composetypes.Services{"redis": {ContainerName: "otto-stack-redis"}},
	}

	toCreate, err := client.ResolveServicesToStart(context.Background(), proj)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(toCreate) != 1 || toCreate[0] != "redis" {
		t.Errorf("Expected [redis] to create, got %v", toCreate)
	}
}

func TestClient_ResolveServicesToStart_StoppedContainer(t *testing.T) {
	startCalled := false
	mockDocker := &testhelpers.MockDockerClient{
		ContainerInspectFunc: func(ctx context.Context, containerID string) (container.InspectResponse, error) {
			return testhelpers.MockContainerJSON(containerID, containerID, "redis:7-alpine", "shared", false), nil
		},
		ContainerStartFunc: func(ctx context.Context, containerID string, options container.StartOptions) error {
			startCalled = true
			return nil
		},
	}
	client := NewClientWithDependencies(mockDocker, nil, testhelpers.MockLogger())

	proj := &composetypes.Project{
		Name:     "shared",
		Services: composetypes.Services{"redis": {ContainerName: "otto-stack-redis"}},
	}

	toCreate, err := client.ResolveServicesToStart(context.Background(), proj)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(toCreate) != 0 {
		t.Errorf("Expected no services to create after direct start, got %v", toCreate)
	}
	if !startCalled {
		t.Error("Expected ContainerStart to be called for stopped container")
	}
}

func TestClient_ResolveServicesToStart_FallbackContainerName(t *testing.T) {
	inspectedName := ""
	mockDocker := &testhelpers.MockDockerClient{
		ContainerInspectFunc: func(ctx context.Context, containerID string) (container.InspectResponse, error) {
			inspectedName = containerID
			return container.InspectResponse{}, errors.New("No such container")
		},
	}
	client := NewClientWithDependencies(mockDocker, nil, testhelpers.MockLogger())

	// Service has no explicit ContainerName — should fall back to "{project}-{service}"
	proj := &composetypes.Project{
		Name:     "shared",
		Services: composetypes.Services{"postgres": {}},
	}

	toCreate, err := client.ResolveServicesToStart(context.Background(), proj)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if inspectedName != "shared-postgres" {
		t.Errorf("Expected container name 'shared-postgres', got %q", inspectedName)
	}
	if len(toCreate) != 1 || toCreate[0] != "postgres" {
		t.Errorf("Expected [postgres] to create, got %v", toCreate)
	}
}

func TestClient_ResolveServicesToStart_Mixed(t *testing.T) {
	mockDocker := &testhelpers.MockDockerClient{
		ContainerInspectFunc: func(ctx context.Context, containerID string) (container.InspectResponse, error) {
			if containerID == "otto-stack-redis" {
				return testhelpers.MockContainerJSON(containerID, containerID, "redis:7-alpine", "shared", true), nil
			}
			return container.InspectResponse{}, errors.New("No such container")
		},
	}
	client := NewClientWithDependencies(mockDocker, nil, testhelpers.MockLogger())

	proj := &composetypes.Project{
		Name: "shared",
		Services: composetypes.Services{
			"redis":    {ContainerName: "otto-stack-redis"},
			"postgres": {ContainerName: "otto-stack-postgres"},
		},
	}

	toCreate, err := client.ResolveServicesToStart(context.Background(), proj)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(toCreate) != 1 || toCreate[0] != "postgres" {
		t.Errorf("Expected only [postgres] to create, got %v", toCreate)
	}
}

func TestClient_GetServiceStatus_Unit(t *testing.T) {
	mockDocker := &testhelpers.MockDockerClient{
		ContainerListFunc: func(ctx context.Context, options container.ListOptions) ([]container.Summary, error) {
			return []container.Summary{
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
		ContainerInspectFunc: func(ctx context.Context, containerID string) (container.InspectResponse, error) {
			json := testhelpers.MockContainerJSON(containerID, "/test", "test:latest", "test-project", false)
			json.Created = "2024-01-01T00:00:00Z"
			json.State.StartedAt = "2024-01-01T00:00:01Z"
			json.NetworkSettings = &container.NetworkSettings{}
			return json, nil
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
