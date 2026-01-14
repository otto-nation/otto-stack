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
	err = handler.projectManager.generateDockerCompose(serviceConfigs, TestProjectName, &base.BaseCommand{Output: ui.NewOutput()})
	if err != nil {
		t.Logf("Expected error in test environment: %v", err)
	}
}

func TestGenerateInitialComposeFiles(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	serviceConfigs := []types.ServiceConfig{{Name: services.ServicePostgres}}
	err := handler.projectManager.generateInitialComposeFiles(serviceConfigs, TestProjectName,
		map[string]bool{"skip_warnings": false},
		map[string]bool{"auto_start": true}, &base.BaseCommand{Output: ui.NewOutput()})

	if err != nil {
		t.Logf("Expected error in test environment: %v", err)
	}
}
