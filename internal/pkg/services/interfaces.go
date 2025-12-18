package services

import (
	"context"

	"github.com/otto-nation/otto-stack/internal/pkg/config"
)

// StackService defines the interface for stack operations
type StackService interface {
	// StartStack starts the specified services
	StartStack(ctx context.Context, services []string, options StartOptions) error

	// StopStack stops the stack services
	StopStack(ctx context.Context, preserveVolumes bool) error

	// GetStackStatus returns the current status of all services
	GetStackStatus(ctx context.Context) (*StackStatus, error)

	// RestartStack restarts the specified services
	RestartStack(ctx context.Context, services []string) error
}

// ConfigService defines the interface for configuration operations
type ConfigService interface {
	// LoadConfig loads the project configuration
	LoadConfig() (*config.Config, error)

	// SaveConfig saves the project configuration
	SaveConfig(cfg *config.Config) error

	// ValidateConfig validates the configuration
	ValidateConfig(cfg *config.Config) error

	// GetConfigHash returns a hash of the current configuration
	GetConfigHash(cfg *config.Config) (string, error)
}

// DockerService defines the interface for Docker operations
type DockerService interface {
	// CreateContainer creates a new container
	CreateContainer(ctx context.Context, config ContainerConfig) (string, error)

	// StartContainer starts a container by ID
	StartContainer(ctx context.Context, containerID string) error

	// StopContainer stops a container by ID
	StopContainer(ctx context.Context, containerID string) error

	// GetContainerStatus returns the status of a container
	GetContainerStatus(ctx context.Context, containerID string) (*ContainerStatus, error)

	// ListContainers returns all containers for the project
	ListContainers(ctx context.Context, projectName string) ([]ContainerStatus, error)
}

// StartOptions contains options for starting services
type StartOptions struct {
	Build          bool
	ForceRecreate  bool
	Detach         bool
	NoDeps         bool
	ResolveDeps    bool
	CheckConflicts bool
	RemoveOrphans  bool
}

// StackStatus represents the status of the entire stack
type StackStatus struct {
	Services []ServiceStatus `json:"services"`
	Overall  string          `json:"overall"`
}

// ServiceStatus represents the status of a single service
type ServiceStatus struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	Health    string `json:"health"`
	Ports     []Port `json:"ports"`
	Uptime    string `json:"uptime"`
	UpdatedAt string `json:"updated_at"`
}

// ContainerConfig represents container configuration
type ContainerConfig struct {
	Name          string
	Image         string
	Ports         []Port
	Environment   map[string]string
	Volumes       []string
	Networks      []string
	HealthCheck   *HealthCheck
	RestartPolicy string
}

// ContainerStatus represents the status of a container
type ContainerStatus struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Status  string            `json:"status"`
	Health  string            `json:"health"`
	Ports   []Port            `json:"ports"`
	Labels  map[string]string `json:"labels"`
	Created string            `json:"created"`
	Started string            `json:"started"`
}

// Port represents a port mapping
type Port struct {
	External string `json:"external"`
	Internal string `json:"internal"`
	Protocol string `json:"protocol"`
}

// HealthCheck represents a health check configuration
type HealthCheck struct {
	Test     []string `json:"test"`
	Interval string   `json:"interval"`
	Timeout  string   `json:"timeout"`
	Retries  int      `json:"retries"`
}
