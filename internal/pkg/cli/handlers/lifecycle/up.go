package lifecycle

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/registry"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/validation"
)

// UpHandler handles the up command
type UpHandler struct{}

// NewUpHandler creates a new up handler
func NewUpHandler() *UpHandler {
	return &UpHandler{}
}

// Handle executes the up command
func (h *UpHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	detector, err := clicontext.NewDetector()
	if err != nil {
		return err
	}

	execCtx, err := detector.Detect()
	if err != nil {
		return err
	}

	if execCtx.Type == clicontext.Global {
		return h.handleGlobalContext(ctx, cmd, args, base, execCtx)
	}

	return h.handleProjectContext(ctx, cmd, args, base, execCtx)
}

func (h *UpHandler) handleProjectContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand, execCtx *clicontext.ExecutionContext) error {
	base.Output.Header("%s", core.MsgLifecycle_starting)

	setup, cleanup, err := common.SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	serviceConfigs, err := common.ResolveServiceConfigs(args, setup)
	if err != nil {
		return err
	}

	service, err := common.NewServiceManager(false)
	if err != nil {
		return pkgerrors.NewServiceError(common.ComponentStack, common.ActionCreateService, err)
	}

	force, _ := cmd.Flags().GetBool(core.FlagForce)
	startRequest := services.StartRequest{
		Project:        setup.Config.Project.Name,
		ServiceConfigs: serviceConfigs,
		Build:          force,
		ForceRecreate:  false,
	}

	if err = service.Start(ctx, startRequest); err != nil {
		return pkgerrors.NewServiceError(common.ComponentStack, common.ActionStartServices, err)
	}

	base.Output.Success("Services started successfully")
	base.Output.Info("Project: %s", setup.Config.Project.Name)
	for _, svc := range serviceConfigs {
		base.Output.Info("  %s %s", display.StatusSuccess, svc.Name)
	}

	return nil
}

func (h *UpHandler) handleGlobalContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand, execCtx *clicontext.ExecutionContext) error {
	base.Output.Header("Starting shared containers")

	if len(args) == 0 {
		return pkgerrors.NewValidationError("", "No services specified. Use 'otto-stack up <service>' to start shared containers", nil)
	}

	// Resolve service configs
	serviceUtils := services.NewServiceUtils()
	var serviceConfigs []types.ServiceConfig
	for _, serviceName := range args {
		cfg, err := serviceUtils.LoadServiceConfig(serviceName)
		if err != nil {
			return pkgerrors.NewServiceError("service", "load", err)
		}
		if cfg != nil {
			serviceConfigs = append(serviceConfigs, *cfg)
		}
	}

	// Register shared containers
	reg := registry.NewManager(execCtx.Shared.Root)
	regData, err := reg.Load()
	if err != nil {
		return err
	}

	for _, svc := range serviceConfigs {
		containerName := "otto-shared-" + svc.Name
		if err := reg.Register(svc.Name, containerName, "global"); err != nil {
			base.Output.Warning("Failed to register %s: %v", svc.Name, err)
		}
	}

	if err := reg.Save(regData); err != nil {
		base.Output.Warning("Failed to save registry: %v", err)
	}

	base.Output.Success("Shared containers registered")
	for _, svc := range serviceConfigs {
		base.Output.Info("  %s %s (shared)", display.StatusSuccess, svc.Name)
	}

	return nil
}

// ValidateArgs validates the command arguments
func (h *UpHandler) ValidateArgs(args []string) error {
	return validation.ValidateUpArgs(args)
}

// GetRequiredFlags returns required flags for this command
func (h *UpHandler) GetRequiredFlags() []string {
	return []string{}
}
