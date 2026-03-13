package middleware

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
)

type execContextKeyType struct{}

var execContextKey = execContextKeyType{}

// WithExecContext runs DetectExecutionContext once and injects the result into the
// request context. Handlers use ExecContextOrDetect instead of calling
// common.DetectExecutionContext directly, eliminating the repeated 4-line pattern.
func WithExecContext() Middleware {
	return func(next base.CommandHandler) base.CommandHandler {
		return HandlerFunc(func(ctx context.Context, cmd *cobra.Command, args []string, b *base.BaseCommand) error {
			execCtx, err := common.DetectExecutionContext()
			if err != nil {
				return err
			}
			return next.Handle(context.WithValue(ctx, execContextKey, execCtx), cmd, args, b)
		})
	}
}

// ExecContextFromCtx retrieves the execution context injected by WithExecContext.
// Returns (nil, false) if not present.
func ExecContextFromCtx(ctx context.Context) (clicontext.ExecutionMode, bool) {
	v, ok := ctx.Value(execContextKey).(clicontext.ExecutionMode)
	return v, ok
}

// ExecContextOrDetect returns the injected execution context if available, otherwise
// calls common.DetectExecutionContext. Handlers call this instead of DetectExecutionContext
// directly so they work both with and without the middleware chain.
func ExecContextOrDetect(ctx context.Context) (clicontext.ExecutionMode, error) {
	if v, ok := ExecContextFromCtx(ctx); ok {
		return v, nil
	}
	return common.DetectExecutionContext()
}
