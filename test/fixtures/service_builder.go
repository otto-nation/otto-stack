package fixtures

import (
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// ServiceConfigBuilder builds ServiceConfig instances for tests
type ServiceConfigBuilder struct {
	config types.ServiceConfig
}

// NewServiceConfig creates a builder with minimal required fields
func NewServiceConfig(name string) *ServiceConfigBuilder {
	return &ServiceConfigBuilder{
		config: types.ServiceConfig{
			Name:      name,
			Container: types.ContainerSpec{},
		},
	}
}

func (b *ServiceConfigBuilder) WithImage(image string) *ServiceConfigBuilder {
	b.config.Container.Image = image
	return b
}

func (b *ServiceConfigBuilder) WithPort(external, internal string) *ServiceConfigBuilder {
	b.config.Container.Ports = append(b.config.Container.Ports, types.PortSpec{
		External: external,
		Internal: internal,
	})
	return b
}

func (b *ServiceConfigBuilder) WithEnv(key, value string) *ServiceConfigBuilder {
	if b.config.Container.Environment == nil {
		b.config.Container.Environment = make(map[string]string)
	}
	b.config.Container.Environment[key] = value
	return b
}

func (b *ServiceConfigBuilder) WithVolume(mount string) *ServiceConfigBuilder {
	b.config.Container.Volumes = append(b.config.Container.Volumes, types.VolumeSpec{Mount: mount})
	return b
}

func (b *ServiceConfigBuilder) WithCategory(category string) *ServiceConfigBuilder {
	b.config.Category = category
	return b
}

func (b *ServiceConfigBuilder) WithServiceType(serviceType types.ServiceType) *ServiceConfigBuilder {
	b.config.ServiceType = serviceType
	return b
}

func (b *ServiceConfigBuilder) WithHealthCheck(test []string, interval, timeout, retries int) *ServiceConfigBuilder {
	b.config.Container.HealthCheck = &types.HealthCheckSpec{
		Test:     test,
		Interval: time.Duration(interval) * time.Second,
		Timeout:  time.Duration(timeout) * time.Second,
		Retries:  retries,
	}
	return b
}

func (b *ServiceConfigBuilder) WithRestart(policy string) *ServiceConfigBuilder {
	b.config.Container.Restart = types.RestartPolicy(policy)
	return b
}

func (b *ServiceConfigBuilder) WithCommand(command []string) *ServiceConfigBuilder {
	b.config.Container.Command = command
	return b
}

func (b *ServiceConfigBuilder) WithMemoryLimit(limit string) *ServiceConfigBuilder {
	b.config.Container.MemoryLimit = limit
	return b
}

func (b *ServiceConfigBuilder) Build() types.ServiceConfig {
	return b.config
}
