package types

import (
	"context"
	"log/slog"
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/spf13/cobra"
)

// Project represents a development stack project
type Project struct {
	Name        string    `yaml:"name" json:"name"`
	Type        string    `yaml:"type" json:"type"`
	Environment string    `yaml:"environment" json:"environment"`
	Services    []string  `yaml:"services" json:"services"`
	CreatedAt   time.Time `yaml:"created_at" json:"created_at"`
	UpdatedAt   time.Time `yaml:"updated_at" json:"updated_at"`
}

// Service represents a service within a development stack
type Service struct {
	Name        string            `yaml:"name" json:"name"`
	Type        string            `yaml:"type" json:"type"`
	Image       string            `yaml:"image,omitempty" json:"image,omitempty"`
	Ports       []PortMapping     `yaml:"ports,omitempty" json:"ports,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`
	Volumes     []VolumeMapping   `yaml:"volumes,omitempty" json:"volumes,omitempty"`
	DependsOn   []string          `yaml:"depends_on,omitempty" json:"depends_on,omitempty"`
	HealthCheck *HealthCheck      `yaml:"health_check,omitempty" json:"health_check,omitempty"`
}

// PortMapping represents a port mapping between host and container
type PortMapping struct {
	Host      string `yaml:"host" json:"host"`
	Container string `yaml:"container" json:"container"`
	Protocol  string `yaml:"protocol,omitempty" json:"protocol,omitempty"`
}

// VolumeMapping represents a volume mapping between host and container
type VolumeMapping struct {
	Host      string `yaml:"host" json:"host"`
	Container string `yaml:"container" json:"container"`
	Mode      string `yaml:"mode,omitempty" json:"mode,omitempty"`
}

// HealthCheck represents health check configuration
type HealthCheck struct {
	Test        []string      `yaml:"test" json:"test"`
	Interval    time.Duration `yaml:"interval,omitempty" json:"interval,omitempty"`
	Timeout     time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Retries     int           `yaml:"retries,omitempty" json:"retries,omitempty"`
	StartPeriod time.Duration `yaml:"start_period,omitempty" json:"start_period,omitempty"`
}

// Config represents the main application configuration
type Config struct {
	Project              ProjectInfo              `yaml:"project" json:"project"`
	Stack                StackConfig              `yaml:"stack" json:"stack"`
	ServiceConfiguration map[string]any           `yaml:"service_configuration,omitempty" json:"service_configuration,omitempty"`
	Global               GlobalConfig             `yaml:"global" json:"global"`
	Projects             map[string]ProjectConfig `yaml:"projects,omitempty" json:"projects,omitempty"`
}

// GlobalConfig represents global application settings
type GlobalConfig struct {
	DefaultProjectType string `yaml:"default_project_type" json:"default_project_type"`
	LogLevel           string `yaml:"log_level" json:"log_level"`
	ColorOutput        bool   `yaml:"color_output" json:"color_output"`
	CheckUpdates       bool   `yaml:"check_updates" json:"check_updates"`
}

// OttoStackConfig represents the main otto-stack-config.yml structure
type OttoStackConfig struct {
	Project              ProjectInfo    `yaml:"project" json:"project"`
	Stack                StackConfig    `yaml:"stack" json:"stack"`
	ServiceConfiguration map[string]any `yaml:"service_configuration,omitempty" json:"service_configuration,omitempty"`
}

// ProjectInfo represents project information
type ProjectInfo struct {
	Name        string `yaml:"name" json:"name"`
	Environment string `yaml:"environment" json:"environment"`
}

// StackConfig represents stack configuration
type StackConfig struct {
	Enabled []string `yaml:"enabled" json:"enabled"`
}

// ProjectConfig represents project-level configuration
type ProjectConfig struct {
	Version string `yaml:"version" json:"version"`
}

// ServiceInfo represents service information
type ServiceInfo struct {
	Name         string
	Description  string
	Category     string
	Status       string
	Type         string
	Visibility   string
	Components   []string
	Dependencies []string
}

// ServiceConfig represents service configuration
type ServiceConfig struct {
	Name         string              `yaml:"name"`
	Description  string              `yaml:"description"`
	Category     string              `yaml:"category"`
	Image        string              `yaml:"image,omitempty"`
	Ports        []string            `yaml:"ports,omitempty"`
	Volumes      []VolumeConfig      `yaml:"volumes,omitempty"`
	Environment  map[string]string   `yaml:"environment,omitempty"`
	Type         string              `yaml:"type,omitempty"`
	Visibility   string              `yaml:"visibility,omitempty"`
	Components   []string            `yaml:"components,omitempty"`
	Dependencies ServiceDependencies `yaml:"dependencies,omitempty"`
	Docker       map[string]any      `yaml:"docker,omitempty"`
	Management   map[string]any      `yaml:"management,omitempty"`

	// Flexible fields for documentation and other unknown fields
	Connection           map[string]any `yaml:"connection,omitempty"`
	Documentation        map[string]any `yaml:"documentation,omitempty"`
	ServiceConfiguration []any          `yaml:"service_configuration,omitempty"`

	// Catch-all for any other fields
	Extra map[string]any `yaml:",inline"`
}

type VolumeConfig struct {
	Name        string `yaml:"name"`
	Mount       string `yaml:"mount"`
	Description string `yaml:"description,omitempty"`
}

type ServiceDependencies struct {
	Required  []string `yaml:"required,omitempty"`
	Soft      []string `yaml:"soft,omitempty"`
	Conflicts []string `yaml:"conflicts,omitempty"`
	Provides  []string `yaml:"provides,omitempty"`
}

// Helper methods for ServiceConfig backward compatibility
func (sc *ServiceConfig) GetName() string {
	return sc.Name
}

func (sc *ServiceConfig) GetDescription() string {
	return sc.Description
}

func (sc *ServiceConfig) GetCategory() string {
	return sc.Category
}

func (sc *ServiceConfig) GetVisibility() string {
	return sc.Visibility
}

func (sc *ServiceConfig) GetPorts() []string {
	return sc.Ports
}

func (sc *ServiceConfig) GetEnvironment() map[string]string {
	return sc.Environment
}

func (sc *ServiceConfig) GetDocker() map[string]any {
	return sc.Docker
}

func (sc *ServiceConfig) GetConnection() map[string]any {
	return sc.Connection
}

func (sc *ServiceConfig) GetDependencies() ServiceDependencies {
	return sc.Dependencies
}

func (sc *ServiceConfig) GetManagement() map[string]any {
	return sc.Management
}

// ServiceStatus represents the runtime status of a service
type ServiceStatus struct {
	Name      string        `json:"name"`
	State     ServiceState  `json:"state"`
	Health    HealthStatus  `json:"health"`
	Uptime    time.Duration `json:"uptime"`
	CPUUsage  float64       `json:"cpu_usage"`
	Memory    MemoryUsage   `json:"memory"`
	Ports     []PortMapping `json:"ports"`
	CreatedAt time.Time     `json:"created_at"`
	StartedAt *time.Time    `json:"started_at,omitempty"`
}

// MemoryUsage represents memory usage statistics
type MemoryUsage struct {
	Used  uint64 `json:"used"`
	Limit uint64 `json:"limit"`
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
	ShellTypeBash       ShellType = constants.ShellBash
	ShellTypeZsh        ShellType = constants.ShellZsh
	ShellTypeFish       ShellType = constants.ShellFish
	ShellTypePowerShell ShellType = constants.ShellPowerShell
)

// String returns the string representation of the shell type
func (s ShellType) String() string {
	return string(s)
}

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

// Output interface for UI output (compatibility)
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
