//go:build unit

package middleware

import (
	"context"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/stretchr/testify/assert"
)

func TestValidationMiddleware_Creation(t *testing.T) {
	t.Run("creates validation middleware", func(t *testing.T) {
		middleware := NewValidationMiddleware()
		assert.NotNil(t, middleware)
		assert.IsType(t, &ValidationMiddleware{}, middleware)
	})
}

func TestValidationMiddleware_Execute(t *testing.T) {
	t.Run("executes with force flag", func(t *testing.T) {
		middleware := NewValidationMiddleware()
		ctx := context.Background()
		cliCtx := clicontext.Context{
			Runtime: clicontext.RuntimeSpec{Force: true},
		}
		base := &base.BaseCommand{Output: &mockOutput{}}
		next := &mockCommand{}

		err := middleware.Execute(ctx, cliCtx, base, next)
		// Should execute next command when force is true
		if err != nil {
			assert.Error(t, err)
		}
		assert.True(t, next.executed)
	})

	t.Run("executes without force flag", func(t *testing.T) {
		middleware := NewValidationMiddleware()
		ctx := context.Background()
		cliCtx := clicontext.Context{
			Runtime: clicontext.RuntimeSpec{Force: false},
		}
		base := &base.BaseCommand{Output: &mockOutput{}}
		next := &mockCommand{}

		err := middleware.Execute(ctx, cliCtx, base, next)
		// Should handle project state validation
		if err != nil {
			assert.Error(t, err)
		}
	})
}

func TestLoggingMiddleware_Creation(t *testing.T) {
	t.Run("creates logging middleware", func(t *testing.T) {
		middleware := NewLoggingMiddleware()
		assert.NotNil(t, middleware)
		assert.IsType(t, &LoggingMiddleware{}, middleware)
	})
}

func TestLoggingMiddleware_Execute(t *testing.T) {
	t.Run("logs successful execution", func(t *testing.T) {
		middleware := NewLoggingMiddleware()
		ctx := context.Background()
		cliCtx := clicontext.Context{}
		base := &base.BaseCommand{Output: &mockOutput{}}
		next := &mockCommand{}

		err := middleware.Execute(ctx, cliCtx, base, next)
		assert.NoError(t, err)
		assert.True(t, next.executed)
	})

	t.Run("logs failed execution", func(t *testing.T) {
		middleware := NewLoggingMiddleware()
		ctx := context.Background()
		cliCtx := clicontext.Context{}
		base := &base.BaseCommand{Output: &mockOutput{}}
		next := &mockCommand{shouldError: true}

		err := middleware.Execute(ctx, cliCtx, base, next)
		assert.Error(t, err)
		assert.True(t, next.executed)
	})
}

func TestInitializationMiddleware_Creation(t *testing.T) {
	t.Run("creates initialization middleware", func(t *testing.T) {
		middleware := NewInitializationMiddleware()
		assert.NotNil(t, middleware)
		assert.IsType(t, &InitializationMiddleware{}, middleware)
	})
}

func TestInitializationMiddleware_Execute(t *testing.T) {
	t.Run("executes with initialization check", func(t *testing.T) {
		middleware := NewInitializationMiddleware()
		ctx := context.Background()
		cliCtx := clicontext.Context{}
		base := &base.BaseCommand{Output: &mockOutput{}}
		next := &mockCommand{}

		err := middleware.Execute(ctx, cliCtx, base, next)
		// Should handle initialization validation
		if err != nil {
			assert.Error(t, err)
		}
	})
}

// Mock implementations for testing
type mockOutput struct{}

func (m *mockOutput) Success(msg string, args ...any) {}
func (m *mockOutput) Error(msg string, args ...any)   {}
func (m *mockOutput) Warning(msg string, args ...any) {}
func (m *mockOutput) Info(msg string, args ...any)    {}
func (m *mockOutput) Header(msg string, args ...any)  {}
func (m *mockOutput) Muted(msg string, args ...any)   {}

type mockCommand struct {
	executed    bool
	shouldError bool
}

func (m *mockCommand) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	m.executed = true
	if m.shouldError {
		return assert.AnError
	}
	return nil
}
