package operations

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestGetContainerName(t *testing.T) {
	tests := []struct {
		name     string
		config   types.ServiceConfig
		expected string
	}{
		{
			name:     "hidden service returns name",
			config:   types.ServiceConfig{Name: "test-service", Hidden: true},
			expected: "test-service",
		},
		{
			name: "service with dependency returns dependency",
			config: types.ServiceConfig{
				Name: "test-service",
				Service: types.ServiceSpec{
					Dependencies: types.DependenciesSpec{
						Required: []string{"dep-service"},
					},
				},
			},
			expected: "dep-service",
		},
		{
			name:     "service without dependency returns name",
			config:   types.ServiceConfig{Name: "test-service"},
			expected: "test-service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getContainerName(tt.config)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilterInitContainers(t *testing.T) {
	configs := []types.ServiceConfig{
		{Name: "service1", Container: types.ContainerSpec{Restart: types.RestartPolicyAlways}},
		{Name: "service2", Container: types.ContainerSpec{Restart: types.RestartPolicyNo}},
		{Name: "service3", Container: types.ContainerSpec{Restart: types.RestartPolicyOnFailure}},
	}

	result := filterInitContainers(configs)
	assert.Len(t, result, 2)
	assert.Contains(t, result, "service1")
	assert.Contains(t, result, "service3")
	assert.NotContains(t, result, "service2")
}
