package docker

import (
	"time"

	"github.com/docker/compose/v5/pkg/api"
	"github.com/otto-nation/otto-stack/internal/config"
	"gopkg.in/yaml.v3"
)

// UpOptions represents options for compose up operations
type UpOptions struct {
	Build         bool
	ForceRecreate bool
	RemoveOrphans bool
	NoDeps        bool
	Detach        bool
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
		Services:      o.Services,
		RemoveOrphans: o.RemoveOrphans,
		Timeout:       o.Timeout,
	}
	if o.Build {
		createOpts.Build = &api.BuildOptions{}
	}
	if o.ForceRecreate {
		createOpts.Recreate = api.RecreateForce
	}
	if o.NoDeps {
		createOpts.RecreateDependencies = api.RecreateNever
	}

	startOpts := api.StartOptions{
		Services: o.Services,
		Wait:     !o.Detach, // If detached, don't wait
	}

	return api.UpOptions{
		Create: createOpts,
		Start:  startOpts,
	}
}

// ToSDK converts DownOptions to Docker Compose API options
func (o DownOptions) ToSDK() api.DownOptions {
	return api.DownOptions{
		Services:      o.Services,
		Volumes:       o.RemoveVolumes,
		RemoveOrphans: o.RemoveOrphans,
		Timeout:       o.Timeout,
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

// ServiceCharacteristicsConfig defines Docker behaviors for service characteristics
type ServiceCharacteristicsConfig struct {
	ServiceCharacteristics map[string]ServiceCharacteristic `yaml:"service_characteristics"`
}

// ServiceCharacteristic defines flags for different Docker operations
type ServiceCharacteristic struct {
	ComposeUpFlags   []string `yaml:"compose_up_flags"`
	ComposeDownFlags []string `yaml:"compose_down_flags"`
	RunFlags         []string `yaml:"run_flags"`
}

// ServiceCharacteristicsResolver resolves Docker flags based on service characteristics
type ServiceCharacteristicsResolver struct {
	config *ServiceCharacteristicsConfig
}

// NewServiceCharacteristicsResolver creates a new service characteristics resolver
func NewServiceCharacteristicsResolver() (*ServiceCharacteristicsResolver, error) {
	serviceConfig, err := loadServiceCharacteristicsConfig()
	if err != nil {
		return nil, err
	}

	return &ServiceCharacteristicsResolver{
		config: serviceConfig,
	}, nil
}

// ResolveComposeUpFlags resolves flags for compose up based on service characteristics
func (scr *ServiceCharacteristicsResolver) ResolveComposeUpFlags(characteristics []string) []string {
	flags := []string{}

	// Add characteristic-based flags
	for _, characteristic := range characteristics {
		if serviceChar, exists := scr.config.ServiceCharacteristics[characteristic]; exists {
			flags = append(flags, serviceChar.ComposeUpFlags...)
		}
	}

	return flags
}

// ResolveComposeDownFlags resolves flags for compose down based on service characteristics
func (scr *ServiceCharacteristicsResolver) ResolveComposeDownFlags(characteristics []string) []string {
	flags := []string{}

	// Add characteristic-based flags
	for _, characteristic := range characteristics {
		if serviceChar, exists := scr.config.ServiceCharacteristics[characteristic]; exists {
			flags = append(flags, serviceChar.ComposeDownFlags...)
		}
	}

	return flags
}

// loadServiceCharacteristicsConfig loads the service characteristics configuration
func loadServiceCharacteristicsConfig() (*ServiceCharacteristicsConfig, error) {
	var serviceConfig ServiceCharacteristicsConfig
	if err := yaml.Unmarshal(config.EmbeddedServiceCharacteristicsYAML, &serviceConfig); err != nil {
		return nil, err
	}
	return &serviceConfig, nil
}
