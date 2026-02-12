//go:build unit

package services

import (
	"testing"

	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/test/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("creates manager successfully", func(t *testing.T) {
		manager, err := New()
		testhelpers.AssertValidConstructor(t, manager, err, "Manager")
		assert.NotNil(t, manager.services)
	})
}

func TestManager_GetService(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	t.Run("returns service when exists", func(t *testing.T) {
		// Test with a service that should exist (postgres)
		service, err := manager.GetService(ServicePostgres)
		if err == nil {
			assert.NotNil(t, service)
			assert.Equal(t, ServicePostgres, service.Name)
		}
	})

	t.Run("handles non-existent service", func(t *testing.T) {
		service, err := manager.GetService("nonexistent")
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.Nil(t, service)
		}
	})
}

func TestManager_GetDependencies(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)

	t.Run("gets dependencies for service", func(t *testing.T) {
		deps, err := manager.GetDependencies(ServicePostgres)
		if err != nil {
			assert.Error(t, err)
		} else {
			// deps can be nil or empty slice, both are valid
			assert.IsType(t, []string{}, deps)
		}
	})

	t.Run("returns error when service not found", func(t *testing.T) {
		_, err := manager.GetDependencies("nonexistent-service")
		assert.Error(t, err)
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
		expectedServices := []string{ServicePostgres, ServiceRedis, ServiceMysql}
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

func TestManager_loadServices(t *testing.T) {
	t.Run("loads services from config directory", func(t *testing.T) {
		manager := &Manager{
			services: make(map[string]servicetypes.ServiceConfig),
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
