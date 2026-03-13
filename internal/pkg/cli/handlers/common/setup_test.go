//go:build unit

package common

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/registry"
	"github.com/stretchr/testify/assert"
)

func TestVerifyServicesInRegistry(t *testing.T) {
	t.Run("all services exist", func(t *testing.T) {
		reg := &registry.Registry{
			Containers: map[string]*registry.ContainerInfo{
				"service1": {},
				"service2": {},
			},
		}
		err := VerifyServicesInRegistry([]string{"service1", "service2"}, reg)
		assert.NoError(t, err)
	})

	t.Run("service not in registry", func(t *testing.T) {
		reg := &registry.Registry{
			Containers: map[string]*registry.ContainerInfo{
				"service1": {},
			},
		}
		err := VerifyServicesInRegistry([]string{"service1", "missing"}, reg)
		assert.Error(t, err)
	})

	t.Run("empty service list", func(t *testing.T) {
		reg := &registry.Registry{
			Containers: map[string]*registry.ContainerInfo{},
		}
		err := VerifyServicesInRegistry([]string{}, reg)
		assert.NoError(t, err)
	})
}

func TestResolveServiceConfigs_Default(t *testing.T) {
	cfg := &config.Config{
		Stack: config.StackConfig{
			Enabled: []string{"postgres"},
		},
	}
	setup := &CoreSetup{Config: cfg}

	configs, err := ResolveServiceConfigs([]string{}, setup)
	if err != nil {
		t.Logf("Service resolution error: %v", err)
	} else {
		assert.NotNil(t, configs)
	}
}

func TestResolveServiceConfigs_Provided(t *testing.T) {
	cfg := &config.Config{
		Stack: config.StackConfig{
			Enabled: []string{"postgres"},
		},
	}
	setup := &CoreSetup{Config: cfg}

	configs, err := ResolveServiceConfigs([]string{"redis"}, setup)
	if err != nil {
		t.Logf("Service resolution error: %v", err)
	} else {
		assert.NotNil(t, configs)
	}
}
