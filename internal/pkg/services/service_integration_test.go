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
