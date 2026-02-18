//go:build unit

package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"io"

	composetypes "github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/system"
	"github.com/docker/docker/api/types/volume"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/test/fixtures"
	"github.com/otto-nation/otto-stack/test/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockResolver struct{}

func (m *mockResolver) ResolveUpOptions(_ []string, _ []servicetypes.ServiceConfig, base docker.UpOptions) docker.UpOptions {
	return base
}

func (m *mockResolver) ResolveDownOptions(_ []string, _ []servicetypes.ServiceConfig, base docker.DownOptions) docker.DownOptions {
	return base
}

func (m *mockResolver) ResolveStopOptions(_ []string, _ []servicetypes.ServiceConfig, base docker.StopOptions) docker.StopOptions {
	return base
}

type mockDockerClient struct{}

func (m *mockDockerClient) ContainerList(ctx context.Context, options container.ListOptions) ([]container.Summary, error) {
	return []container.Summary{}, nil
}

func (m *mockDockerClient) ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error {
	return nil
}

func (m *mockDockerClient) ContainerInspect(ctx context.Context, containerID string) (container.InspectResponse, error) {
	return container.InspectResponse{}, nil
}

func (m *mockDockerClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, containerName string) (container.CreateResponse, error) {
	return container.CreateResponse{}, nil
}

func (m *mockDockerClient) ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error {
	return nil
}

func (m *mockDockerClient) ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error {
	return nil
}

func (m *mockDockerClient) ContainerRestart(ctx context.Context, containerID string, options container.StopOptions) error {
	return nil
}

func (m *mockDockerClient) ContainerWait(ctx context.Context, containerID string, condition container.WaitCondition) (<-chan container.WaitResponse, <-chan error) {
	respChan := make(chan container.WaitResponse, 1)
	errChan := make(chan error, 1)
	close(respChan)
	close(errChan)
	return respChan, errChan
}

func (m *mockDockerClient) ContainerLogs(ctx context.Context, container string, options container.LogsOptions) (io.ReadCloser, error) {
	return io.NopCloser(nil), nil
}

func (m *mockDockerClient) VolumeList(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error) {
	return volume.ListResponse{}, nil
}

func (m *mockDockerClient) VolumeRemove(ctx context.Context, volumeID string, force bool) error {
	return nil
}

func (m *mockDockerClient) NetworkList(ctx context.Context, options network.ListOptions) ([]network.Summary, error) {
	return []network.Summary{}, nil
}

func (m *mockDockerClient) NetworkRemove(ctx context.Context, networkID string) error {
	return nil
}

func (m *mockDockerClient) ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error) {
	return []image.Summary{}, nil
}

func (m *mockDockerClient) ImageRemove(ctx context.Context, imageID string, options image.RemoveOptions) ([]image.DeleteResponse, error) {
	return []image.DeleteResponse{}, nil
}

func (m *mockDockerClient) Info(ctx context.Context) (system.Info, error) {
	return system.Info{}, nil
}

func (m *mockDockerClient) Ping(ctx context.Context) (types.Ping, error) {
	return types.Ping{}, nil
}

func (m *mockDockerClient) Close() error {
	return nil
}

func setupTestEnv(t *testing.T) (string, func()) {
	tmpDir := t.TempDir()
	ottoDir := filepath.Join(tmpDir, ".otto-stack")
	require.NoError(t, os.MkdirAll(ottoDir, 0755))

	composeFile := filepath.Join(ottoDir, "docker-compose.yml")
	composeContent := `version: "3.8"
services:
  postgres:
    image: postgres:latest
`
	require.NoError(t, os.WriteFile(composeFile, []byte(composeContent), 0644))

	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)

	return tmpDir, func() { os.Chdir(oldDir) }
}

