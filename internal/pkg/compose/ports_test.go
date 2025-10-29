package compose

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/stretchr/testify/assert"
)

func TestGetRequiredPorts(t *testing.T) {
	// Create a mock service registry
	registry := &services.ServiceRegistry{}

	// Create generator
	generator := &Generator{
		projectName: "test",
		registry:    registry,
	}

	// Test with empty services
	ports, err := generator.GetRequiredPorts([]string{})
	assert.NoError(t, err)
	assert.Empty(t, ports)
}
