package operations

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/middleware"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/registry"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

// LogsHandler handles the logs command
type LogsHandler struct{}

// NewLogsHandler creates a new logs handler
func NewLogsHandler() *LogsHandler {
	return &LogsHandler{}
}

// Handle executes the logs command
func (h *LogsHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	execCtx, err := middleware.ExecContextOrDetect(ctx)
	if err != nil {
		return err
	}

	switch mode := execCtx.(type) {
	case *clicontext.ProjectMode:
		return h.handleProjectContext(ctx, cmd, args, base)
	case *clicontext.SharedMode:
		return h.handleSharedContext(ctx, cmd, args, base, mode)
	default:
		return pkgerrors.NewSystemErrorf(pkgerrors.ErrCodeInternal, messages.ErrorsContextUnknownMode, execCtx)
	}
}

func (h *LogsHandler) handleProjectContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	setup, cleanup, err := middleware.CoreSetupOrCreate(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	serviceConfigs, err := common.ResolveServiceConfigs(args, setup)
	if err != nil {
		return err
	}

	stackService, err := common.NewServiceManager(false)
	if err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, messages.ErrorsStackCreateFailed, err)
	}

	flags, err := core.ParseLogsFlags(cmd)
	if err != nil {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldFlags, messages.ValidationFailedParseFlags, err)
	}

	tail := flags.Tail
	if tail == "" {
		tail = core.DefaultLogTailLines
	}

	logReq := services.LogRequest{
		Project:        setup.Config.Project.Name,
		ServiceConfigs: serviceConfigs,
		Follow:         flags.Follow,
		Timestamps:     flags.Timestamps,
		Tail:           tail,
		Since:          flags.Since,
		NoColor:        base.Output.GetNoColor(),
		Writer:         base.Output.Writer(),
	}

	return stackService.Logs(ctx, logReq)
}

func (h *LogsHandler) handleSharedContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand, mode *clicontext.SharedMode) error {
	if len(args) == 0 {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationServiceNameRequired, nil)
	}

	reg, err := registry.NewManager(mode.Shared.Root).Load()
	if err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsRegistryLoadFailed, err)
	}
	if err := common.VerifyServicesInRegistry(args, reg); err != nil {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ErrorsServiceNotInRegistry, err)
	}

	composeManager, err := docker.NewManager()
	if err != nil {
		return pkgerrors.NewDockerError(pkgerrors.ErrCodeOperationFail, messages.ErrorsDockerManagerCreateFailed, err)
	}

	flags, err := core.ParseLogsFlags(cmd)
	if err != nil {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldFlags, messages.ValidationFailedParseFlags, err)
	}

	tail := flags.Tail
	if tail == "" {
		tail = core.DefaultLogTailLines
	}

	options := docker.LogOptions{
		Services:   args,
		Follow:     flags.Follow,
		Timestamps: flags.Timestamps,
		Tail:       tail,
		Since:      flags.Since,
	}

	consumer := docker.NewServiceLogConsumer(base.Output.Writer(), base.Output.GetNoColor(), len(args))
	return composeManager.Logs(ctx, core.SharedDir, consumer, options.ToSDK())
}

// ValidateArgs validates the command arguments
func (h *LogsHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *LogsHandler) GetRequiredFlags() []string {
	return []string{}
}
