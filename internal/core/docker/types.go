package docker

import (
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
)

// StartOptions defines options for starting services
type StartOptions struct {
	Build          bool
	ForceRecreate  bool
	NoDeps         bool
	Detach         bool
	Timeout        time.Duration
	ResolveDeps    bool
	CheckConflicts bool
}

// StopOptions defines options for stopping services
type StopOptions struct {
	Timeout       int
	Remove        bool
	RemoveVolumes bool
	RemoveOrphans bool
	RemoveImages  string
}

// LogOptions defines options for retrieving container logs
type LogOptions struct {
	Follow     bool
	Timestamps bool
	Tail       string
	Since      string
}

// CleanupOptions defines options for cleaning up resources
type CleanupOptions struct {
	RemoveVolumes  bool
	RemoveImages   bool
	RemoveNetworks bool
	All            bool
	DryRun         bool
}

// ExecOptions defines options for executing commands in containers
type ExecOptions struct {
	User        string
	WorkingDir  string
	Env         []string
	Interactive bool
	TTY         bool
	Detach      bool
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
	DockerServiceStateRunning DockerServiceState = constants.StateRunning
	DockerServiceStateStopped DockerServiceState = constants.StateStopped
	DockerServiceStateCreated DockerServiceState = constants.StateCreated
)

// IsRunning returns true if the service is running
func (s DockerServiceState) IsRunning() bool {
	return s == DockerServiceStateRunning
}

// DockerHealthStatus represents the health status of a service
type DockerHealthStatus string

const (
	DockerHealthStatusHealthy   DockerHealthStatus = constants.HealthHealthy
	DockerHealthStatusUnhealthy DockerHealthStatus = constants.HealthUnhealthy
	DockerHealthStatusStarting  DockerHealthStatus = constants.HealthStarting
	DockerHealthStatusNone      DockerHealthStatus = constants.HealthNone
)
