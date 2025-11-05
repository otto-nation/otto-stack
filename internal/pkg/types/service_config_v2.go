package types

import (
	"fmt"
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
)

// ServiceConfigV2 represents the new structured service configuration format
type ServiceConfigV2 struct {
	// Metadata
	Name        string                `yaml:"name" validate:"required,min=1"`
	Description string                `yaml:"description" validate:"required,min=1"`
	Hidden      bool                  `yaml:"hidden,omitempty"`
	Type        constants.ServiceType `yaml:"type,omitempty"`

	// Configuration sections
	Runtime     RuntimeSpec     `yaml:"runtime" validate:"required"`
	Integration IntegrationSpec `yaml:"integration,omitempty"`

	// Documentation and parameters
	Documentation DocumentationSpec `yaml:"documentation,omitempty"`
	Parameters    ParametersSpec    `yaml:"parameters,omitempty"`
}

// RuntimeSpec defines how the service runs
type RuntimeSpec struct {
	Image       string            `yaml:"image,omitempty"`
	Ports       []PortSpec        `yaml:"ports,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Container   ContainerSpec     `yaml:"container,omitempty"`
	Volumes     []VolumeSpec      `yaml:"volumes,omitempty"`
}

// PortSpec represents a port mapping
type PortSpec struct {
	Host      string `yaml:"host"`
	Container string `yaml:"container"`
	Protocol  string `yaml:"protocol,omitempty"` // tcp, udp
}

// ContainerSpec defines container-specific configuration
type ContainerSpec struct {
	Restart     constants.RestartPolicy `yaml:"restart,omitempty"`
	Networks    []string                `yaml:"networks,omitempty"`
	MemoryLimit string                  `yaml:"memory_limit,omitempty"`
	CPULimit    string                  `yaml:"cpu_limit,omitempty"`
	Command     []string                `yaml:"command,omitempty"`
	Environment []string                `yaml:"environment,omitempty"`
	HealthCheck *HealthCheckSpec        `yaml:"health_check,omitempty"`
}

// HealthCheckSpec defines container health check
type HealthCheckSpec struct {
	Test        []string      `yaml:"test" validate:"required,min=1"`
	Interval    time.Duration `yaml:"interval,omitempty"`
	Timeout     time.Duration `yaml:"timeout,omitempty"`
	Retries     int           `yaml:"retries,omitempty"`
	StartPeriod time.Duration `yaml:"start_period,omitempty"`
}

// VolumeSpec represents a volume mount
type VolumeSpec struct {
	Name        string `yaml:"name" validate:"required"`
	Mount       string `yaml:"mount" validate:"required"`
	Description string `yaml:"description,omitempty"`
	ReadOnly    bool   `yaml:"read_only,omitempty"`
}

// IntegrationSpec defines how the service integrates with others
type IntegrationSpec struct {
	Connection   *ConnectionSpec  `yaml:"connection,omitempty"`
	Dependencies DependenciesSpec `yaml:"dependencies,omitempty"`
	Management   *ManagementSpec  `yaml:"management,omitempty"`
}

// ConnectionSpec defines how to connect to the service
type ConnectionSpec struct {
	Type        constants.ConnectionType `yaml:"type" validate:"required"`
	DefaultPort int                      `yaml:"default_port,omitempty"`
	DefaultUser string                   `yaml:"default_user,omitempty"`

	// CLI-specific fields
	Client     string   `yaml:"client,omitempty"`
	HostFlag   string   `yaml:"host_flag,omitempty"`
	PortFlag   string   `yaml:"port_flag,omitempty"`
	UserFlag   string   `yaml:"user_flag,omitempty"`
	DBFlag     string   `yaml:"database_flag,omitempty"`
	ExtraFlags []string `yaml:"extra_flags,omitempty"`

	// Web-specific fields
	URLPattern string `yaml:"url_pattern,omitempty"`
}

// DependenciesSpec defines service dependencies (alias for backward compatibility)
type DependenciesSpec = ServiceDependencies

// ManagementSpec defines management operations
type ManagementSpec struct {
	Connect *OperationSpec            `yaml:"connect,omitempty"`
	Backup  *OperationSpec            `yaml:"backup,omitempty"`
	Restore *OperationSpec            `yaml:"restore,omitempty"`
	Custom  map[string]*OperationSpec `yaml:"custom,omitempty"`
}

// OperationSpec defines a management operation
type OperationSpec struct {
	Type        string              `yaml:"type,omitempty"` // command, script, api
	Command     []string            `yaml:"command,omitempty"`
	Args        map[string][]string `yaml:"args,omitempty"`
	Defaults    map[string]string   `yaml:"defaults,omitempty"`
	PreCommands map[string][]string `yaml:"pre_commands,omitempty"`
	Extension   string              `yaml:"extension,omitempty"`
}

// DocumentationSpec defines service documentation
type DocumentationSpec struct {
	Examples      []string      `yaml:"examples,omitempty"`
	UsageNotes    string        `yaml:"usage_notes,omitempty"`
	Links         []string      `yaml:"links,omitempty"`
	UseCases      []string      `yaml:"use_cases,omitempty"`
	Docs          []DocLinkSpec `yaml:"docs,omitempty"`
	WebInterfaces []WebUISpec   `yaml:"web_interfaces,omitempty"`
	SpringConfig  *SpringConfig `yaml:"spring_config,omitempty"`
}

// DocLinkSpec represents a documentation link
type DocLinkSpec struct {
	Name string `yaml:"name" validate:"required"`
	URL  string `yaml:"url" validate:"required,url"`
}

// WebUISpec represents a web interface
type WebUISpec struct {
	Name        string `yaml:"name" validate:"required"`
	URL         string `yaml:"url" validate:"required"`
	Description string `yaml:"description,omitempty"`
}

// SpringConfig represents Spring Boot configuration
type SpringConfig struct {
	Properties []string `yaml:"properties,omitempty"`
}

// ParametersSpec defines configurable parameters
type ParametersSpec struct {
	Parameters []ParameterSpec `yaml:"parameters,omitempty"`
}

// ParameterSpec defines a configurable parameter
type ParameterSpec struct {
	Name        string                  `yaml:"name" validate:"required"`
	Type        constants.ParameterType `yaml:"type" validate:"required"`
	Description string                  `yaml:"description" validate:"required"`
	Default     string                  `yaml:"default,omitempty"`
	Example     string                  `yaml:"example,omitempty"`
	Required    bool                    `yaml:"required,omitempty"`
}

// Validate validates the entire service configuration
func (sc *ServiceConfigV2) Validate() error {
	if sc.Name == "" {
		return fmt.Errorf("service name is required")
	}
	if sc.Description == "" {
		return fmt.Errorf("service description is required")
	}

	if err := sc.Type.Validate(); err != nil {
		return fmt.Errorf("invalid service type: %w", err)
	}

	if err := sc.Runtime.Validate(); err != nil {
		return fmt.Errorf("invalid runtime config: %w", err)
	}

	if err := sc.Integration.Validate(); err != nil {
		return fmt.Errorf("invalid integration config: %w", err)
	}

	return nil
}

// Validate validates the runtime specification
func (rs *RuntimeSpec) Validate() error {
	if err := rs.Container.Validate(); err != nil {
		return fmt.Errorf("invalid container config: %w", err)
	}

	for i, volume := range rs.Volumes {
		if volume.Name == "" {
			return fmt.Errorf("volume %d: name is required", i)
		}
		if volume.Mount == "" {
			return fmt.Errorf("volume %d: mount path is required", i)
		}
	}

	return nil
}

// Validate validates the container specification
func (cs *ContainerSpec) Validate() error {
	return cs.Restart.Validate()
}

// Validate validates the integration specification
func (is *IntegrationSpec) Validate() error {
	if is.Connection != nil {
		if err := is.Connection.Validate(); err != nil {
			return fmt.Errorf("invalid connection config: %w", err)
		}
	}
	return nil
}

// Validate validates the connection specification
func (cs *ConnectionSpec) Validate() error {
	return cs.Type.Validate()
}