func TestService_StartWithMocks(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	mockCompose := testhelpers.NewMockCompose().WithUpSuccess()
	mockLoader := testhelpers.NewMockProjectLoader().WithLoadSuccess("test-project")

	resolver := &mockResolver{}
	mockDocker := &mockDockerClient{}
	dockerClient := docker.NewClientWithDependencies(mockDocker, nil, nil)
	service := NewServiceWithClient(mockCompose, resolver, mockLoader, dockerClient)

	req := StartRequest{
		Project:        "test-project",
		ServiceConfigs: []servicetypes.ServiceConfig{fixtures.NewServiceConfig(ServicePostgres).Build()},
		Timeout:        10 * time.Second,
	}

	err := service.Start(context.Background(), req)
	assert.NoError(t, err)
}

func TestService_StopWithMocks(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	mockCompose := testhelpers.NewMockCompose().WithStopSuccess()
	mockLoader := testhelpers.NewMockProjectLoader().WithLoadSuccess("test-project")

	resolver := &mockResolver{}
	mockDocker := &mockDockerClient{}
	dockerClient := docker.NewClientWithDependencies(mockDocker, nil, nil)
	service := NewServiceWithClient(mockCompose, resolver, mockLoader, dockerClient)

	req := StopRequest{
		Project:        "test-project",
		ServiceConfigs: []servicetypes.ServiceConfig{fixtures.NewServiceConfig(ServicePostgres).Build()},
		Timeout:        10 * time.Second,
	}

	err := service.Stop(context.Background(), req)
	assert.NoError(t, err)
}

func TestService_StatusWithMocks(t *testing.T) {
	mockCompose := testhelpers.NewMockCompose().WithPsResult([]api.ContainerSummary{
		{Name: ServicePostgres, State: "running"},
	})

	mockLoader := testhelpers.NewMockProjectLoader()
	resolver := &mockResolver{}
	mockDocker := &mockDockerClient{}
	dockerClient := docker.NewClientWithDependencies(mockDocker, nil, nil)
	service := NewServiceWithClient(mockCompose, resolver, mockLoader, dockerClient)

	req := StatusRequest{
		Project:  "test-project",
		Services: []string{ServicePostgres},
	}

	statuses, err := service.Status(context.Background(), req)
	assert.NoError(t, err)
	assert.Len(t, statuses, 1)
}

func TestService_LogsWithMocks(t *testing.T) {
	mockCompose := testhelpers.NewMockCompose().WithLogsSuccess()
	mockLoader := testhelpers.NewMockProjectLoader()
	resolver := &mockResolver{}
	mockDocker := &mockDockerClient{}
	dockerClient := docker.NewClientWithDependencies(mockDocker, nil, nil)
	service := NewServiceWithClient(mockCompose, resolver, mockLoader, dockerClient)

	req := LogRequest{
		Project:        "test-project",
		ServiceConfigs: []servicetypes.ServiceConfig{fixtures.NewServiceConfig(ServicePostgres).Build()},
		Follow:         false,
	}

	err := service.Logs(context.Background(), req)
	assert.NoError(t, err)
}

func TestService_ExecWithMocks(t *testing.T) {
	mockCompose := testhelpers.NewMockCompose().WithExecSuccess()
	mockLoader := testhelpers.NewMockProjectLoader()
	resolver := &mockResolver{}
	mockDocker := &mockDockerClient{}
	dockerClient := docker.NewClientWithDependencies(mockDocker, nil, nil)
	service := NewServiceWithClient(mockCompose, resolver, mockLoader, dockerClient)

	req := ExecRequest{
		Project: "test-project",
		Service: ServicePostgres,
		Command: []string{"psql", "--version"},
	}

	err := service.Exec(context.Background(), req)
	assert.NoError(t, err)
}

func TestService_CleanupWithMocks(t *testing.T) {
	mockCompose := testhelpers.NewMockCompose().WithDownSuccess()
	mockLoader := testhelpers.NewMockProjectLoader().WithLoadSuccess("test-project")

	resolver := &mockResolver{}
	mockDocker := &mockDockerClient{}
	dockerClient := docker.NewClientWithDependencies(mockDocker, nil, nil)
	service := NewServiceWithClient(mockCompose, resolver, mockLoader, dockerClient)

	req := CleanupRequest{
		Project:       "test-project",
		RemoveVolumes: true,
		RemoveImages:  false,
	}

	err := service.Cleanup(context.Background(), req)
	assert.NoError(t, err)
}

