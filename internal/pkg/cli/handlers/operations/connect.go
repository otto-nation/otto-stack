package operations

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/registry"
)

// ConnectHandler handles the connect command
type ConnectHandler struct{}

// NewConnectHandler creates a new connect handler
func NewConnectHandler() *ConnectHandler {
	return &ConnectHandler{}
}

// ValidateArgs validates the command arguments
func (h *ConnectHandler) ValidateArgs(args []string) error {
	if len(args) < 1 {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationServiceNameRequired, nil)
	}
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *ConnectHandler) GetRequiredFlags() []string {
	return []string{}
}

// Handle executes the connect command
func (h *ConnectHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	execCtx, err := h.detectContext()
	if err != nil {
		return err
	}

	switch mode := execCtx.(type) {
	case *clicontext.ProjectMode:
		return h.handleProjectContext(ctx, args, base)
	case *clicontext.SharedMode:
		return h.handleSharedContext(args, base, mode)
	default:
		return fmt.Errorf("unknown execution mode: %T", execCtx)
	}
}

func (h *ConnectHandler) detectContext() (clicontext.ExecutionMode, error) {
	detector, err := clicontext.NewDetector()
	if err != nil {
		return nil, err
	}
	return detector.DetectContext()
}

func (h *ConnectHandler) handleProjectContext(ctx context.Context, args []string, base *base.BaseCommand) error {
	setup, cleanup, err := common.SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	base.Output.Header(messages.InfoConnectingToService)

	serviceConfigs, err := common.ResolveServiceConfigs(args, setup)
	if err != nil {
		return err
	}

	if len(serviceConfigs) > 0 {
		base.Output.Info(messages.InfoServiceInfo, serviceConfigs[0].Name)
	}
	base.Output.Success(messages.SuccessConnected)
	base.Output.Info(messages.InfoProjectInfo, setup.Config.Project.Name)

	return nil
}

func (h *ConnectHandler) handleSharedContext(args []string, base *base.BaseCommand, mode *clicontext.SharedMode) error {
	serviceName := args[0]

	if err := h.verifyServiceInRegistry(serviceName, mode); err != nil {
		return err
	}

	base.Output.Header(messages.InfoConnectingToService)
	base.Output.Info(messages.InfoServiceInfo, serviceName)
	base.Output.Success(messages.SuccessConnected)
	base.Output.Info(messages.InfoContextInfo, "shared")

	return nil
}

func (h *ConnectHandler) verifyServiceInRegistry(serviceName string, mode *clicontext.SharedMode) error {
	reg := registry.NewManager(mode.Shared.Root)
	registryData, err := reg.Load()
	if err != nil {
		return err
	}

	if _, exists := registryData.Containers[serviceName]; !exists {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, fmt.Sprintf(messages.SharedServiceNotInRegistry, serviceName), nil)
	}
	return nil
}
