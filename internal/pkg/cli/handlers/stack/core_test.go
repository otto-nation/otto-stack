package stack

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLogger implements the Logger interface for testing
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, args ...any) {
	m.Called(msg, args)
}

func (m *MockLogger) Error(msg string, args ...any) {
	m.Called(msg, args)
}

func (m *MockLogger) Debug(msg string, args ...any) {
	m.Called(msg, args)
}

func (m *MockLogger) SlogLogger() *slog.Logger {
	return slog.Default()
}

// MockOutput implements Output for testing
type MockOutput struct{}

func (m *MockOutput) Success(msg string, args ...any) {}
func (m *MockOutput) Error(msg string, args ...any)   {}
func (m *MockOutput) Warning(msg string, args ...any) {}
func (m *MockOutput) Info(msg string, args ...any)    {}
func (m *MockOutput) Header(msg string, args ...any)  {}
func (m *MockOutput) Muted(msg string, args ...any)   {}

func TestNewUpHandler(t *testing.T) {
	handler := NewUpHandler()
	assert.NotNil(t, handler)
	assert.IsType(t, &UpHandler{}, handler)
}

func TestUpHandler_ValidateArgs(t *testing.T) {
	handler := NewUpHandler()

	t.Run("no args", func(t *testing.T) {
		err := handler.ValidateArgs([]string{})
		assert.NoError(t, err)
	})

	t.Run("with args", func(t *testing.T) {
		err := handler.ValidateArgs([]string{"service1", "service2"})
		assert.NoError(t, err)
	})
}

func TestUpHandler_GetRequiredFlags(t *testing.T) {
	handler := NewUpHandler()
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags)
}

func TestUpHandler_Handle_ConfigNotFound(t *testing.T) {
	handler := NewUpHandler()
	mockLogger := &MockLogger{}
	mockOutput := &MockOutput{}

	cmd := &cobra.Command{}
	cmd.Flags().Bool("build", false, "")
	cmd.Flags().Bool("force-recreate", false, "")

	base := &base.BaseCommand{
		Logger: mockLogger,
		Output: mockOutput,
	}

	err := handler.Handle(context.Background(), cmd, []string{}, base)
	assert.Error(t, err)
	// The error might be empty or contain initialization-related message
	if err.Error() != "" {
		assert.True(t,
			strings.Contains(err.Error(), "not initialized") ||
				strings.Contains(err.Error(), "config") ||
				strings.Contains(err.Error(), "failed"),
			"Expected initialization error, got: %s", err.Error())
	}
}

func TestLoadProjectConfig(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("valid config", func(t *testing.T) {
		configContent := `project:
  name: test-project
stack:
  enabled:
    - redis
    - postgres`

		configPath := filepath.Join(tmpDir, "config.yml")
		err := os.WriteFile(configPath, []byte(configContent), core.FilePermReadWrite)
		assert.NoError(t, err)

		cfg, err := LoadProjectConfig(configPath)
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "test-project", cfg.Project.Name)
		assert.Equal(t, []string{"redis", "postgres"}, cfg.Stack.Enabled)
	})

	t.Run("file not found", func(t *testing.T) {
		cfg, err := LoadProjectConfig("/nonexistent/config.yml")
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})

	t.Run("invalid YAML", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, "invalid.yml")
		err := os.WriteFile(configPath, []byte("invalid: yaml: [unclosed"), core.FilePermReadWrite)
		assert.NoError(t, err)

		cfg, err := LoadProjectConfig(configPath)
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.True(t,
			strings.Contains(err.Error(), "failed to parse config") ||
				strings.Contains(err.Error(), "yaml:"),
			"Expected YAML parse error, got: %s", err.Error())
	})

	t.Run("empty config", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, "empty.yml")
		err := os.WriteFile(configPath, []byte(""), core.FilePermReadWrite)
		assert.NoError(t, err)

		cfg, err := LoadProjectConfig(configPath)
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Empty(t, cfg.Project.Name)
		assert.Empty(t, cfg.Stack.Enabled)
	})
}

func TestProjectConfig_Structure(t *testing.T) {
	cfg := config.Config{}

	// Test that struct fields exist and can be set
	cfg.Project.Name = "test"
	cfg.Stack.Enabled = []string{"service1"}

	assert.Equal(t, "test", cfg.Project.Name)
	assert.Equal(t, []string{"service1"}, cfg.Stack.Enabled)
}

func TestMockLogger(t *testing.T) {
	mockLogger := &MockLogger{}

	mockLogger.On("Info", "test info", mock.Anything).Return()
	mockLogger.On("Error", "test error", mock.Anything).Return()
	mockLogger.On("Debug", "test debug", mock.Anything).Return()

	// Test interface methods
	mockLogger.Info("test info", "key", "value")
	mockLogger.Error("test error", "key", "value")
	mockLogger.Debug("test debug", "key", "value")

	mockLogger.AssertExpectations(t)
}
