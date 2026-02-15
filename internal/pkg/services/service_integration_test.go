//go:build unit

package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	composetypes "github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/test/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestEnv(t *testing.T) (string, func()) {
	tmpDir := t.TempDir()
	ottoDir := filepath.Join(tmpDir, ".otto-stack")
	require.NoError(t, os.MkdirAll(ottoDir, 0755))

	composeFile := filepath.Join(ottoDir, "docker-compose.yml")
	composeContent := `version: "3.8"
services:
  postgres:
    image: postgres:latest
`
	require.NoError(t, os.WriteFile(composeFile, []byte(composeContent), 0644))

	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)

	return tmpDir, func() { os.Chdir(oldDir) }
}

func TestService_StartWithMocks(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	t.Run("starts services successfully", func(t *testing.T) {
		mockCompose := &testhelpers.MockComposeAPI{
			UpFunc: func(ctx context.Context, project *composetypes.Project, options api.UpOptions) error {
				return nil
			},
		}

		mockLoader := &testhelpers.MockProjectLoader{
			LoadFunc: func(projectName string) (*composetypes.Project, error) {
				return &composetypes.Project{Name: projectName}, nil
			},
		}

		resolver, _ := NewDefaultCharacteristicsResolver()
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := StartRequest{
			Project: "test-project",
			ServiceConfigs: []types.ServiceConfig{
				{Name: "postgres"},
			},
			Timeout: 10 * time.Second,
		}

		err = service.Start(context.Background(), req)
		assert.NoError(t, err)
	})
}

func TestService_StopWithMocks(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	t.Run("stops services successfully", func(t *testing.T) {
		mockCompose := &testhelpers.MockComposeAPI{
			StopFunc: func(ctx context.Context, projectName string, options api.StopOptions) error {
				return nil
			},
		}

		mockLoader := &testhelpers.MockProjectLoader{
			LoadFunc: func(projectName string) (*composetypes.Project, error) {
				return &composetypes.Project{Name: projectName}, nil
			},
		}

		resolver, _ := NewDefaultCharacteristicsResolver()
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := StopRequest{
			Project: "test-project",
			ServiceConfigs: []types.ServiceConfig{
				{Name: "postgres"},
			},
			Timeout: 10 * time.Second,
		}

		err = service.Stop(context.Background(), req)
		assert.NoError(t, err)
	})
}

func TestService_StatusWithMocks(t *testing.T) {
	t.Run("gets service status", func(t *testing.T) {
		mockCompose := &testhelpers.MockComposeAPI{
			PsFunc: func(ctx context.Context, projectName string, options api.PsOptions) ([]api.ContainerSummary, error) {
				return []api.ContainerSummary{
					{Name: "postgres", State: "running"},
				}, nil
			},
		}

		mockLoader := &testhelpers.MockProjectLoader{}
		resolver, _ := NewDefaultCharacteristicsResolver()
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := StatusRequest{
			Project:  "test-project",
			Services: []string{"postgres"},
		}

		statuses, err := service.Status(context.Background(), req)
		assert.NoError(t, err)
		assert.Len(t, statuses, 1)
	})
}

func TestService_LogsWithMocks(t *testing.T) {
	t.Run("streams logs successfully", func(t *testing.T) {
		mockCompose := &testhelpers.MockComposeAPI{
			LogsFunc: func(ctx context.Context, projectName string, consumer api.LogConsumer, options api.LogOptions) error {
				return nil
			},
		}

		mockLoader := &testhelpers.MockProjectLoader{}
		resolver, _ := NewDefaultCharacteristicsResolver()
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := LogRequest{
			Project: "test-project",
			ServiceConfigs: []types.ServiceConfig{
				{Name: "postgres"},
			},
			Follow: false,
		}

		err = service.Logs(context.Background(), req)
		assert.NoError(t, err)
	})
}

func TestService_ExecWithMocks(t *testing.T) {
	t.Run("executes command successfully", func(t *testing.T) {
		mockCompose := &testhelpers.MockComposeAPI{
			ExecFunc: func(ctx context.Context, projectName string, options api.RunOptions) (int, error) {
				return 0, nil
			},
		}

		mockLoader := &testhelpers.MockProjectLoader{}
		resolver, _ := NewDefaultCharacteristicsResolver()
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := ExecRequest{
			Project: "test-project",
			Service: "postgres",
			Command: []string{"psql", "--version"},
		}

		err = service.Exec(context.Background(), req)
		assert.NoError(t, err)
	})
}

