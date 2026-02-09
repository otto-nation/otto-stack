package operations

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
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
	base.Output.Header(messages.LifecycleLogs)

	execCtx, err := h.detectContext()
	if err != nil {
		return err
	}

	if execCtx.Type == clicontext.Shared {
		return h.handleSharedContext(ctx, cmd, args, execCtx)
	}

	return h.handleProjectContext(ctx, cmd, args, base)
}

func (h *LogsHandler) detectContext() (*clicontext.ExecutionContext, error) {
	detector, err := clicontext.NewDetector()
	if err != nil {
		return nil, err
	}
	return detector.Detect()
}

func (h *LogsHandler) handleProjectContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	setup, cleanup, err := common.SetupCoreCommand(ctx, base)
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

	logReq := services.LogRequest{
		Project:        setup.Config.Project.Name,
		ServiceConfigs: serviceConfigs,
		Follow:         h.getFlag(cmd, docker.FlagFollow),
		Timestamps:     h.getFlag(cmd, docker.FlagTimestamps),
		Tail:           h.getTailFlag(cmd),
	}

	return stackService.Logs(ctx, logReq)
}

func (h *LogsHandler) handleSharedContext(ctx context.Context, cmd *cobra.Command, args []string, execCtx *clicontext.ExecutionContext) error {
	if len(args) == 0 {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationServiceNameRequired, nil)
	}

	if err := h.verifyServicesInRegistry(args, execCtx); err != nil {
		return err
	}

	composeManager, err := docker.NewManager()
	if err != nil {
		return err
	}

	options := docker.LogOptions{
		Services:   args,
		Follow:     h.getFlag(cmd, docker.FlagFollow),
		Timestamps: h.getFlag(cmd, docker.FlagTimestamps),
		Tail:       h.getTailFlag(cmd),
	}

	consumer := &docker.SimpleLogConsumer{}
	return composeManager.Logs(ctx, "shared", consumer, options.ToSDK())
}

func (h *LogsHandler) verifyServicesInRegistry(serviceNames []string, execCtx *clicontext.ExecutionContext) error {
	reg := registry.NewManager(execCtx.SharedContainers.Root)
	registryData, err := reg.Load()
	if err != nil {
		return err
	}

	for _, serviceName := range serviceNames {
		if _, exists := registryData.Containers[serviceName]; !exists {
			return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, fmt.Sprintf(messages.SharedServiceNotInRegistry, serviceName), nil)
		}
	}
	return nil
}

func (h *LogsHandler) getFlag(cmd *cobra.Command, flag string) bool {
	val, _ := cmd.Flags().GetBool(flag)
	return val
}

func (h *LogsHandler) getTailFlag(cmd *cobra.Command) string {
	tail, _ := cmd.Flags().GetString(docker.FlagTail)
	if tail == "" {
		return core.DefaultLogTailLines
	}
	return tail
}

// ValidateArgs validates the command arguments
func (h *LogsHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *LogsHandler) GetRequiredFlags() []string {
	return []string{}
}
