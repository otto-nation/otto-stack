// Package middleware provides a composable middleware chain for command handlers.
// Apply Chain() in handlers.Get() so every handler gets logging, execution-context
// detection, and (for project-mode commands) config + Docker client setup without
// repeating the boilerplate inside each Handle() implementation.
package middleware

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
)

// Middleware wraps a CommandHandler with cross-cutting behavior.
type Middleware func(base.CommandHandler) base.CommandHandler

// HandlerFunc is a function that satisfies base.CommandHandler.
type HandlerFunc func(context.Context, *cobra.Command, []string, *base.BaseCommand) error

func (f HandlerFunc) Handle(ctx context.Context, cmd *cobra.Command, args []string, b *base.BaseCommand) error {
	return f(ctx, cmd, args, b)
}

// Chain applies middlewares to a handler. The first middleware listed is the outermost
// wrapper (runs first on the way in, last on the way out).
func Chain(h base.CommandHandler, middlewares ...Middleware) base.CommandHandler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// Logging wraps a handler with structured start/completion logging using the actual
// cobra command name. This replaces per-handler logger.Info("starting operation") calls
// that previously hardcoded the operation name.
func Logging() Middleware {
	return func(next base.CommandHandler) base.CommandHandler {
		return HandlerFunc(func(ctx context.Context, cmd *cobra.Command, args []string, b *base.BaseCommand) error {
			start := time.Now()
			logger.Info(logger.LogMsgStartingOperation, logger.LogFieldOperation, cmd.Name())
			err := next.Handle(ctx, cmd, args, b)
			ms := time.Since(start).Milliseconds()
			if err != nil {
				logger.Error(logger.LogMsgOperationFailed, logger.LogFieldOperation, cmd.Name(), logger.LogFieldError, err, "duration_ms", ms)
			} else {
				logger.Debug("command completed", logger.LogFieldOperation, cmd.Name(), "duration_ms", ms)
			}
			return err
		})
	}
}
