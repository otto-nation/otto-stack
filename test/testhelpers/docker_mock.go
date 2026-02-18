package testhelpers

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/system"
	"github.com/docker/docker/api/types/volume"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// MockDockerClient is a mock implementation of docker.DockerClient for testing
type MockDockerClient struct {
	ContainerListFunc    func(ctx context.Context, options container.ListOptions) ([]types.Container, error)
	ContainerRemoveFunc  func(ctx context.Context, containerID string, options container.RemoveOptions) error
	ContainerInspectFunc func(ctx context.Context, containerID string) (types.ContainerJSON, error)
	ContainerCreateFunc  func(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, containerName string) (container.CreateResponse, error)
	ContainerStartFunc   func(ctx context.Context, containerID string, options container.StartOptions) error
	ContainerStopFunc    func(ctx context.Context, containerID string, options container.StopOptions) error
	ContainerRestartFunc func(ctx context.Context, containerID string, options container.StopOptions) error
	ContainerWaitFunc    func(ctx context.Context, containerID string, condition container.WaitCondition) (<-chan container.WaitResponse, <-chan error)
	ContainerLogsFunc    func(ctx context.Context, container string, options container.LogsOptions) (io.ReadCloser, error)
	VolumeListFunc       func(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error)
	VolumeRemoveFunc     func(ctx context.Context, volumeID string, force bool) error
	NetworkListFunc      func(ctx context.Context, options network.ListOptions) ([]network.Summary, error)
	NetworkRemoveFunc    func(ctx context.Context, networkID string) error
	ImageListFunc        func(ctx context.Context, options image.ListOptions) ([]image.Summary, error)
	ImageRemoveFunc      func(ctx context.Context, imageID string, options image.RemoveOptions) ([]image.DeleteResponse, error)
	InfoFunc             func(ctx context.Context) (system.Info, error)
	PingFunc             func(ctx context.Context) (types.Ping, error)
	CloseFunc            func() error
}

func (m *MockDockerClient) ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
	if m.ContainerListFunc != nil {
		return m.ContainerListFunc(ctx, options)
	}
	return []types.Container{}, nil
}

func (m *MockDockerClient) ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error {
	if m.ContainerRemoveFunc != nil {
		return m.ContainerRemoveFunc(ctx, containerID, options)
	}
	return nil
}

func (m *MockDockerClient) ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	if m.ContainerInspectFunc != nil {
		return m.ContainerInspectFunc(ctx, containerID)
	}
	return types.ContainerJSON{}, nil
}

func (m *MockDockerClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, containerName string) (container.CreateResponse, error) {
	if m.ContainerCreateFunc != nil {
		return m.ContainerCreateFunc(ctx, config, hostConfig, networkingConfig, platform, containerName)
	}
	return container.CreateResponse{}, nil
}

func (m *MockDockerClient) ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error {
	if m.ContainerStartFunc != nil {
		return m.ContainerStartFunc(ctx, containerID, options)
	}
	return nil
}

func (m *MockDockerClient) ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error {
	if m.ContainerStopFunc != nil {
		return m.ContainerStopFunc(ctx, containerID, options)
	}
	return nil
}

func (m *MockDockerClient) ContainerRestart(ctx context.Context, containerID string, options container.StopOptions) error {
	if m.ContainerRestartFunc != nil {
		return m.ContainerRestartFunc(ctx, containerID, options)
	}
	return nil
}

func (m *MockDockerClient) ContainerWait(ctx context.Context, containerID string, condition container.WaitCondition) (<-chan container.WaitResponse, <-chan error) {
	if m.ContainerWaitFunc != nil {
		return m.ContainerWaitFunc(ctx, containerID, condition)
	}
	ch := make(chan container.WaitResponse, 1)
	errCh := make(chan error, 1)
	ch <- container.WaitResponse{StatusCode: 0}
	close(ch)
	close(errCh)
	return ch, errCh
}

func (m *MockDockerClient) ContainerLogs(ctx context.Context, container string, options container.LogsOptions) (io.ReadCloser, error) {
	if m.ContainerLogsFunc != nil {
		return m.ContainerLogsFunc(ctx, container, options)
	}
	return io.NopCloser(nil), nil
}

func (m *MockDockerClient) VolumeList(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error) {
	if m.VolumeListFunc != nil {
		return m.VolumeListFunc(ctx, options)
	}
	return volume.ListResponse{}, nil
}

func (m *MockDockerClient) VolumeRemove(ctx context.Context, volumeID string, force bool) error {
	if m.VolumeRemoveFunc != nil {
		return m.VolumeRemoveFunc(ctx, volumeID, force)
	}
	return nil
}

func (m *MockDockerClient) NetworkList(ctx context.Context, options network.ListOptions) ([]network.Summary, error) {
	if m.NetworkListFunc != nil {
		return m.NetworkListFunc(ctx, options)
	}
	return []network.Summary{}, nil
}

func (m *MockDockerClient) NetworkRemove(ctx context.Context, networkID string) error {
	if m.NetworkRemoveFunc != nil {
		return m.NetworkRemoveFunc(ctx, networkID)
	}
	return nil
}

func (m *MockDockerClient) ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error) {
	if m.ImageListFunc != nil {
		return m.ImageListFunc(ctx, options)
	}
	return []image.Summary{}, nil
}

func (m *MockDockerClient) ImageRemove(ctx context.Context, imageID string, options image.RemoveOptions) ([]image.DeleteResponse, error) {
	if m.ImageRemoveFunc != nil {
		return m.ImageRemoveFunc(ctx, imageID, options)
	}
	return []image.DeleteResponse{}, nil
}

func (m *MockDockerClient) Info(ctx context.Context) (system.Info, error) {
	if m.InfoFunc != nil {
		return m.InfoFunc(ctx)
	}
	return system.Info{}, nil
}

func (m *MockDockerClient) Ping(ctx context.Context) (types.Ping, error) {
	if m.PingFunc != nil {
		return m.PingFunc(ctx)
	}
	return types.Ping{}, nil
}

func (m *MockDockerClient) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

// MockContainerJSON creates a ContainerJSON for testing
func MockContainerJSON(id, name, image, project string, running bool) types.ContainerJSON {
	return types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID:   id,
			Name: name,
			State: &types.ContainerState{
				Running: running,
				Status:  map[bool]string{true: "running", false: "exited"}[running],
			},
		},
		Config: &container.Config{
			Image:  image,
			Labels: map[string]string{"com.docker.compose.project": project},
		},
	}
}

// MockContainerJSONWithHealth creates a ContainerJSON with health status for testing
func MockContainerJSONWithHealth(id string, running bool, healthStatus string) types.ContainerJSON {
	return types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID: id,
			State: &types.ContainerState{
				Running: running,
				Health: &types.Health{
					Status: healthStatus,
				},
			},
		},
	}
}
