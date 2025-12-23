package stack

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/stretchr/testify/assert"
)

func TestRestartHandler_ResolveServiceNames(t *testing.T) {
	handler := NewRestartHandler()

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
			result := handler.resolveServiceNames(tt.args, tt.enabledServices)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewRestartHandler(t *testing.T) {
	handler := NewRestartHandler()
	assert.NotNil(t, handler)
}
