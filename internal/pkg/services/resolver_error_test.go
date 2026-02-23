//go:build unit

package services

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceResolver_ErrorMessages_NoDoubleWrapping(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)
	resolver := NewServiceResolver(manager)

	tests := []struct {
		name             string
		serviceNames     []string
		shouldContain    string
		shouldNotContain string
	}{
		{
			name:             "invalid service returns clean error",
			serviceNames:     []string{"nonexistent-service"},
			shouldContain:    "nonexistent-service",
			shouldNotContain: "%v",
		},
		{
			name:             "empty list returns clean error",
			serviceNames:     []string{},
			shouldContain:    messages.ValidationNoServicesSelected,
			shouldNotContain: "%v",
		},
		{
			name:             "multiple invalid services",
			serviceNames:     []string{"invalid1", "invalid2"},
			shouldContain:    "invalid1",
			shouldNotContain: "%v",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := resolver.ResolveServices(tt.serviceNames)
			require.Error(t, err)

			errMsg := err.Error()

			if tt.shouldNotContain != "" {
				assert.NotContains(t, errMsg, tt.shouldNotContain)
			}

			if tt.shouldContain != "" {
				assert.Contains(t, errMsg, tt.shouldContain)
			}

			// Check for duplicate error wrapping patterns
			assert.NotContains(t, errMsg, "error: error:")
			assert.NotContains(t, errMsg, "failed: failed:")
		})
	}
}

func TestServiceResolver_ErrorMessages_ValidateAgainstMessagesYAML(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)
	resolver := NewServiceResolver(manager)

	t.Run("unknown service error uses correct message", func(t *testing.T) {
		_, err := resolver.ResolveServices([]string{"unknown-service-xyz"})
		require.Error(t, err)

		errMsg := err.Error()
		assert.Contains(t, errMsg, "unknown-service-xyz")
		assert.NotContains(t, errMsg, "%v")
		assert.NotContains(t, errMsg, "%s")
	})

	t.Run("empty services error uses correct message", func(t *testing.T) {
		_, err := resolver.ResolveServices([]string{})
		require.Error(t, err)

		errMsg := err.Error()
		assert.Contains(t, errMsg, "at least one service must be selected")
		assert.NotContains(t, errMsg, "%v")
	})
}
