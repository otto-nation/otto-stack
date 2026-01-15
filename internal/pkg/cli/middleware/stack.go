package middleware

import (
	"context"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/command"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/validation"
)

// InitializationMiddleware checks if project is initialized
type InitializationMiddleware struct{}

// NewInitializationMiddleware creates a new initialization middleware
func NewInitializationMiddleware() *InitializationMiddleware {
	return &InitializationMiddleware{}
}

// Execute checks if the project is initialized
func (m *InitializationMiddleware) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand, next command.Command) error {
	if err := validation.CheckInitialization(); err != nil {
		return err
	}
	return next.Execute(ctx, cliCtx, base)
}
