package lifecycle

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/registry"
	"github.com/stretchr/testify/assert"
)

func TestRestartHandler_ValidateArgs(t *testing.T) {
	handler := &RestartHandler{}
	assert.NoError(t, handler.ValidateArgs([]string{}))
}

func TestRestartHandler_GetRequiredFlags(t *testing.T) {
	handler := &RestartHandler{}
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags)
}

func TestRestartHandler_verifyServicesInRegistry(t *testing.T) {
	handler := &RestartHandler{}

	t.Run("all services exist", func(t *testing.T) {
		reg := &registry.Registry{
			Containers: map[string]*registry.ContainerInfo{
				"service1": {},
				"service2": {},
			},
		}
		err := handler.verifyServicesInRegistry([]string{"service1", "service2"}, reg)
		assert.NoError(t, err)
	})

	t.Run("service not in registry", func(t *testing.T) {
		reg := &registry.Registry{
			Containers: map[string]*registry.ContainerInfo{
				"service1": {},
			},
		}
		err := handler.verifyServicesInRegistry([]string{"service1", "missing"}, reg)
		assert.Error(t, err)
	})

	t.Run("empty service list", func(t *testing.T) {
		reg := &registry.Registry{
			Containers: map[string]*registry.ContainerInfo{},
		}
		err := handler.verifyServicesInRegistry([]string{}, reg)
		assert.NoError(t, err)
	})
}
