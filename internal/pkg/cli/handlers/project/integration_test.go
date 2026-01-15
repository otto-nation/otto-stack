//go:build unit

package project

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestHandle_ValidationFailure(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	createTestConfig(t)

	handler := NewInitHandler()
	cmd := &cobra.Command{}
	cmd.Flags().Bool("force", false, "force initialization")
	cmd.Flags().Bool("non-interactive", true, "non-interactive mode")
	cmd.Flag("non-interactive").Value.Set("true")

	base := &base.BaseCommand{
		Logger: &MockLogger{},
		Output: &MockOutput{},
	}

	err := handler.Handle(context.Background(), cmd, []string{}, base)
	assert.Error(t, err)
	assert.True(t,
		strings.Contains(err.Error(), "Non-interactive mode requires explicit configuration") ||
			strings.Contains(err.Error(), "non-interactive mode requires") ||
			strings.Contains(err.Error(), ActionValidation) ||
			strings.Contains(err.Error(), "already initialized"),
		"Expected validation or initialization error, got: %s", err.Error())
}

func TestHandle_WithForceFlag(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	handler := NewInitHandler()
	cmd := &cobra.Command{}
	cmd.Flags().Bool("force", true, "force initialization")

	base := &base.BaseCommand{
		Logger: &MockLogger{},
		Output: &MockOutput{},
	}

	err := handler.Handle(context.Background(), cmd, []string{}, base)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), MsgFailedToGetProjectDetails)
}

func TestCreateDirectoryStructure(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	err := handler.projectManager.directoryManager.CreateDirectoryStructure()
	assert.NoError(t, err)

	_, err = os.Stat(core.OttoStackDir)
	assert.NoError(t, err)
}

func TestCreateConfigFile(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	err := handler.projectManager.directoryManager.CreateDirectoryStructure()
	assert.NoError(t, err)

	originalServiceNames := []string{services.ServicePostgres}
	err = handler.projectManager.configManager.CreateConfigFile(TestProjectName, originalServiceNames, nil,
		&base.BaseCommand{Output: ui.NewOutput()})
	assert.NoError(t, err)

	_, err = os.Stat(TestConfigFilePath)
	assert.NoError(t, err)
}

func TestCreateGitignoreEntries(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	// Create directory structure first
	err := handler.projectManager.directoryManager.CreateDirectoryStructure()
	assert.NoError(t, err)

	err = handler.projectManager.createGitignoreEntries(&base.BaseCommand{Output: ui.NewOutput()})
	assert.NoError(t, err)

	gitignorePath := filepath.Join(core.OttoStackDir, core.GitIgnoreFileName)
	_, err = os.Stat(gitignorePath)
	assert.NoError(t, err)
}

func TestCreateReadme(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	err := handler.projectManager.directoryManager.CreateDirectoryStructure()
	assert.NoError(t, err)

	serviceConfigs := []types.ServiceConfig{{Name: services.ServicePostgres}, {Name: services.ServiceRedis}}
	err = handler.projectManager.createReadme(TestProjectName, serviceConfigs, &base.BaseCommand{Output: ui.NewOutput()})
	assert.NoError(t, err)

	readmePath := filepath.Join(core.OttoStackDir, core.ReadmeFileName)
	_, err = os.Stat(readmePath)
	assert.NoError(t, err)
}
