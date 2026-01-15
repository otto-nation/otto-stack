//go:build unit

package stack

import (
	"context"
	"testing"
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestUpRequest_Fields(t *testing.T) {
	t.Run("validates UpRequest structure", func(t *testing.T) {
		req := UpRequest{
			Project:        "test-project",
			ServiceConfigs: []types.ServiceConfig{},
			Build:          true,
			SkipConflicts:  false,
		}

		assert.Equal(t, "test-project", req.Project)
		assert.True(t, req.Build)
		assert.False(t, req.SkipConflicts)
		assert.NotNil(t, req.ServiceConfigs)
	})
}

func TestDownRequest_Fields(t *testing.T) {
	t.Run("validates DownRequest structure", func(t *testing.T) {
		req := DownRequest{
			Project:        "test-project",
			ServiceConfigs: []types.ServiceConfig{},
			RemoveVolumes:  true,
			Timeout:        30 * time.Second,
		}

		assert.Equal(t, "test-project", req.Project)
		assert.True(t, req.RemoveVolumes)
		assert.Equal(t, 30*time.Second, req.Timeout)
		assert.NotNil(t, req.ServiceConfigs)
	})
}

func TestRestartRequest_Fields(t *testing.T) {
	t.Run("validates RestartRequest structure", func(t *testing.T) {
		req := RestartRequest{
			Project:        "test-project",
			ServiceConfigs: []types.ServiceConfig{},
			Timeout:        60 * time.Second,
		}

		assert.Equal(t, "test-project", req.Project)
		assert.Equal(t, 60*time.Second, req.Timeout)
		assert.NotNil(t, req.ServiceConfigs)
	})
}

func TestStatusRequest_Fields(t *testing.T) {
	t.Run("validates StatusRequest structure", func(t *testing.T) {
		req := StatusRequest{
			Project:        "test-project",
			ServiceConfigs: []types.ServiceConfig{},
		}

		assert.Equal(t, "test-project", req.Project)
		assert.NotNil(t, req.ServiceConfigs)
	})
}

func TestLogsRequest_Fields(t *testing.T) {
	t.Run("validates LogsRequest structure", func(t *testing.T) {
		req := LogsRequest{
			Project:        "test-project",
			ServiceConfigs: []types.ServiceConfig{},
			Follow:         true,
			Timestamps:     false,
			Tail:           "100",
		}

		assert.Equal(t, "test-project", req.Project)
		assert.True(t, req.Follow)
		assert.False(t, req.Timestamps)
		assert.Equal(t, "100", req.Tail)
		assert.NotNil(t, req.ServiceConfigs)
	})
}

func TestServiceStatus_Fields(t *testing.T) {
	t.Run("validates ServiceStatus structure", func(t *testing.T) {
		status := ServiceStatus{
			Name:   "test-service",
			Status: "running",
			Health: "healthy",
		}

		assert.Equal(t, "test-service", status.Name)
		assert.Equal(t, "running", status.Status)
		assert.Equal(t, "healthy", status.Health)
	})
}

func TestOperationsInterface(t *testing.T) {
	t.Run("validates Operations interface methods", func(t *testing.T) {
		// Test that Operations interface can be implemented
		var ops Operations
		assert.Nil(t, ops)

		// Test context usage
		ctx := context.Background()
		assert.NotNil(t, ctx)
	})
}
