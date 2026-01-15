//go:build unit

package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"github.com/stretchr/testify/assert"
)

func TestNewInitHandler(t *testing.T) {
	handler := NewInitHandler()
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.projectManager)
}

func TestInitHandler_ValidateArgs(t *testing.T) {
	handler := NewInitHandler()

	tests := []struct {
		name string
		args []string
	}{
		{"no args", []string{}},
		{"with args", []string{"arg1", "arg2"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.ValidateArgs(tt.args)
			assert.NoError(t, err)
		})
	}
}

func TestInitHandler_GetRequiredFlags(t *testing.T) {
	handler := NewInitHandler()
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags)
}

func TestCreateGitignoreEntries_ExistingContent(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	// Create directory structure first
	err := handler.projectManager.directoryManager.CreateDirectoryStructure()
	assert.NoError(t, err)

	// Create .gitignore with existing content
	createTestFile(t, core.GitIgnoreFileName, TestGitignoreContent)

	err = handler.projectManager.createGitignoreEntries(&base.BaseCommand{Output: ui.NewOutput()})
	assert.NoError(t, err)

	gitignorePath := filepath.Join(core.OttoStackDir, core.GitIgnoreFileName)
	content, err := os.ReadFile(gitignorePath)
	assert.NoError(t, err)
	assert.Contains(t, string(content), core.OttoStackDir+"/")
}

func TestCreateReadme_WithServices(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	err := handler.projectManager.directoryManager.CreateDirectoryStructure()
	assert.NoError(t, err)

	serviceConfigs := []types.ServiceConfig{{Name: services.ServicePostgres}, {Name: services.ServiceRedis}}
	err = handler.projectManager.createReadme(TestProjectName, serviceConfigs, &base.BaseCommand{Output: ui.NewOutput()})
	assert.NoError(t, err)

	readmePath := filepath.Join(core.OttoStackDir, core.ReadmeFileName)
	content, err := os.ReadFile(readmePath)
	assert.NoError(t, err)
	assert.Contains(t, string(content), TestProjectName)
	assert.Contains(t, string(content), services.ServicePostgres)
	assert.Contains(t, string(content), services.ServiceRedis)
}
