//go:build unit

package project

import (
	"testing"

	"github.com/stretchr/testify/assert"

	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
)

func TestConflictsHandler_ValidateArgs(t *testing.T) {
	handler := &ConflictsHandler{}

	t.Run("accepts no arguments", func(t *testing.T) {
		err := handler.ValidateArgs([]string{})
		assert.NoError(t, err)
	})

	t.Run("accepts service names", func(t *testing.T) {
		err := handler.ValidateArgs([]string{"postgres", "redis"})
		assert.NoError(t, err)
	})
}

func TestConflictsHandler_GetRequiredFlags(t *testing.T) {
	handler := &ConflictsHandler{}
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags)
}

func TestConflictsHandler_ParsePort(t *testing.T) {
	handler := &ConflictsHandler{}

	t.Run("parses valid port", func(t *testing.T) {
		port := handler.parsePort("8080")
		assert.Equal(t, 8080, port)
	})

	t.Run("returns 0 for invalid port", func(t *testing.T) {
		port := handler.parsePort("invalid")
		assert.Equal(t, 0, port)
	})

	t.Run("returns 0 for empty string", func(t *testing.T) {
		port := handler.parsePort("")
		assert.Equal(t, 0, port)
	})
}

func TestConflictsHandler_ExtractPortsFromService(t *testing.T) {
	handler := &ConflictsHandler{}

	t.Run("extracts ports from service", func(t *testing.T) {
		service := &servicetypes.ServiceConfig{
			Container: servicetypes.ContainerSpec{
				Ports: []servicetypes.PortSpec{
					{External: "8080"},
					{External: "9090"},
				},
			},
		}
		ports := handler.extractPortsFromService(service)
		assert.Len(t, ports, 2)
		assert.Contains(t, ports, 8080)
		assert.Contains(t, ports, 9090)
	})

	t.Run("skips invalid ports", func(t *testing.T) {
		service := &servicetypes.ServiceConfig{
			Container: servicetypes.ContainerSpec{
				Ports: []servicetypes.PortSpec{
					{External: "8080"},
					{External: "invalid"},
				},
			},
		}
		ports := handler.extractPortsFromService(service)
		assert.Len(t, ports, 1)
		assert.Contains(t, ports, 8080)
	})

	t.Run("returns empty for no ports", func(t *testing.T) {
		service := &servicetypes.ServiceConfig{}
		ports := handler.extractPortsFromService(service)
		assert.Empty(t, ports)
	})
}

func TestDepsHandler_ValidateArgs(t *testing.T) {
	handler := &DepsHandler{}

	t.Run("accepts no arguments", func(t *testing.T) {
		err := handler.ValidateArgs([]string{})
		assert.NoError(t, err)
	})

	t.Run("accepts service names", func(t *testing.T) {
		err := handler.ValidateArgs([]string{"postgres"})
		assert.NoError(t, err)
	})
}

func TestDepsHandler_GetRequiredFlags(t *testing.T) {
	handler := &DepsHandler{}
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags)
}
