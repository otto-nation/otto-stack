package types

import (
	"time"
)

// ContainerSpec defines container configuration
type ContainerSpec struct {
	Image       string            `yaml:"image,omitempty"`
	Entrypoint  []string          `yaml:"entrypoint,omitempty"`
	Ports       []PortSpec        `yaml:"ports,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Volumes     []VolumeSpec      `yaml:"volumes,omitempty"`
	Restart     RestartPolicy     `yaml:"restart,omitempty"`
	Command     []string          `yaml:"command,omitempty"`
	Networks    []string          `yaml:"networks,omitempty"`
	MemoryLimit string            `yaml:"memory_limit,omitempty"`
	HealthCheck *HealthCheckSpec  `yaml:"health_check,omitempty"`
}

// PortSpec defines port mapping
type PortSpec struct {
	External string `yaml:"external,omitempty"`
	Internal string `yaml:"internal,omitempty"`
	Protocol string `yaml:"protocol,omitempty"`
}

// HealthCheckSpec defines health check configuration
type HealthCheckSpec struct {
	Test        []string      `yaml:"test,omitempty"`
	Interval    time.Duration `yaml:"interval,omitempty"`
	Timeout     time.Duration `yaml:"timeout,omitempty"`
	Retries     int           `yaml:"retries,omitempty"`
	StartPeriod time.Duration `yaml:"start_period,omitempty"`
}

// VolumeSpec defines volume configuration
type VolumeSpec struct {
	Name     string `yaml:"name,omitempty"`
	Mount    string `yaml:"mount,omitempty"`
	ReadOnly bool   `yaml:"read_only,omitempty"`
}

// ServiceSpec defines service integration
type ServiceSpec struct {
	Connection   *ConnectionSpec  `yaml:"connection,omitempty"`
	Dependencies DependenciesSpec `yaml:"dependencies,omitempty"`
	Management   *ManagementSpec  `yaml:"management,omitempty"`
}

// ConnectionSpec defines connection configuration
type ConnectionSpec struct {
	Type        ConnectionType `yaml:"type,omitempty"`
	DefaultPort int            `yaml:"default_port,omitempty"`
	DefaultUser string         `yaml:"default_user,omitempty"`
	Client      string         `yaml:"client,omitempty"`
	HostFlag    string         `yaml:"host_flag,omitempty"`
	PortFlag    string         `yaml:"port_flag,omitempty"`
	UserFlag    string         `yaml:"user_flag,omitempty"`
	DBFlag      string         `yaml:"database_flag,omitempty"`
	ExtraFlags  []string       `yaml:"extra_flags,omitempty"`
	URLPattern  string         `yaml:"url_pattern,omitempty"`
}

// DependenciesSpec defines service dependencies
type DependenciesSpec struct {
	Required  []string `yaml:"required,omitempty"`
	Soft      []string `yaml:"soft,omitempty"`
	Conflicts []string `yaml:"conflicts,omitempty"`
	Provides  []string `yaml:"provides,omitempty"`
}

// ManagementSpec defines management operations
type ManagementSpec struct {
	Connect *OperationSpec            `yaml:"connect,omitempty"`
	Backup  *OperationSpec            `yaml:"backup,omitempty"`
	Restore *OperationSpec            `yaml:"restore,omitempty"`
	Custom  map[string]*OperationSpec `yaml:"custom,omitempty"`
}

// OperationSpec defines a management operation
type OperationSpec struct {
	Type        string              `yaml:"type,omitempty"`
	Command     []string            `yaml:"command,omitempty"`
	Args        map[string][]string `yaml:"args,omitempty"`
	Defaults    map[string]string   `yaml:"defaults,omitempty"`
	PreCommands map[string][]string `yaml:"pre_commands,omitempty"`
	Extension   string              `yaml:"extension,omitempty"`
}

// DocumentationSpec defines service documentation
type DocumentationSpec struct {
	Examples      []string       `yaml:"examples,omitempty"`
	UsageNotes    string         `yaml:"usage_notes,omitempty"`
	Links         []string       `yaml:"links,omitempty"`
	UseCases      []string       `yaml:"use_cases,omitempty"`
	WebInterfaces []WebInterface `yaml:"web_interfaces,omitempty"`
}

// WebInterface defines a web interface for a service
type WebInterface struct {
	Name        string `yaml:"name"`
	URL         string `yaml:"url"`
	Description string `yaml:"description"`
}

// ParametersSpec defines service parameters
type ParametersSpec struct {
	Required []ParameterSpec `yaml:"required,omitempty"`
	Optional []ParameterSpec `yaml:"optional,omitempty"`
}

// ParameterSpec defines a single parameter
type ParameterSpec struct {
	Name        string        `yaml:"name,omitempty"`
	Type        ParameterType `yaml:"type,omitempty"`
	Description string        `yaml:"description,omitempty"`
	Default     string        `yaml:"default,omitempty"`
}

// Legacy types for backward compatibility
type ConnectionConfig = ConnectionSpec
type ServiceDependencies = DependenciesSpec

// Service Types
type ServiceType string

const (
	ServiceTypeContainer                     = "container"
	ServiceTypeComposite                     = "composite"
	ServiceTypeConfiguration                 = "configuration"
	ServiceTypeContainerType     ServiceType = ServiceTypeContainer
	ServiceTypeCompositeType     ServiceType = ServiceTypeComposite
	ServiceTypeConfigurationType ServiceType = ServiceTypeConfiguration
)

// Restart Policies
type RestartPolicy string

const (
	RestartPolicyNo                              = "no"
	RestartPolicyAlways                          = "always"
	RestartPolicyOnFailure                       = "on-failure"
	RestartPolicyUnlessStopped                   = "unless-stopped"
	RestartPolicyNoType            RestartPolicy = RestartPolicyNo
	RestartPolicyAlwaysType        RestartPolicy = RestartPolicyAlways
	RestartPolicyOnFailureType     RestartPolicy = RestartPolicyOnFailure
	RestartPolicyUnlessStoppedType RestartPolicy = RestartPolicyUnlessStopped
)

// Connection Types
type ConnectionType string

const (
	ConnectionClientCLI                = "cli"
	ConnectionTypeCLI   ConnectionType = ConnectionClientCLI
)

// Parameter Types
type ParameterType string

const (
	ParamTypeString                    = "string"
	ParamTypeInteger                   = "integer"
	ParameterTypeString  ParameterType = ParamTypeString
	ParameterTypeInteger ParameterType = ParamTypeInteger
)
