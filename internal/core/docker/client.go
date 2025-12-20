package docker

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"

	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type ResourceType string

const (
	ResourceContainer ResourceType = "container"
	ResourceVolume    ResourceType = "volume"
	ResourceNetwork   ResourceType = "network"
	ResourceImage     ResourceType = "image"
)

// InitContainerConfig holds configuration for init containers
type InitContainerConfig struct {
	Image       string
	Command     []string
	Environment map[string]string
	Volumes     []string
	WorkingDir  string
	Networks    []string
}

type Client struct {
	cli       *client.Client
	logger    *slog.Logger
	resources *ResourceManager
}

func NewClient(logger *slog.Logger) (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, pkgerrors.NewDockerError("create Docker client", "", err)
	}

	dc := &Client{
		cli:    cli,
		logger: logger,
	}
	dc.resources = NewResourceManager(dc)

	return dc, nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}

// Compose operations using docker compose CLI
func (c *Client) ComposeUp(ctx context.Context, project string, services []string, options StartOptions) error {
	builder := NewComposeBuilder().Project(project).Services(services...)

	flags := []string{}
	if options.Build {
		flags = append(flags, FlagBuild)
	}
	if options.ForceRecreate {
		flags = append(flags, FlagForceRecreate)
	}
	if options.RemoveOrphans {
		flags = append(flags, FlagRemoveOrphans)
	}

	// Add characteristic-based flags
	if len(options.Characteristics) > 0 {
		resolver, err := NewServiceCharacteristicsResolver()
		if err == nil {
			charFlags := resolver.ResolveComposeUpFlags(options.Characteristics)
			flags = append(flags, charFlags...)
		}
	}

	if len(flags) > 0 {
		builder = builder.WithFlags(flags...)
	}

	return builder.Up()
}

func (c *Client) ComposeDown(ctx context.Context, project string, options StopOptions) error {
	charFlags := c.getCharacteristicFlags(options.Characteristics)

	if options.Remove {
		return c.composeDownWithRemove(project, options, charFlags)
	}

	return c.composeStop(project, options)
}

func (c *Client) getCharacteristicFlags(characteristics []string) []string {
	if len(characteristics) == 0 {
		return nil
	}

	resolver, err := NewServiceCharacteristicsResolver()
	if err != nil {
		return nil
	}

	return resolver.ResolveComposeDownFlags(characteristics)
}

func (c *Client) composeDownWithRemove(project string, options StopOptions, charFlags []string) error {
	if options.RemoveOrphans && !options.RemoveVolumes {
		return c.composeDownCustom(project, options.Services, charFlags)
	}

	builder := NewComposeBuilder().Project(project)
	if len(options.Services) > 0 {
		builder = builder.Services(options.Services...)
	}
	if len(charFlags) > 0 {
		builder = builder.WithFlags(charFlags...)
	}
	return builder.Down()
}

func (c *Client) composeDownCustom(project string, services []string, charFlags []string) error {
	cmd := NewCommand(DockerCmd).
		Subcommand(DockerComposeCmd).
		Flag(FlagProjectName, project).
		Args(ComposeDownCmd)

	if len(services) > 0 {
		cmd = cmd.Args(services...)
	}

	cmd = cmd.BoolFlag(FlagRemoveOrphans)

	for _, flag := range charFlags {
		cmd = cmd.BoolFlag(flag)
	}

	builtCmd := cmd.Build()
	builtCmd.Stdout = os.Stdout
	builtCmd.Stderr = os.Stderr
	return builtCmd.Run()
}

func (c *Client) composeStop(project string, options StopOptions) error {
	cmd := NewCommand(DockerCmd).
		Subcommand(DockerComposeCmd).
		Flag(FlagProjectName, project).
		Args(ComposeStopCmd)

	if len(options.Services) > 0 {
		cmd = cmd.Args(options.Services...)
	}

	if options.Timeout > 0 {
		cmd = cmd.Flag(FlagTimeout, fmt.Sprintf("%d", options.Timeout))
	}

	builtCmd := cmd.Build()
	builtCmd.Stdout = os.Stdout
	builtCmd.Stderr = os.Stderr
	return builtCmd.Run()
}

func (c *Client) ComposeLogs(ctx context.Context, project string, services []string, options LogOptions) error {
	builder := NewComposeBuilder().Project(project).Services(services...)

	flags := []string{}
	if options.Follow {
		flags = append(flags, FlagFollow)
	}
	if options.Timestamps {
		flags = append(flags, "timestamps")
	}
	if options.Tail != "" {
		flags = append(flags, FlagTail, options.Tail)
	}

	if len(flags) > 0 {
		builder = builder.WithFlags(flags...)
	}

	return builder.Logs()
}

// Resource management
func (c *Client) ListResources(ctx context.Context, resourceType ResourceType, project string) ([]string, error) {
	filter := NewProjectFilter(project)
	return c.resources.List(ctx, resourceType, filter)
}

func (c *Client) RemoveResources(ctx context.Context, resourceType ResourceType, project string) error {
	names, err := c.ListResources(ctx, resourceType, project)
	if err != nil {
		return err
	}
	return c.resources.Remove(ctx, resourceType, names)
}

// Container status
func (c *Client) GetDockerServiceStatus(ctx context.Context, project string, services []string) ([]DockerServiceStatus, error) {
	filter := NewProjectFilter(project)
	containers, err := c.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filter})
	if err != nil {
		return nil, err
	}

	var statuses []DockerServiceStatus
	for _, container := range containers {
		serviceName := container.Labels[ComposeServiceLabel]
		if len(services) > 0 && !contains(services, serviceName) {
			continue
		}

		status := DockerServiceStatus{
			Name:   serviceName,
			State:  DockerServiceState(container.State),
			Health: c.getContainerHealth(container),
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

func (c *Client) getContainerHealth(cont container.Summary) DockerHealthStatus {
	switch cont.State {
	case StateRunning:
		return DockerHealthStatusHealthy
	case StateStopped:
		return DockerHealthStatusUnhealthy
	case StateStarting:
		return DockerHealthStatusStarting
	default:
		return DockerHealthStatusNone
	}
}

// RunCommand executes a docker command with the given arguments
func (c *Client) RunCommand(ctx context.Context, args ...string) error {
	cmd := NewCommand(DockerCmd).Args(args...).Context(ctx).Build()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunInitContainer runs an init container and waits for completion
func (c *Client) RunInitContainer(ctx context.Context, containerName string, config InitContainerConfig) error {
	builder := NewCommand(DockerCmd).
		Subcommand(DockerRunCmd).
		BoolFlag(FlagRm).
		Flag(FlagName, containerName).
		Context(ctx)

	// Add environment variables
	for key, value := range config.Environment {
		builder = builder.Flag(FlagEnv, fmt.Sprintf("%s=%s", key, value))
	}

	// Add volumes
	for _, volume := range config.Volumes {
		builder = builder.Flag(FlagVolume, volume)
	}

	// Add working directory
	if config.WorkingDir != "" {
		builder = builder.Flag(FlagWorkingDir, config.WorkingDir)
	}

	// Add networks
	for _, network := range config.Networks {
		builder = builder.Flag(FlagNetwork, network)
	}

	// Add image and command
	cmd := builder.Args(config.Image).Args(config.Command...).Build()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}
