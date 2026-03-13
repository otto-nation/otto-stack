package middleware

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
)

type coreSetupKeyType struct{}

var coreSetupKey = coreSetupKeyType{}

// WithProjectSetup runs SetupCoreCommand when the detected execution context is
// ProjectMode, injects the result into ctx, and defers Docker client cleanup.
// For SharedMode commands this is a no-op so global-only commands are unaffected.
//
// The --global and --project flags bypass normal context detection; this middleware
// skips setup in both cases so the handler controls setup directly (the handler
// either runs in SharedMode or chdirs via buildProjectMode before calling
// SetupCoreCommand itself).
func WithProjectSetup() Middleware {
	return func(next base.CommandHandler) base.CommandHandler {
		return HandlerFunc(func(ctx context.Context, cmd *cobra.Command, args []string, b *base.BaseCommand) error {
			// --global forces SharedMode — skip project setup.
			if globalFlag, _ := cmd.Flags().GetBool(docker.FlagGlobal); globalFlag {
				return next.Handle(ctx, cmd, args, b)
			}
			// --project changes cwd inside Handle(); skip setup here so the handler
			// picks up the correct project config after the directory change.
			if projectDir, _ := cmd.Flags().GetString(docker.FlagProject); projectDir != "" {
				return next.Handle(ctx, cmd, args, b)
			}

			execCtx, ok := ExecContextFromCtx(ctx)
			if !ok {
				// WithExecContext did not run — fall through to handler's own setup.
				return next.Handle(ctx, cmd, args, b)
			}
			if _, isProject := execCtx.(*clicontext.ProjectMode); !isProject {
				// SharedMode — no project config to load.
				return next.Handle(ctx, cmd, args, b)
			}

			setup, cleanup, err := common.SetupCoreCommand(ctx, b)
			if err != nil {
				return err
			}
			defer cleanup()
			return next.Handle(context.WithValue(ctx, coreSetupKey, setup), cmd, args, b)
		})
	}
}

// CoreSetupOrCreate returns the CoreSetup injected by WithProjectSetup when available.
// If not present (e.g. the --project flag path or SharedMode), it falls back to
// calling SetupCoreCommand directly and returns its cleanup function.
// Callers must always defer the returned cleanup.
func CoreSetupOrCreate(ctx context.Context, b *base.BaseCommand) (*common.CoreSetup, func(), error) {
	if setup, ok := ctx.Value(coreSetupKey).(*common.CoreSetup); ok {
		return setup, func() {}, nil // cleanup already deferred by WithProjectSetup
	}
	return common.SetupCoreCommand(ctx, b)
}