func TestService_CheckDockerHealthWithMocks(t *testing.T) {
	mockCompose := testhelpers.NewMockCompose()
	mockLoader := testhelpers.NewMockProjectLoader()
	resolver := &mockResolver{}
	mockDocker := &mockDockerClient{}
	dockerClient := docker.NewClientWithDependencies(mockDocker, nil, nil)
	service := NewServiceWithClient(mockCompose, resolver, mockLoader, dockerClient)

	err := service.CheckDockerHealth(context.Background())
	assert.NoError(t, err)
}

func TestService_NewServiceErrorHandling(t *testing.T) {
	mockCompose := testhelpers.NewMockCompose()
	mockLoader := testhelpers.NewMockProjectLoader()
	resolver := &mockResolver{}

	service, err := NewService(mockCompose, resolver, mockLoader)
	assert.NoError(t, err)
	assert.NotNil(t, service)
	assert.NotNil(t, service.DockerClient)
}

func TestService_StartWithBuildFlag(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	t.Run("starts with build flag", func(t *testing.T) {
		buildCalled := false
		mockCompose := &testhelpers.MockComposeAPI{
			UpFunc: func(ctx context.Context, project *composetypes.Project, options api.UpOptions) error {
				buildCalled = true
				return nil
			},
		}

		mockLoader := testhelpers.NewMockProjectLoader().WithLoadSuccess("test-project")

		resolver := &mockResolver{}
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := StartRequest{
			Project: "test-project",
			ServiceConfigs: []servicetypes.ServiceConfig{
				{Name: "postgres"},
			},
			Build:   true,
			Timeout: 10 * time.Second,
		}

		err = service.Start(context.Background(), req)
		assert.NoError(t, err)
		assert.True(t, buildCalled)
	})

	t.Run("starts with force recreate", func(t *testing.T) {
		recreateCalled := false
		mockCompose := &testhelpers.MockComposeAPI{
			UpFunc: func(ctx context.Context, project *composetypes.Project, options api.UpOptions) error {
				recreateCalled = true
				return nil
			},
		}

		mockLoader := testhelpers.NewMockProjectLoader().WithLoadSuccess("test-project")

		resolver := &mockResolver{}
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := StartRequest{
			Project: "test-project",
			ServiceConfigs: []servicetypes.ServiceConfig{
				{Name: "postgres"},
			},
			ForceRecreate: true,
			Timeout:       10 * time.Second,
		}

		err = service.Start(context.Background(), req)
		assert.NoError(t, err)
		assert.True(t, recreateCalled)
	})
}

func TestService_StartErrorPaths(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	t.Run("handles compose up error", func(t *testing.T) {
		mockCompose := testhelpers.NewMockCompose().WithUpError(assert.AnError)
		mockLoader := testhelpers.NewMockProjectLoader().WithLoadSuccess("test-project")

		resolver := &mockResolver{}
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := StartRequest{
			Project: "test-project",
			ServiceConfigs: []servicetypes.ServiceConfig{
				{Name: "postgres"},
			},
			Timeout: 10 * time.Second,
		}

		err = service.Start(context.Background(), req)
		assert.Error(t, err)
	})

	t.Run("handles project load error", func(t *testing.T) {
		mockCompose := testhelpers.NewMockCompose()

		mockLoader := &testhelpers.MockProjectLoader{
			LoadFunc: func(projectName string) (*composetypes.Project, error) {
				return nil, assert.AnError
			},
		}

		resolver := &mockResolver{}
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := StartRequest{
			Project: "test-project",
			ServiceConfigs: []servicetypes.ServiceConfig{
				{Name: "postgres"},
			},
			Timeout: 10 * time.Second,
		}

		err = service.Start(context.Background(), req)
		assert.Error(t, err)
	})
}