func TestService_CleanupWithMocks(t *testing.T) {
	t.Run("cleans up resources successfully", func(t *testing.T) {
		mockCompose := &testhelpers.MockComposeAPI{
			DownFunc: func(ctx context.Context, projectName string, options api.DownOptions) error {
				return nil
			},
		}

		mockLoader := &testhelpers.MockProjectLoader{
			LoadFunc: func(projectName string) (*composetypes.Project, error) {
				return &composetypes.Project{Name: projectName}, nil
			},
		}

		resolver, _ := NewDefaultCharacteristicsResolver()
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := CleanupRequest{
			Project:       "test-project",
			RemoveVolumes: true,
			RemoveImages:  false,
		}

		err = service.Cleanup(context.Background(), req)
		assert.NoError(t, err)
	})
}

func TestService_CheckDockerHealthWithMocks(t *testing.T) {
	t.Run("checks docker health successfully", func(t *testing.T) {
		mockCompose := &testhelpers.MockComposeAPI{}
		mockLoader := &testhelpers.MockProjectLoader{}
		resolver, _ := NewDefaultCharacteristicsResolver()
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		err = service.CheckDockerHealth(context.Background())
		assert.NoError(t, err)
	})
}

func TestService_NewServiceErrorHandling(t *testing.T) {
	t.Run("creates service successfully", func(t *testing.T) {
		mockCompose := &testhelpers.MockComposeAPI{}
		mockLoader := &testhelpers.MockProjectLoader{}
		resolver, _ := NewDefaultCharacteristicsResolver()

		service, err := NewService(mockCompose, resolver, mockLoader)
		assert.NoError(t, err)
		assert.NotNil(t, service)
		assert.NotNil(t, service.DockerClient)
	})
}

func TestService_StartWithBuildFlag(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	t.Run("starts with build flag", func(t *testing.T) {
		buildCalled := false
		mockCompose := &testhelpers.MockComposeAPI{
			UpFunc: func(ctx context.Context, project *composetypes.Project, options api.UpOptions) error {
				buildCalled = true
				return nil
			},
		}

		mockLoader := &testhelpers.MockProjectLoader{
			LoadFunc: func(projectName string) (*composetypes.Project, error) {
				return &composetypes.Project{Name: projectName}, nil
			},
		}

		resolver, _ := NewDefaultCharacteristicsResolver()
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := StartRequest{
			Project: "test-project",
			ServiceConfigs: []types.ServiceConfig{
				{Name: "postgres"},
			},
			Build:   true,
			Timeout: 10 * time.Second,
		}

		err = service.Start(context.Background(), req)
		assert.NoError(t, err)
		assert.True(t, buildCalled)
	})

	t.Run("starts with force recreate", func(t *testing.T) {
		recreateCalled := false
		mockCompose := &testhelpers.MockComposeAPI{
			UpFunc: func(ctx context.Context, project *composetypes.Project, options api.UpOptions) error {
				recreateCalled = true
				return nil
			},
		}

		mockLoader := &testhelpers.MockProjectLoader{
			LoadFunc: func(projectName string) (*composetypes.Project, error) {
				return &composetypes.Project{Name: projectName}, nil
			},
		}

		resolver, _ := NewDefaultCharacteristicsResolver()
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := StartRequest{
			Project: "test-project",
			ServiceConfigs: []types.ServiceConfig{
				{Name: "postgres"},
			},
			ForceRecreate: true,
			Timeout:       10 * time.Second,
		}

		err = service.Start(context.Background(), req)
		assert.NoError(t, err)
		assert.True(t, recreateCalled)
	})
}

func TestService_StartErrorPaths(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	t.Run("handles compose up error", func(t *testing.T) {
		mockCompose := &testhelpers.MockComposeAPI{
			UpFunc: func(ctx context.Context, project *composetypes.Project, options api.UpOptions) error {
				return assert.AnError
			},
		}

		mockLoader := &testhelpers.MockProjectLoader{
			LoadFunc: func(projectName string) (*composetypes.Project, error) {
				return &composetypes.Project{Name: projectName}, nil
			},
		}

		resolver, _ := NewDefaultCharacteristicsResolver()
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := StartRequest{
			Project: "test-project",
			ServiceConfigs: []types.ServiceConfig{
				{Name: "postgres"},
			},
			Timeout: 10 * time.Second,
		}

		err = service.Start(context.Background(), req)
		assert.Error(t, err)
	})

	t.Run("handles project load error", func(t *testing.T) {
		mockCompose := &testhelpers.MockComposeAPI{}

		mockLoader := &testhelpers.MockProjectLoader{
			LoadFunc: func(projectName string) (*composetypes.Project, error) {
				return nil, assert.AnError
			},
		}

		resolver, _ := NewDefaultCharacteristicsResolver()
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := StartRequest{
			Project: "test-project",
			ServiceConfigs: []types.ServiceConfig{
				{Name: "postgres"},
			},
			Timeout: 10 * time.Second,
		}

		err = service.Start(context.Background(), req)
		assert.Error(t, err)
	})
}

