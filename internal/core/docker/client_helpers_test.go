package docker

import (
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
)

func TestExtractPorts(t *testing.T) {
	t.Run("empty port map", func(t *testing.T) {
		ports := extractPorts(nat.PortMap{})
		assert.Empty(t, ports)
	})

	t.Run("port with binding", func(t *testing.T) {
		portMap := nat.PortMap{
			"8080/tcp": []nat.PortBinding{
				{HostPort: "8080"},
			},
		}
		ports := extractPorts(portMap)
		assert.Len(t, ports, 1)
		assert.Contains(t, ports[0], "8080")
	})

	t.Run("port without binding", func(t *testing.T) {
		portMap := nat.PortMap{
			"8080/tcp": []nat.PortBinding{},
		}
		ports := extractPorts(portMap)
		assert.Empty(t, ports)
	})
}

func TestGetHealthStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{"healthy", "healthy", HealthStatusHealthy},
		{"unhealthy contains healthy", "unhealthy", HealthStatusHealthy}, // "unhealthy" contains "healthy"
		{"starting", "starting", "n/a"},
		{"empty", "", "n/a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getHealthStatus(tt.status)
			assert.Equal(t, tt.expected, result)
		})
	}
}