func TestService_StopErrorPaths(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	t.Run("handles stop error", func(t *testing.T) {
		mockCompose := testhelpers.NewMockCompose().WithStopError(assert.AnError)
		mockLoader := testhelpers.NewMockProjectLoader().WithLoadSuccess("test-project")

		resolver := &mockResolver{}
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := StopRequest{
			Project: "test-project",
			ServiceConfigs: []servicetypes.ServiceConfig{
				{Name: "postgres"},
			},
			Timeout: 10 * time.Second,
		}

		err = service.Stop(context.Background(), req)
		assert.Error(t, err)
	})

	t.Run("uses down when remove is true", func(t *testing.T) {
		downCalled := false
		mockCompose := &testhelpers.MockComposeAPI{
			DownFunc: func(ctx context.Context, projectName string, options api.DownOptions) error {
				downCalled = true
				return nil
			},
		}

		mockLoader := testhelpers.NewMockProjectLoader().WithLoadSuccess("test-project")

		resolver := &mockResolver{}
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := StopRequest{
			Project: "test-project",
			ServiceConfigs: []servicetypes.ServiceConfig{
				{Name: "postgres"},
			},
			Remove:  true,
			Timeout: 10 * time.Second,
		}

		err = service.Stop(context.Background(), req)
		assert.NoError(t, err)
		assert.True(t, downCalled)
	})
}

func TestService_GenerateComposeFile(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	t.Run("generates compose file successfully", func(t *testing.T) {
		mockCompose := testhelpers.NewMockCompose()
		mockLoader := testhelpers.NewMockProjectLoader()
		resolver := &mockResolver{}
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		configs := []servicetypes.ServiceConfig{
			{Name: "postgres", Description: "PostgreSQL"},
		}

		err = service.GenerateComposeFile("test-project", configs)
		assert.NoError(t, err)
	})
}

func TestService_LogsWithOptions(t *testing.T) {
	t.Run("logs with follow option", func(t *testing.T) {
		mockCompose := &testhelpers.MockComposeAPI{
			LogsFunc: func(ctx context.Context, projectName string, consumer api.LogConsumer, options api.LogOptions) error {
				assert.True(t, options.Follow)
				return nil
			},
		}

		mockLoader := testhelpers.NewMockProjectLoader()
		resolver := &mockResolver{}
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := LogRequest{
			Project: "test-project",
			ServiceConfigs: []servicetypes.ServiceConfig{
				{Name: "postgres"},
			},
			Follow: true,
		}

		err = service.Logs(context.Background(), req)
		assert.NoError(t, err)
	})

	t.Run("logs with timestamps", func(t *testing.T) {
		mockCompose := &testhelpers.MockComposeAPI{
			LogsFunc: func(ctx context.Context, projectName string, consumer api.LogConsumer, options api.LogOptions) error {
				assert.True(t, options.Timestamps)
				return nil
			},
		}

		mockLoader := testhelpers.NewMockProjectLoader()
		resolver := &mockResolver{}
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := LogRequest{
			Project: "test-project",
			ServiceConfigs: []servicetypes.ServiceConfig{
				{Name: "postgres"},
			},
			Timestamps: true,
		}

		err = service.Logs(context.Background(), req)
		assert.NoError(t, err)
	})
}

func TestService_ExecWithOptions(t *testing.T) {
	t.Run("exec with user option", func(t *testing.T) {
		mockCompose := &testhelpers.MockComposeAPI{
			ExecFunc: func(ctx context.Context, projectName string, options api.RunOptions) (int, error) {
				assert.Equal(t, "admin", options.User)
				return 0, nil
			},
		}

		mockLoader := testhelpers.NewMockProjectLoader()
		resolver := &mockResolver{}
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := ExecRequest{
			Project: "test-project",
			Service: "postgres",
			Command: []string{"psql"},
			User:    "admin",
		}

		err = service.Exec(context.Background(), req)
		assert.NoError(t, err)
	})
}
