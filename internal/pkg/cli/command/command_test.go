//go:build unit

package command

import (
	"context"
	"errors"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/test/testhelpers"
)

// MockCommand for testing
type MockCommand struct {
	executed bool
	err      error
}

func (m *MockCommand) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	m.executed = true
	return m.err
}

// MockMiddleware for testing
type MockMiddleware struct {
	executed bool
	err      error
}

func (m *MockMiddleware) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand, next Command) error {
	m.executed = true
	if m.err != nil {
		return m.err
	}
	return next.Execute(ctx, cliCtx, base)
}

func TestHandler_Execute(t *testing.T) {
	tests := []struct {
		name          string
		command       *MockCommand
		middlewares   []Middleware
		expectedError bool
		commandCalled bool
	}{
		{
			name:          "executes command without middleware",
			command:       &MockCommand{},
			middlewares:   nil,
			expectedError: false,
			commandCalled: true,
		},
		{
			name:          "executes command with single middleware",
			command:       &MockCommand{},
			middlewares:   []Middleware{&MockMiddleware{}},
			expectedError: false,
			commandCalled: true,
		},
		{
			name:          "handles command error",
			command:       &MockCommand{err: errors.New("command error")},
			middlewares:   nil,
			expectedError: true,
			commandCalled: true,
		},
		{
			name:          "handles middleware error",
			command:       &MockCommand{},
			middlewares:   []Middleware{&MockMiddleware{err: errors.New("middleware error")}},
			expectedError: true,
			commandCalled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.command, tt.middlewares...)

			ctx := context.Background()
			cliCtx := clicontext.Context{}
			base := &base.BaseCommand{}

			err := handler.Execute(ctx, cliCtx, base)

			if tt.expectedError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.commandCalled != tt.command.executed {
				t.Errorf("expected command executed=%v, got=%v", tt.commandCalled, tt.command.executed)
			}
		})
	}
}

func TestHandler_MiddlewareChain(t *testing.T) {
	command := &MockCommand{}
	middleware1 := &MockMiddleware{}
	middleware2 := &MockMiddleware{}

	handler := NewHandler(command, middleware1, middleware2)

	ctx := context.Background()
	cliCtx := clicontext.Context{}
	base := &base.BaseCommand{}

	err := handler.Execute(ctx, cliCtx, base)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !middleware1.executed {
		t.Error("middleware1 was not executed")
	}
	if !middleware2.executed {
		t.Error("middleware2 was not executed")
	}
	if !command.executed {
		t.Error("command was not executed")
	}
}

func TestCommand_adapters(t *testing.T) {
	t.Run("new cobra adapter", func(t *testing.T) {
		adapter := NewCobraAdapter(nil)
		testhelpers.AssertNoError(t, nil, "NewCobraAdapter should not error")
		if adapter == nil {
			t.Error("NewCobraAdapter should return adapter")
		}
	})
}
