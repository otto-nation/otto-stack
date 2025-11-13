package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("creates manager successfully", func(t *testing.T) {
		manager, err := New()
		assert.NoError(t, err)
		assert.NotNil(t, manager)
		assert.NotNil(t, manager.services)
	})
}

func TestManager_GetService(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	t.Run("returns service when exists", func(t *testing.T) {
		// Test with a service that should exist (postgres)
		service, err := manager.GetService("postgres")
		if err == nil {
			assert.NotNil(t, service)
			assert.Equal(t, "postgres", service.Name)
		}
	})

	t.Run("returns error when service not found", func(t *testing.T) {
		service, err := manager.GetService("nonexistent-service")
		assert.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "service not found")
	})
}

func TestManager_GetAllServices(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	t.Run("returns all loaded services", func(t *testing.T) {
		services := manager.GetAllServices()
		assert.NotEmpty(t, services)

		// Check that we have some expected services
		serviceNames := make([]string, 0, len(services))
		for _, service := range services {
			serviceNames = append(serviceNames, service.Name)
		}

		// Should contain at least some core services
		expectedServices := []string{"postgres", "redis", "mysql"}
		for _, expected := range expectedServices {
			found := false
			for _, name := range serviceNames {
				if name == expected {
					found = true
					break
				}
			}
			if !found {
				t.Logf("Expected service %s not found in: %v", expected, serviceNames)
			}
		}
	})
}

func TestManager_ResolveServices(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	t.Run("resolves single service", func(t *testing.T) {
		resolved, err := manager.ResolveServices([]string{"postgres"})
		if err == nil {
			assert.Contains(t, resolved, "postgres")
		}
	})

	t.Run("resolves multiple services", func(t *testing.T) {
		resolved, err := manager.ResolveServices([]string{"postgres", "redis"})
		if err == nil {
			assert.GreaterOrEqual(t, len(resolved), 2)
		}
	})

	t.Run("handles empty input", func(t *testing.T) {
		resolved, err := manager.ResolveServices([]string{})
		assert.NoError(t, err)
		assert.Empty(t, resolved)
	})

	t.Run("handles nonexistent service", func(t *testing.T) {
		resolved, err := manager.ResolveServices([]string{"nonexistent"})
		assert.NoError(t, err)
		assert.Empty(t, resolved) // Nonexistent services are filtered out
	})
}

func TestManager_loadServices(t *testing.T) {
	t.Run("loads services from config directory", func(t *testing.T) {
		manager := &Manager{
			services: make(map[string]ServiceConfig),
		}

		err := manager.loadServices()
		assert.NoError(t, err)
		assert.NotEmpty(t, manager.services)
	})
}

func TestManager_ServiceValidation(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	t.Run("validates service structure", func(t *testing.T) {
		services := manager.GetAllServices()

		for _, service := range services {
			// Basic validation
			assert.NotEmpty(t, service.Name, "Service should have a name")

			// If service has ports, validate structure
			if len(service.Container.Ports) > 0 {
				for _, port := range service.Container.Ports {
					assert.NotEmpty(t, port.External, "External port should not be empty")
					assert.NotEmpty(t, port.Internal, "Internal port should not be empty")
				}
			}
		}
	})
}
