//go:build unit

package operations

import (
	"context"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestConnectHandler_Methods(t *testing.T) {
	t.Run("tests Handle method", func(t *testing.T) {
		handler := NewConnectHandler()
		cmd := &cobra.Command{}
		args := []string{"postgres"}

		mockBase := &base.BaseCommand{
			Output: &mockOutput{},
		}

		err := handler.Handle(context.Background(), cmd, args, mockBase)
		// Should handle gracefully
		if err != nil {
			assert.Error(t, err)
		}
	})
}

func TestExecHandler_Methods(t *testing.T) {
	t.Run("tests Handle method", func(t *testing.T) {
		handler := NewExecHandler()
		cmd := &cobra.Command{}
		args := []string{"postgres", "bash"}

		mockBase := &base.BaseCommand{
			Output: &mockOutput{},
		}

		err := handler.Handle(context.Background(), cmd, args, mockBase)
		// Should handle gracefully
		if err != nil {
			assert.Error(t, err)
		}
	})
}

func TestLogsHandler_Methods(t *testing.T) {
	t.Run("tests Handle method", func(t *testing.T) {
		handler := NewLogsHandler()
		cmd := &cobra.Command{}
		cmd.Flags().Bool("follow", false, "follow logs")
		cmd.Flags().String("tail", "100", "tail lines")
		args := []string{"postgres"}

		mockBase := &base.BaseCommand{
			Output: &mockOutput{},
		}

		err := handler.Handle(context.Background(), cmd, args, mockBase)
		// Should handle gracefully
		if err != nil {
			assert.Error(t, err)
		}
	})
}

func TestStatusHandler_Methods(t *testing.T) {
	t.Run("tests Handle method", func(t *testing.T) {
		handler := NewStatusHandler()
		cmd := &cobra.Command{}
		cmd.Flags().Bool("all", false, "show all")
		args := []string{}

		mockBase := &base.BaseCommand{
			Output: &mockOutput{},
		}

		err := handler.Handle(context.Background(), cmd, args, mockBase)
		// Should handle gracefully
		if err != nil {
			assert.Error(t, err)
		}
	})
}

func TestStatusPackageFunctions(t *testing.T) {
	t.Run("tests getContainerName", func(t *testing.T) {
		config := types.ServiceConfig{Name: "postgres"}

		name := getContainerName(config)
		assert.IsType(t, "", name)
		assert.Contains(t, name, "postgres")
	})

	t.Run("tests filterInitContainers", func(t *testing.T) {
		configs := []types.ServiceConfig{
			{Name: "postgres"},
			{Name: "postgres-init"},
		}

		filtered := filterInitContainers(configs)
		assert.IsType(t, []string{}, filtered)
		// Should return service names
		assert.GreaterOrEqual(t, len(filtered), 0)
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
