package middleware

import (
	"context"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/command"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
)

// ValidationMiddleware validates project state before command execution
type ValidationMiddleware struct{}

// NewValidationMiddleware creates a new validation middleware
func NewValidationMiddleware() *ValidationMiddleware {
	return &ValidationMiddleware{}
}

// Execute validates the project state
func (m *ValidationMiddleware) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand, next command.Command) error {
	// Validation moved to init handler for immediate feedback
	return next.Execute(ctx, cliCtx, base)
}

// LoggingMiddleware logs command execution
type LoggingMiddleware struct{}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{}
}

// Execute logs the command execution
func (m *LoggingMiddleware) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand, next command.Command) error {
	logger.Info("starting_operation", "operation", "init")

	err := next.Execute(ctx, cliCtx, base)

	if err != nil {
		logger.Debug("operation_failed", "operation", "init", "error", err)
	} else {
		logger.Info("operation_completed", "operation", "init")
	}

	return err
}
