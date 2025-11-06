package types

import (
	"context"
	"log/slog"
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/spf13/cobra"
)

// HealthCheck represents health check configuration
type HealthCheck struct {
	Test        []string      `yaml:"test" json:"test"`
	Interval    time.Duration `yaml:"interval,omitempty" json:"interval,omitempty"`
	Timeout     time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Retries     int           `yaml:"retries,omitempty" json:"retries,omitempty"`
	StartPeriod time.Duration `yaml:"start_period,omitempty" json:"start_period,omitempty"`
}

// ServiceStatus represents the runtime status of a service
type ServiceStatus struct {
	Name      string        `json:"name"`
	State     ServiceState  `json:"state"`
	Health    HealthStatus  `json:"health"`
	Uptime    time.Duration `json:"uptime"`
	CPUUsage  float64       `json:"cpu_usage"`
	Memory    uint64        `json:"memory"`
	StartedAt *time.Time    `json:"started_at,omitempty"`
	Ports     []string      `json:"ports,omitempty"`
	Image     string        `json:"image,omitempty"`
	ID        string        `json:"id,omitempty"`
}

// StackStatus represents the overall status of a development stack
type StackStatus struct {
	Project     string          `json:"project"`
	Services    []ServiceStatus `json:"services"`
	TotalCPU    float64         `json:"total_cpu"`
	TotalMemory uint64          `json:"total_memory"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// ServiceState represents the state of a service
type ServiceState string

const (
	ServiceStateRunning ServiceState = constants.StateRunning
	ServiceStateStopped ServiceState = constants.StateStopped
	ServiceStateCreated ServiceState = constants.StateCreated
)

// String returns the string representation of the service state
func (s ServiceState) String() string {
	return string(s)
}

// IsRunning returns true if the service is running
func (s ServiceState) IsRunning() bool {
	return s == ServiceStateRunning
}

// IsStopped returns true if the service is stopped
func (s ServiceState) IsStopped() bool {
	return s == ServiceStateStopped
}

// HealthStatus represents the health status of a service
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = constants.HealthHealthy
	HealthStatusUnhealthy HealthStatus = constants.HealthUnhealthy
	HealthStatusStarting  HealthStatus = constants.HealthStarting
	HealthStatusNone      HealthStatus = constants.HealthNone
)

// String returns the string representation of the health status
func (h HealthStatus) String() string {
	return string(h)
}

// IsHealthy returns true if the service is healthy
func (h HealthStatus) IsHealthy() bool {
	return h == HealthStatusHealthy
}

// IsUnhealthy returns true if the service is unhealthy
func (h HealthStatus) IsUnhealthy() bool {
	return h == HealthStatusUnhealthy
}

// IsStarting returns true if the service is starting
func (h HealthStatus) IsStarting() bool {
	return h == HealthStatusStarting
}

// ShellType represents supported shell types for completion
type ShellType string

const (
	ShellTypeBash       ShellType = "bash"
	ShellTypeZsh        ShellType = "zsh"
	ShellTypeFish       ShellType = "fish"
	ShellTypePowerShell ShellType = "powershell"
)

// IsValid returns true if the shell type is supported
func (s ShellType) IsValid() bool {
	switch s {
	case ShellTypeBash, ShellTypeZsh, ShellTypeFish, ShellTypePowerShell:
		return true
	default:
		return false
	}
}

// AllShellTypeStrings returns all supported shell types as strings
func AllShellTypeStrings() []string {
	return []string{
		string(ShellTypeBash),
		string(ShellTypeZsh),
		string(ShellTypeFish),
		string(ShellTypePowerShell),
	}
}

// Output interface for command output
type Output interface {
	Success(msg string, args ...any)
	Error(msg string, args ...any)
	Warning(msg string, args ...any)
	Info(msg string, args ...any)
	Header(msg string, args ...any)
	Muted(msg string, args ...any)
}

// ServiceManagerAdapter interface for service management operations
type ServiceManagerAdapter interface {
	StartServices(ctx context.Context, serviceNames []string, options StartOptions) error
	StopServices(ctx context.Context, serviceNames []string, options StopOptions) error
	GetServiceStatus(ctx context.Context, serviceNames []string) ([]ServiceStatus, error)
}

// LoggerAdapter interface for accessing underlying slog.Logger
type LoggerAdapter interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	SlogLogger() *slog.Logger
}

// BaseCommand provides common dependencies for command handlers
type BaseCommand struct {
	Manager ServiceManagerAdapter
	Logger  LoggerAdapter
	Output  Output
}

// CommandHandler interface for command handlers
type CommandHandler interface {
	Handle(ctx context.Context, cmd *cobra.Command, args []string, base *BaseCommand) error
}