func TestService_StopErrorPaths(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	t.Run("handles stop error", func(t *testing.T) {
		mockCompose := &testhelpers.MockComposeAPI{
			StopFunc: func(ctx context.Context, projectName string, options api.StopOptions) error {
				return assert.AnError
			},
		}

		mockLoader := &testhelpers.MockProjectLoader{
			LoadFunc: func(projectName string) (*composetypes.Project, error) {
				return &composetypes.Project{Name: projectName}, nil
			},
		}

		resolver, _ := NewDefaultCharacteristicsResolver()
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := StopRequest{
			Project: "test-project",
			ServiceConfigs: []types.ServiceConfig{
				{Name: "postgres"},
			},
			Timeout: 10 * time.Second,
		}

		err = service.Stop(context.Background(), req)
		assert.Error(t, err)
	})

	t.Run("uses down when remove is true", func(t *testing.T) {
		downCalled := false
		mockCompose := &testhelpers.MockComposeAPI{
			DownFunc: func(ctx context.Context, projectName string, options api.DownOptions) error {
				downCalled = true
				return nil
			},
		}

		mockLoader := &testhelpers.MockProjectLoader{
			LoadFunc: func(projectName string) (*composetypes.Project, error) {
				return &composetypes.Project{Name: projectName}, nil
			},
		}

		resolver, _ := NewDefaultCharacteristicsResolver()
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := StopRequest{
			Project: "test-project",
			ServiceConfigs: []types.ServiceConfig{
				{Name: "postgres"},
			},
			Remove:  true,
			Timeout: 10 * time.Second,
		}

		err = service.Stop(context.Background(), req)
		assert.NoError(t, err)
		assert.True(t, downCalled)
	})
}

func TestService_GenerateComposeFile(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	t.Run("generates compose file successfully", func(t *testing.T) {
		mockCompose := &testhelpers.MockComposeAPI{}
		mockLoader := &testhelpers.MockProjectLoader{}
		resolver, _ := NewDefaultCharacteristicsResolver()
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		configs := []types.ServiceConfig{
			{Name: "postgres", Description: "PostgreSQL"},
		}

		err = service.GenerateComposeFile("test-project", configs)
		assert.NoError(t, err)
	})
}

func TestService_LogsWithOptions(t *testing.T) {
	t.Run("logs with follow option", func(t *testing.T) {
		mockCompose := &testhelpers.MockComposeAPI{
			LogsFunc: func(ctx context.Context, projectName string, consumer api.LogConsumer, options api.LogOptions) error {
				assert.True(t, options.Follow)
				return nil
			},
		}

		mockLoader := &testhelpers.MockProjectLoader{}
		resolver, _ := NewDefaultCharacteristicsResolver()
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := LogRequest{
			Project: "test-project",
			ServiceConfigs: []types.ServiceConfig{
				{Name: "postgres"},
			},
			Follow: true,
		}

		err = service.Logs(context.Background(), req)
		assert.NoError(t, err)
	})

	t.Run("logs with timestamps", func(t *testing.T) {
		mockCompose := &testhelpers.MockComposeAPI{
			LogsFunc: func(ctx context.Context, projectName string, consumer api.LogConsumer, options api.LogOptions) error {
				assert.True(t, options.Timestamps)
				return nil
			},
		}

		mockLoader := &testhelpers.MockProjectLoader{}
		resolver, _ := NewDefaultCharacteristicsResolver()
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := LogRequest{
			Project: "test-project",
			ServiceConfigs: []types.ServiceConfig{
				{Name: "postgres"},
			},
			Timestamps: true,
		}

		err = service.Logs(context.Background(), req)
		assert.NoError(t, err)
	})
}

func TestService_ExecWithOptions(t *testing.T) {
	t.Run("exec with user option", func(t *testing.T) {
		mockCompose := &testhelpers.MockComposeAPI{
			ExecFunc: func(ctx context.Context, projectName string, options api.RunOptions) (int, error) {
				assert.Equal(t, "admin", options.User)
				return 0, nil
			},
		}

		mockLoader := &testhelpers.MockProjectLoader{}
		resolver, _ := NewDefaultCharacteristicsResolver()
		service, err := NewService(mockCompose, resolver, mockLoader)
		require.NoError(t, err)

		req := ExecRequest{
			Project: "test-project",
			Service: "postgres",
			Command: []string{"psql"},
			User:    "admin",
		}

		err = service.Exec(context.Background(), req)
		assert.NoError(t, err)
	})
}
