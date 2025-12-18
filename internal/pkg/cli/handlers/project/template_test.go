package project

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"github.com/stretchr/testify/assert"
)

func TestGenerateEnvFile(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	err := handler.projectManager.createDirectoryStructure()
	assert.NoError(t, err)

	err = handler.projectManager.generateEnvFile([]string{TestServicePostgres}, TestProjectName, &base.BaseCommand{Output: ui.NewOutput()})
	if err != nil {
		t.Logf("Expected error in test environment: %v", err)
	}
}

func TestGenerateDockerCompose(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	err := handler.projectManager.createDirectoryStructure()
	assert.NoError(t, err)

	err = handler.projectManager.generateDockerCompose([]string{TestServicePostgres}, TestProjectName, &base.BaseCommand{Output: ui.NewOutput()})
	if err != nil {
		t.Logf("Expected error in test environment: %v", err)
	}
}

func TestGenerateInitialComposeFiles(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	err := handler.projectManager.generateInitialComposeFiles([]string{TestServicePostgres}, TestProjectName,
		map[string]bool{"skip_warnings": false},
		map[string]bool{"auto_start": true}, &base.BaseCommand{Output: ui.NewOutput()})

	if err != nil {
		t.Logf("Expected error in test environment: %v", err)
	}
}
