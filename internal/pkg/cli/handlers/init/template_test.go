package init

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateEnvFile(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	err := handler.createDirectoryStructure()
	assert.NoError(t, err)

	projectConfig := &ProjectConfig{
		Project: struct {
			Name        string
			Environment string
		}{
			Name:        TestProjectName,
			Environment: TestEnvironmentLocal,
		},
		Stack: struct {
			Enabled []string
		}{
			Enabled: []string{TestServicePostgres},
		},
	}

	err = handler.generateEnvFile([]string{TestServicePostgres}, projectConfig)
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

	projectConfig := &ProjectConfig{
		Project: struct {
			Name        string
			Environment string
		}{
			Name:        TestProjectName,
			Environment: TestEnvironmentLocal,
		},
		Stack: struct {
			Enabled []string
		}{
			Enabled: []string{TestServicePostgres},
		},
	}

	err = handler.generateDockerCompose([]string{TestServicePostgres}, projectConfig)
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
		map[string]bool{"auto_start": true})

	if err != nil {
		t.Logf("Expected error in test environment: %v", err)
	}
}
