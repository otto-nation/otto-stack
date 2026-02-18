package middleware

import (
	"context"
	"errors"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/stretchr/testify/assert"
)

type mockCommand struct {
	executeFunc func(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error
}

func (m *mockCommand) Execute(ctx context.Context, cliCtx clicontext.Context, baseCmd *base.BaseCommand) error {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, cliCtx, baseCmd)
	}
	return nil
}

func TestNewValidationMiddleware(t *testing.T) {
	m := NewValidationMiddleware()
	assert.NotNil(t, m)
}

func TestValidationMiddleware_Execute(t *testing.T) {
	m := NewValidationMiddleware()
	ctx := context.Background()
	cliCtx := clicontext.Context{}
	baseCmd := &base.BaseCommand{}

	next := &mockCommand{
		executeFunc: func(ctx context.Context, cliCtx clicontext.Context, baseCmd *base.BaseCommand) error {
			return nil
		},
	}

	err := m.Execute(ctx, cliCtx, baseCmd, next)
	assert.NoError(t, err)
}

func TestNewLoggingMiddleware(t *testing.T) {
	m := NewLoggingMiddleware()
	assert.NotNil(t, m)
}

func TestLoggingMiddleware_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m := NewLoggingMiddleware()
		ctx := context.Background()
		cliCtx := clicontext.Context{}
		baseCmd := &base.BaseCommand{}

		next := &mockCommand{
			executeFunc: func(ctx context.Context, cliCtx clicontext.Context, baseCmd *base.BaseCommand) error {
				return nil
			},
		}

		err := m.Execute(ctx, cliCtx, baseCmd, next)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		m := NewLoggingMiddleware()
		ctx := context.Background()
		cliCtx := clicontext.Context{}
		baseCmd := &base.BaseCommand{}

		expectedErr := errors.New("test error")
		next := &mockCommand{
			executeFunc: func(ctx context.Context, cliCtx clicontext.Context, baseCmd *base.BaseCommand) error {
				return expectedErr
			},
		}

		err := m.Execute(ctx, cliCtx, baseCmd, next)
		assert.Equal(t, expectedErr, err)
	})
}
