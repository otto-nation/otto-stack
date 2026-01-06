package stack

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/stretchr/testify/assert"
)

func TestRestartHandler_ResolveServiceNames(t *testing.T) {

	tests := []struct {
		name            string
		args            []string
		enabledServices []string
		expected        []string
	}{
		{
			name:            "use args when provided",
			args:            []string{services.ServicePostgres, services.ServiceRedis},
			enabledServices: []string{services.ServicePostgres, services.ServiceRedis, services.ServiceKafka},
			expected:        []string{services.ServicePostgres, services.ServiceRedis},
		},
		{
			name:            "use enabled services when no args",
			args:            []string{},
			enabledServices: []string{services.ServicePostgres, services.ServiceRedis},
			expected:        []string{services.ServicePostgres, services.ServiceRedis},
		},
		{
			name:            "use enabled services when nil args",
			args:            nil,
			enabledServices: []string{services.ServicePostgres},
			expected:        []string{services.ServicePostgres},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the shared ResolveServiceConfigs function instead
			setup := &CoreSetup{Config: &config.Config{Stack: config.StackConfig{Enabled: tt.enabledServices}}}
			configs, err := ResolveServiceConfigs(tt.args, setup)
			assert.NoError(t, err)

			result := make([]string, len(configs))
			for i, config := range configs {
				result[i] = config.Name
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewRestartHandler(t *testing.T) {
	handler := NewRestartHandler()
	assert.NotNil(t, handler)
}
