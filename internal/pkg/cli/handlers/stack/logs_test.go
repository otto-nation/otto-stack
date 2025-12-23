package stack

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/stretchr/testify/assert"
)

func TestLogsHandler_ResolveServiceNames(t *testing.T) {
	handler := NewLogsHandler()

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

func TestLogsHandler_ValidateArgs(t *testing.T) {
	handler := NewLogsHandler()

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "no args should be valid",
			args: []string{},
		},
		{
			name: "single service should be valid",
			args: []string{services.ServicePostgres},
		},
		{
			name: "multiple services should be valid",
			args: []string{services.ServicePostgres, services.ServiceRedis},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.ValidateArgs(tt.args)
			assert.NoError(t, err)
		})
	}
}

func TestLogsHandler_GetRequiredFlags(t *testing.T) {
	handler := NewLogsHandler()
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags, "logs handler should not require any flags")
}

func TestNewLogsHandler(t *testing.T) {
	handler := NewLogsHandler()
	assert.NotNil(t, handler)
}
