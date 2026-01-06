package docker

import (
	"time"

	"github.com/docker/compose/v5/pkg/api"
)

// UpOptions represents options for compose up operations
type UpOptions struct {
	Build         bool
	ForceRecreate bool
	RemoveOrphans bool
	Services      []string
	Timeout       *time.Duration
}

// DownOptions represents options for compose down operations
type DownOptions struct {
	RemoveVolumes bool
	RemoveOrphans bool
	Timeout       *time.Duration
	Services      []string
}

// StopOptions represents options for compose stop operations
type StopOptions struct {
	Timeout  *time.Duration
	Services []string
}

// LogOptions represents options for compose logs operations
type LogOptions struct {
	Services   []string
	Follow     bool
	Timestamps bool
	Tail       string
}

// ToSDK converts UpOptions to Docker Compose API options
func (o UpOptions) ToSDK() api.UpOptions {
	createOpts := api.CreateOptions{
		Services: o.Services,
	}
	if o.Build {
		createOpts.Build = &api.BuildOptions{}
	}
	if o.ForceRecreate {
		createOpts.Recreate = api.RecreateForce
	}
	return api.UpOptions{Create: createOpts}
}

// ToSDK converts DownOptions to Docker Compose API options
func (o DownOptions) ToSDK() api.DownOptions {
	return api.DownOptions{
		Services: o.Services,
		Volumes:  o.RemoveVolumes,
		Timeout:  o.Timeout,
	}
}

// ToSDK converts StopOptions to Docker Compose API options
func (o StopOptions) ToSDK() api.StopOptions {
	return api.StopOptions{
		Services: o.Services,
		Timeout:  o.Timeout,
	}
}

// ToSDK converts LogOptions to Docker Compose API options
func (o LogOptions) ToSDK() api.LogOptions {
	return api.LogOptions{
		Services:   o.Services,
		Follow:     o.Follow,
		Timestamps: o.Timestamps,
		Tail:       o.Tail,
	}
}

// DockerServiceStatus represents the runtime status of a service
type DockerServiceStatus struct {
	Name      string             `json:"name"`
	State     DockerServiceState `json:"state"`
	Health    DockerHealthStatus `json:"health"`
	Uptime    time.Duration      `json:"uptime"`
	CPUUsage  float64            `json:"cpu_usage"`
	Memory    uint64             `json:"memory"`
	StartedAt *time.Time         `json:"started_at,omitempty"`
	Ports     []string           `json:"ports,omitempty"`
	Image     string             `json:"image,omitempty"`
	ID        string             `json:"id,omitempty"`
}

// DockerServiceState represents the state of a service
type DockerServiceState string

const (
	DockerServiceStateRunning DockerServiceState = StateRunning
	DockerServiceStateStopped DockerServiceState = StateStopped
	DockerServiceStateCreated DockerServiceState = StateCreated
)

// IsRunning returns true if the service is running
func (s DockerServiceState) IsRunning() bool {
	return s == DockerServiceStateRunning
}

// DockerHealthStatus represents the health status of a service
type DockerHealthStatus string

const (
	DockerHealthStatusHealthy   DockerHealthStatus = HealthHealthy
	DockerHealthStatusUnhealthy DockerHealthStatus = HealthUnhealthy
	DockerHealthStatusStarting  DockerHealthStatus = HealthStarting
	DockerHealthStatusNone      DockerHealthStatus = HealthNone
)
