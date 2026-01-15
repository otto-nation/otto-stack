package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/system"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// DockerClient defines the interface for Docker operations
type DockerClient interface {
	ContainerList(ctx context.Context, options container.ListOptions) ([]container.Summary, error)
	ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error
	ContainerInspect(ctx context.Context, containerID string) (container.InspectResponse, error)
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, containerName string) (container.CreateResponse, error)
	ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error
	ContainerWait(ctx context.Context, containerID string, condition container.WaitCondition) (<-chan container.WaitResponse, <-chan error)
	ContainerLogs(ctx context.Context, container string, options container.LogsOptions) (io.ReadCloser, error)
	VolumeList(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error)
	VolumeRemove(ctx context.Context, volumeID string, force bool) error
	NetworkList(ctx context.Context, options network.ListOptions) ([]network.Summary, error)
	NetworkRemove(ctx context.Context, networkID string) error
	ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error)
	ImageRemove(ctx context.Context, imageID string, options image.RemoveOptions) ([]image.DeleteResponse, error)
	Info(ctx context.Context) (system.Info, error)
	Ping(ctx context.Context) (types.Ping, error)
	Close() error
}

// dockerClientAdapter adapts the real Docker client to our interface
type dockerClientAdapter struct {
	client *client.Client
}

// NewDockerClientAdapter wraps a real Docker client
func NewDockerClientAdapter(cli *client.Client) DockerClient {
	return &dockerClientAdapter{client: cli}
}

func (a *dockerClientAdapter) ContainerList(ctx context.Context, options container.ListOptions) ([]container.Summary, error) {
	return a.client.ContainerList(ctx, options)
}

func (a *dockerClientAdapter) ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error {
	return a.client.ContainerRemove(ctx, containerID, options)
}

func (a *dockerClientAdapter) ContainerInspect(ctx context.Context, containerID string) (container.InspectResponse, error) {
	return a.client.ContainerInspect(ctx, containerID)
}

func (a *dockerClientAdapter) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, containerName string) (container.CreateResponse, error) {
	return a.client.ContainerCreate(ctx, config, hostConfig, networkingConfig, platform, containerName)
}

func (a *dockerClientAdapter) ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error {
	return a.client.ContainerStart(ctx, containerID, options)
}

func (a *dockerClientAdapter) ContainerWait(ctx context.Context, containerID string, condition container.WaitCondition) (<-chan container.WaitResponse, <-chan error) {
	return a.client.ContainerWait(ctx, containerID, condition)
}

func (a *dockerClientAdapter) ContainerLogs(ctx context.Context, container string, options container.LogsOptions) (io.ReadCloser, error) {
	return a.client.ContainerLogs(ctx, container, options)
}

func (a *dockerClientAdapter) VolumeList(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error) {
	return a.client.VolumeList(ctx, options)
}

func (a *dockerClientAdapter) VolumeRemove(ctx context.Context, volumeID string, force bool) error {
	return a.client.VolumeRemove(ctx, volumeID, force)
}

func (a *dockerClientAdapter) NetworkList(ctx context.Context, options network.ListOptions) ([]network.Summary, error) {
	return a.client.NetworkList(ctx, options)
}

func (a *dockerClientAdapter) NetworkRemove(ctx context.Context, networkID string) error {
	return a.client.NetworkRemove(ctx, networkID)
}

func (a *dockerClientAdapter) ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error) {
	return a.client.ImageList(ctx, options)
}

func (a *dockerClientAdapter) ImageRemove(ctx context.Context, imageID string, options image.RemoveOptions) ([]image.DeleteResponse, error) {
	return a.client.ImageRemove(ctx, imageID, options)
}

func (a *dockerClientAdapter) Info(ctx context.Context) (system.Info, error) {
	return a.client.Info(ctx)
}

func (a *dockerClientAdapter) Ping(ctx context.Context) (types.Ping, error) {
	return a.client.Ping(ctx)
}

func (a *dockerClientAdapter) Close() error {
	return a.client.Close()
}
