package compose

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestLabelGeneration(t *testing.T) {
	manager, err := services.New()
	require.NoError(t, err)

	gen, err := NewGenerator("test-project", "", manager)
	require.NoError(t, err)

	// Generate YAML for redis service
	yamlBytes, err := gen.GenerateYAML([]string{services.ServiceRedis})
	require.NoError(t, err)

	// Parse the YAML
	var compose map[string]any
	err = yaml.Unmarshal(yamlBytes, &compose)
	require.NoError(t, err)

	// Get services
	services, ok := compose["services"].(map[string]any)
	require.True(t, ok, "services should be a map")

	// Get redis service
	redis, ok := services["redis"].(map[string]any)
	require.True(t, ok, "redis service should exist")

	// Check labels
	labels, ok := redis["labels"].(map[string]any)
	require.True(t, ok, "labels should exist")

	// Verify Otto Stack labels
	assert.Equal(t, "true", labels["io.otto-stack.managed"])
	assert.Equal(t, "test-project", labels["io.otto-stack.project"])
	assert.Equal(t, "redis", labels["io.otto-stack.service"])
	assert.Equal(t, "isolated", labels["io.otto-stack.sharing-mode"])
	assert.NotEmpty(t, labels["io.otto-stack.version"])
}
