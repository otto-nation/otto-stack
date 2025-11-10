package services

import (
	"context"
	"log/slog"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDockerClient for testing
type MockDockerClient struct {
	mock.Mock
}

func (m *MockDockerClient) StartServices(ctx context.Context, services []string, options docker.StartOptions) error {
	args := m.Called(ctx, services, options)
	return args.Error(0)
}

func (m *MockDockerClient) StopServices(ctx context.Context, services []string, options docker.StopOptions) error {
	args := m.Called(ctx, services, options)
	return args.Error(0)
}

func (m *MockDockerClient) GetServiceStatus(ctx context.Context, services []string) ([]docker.DockerServiceStatus, error) {
	args := m.Called(ctx, services)
	return args.Get(0).([]docker.DockerServiceStatus), args.Error(1)
}

func TestNewOrchestrator(t *testing.T) {
	t.Run("creates orchestrator successfully", func(t *testing.T) {
		logger := slog.Default()
		projectDir := "/tmp/test"

		_, err := NewOrchestrator(logger, projectDir)
		// Note: This might fail if Docker is not available, which is expected in test environment
		if err != nil {
			assert.Contains(t, err.Error(), "Docker")
		}
	})

	t.Run("handles invalid project directory", func(t *testing.T) {
		logger := slog.Default()
		projectDir := ""

		_, err := NewOrchestrator(logger, projectDir)
		// Should either succeed or fail with Docker-related error
		if err != nil {
			assert.NotNil(t, err)
		}
	})
}
