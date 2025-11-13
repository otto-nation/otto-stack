package project

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
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
			strings.Contains(err.Error(), "validation failed") ||
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
	assert.Contains(t, err.Error(), "failed to get project details")
}

func TestCreateDirectoryStructure(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	err := handler.createDirectoryStructure()
	assert.NoError(t, err)

	_, err = os.Stat(core.OttoStackDir)
	assert.NoError(t, err)
}

func TestCreateConfigFile(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	err := handler.createDirectoryStructure()
	assert.NoError(t, err)

	err = handler.createConfigFile(TestProjectName, []string{TestServicePostgres}, nil,
		&base.BaseCommand{Output: ui.NewOutput()})
	assert.NoError(t, err)

	_, err = os.Stat(TestConfigFilePath)
	assert.NoError(t, err)
}

func TestCreateGitignoreEntries(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	err := handler.createGitignoreEntries(&base.BaseCommand{Output: ui.NewOutput()})
	assert.NoError(t, err)

	_, err = os.Stat(core.GitIgnoreFileName)
	assert.NoError(t, err)
}

func TestCreateReadme(t *testing.T) {
	handler := NewInitHandler()
	cleanup := setupTestDir(t)
	defer cleanup()

	err := handler.createDirectoryStructure()
	assert.NoError(t, err)

	err = handler.createReadme(TestProjectName, []string{TestServicePostgres, TestServiceRedis}, &base.BaseCommand{Output: ui.NewOutput()})
	assert.NoError(t, err)

	readmePath := filepath.Join(core.OttoStackDir, core.ReadmeFileName)
	_, err = os.Stat(readmePath)
	assert.NoError(t, err)
}

func TestGenerateConfig(t *testing.T) {
	handler := NewInitHandler()

	config := handler.generateConfig(TestProjectName, []string{TestServicePostgres}, nil)

	assert.Contains(t, config, TestProjectName)
	assert.Contains(t, config, TestServicePostgres)
}
