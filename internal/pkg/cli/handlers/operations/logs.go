package operations

import (
	"context"

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

	// Detect execution context
	detector, err := clicontext.NewDetector()
	if err != nil {
		return err
	}

	execCtx, err := detector.Detect()
	if err != nil {
		return err
	}

	// Handle shared context
	if execCtx.Type == clicontext.Shared {
		return h.handleSharedContext(ctx, cmd, args, execCtx)
	}

	// Handle project context
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

	// Parse log flags
	follow, _ := cmd.Flags().GetBool(docker.FlagFollow)
	timestamps, _ := cmd.Flags().GetBool(docker.FlagTimestamps)
	tail, _ := cmd.Flags().GetString(docker.FlagTail)
	if tail == "" {
		tail = core.DefaultLogTailLines
	}

	// Use Service.Logs method
	logReq := services.LogRequest{
		Project:        setup.Config.Project.Name,
		ServiceConfigs: serviceConfigs,
		Follow:         follow,
		Timestamps:     timestamps,
		Tail:           tail,
	}

	return stackService.Logs(ctx, logReq)
}

func (h *LogsHandler) handleSharedContext(ctx context.Context, cmd *cobra.Command, args []string, execCtx *clicontext.ExecutionContext) error {
	if len(args) == 0 {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationServiceNameRequired, nil)
	}

	reg := registry.NewManager(execCtx.SharedContainers.Root)
	registryData, err := reg.Load()
	if err != nil {
		return err
	}

	// Verify services exist in registry
	for _, serviceName := range args {
		if _, exists := registryData.Containers[serviceName]; !exists {
			return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, "service '"+serviceName+"' not found in shared registry", nil)
		}
	}

	// Parse log flags
	follow, _ := cmd.Flags().GetBool(docker.FlagFollow)
	timestamps, _ := cmd.Flags().GetBool(docker.FlagTimestamps)
	tail, _ := cmd.Flags().GetString(docker.FlagTail)
	if tail == "" {
		tail = core.DefaultLogTailLines
	}

	// Use compose manager for logs
	composeManager, err := docker.NewManager()
	if err != nil {
		return err
	}

	options := docker.LogOptions{
		Services:   args,
		Follow:     follow,
		Timestamps: timestamps,
		Tail:       tail,
	}

	consumer := &docker.SimpleLogConsumer{}
	return composeManager.Logs(ctx, "shared", consumer, options.ToSDK())
}

// ValidateArgs validates the command arguments
func (h *LogsHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *LogsHandler) GetRequiredFlags() []string {
	return []string{}
}
