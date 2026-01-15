//go:build unit

package compose

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/test/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLabelGeneration(t *testing.T) {
	manager, err := services.New()
	testhelpers.AssertValidConstructor(t, manager, err, "Services Manager")

	gen, err := NewGenerator("test-project", "", manager)
	testhelpers.AssertValidConstructor(t, gen, err, "Generator")

	// Test that labels are properly generated in compose structure
	compose, err := gen.buildComposeStructure([]types.ServiceConfig{})
	require.NoError(t, err)

	// Check network labels exist
	networks, ok := compose["networks"].(map[string]any)
	require.True(t, ok, "networks should be a map")

	defaultNet, ok := networks["default"].(map[string]any)
	require.True(t, ok, "default network should exist")

	labels, ok := defaultNet["labels"].(map[string]string)
	require.True(t, ok, "labels should exist")

	// Verify Otto Stack labels on network
	assert.Equal(t, "true", labels["io.otto-stack.managed"])
	assert.Equal(t, "test-project", labels["io.otto-stack.project"])
	assert.Equal(t, "network", labels["io.otto-stack.service"])
	assert.Equal(t, "isolated", labels["io.otto-stack.sharing-mode"])
	assert.NotEmpty(t, labels["io.otto-stack.version"])
}
