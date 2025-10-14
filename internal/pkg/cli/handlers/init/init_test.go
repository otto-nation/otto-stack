package init

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func TestNewInitHandler(t *testing.T) {
	handler := NewInitHandler()
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.serviceUtils)
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

	// Create .gitignore with existing content
	createTestFile(t, constants.GitignoreFileName, TestGitignoreContent)

	err := handler.createGitignoreEntries()
	assert.NoError(t, err)

	content, err := os.ReadFile(constants.GitignoreFileName)
	assert.NoError(t, err)
	assert.Contains(t, string(content), constants.DevStackDir+"/")
}

func TestCreateReadme_WithServices(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	err := handler.createDirectoryStructure()
	assert.NoError(t, err)

	err = handler.createReadme(TestProjectName, []string{TestServicePostgres, TestServiceRedis})
	assert.NoError(t, err)

	readmePath := filepath.Join(constants.DevStackDir, constants.ReadmeFileName)
	content, err := os.ReadFile(readmePath)
	assert.NoError(t, err)
	assert.Contains(t, string(content), TestProjectName)
	assert.Contains(t, string(content), TestServicePostgres)
	assert.Contains(t, string(content), TestServiceRedis)
}
