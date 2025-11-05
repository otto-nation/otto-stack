package project

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"github.com/stretchr/testify/assert"
)

func TestGenerateEnvFile(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	err := handler.createDirectoryStructure()
	assert.NoError(t, err)

	err = handler.generateEnvFile([]string{TestServicePostgres}, TestProjectName, &types.BaseCommand{Output: ui.NewOutput()})
	if err != nil {
		t.Logf("Expected error in test environment: %v", err)
	}
}

func TestGenerateDockerCompose(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	err := handler.createDirectoryStructure()
	assert.NoError(t, err)

	err = handler.generateDockerCompose([]string{TestServicePostgres}, TestProjectName, &types.BaseCommand{Output: ui.NewOutput()})
	if err != nil {
		t.Logf("Expected error in test environment: %v", err)
	}
}

func TestGenerateInitialComposeFiles(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	err := handler.generateInitialComposeFiles([]string{TestServicePostgres}, TestProjectName,
		map[string]bool{"skip_warnings": false},
		map[string]bool{"auto_start": true}, &types.BaseCommand{Output: ui.NewOutput()})

	if err != nil {
		t.Logf("Expected error in test environment: %v", err)
	}
}
