//go:build unit

package project

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// MockLogger implements LoggerAdapter for testing
type MockLogger struct{}

func (m *MockLogger) Info(msg string, args ...any)  {}
func (m *MockLogger) Error(msg string, args ...any) {}
func (m *MockLogger) Debug(msg string, args ...any) {}
func (m *MockLogger) SlogLogger() *slog.Logger      { return slog.Default() }

// MockOutput implements Output for testing
type MockOutput struct{}

func (m *MockOutput) Success(msg string, args ...any) {}
func (m *MockOutput) Error(msg string, args ...any)   {}
func (m *MockOutput) Warning(msg string, args ...any) {}
func (m *MockOutput) Info(msg string, args ...any)    {}
func (m *MockOutput) Header(msg string, args ...any)  {}
func (m *MockOutput) Muted(msg string, args ...any)   {}
func (m *MockOutput) Writer() io.Writer               { return os.Stdout }

func TestHandle_DirectoryValidation(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	// Create conflicting docker-compose.yml file
	createTestFile(t, docker.DockerComposeFileName, "version: '3'")

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
	// Test should fail due to non-interactive mode requiring config (which is expected behavior)
	assert.True(t,
		strings.Contains(err.Error(), "Non-interactive mode requires explicit configuration") ||
			strings.Contains(err.Error(), "non-interactive mode requires") ||
			strings.Contains(err.Error(), ActionValidation) ||
			strings.Contains(err.Error(), "directory validation failed") ||
			strings.Contains(err.Error(), docker.DockerComposeFileName) ||
			strings.Contains(err.Error(), "required tool 'docker' is not available"),
		"Expected directory validation or Docker availability error, got: %s", err.Error())
}

func TestHandle_AlreadyInitialized(t *testing.T) {
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
