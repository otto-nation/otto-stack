package services

import (
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// Service represents a service definition (V1 format)
type Service struct {
	Name         string            `yaml:"name"`
	Description  string            `yaml:"description"`
	Category     string            `yaml:"category"`
	Type         string            `yaml:"type,omitempty"`
	Hidden       bool              `yaml:"hidden,omitempty"`
	Docker       DockerConfig      `yaml:"docker,omitempty"`
	Connection   ConnectionConfig  `yaml:"connection,omitempty"`
	Dependencies DependenciesV1    `yaml:"dependencies,omitempty"`
	Environment  map[string]string `yaml:"environment,omitempty"`
	Management   *ManagementV1     `yaml:"management,omitempty"`
}

// DependenciesV1 represents V1 dependencies structure
type DependenciesV1 struct {
	Required  []string `yaml:"required,omitempty"`
	Soft      []string `yaml:"soft,omitempty"`
	Conflicts []string `yaml:"conflicts,omitempty"`
	Provides  []string `yaml:"provides,omitempty"`
}

// DockerConfig represents Docker configuration (V1 format)
type DockerConfig struct {
	Image         string             `yaml:"image,omitempty"`
	Ports         []string           `yaml:"ports,omitempty"`
	Environment   []string           `yaml:"environment,omitempty"`
	Volumes       []any              `yaml:"volumes,omitempty"` // Can be strings or maps
	SimpleVolumes []string           `yaml:"simple_volumes,omitempty"`
	Command       []string           `yaml:"command,omitempty"`
	Restart       string             `yaml:"restart,omitempty"`
	DependsOn     []string           `yaml:"depends_on,omitempty"`
	Networks      []string           `yaml:"networks,omitempty"`
	MemoryLimit   string             `yaml:"memory_limit,omitempty"`
	HealthCheck   *HealthCheckConfig `yaml:"health_check,omitempty"`
}

// HealthCheckConfig represents health check configuration (V1 format)
type HealthCheckConfig struct {
	Test        []string `yaml:"test,omitempty"`
	Interval    string   `yaml:"interval,omitempty"`
	Timeout     string   `yaml:"timeout,omitempty"`
	Retries     int      `yaml:"retries,omitempty"`
	StartPeriod string   `yaml:"start_period,omitempty"`
}

// ManagementV1 represents management operations (V1 format)
type ManagementV1 struct {
	Connect *OperationV1            `yaml:"connect,omitempty"`
	Backup  *OperationV1            `yaml:"backup,omitempty"`
	Restore *OperationV1            `yaml:"restore,omitempty"`
	Custom  map[string]*OperationV1 `yaml:"custom,omitempty"`
}

// OperationV1 represents a management operation (V1 format)
type OperationV1 struct {
	Type        string              `yaml:"type,omitempty"`
	Command     []string            `yaml:"command,omitempty"`
	Args        map[string][]string `yaml:"args,omitempty"`
	Defaults    map[string]string   `yaml:"defaults,omitempty"`
	PreCommands map[string][]string `yaml:"pre_commands,omitempty"`
	Extension   string              `yaml:"extension,omitempty"`
}

// ConnectionConfig represents connection configuration (V1 format)
type ConnectionConfig struct {
	Client      string   `yaml:"client,omitempty"`
	DefaultUser string   `yaml:"default_user,omitempty"`
	DefaultPort int      `yaml:"default_port,omitempty"`
	UserFlag    string   `yaml:"user_flag,omitempty"`
	HostFlag    string   `yaml:"host_flag,omitempty"`
	PortFlag    string   `yaml:"port_flag,omitempty"`
	DBFlag      string   `yaml:"database_flag,omitempty"`
	ExtraFlags  []string `yaml:"extra_flags,omitempty"`
}

// ServiceV2 is an alias for the V2 format
type ServiceV2 = types.ServiceConfigV2
