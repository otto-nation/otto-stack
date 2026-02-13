package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindOrphans(t *testing.T) {
	tempDir := t.TempDir()

	manager := NewManager(tempDir)

	t.Run("finds orphans with empty projects", func(t *testing.T) {
		// Create a registry with a container that has no projects
		registry := &Registry{
			Containers: map[string]*ContainerInfo{
				"test-service": {
					Name:     "otto-stack-test",
					Projects: []string{}, // Empty projects = orphan
				},
			},
		}
		err := manager.Save(registry)
		require.NoError(t, err)

		// Find orphans
		orphans, err := manager.FindOrphans()
		require.NoError(t, err)
		assert.Len(t, orphans, 1)
		assert.Equal(t, "test-service", orphans[0].Service)
	})

	t.Run("no orphans when all have projects", func(t *testing.T) {
		registry := &Registry{
			Containers: map[string]*ContainerInfo{
				"test-service": {
					Name:     "otto-stack-test",
					Projects: []string{"project1"},
				},
			},
		}
		err := manager.Save(registry)
		require.NoError(t, err)

		orphans, err := manager.FindOrphans()
		require.NoError(t, err)
		assert.Empty(t, orphans)
	})
}

func TestIsSharedContainer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"shared container", "otto-stack-redis", true},
		{"shared container with suffix", "otto-stack-postgres-1", true},
		{"non-shared container", "my-redis", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSharedContainer(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractServiceName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"basic service", "otto-stack-redis", "redis"},
		{"service with suffix", "otto-stack-postgres-1", "postgres-1"},
		{"service with multiple dashes", "otto-stack-localstack-sns", "localstack-sns"},
		{"non-shared container", "my-redis", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractServiceName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProjectExists(t *testing.T) {
	t.Run("returns false for non-existent project", func(t *testing.T) {
		exists := projectExists("/nonexistent/path/to/project")
		assert.False(t, exists)
	})
}
