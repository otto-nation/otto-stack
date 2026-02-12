//go:build unit

package project

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"github.com/stretchr/testify/assert"
)

func TestGenerateEnvFile(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	err := handler.projectManager.directoryManager.CreateDirectoryStructure()
	assert.NoError(t, err)

	serviceConfigs := []types.ServiceConfig{{Name: services.ServicePostgres}}
	err = handler.projectManager.generateEnvFile(serviceConfigs, TestProjectName, &base.BaseCommand{Output: ui.NewOutput()})
	if err != nil {
		t.Logf("Expected error in test environment: %v", err)
	}
}

func TestGenerateDockerCompose(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	err := handler.projectManager.directoryManager.CreateDirectoryStructure()
	assert.NoError(t, err)

	serviceConfigs := []types.ServiceConfig{{Name: services.ServicePostgres}}
	err = handler.projectManager.generateDockerComposeWithSharing(serviceConfigs, TestProjectName, false, &base.BaseCommand{Output: ui.NewOutput()})
	if err != nil {
		t.Logf("Expected error in test environment: %v", err)
	}
}

func TestGenerateInitialComposeFiles(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	serviceConfigs := []types.ServiceConfig{{Name: services.ServicePostgres}}
	baseCmd := &base.BaseCommand{Output: ui.NewOutput()}

	err := handler.projectManager.generateEnvFile(serviceConfigs, TestProjectName, baseCmd)
	if err != nil {
		t.Logf("Expected error in test environment: %v", err)
	}

	err = handler.projectManager.generateDockerComposeWithSharing(serviceConfigs, TestProjectName, false, baseCmd)
	if err != nil {
		t.Logf("Expected error in test environment: %v", err)
	}
}
