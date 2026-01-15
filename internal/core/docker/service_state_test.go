//go:build unit

package docker

import (
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestToSDK_Methods(t *testing.T) {
	t.Run("validates UpOptions ToSDK conversion", func(t *testing.T) {
		opts := UpOptions{
			Build:         true,
			ForceRecreate: true,
			Services:      []string{"postgres"},
		}

		result := opts.ToSDK()
		assert.NotNil(t, result)
	})

	t.Run("validates DownOptions ToSDK conversion", func(t *testing.T) {
		opts := DownOptions{
			Services:      []string{"postgres"},
			RemoveVolumes: true,
		}

		result := opts.ToSDK()
		assert.NotNil(t, result)
	})

	t.Run("validates StopOptions ToSDK conversion", func(t *testing.T) {
		opts := StopOptions{
			Services: []string{"postgres"},
		}

		result := opts.ToSDK()
		assert.NotNil(t, result)
	})

	t.Run("validates LogOptions ToSDK conversion", func(t *testing.T) {
		opts := LogOptions{
			Services:   []string{"postgres"},
			Follow:     true,
			Timestamps: true,
		}

		result := opts.ToSDK()
		assert.NotNil(t, result)
	})
}

func TestDockerServiceState_IsRunning(t *testing.T) {
	testCases := []struct {
		name     string
		state    DockerServiceState
		expected bool
	}{
		{"running state", DockerServiceStateRunning, true},
		{"stopped state", DockerServiceStateStopped, false},
		{"created state", DockerServiceStateCreated, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.state.IsRunning()
			assert.Equal(t, tc.expected, result, "IsRunning result for %v", tc.state)
		})
	}
}

func TestNewProjectFilter(t *testing.T) {
	t.Run("creates project filter", func(t *testing.T) {
		filter := NewProjectFilter("test-project")
		assert.NotNil(t, filter)
	})
}

func TestNewServiceFilter(t *testing.T) {
	t.Run("creates service filter", func(t *testing.T) {
		filter := NewServiceFilter("test-project", "test-service")
		assert.NotNil(t, filter)
	})
}

func TestServiceCharacteristicsResolver(t *testing.T) {
	t.Run("creates service characteristics resolver", func(t *testing.T) {
		resolver, err := NewServiceCharacteristicsResolver()
		testhelpers.AssertValidConstructor(t, resolver, err, "ServiceCharacteristicsResolver")
	})

	t.Run("resolves compose up flags", func(t *testing.T) {
		resolver, err := NewServiceCharacteristicsResolver()
		testhelpers.AssertValidConstructor(t, resolver, err, "ServiceCharacteristicsResolver")

		flags := resolver.ResolveComposeUpFlags([]string{"postgres"})
		assert.NotNil(t, flags)
	})

	t.Run("resolves compose down flags", func(t *testing.T) {
		resolver, err := NewServiceCharacteristicsResolver()
		testhelpers.AssertValidConstructor(t, resolver, err, "ServiceCharacteristicsResolver")

		flags := resolver.ResolveComposeDownFlags([]string{"postgres"})
		assert.NotNil(t, flags)
	})
}
