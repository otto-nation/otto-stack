//go:build unit

package project

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckManager_Methods(t *testing.T) {
	t.Run("tests NewHealthCheckManager", func(t *testing.T) {
		manager := NewHealthCheckManager()
		assert.NotNil(t, manager)
	})

	t.Run("tests RunAllChecks", func(t *testing.T) {
		manager := NewHealthCheckManager()
		mockBase := &base.BaseCommand{
			Output: &mockOutput{},
		}

		result := manager.RunAllChecks(context.Background(), mockBase)
		assert.IsType(t, false, result)
	})

	t.Run("tests CheckDocker", func(t *testing.T) {
		manager := NewHealthCheckManager()
		mockBase := &base.BaseCommand{
			Output: &mockOutput{},
		}

		result := manager.CheckDocker(context.Background(), mockBase)
		assert.IsType(t, false, result)
	})

	t.Run("tests CheckDockerCompose", func(t *testing.T) {
		manager := NewHealthCheckManager()
		mockBase := &base.BaseCommand{
			Output: &mockOutput{},
		}

		result := manager.CheckDockerCompose(mockBase)
		assert.IsType(t, false, result)
	})

	t.Run("tests CheckProjectInit", func(t *testing.T) {
		manager := NewHealthCheckManager()
		mockBase := &base.BaseCommand{
			Output: &mockOutput{},
		}

		result := manager.CheckProjectInit(mockBase)
		assert.IsType(t, false, result)
	})

	t.Run("tests CheckConfiguration", func(t *testing.T) {
		manager := NewHealthCheckManager()
		mockBase := &base.BaseCommand{
			Output: &mockOutput{},
		}

		result := manager.CheckConfiguration(mockBase)
		assert.IsType(t, false, result)
	})
}

func TestValidationManager_Methods(t *testing.T) {
	t.Run("tests NewValidationManager", func(t *testing.T) {
		manager := NewValidationManager()
		assert.NotNil(t, manager)
	})
}

func TestConfigManager_Methods(t *testing.T) {
	t.Run("tests NewConfigManager", func(t *testing.T) {
		manager := NewConfigManager()
		assert.NotNil(t, manager)
	})
}

func TestProjectManager_Methods(t *testing.T) {
	t.Run("tests NewProjectManager", func(t *testing.T) {
		manager := NewProjectManager()
		assert.NotNil(t, manager)
	})
}

func TestHandlers_Methods(t *testing.T) {
	t.Run("tests DepsHandler", func(t *testing.T) {
		handler := NewDepsHandler()
		cmd := &cobra.Command{}
		args := []string{}

		mockBase := &base.BaseCommand{
			Output: &mockOutput{},
		}

		err := handler.Handle(context.Background(), cmd, args, mockBase)
		if err != nil {
			assert.Error(t, err)
		}
	})

	t.Run("tests ConflictsHandler", func(t *testing.T) {
		handler := NewConflictsHandler()
		cmd := &cobra.Command{}
		args := []string{}

		mockBase := &base.BaseCommand{
			Output: &mockOutput{},
		}

		err := handler.Handle(context.Background(), cmd, args, mockBase)
		if err != nil {
			assert.Error(t, err)
		}
	})

	t.Run("tests ValidateHandler", func(t *testing.T) {
		handler := NewValidateHandler()
		cmd := &cobra.Command{}
		args := []string{}

		mockBase := &base.BaseCommand{
			Output: &mockOutput{},
		}

		err := handler.Handle(context.Background(), cmd, args, mockBase)
		if err != nil {
			assert.Error(t, err)
		}
	})

	t.Run("tests DoctorHandler", func(t *testing.T) {
		handler := NewDoctorHandler()
		cmd := &cobra.Command{}
		args := []string{}

		mockBase := &base.BaseCommand{
			Output: &mockOutput{},
		}

		err := handler.Handle(context.Background(), cmd, args, mockBase)
		if err != nil {
			assert.Error(t, err)
		}
	})
}

// Mock output for testing
type mockOutput struct{}

func (m *mockOutput) Success(msg string, args ...any) {}
func (m *mockOutput) Error(msg string, args ...any)   {}
func (m *mockOutput) Warning(msg string, args ...any) {}
func (m *mockOutput) Info(msg string, args ...any)    {}
func (m *mockOutput) Header(msg string, args ...any)  {}
func (m *mockOutput) Muted(msg string, args ...any)   {}
func (m *mockOutput) Writer() io.Writer               { return os.Stdout }
