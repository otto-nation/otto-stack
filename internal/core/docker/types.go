package docker

import (
	"time"
)

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
