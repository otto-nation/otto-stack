package lifecycle

import (
	"context"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/registry"
	"github.com/otto-nation/otto-stack/internal/pkg/services"

	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/spf13/cobra"
)

// CleanupHandler handles the cleanup command
type CleanupHandler struct{}

// NewCleanupHandler creates a new cleanup handler
func NewCleanupHandler() *CleanupHandler {
	return &CleanupHandler{}
}

// Handle executes the cleanup command
func (h *CleanupHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	flags, err := core.ParseCleanupFlags(cmd)
	if err != nil {
		return err
	}

	ciFlags := ci.GetFlags(cmd)
	if !ciFlags.Quiet {
		base.Output.Header(messages.LifecycleCleaning)
	}

	setup, cleanup, err := common.SetupCoreCommand(ctx, base)
	if err != nil {
		return ci.FormatError(ciFlags, err)
	}
	defer cleanup()

	// Check for orphans (informational only)
	_ = h.checkOrphans(setup, base, &ciFlags)

	// Execute cleanup based on flags
	if flags.Orphans {
		return h.handleOrphanCleanup(setup, base, &ciFlags, flags.Force)
	}

	if flags.All || flags.Volumes || flags.Images || flags.Networks {
		return h.handleResourceCleanup(ctx, setup, flags, &ciFlags, base)
	}

	if !ciFlags.Quiet {
		base.Output.Info(messages.InfoNoCleanupOptions)
	}
	return nil
}

func (h *CleanupHandler) handleOrphanCleanup(setup *common.CoreSetup, base *base.BaseCommand, ciFlags *ci.Flags, force bool) error {
	if err := h.cleanOrphans(setup, base, ciFlags, force); err != nil {
		return ci.FormatError(*ciFlags, err)
	}
	if !ciFlags.Quiet {
		base.Output.Success(messages.OrphanCleanupSuccess)
	}
	return nil
}

func (h *CleanupHandler) handleResourceCleanup(ctx context.Context, setup *common.CoreSetup, flags *core.CleanupFlags, ciFlags *ci.Flags, base *base.BaseCommand) error {
	if flags.All {
		flags.Volumes, flags.Images, flags.Networks = true, true, true
	}

	if !flags.Force && !ciFlags.NonInteractive && !h.confirmCleanup(base) {
		return nil
	}

	return h.performCleanup(ctx, setup, flags, ciFlags, base)
}

// performCleanup executes the actual cleanup operations
func (h *CleanupHandler) performCleanup(ctx context.Context, setup *common.CoreSetup, flags *core.CleanupFlags, ciFlags *ci.Flags, base *base.BaseCommand) error {
	projectName := flags.Project
	if projectName == "" {
		projectName = setup.Config.Project.Name
	}

	stackService, err := common.NewServiceManager(false)
	if err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, messages.ErrorsStackCreateFailed, err)
	}

	if !ciFlags.Quiet {
		base.Output.Info(messages.InfoCleaningProject, projectName)
	}

	err = stackService.Cleanup(ctx, services.CleanupRequest{
		Project:       projectName,
		Force:         flags.Force,
		RemoveVolumes: flags.Volumes,
		RemoveImages:  flags.Images,
	})
	if err != nil {
		return err
	}

	if !ciFlags.Quiet {
		base.Output.Success(messages.SuccessCleanupCompleted)
	}

	return nil
}

// confirmCleanup asks for user confirmation before cleaning
func (h *CleanupHandler) confirmCleanup(base *base.BaseCommand) bool {
	base.Output.Warning(messages.WarningsCleanupWarning)
	// TODO: Implement proper confirmation with base.Output
	return true
}

// checkOrphans checks for and reports orphaned shared containers
func (h *CleanupHandler) checkOrphans(_ *common.CoreSetup, base *base.BaseCommand, ciFlags *ci.Flags) error {
	reg, err := h.getRegistry()
	if err != nil {
		return err
	}

	orphans, err := reg.FindOrphans()
	if err != nil || len(orphans) == 0 || ciFlags.Quiet {
		return err
	}

	base.Output.Warning(messages.OrphanFound, len(orphans))
	for _, orphan := range orphans {
		base.Output.Info("  - %s (%s)", orphan.Service, orphan.Reason)
	}
	base.Output.Info(messages.OrphanRunCleanupHint)
	return nil
}

// cleanOrphans removes orphaned shared containers
func (h *CleanupHandler) cleanOrphans(_ *common.CoreSetup, base *base.BaseCommand, ciFlags *ci.Flags, force bool) error {
	reg, err := h.getRegistry()
	if err != nil {
		return err
	}

	orphans, err := reg.FindOrphans()
	if err != nil {
		return err
	}

	if len(orphans) == 0 {
		if !ciFlags.Quiet {
			base.Output.Info(messages.OrphanNoneFound)
		}
		return nil
	}

	if !force && !ciFlags.NonInteractive && !h.confirmOrphanCleanup(base, orphans) {
		base.Output.Info(messages.OrphanCleanupCancelled)
		return nil
	}

	cleaned, err := reg.CleanOrphans()
	if err != nil {
		return err
	}

	if !ciFlags.Quiet {
		base.Output.Success(messages.OrphanRemovedFromRegistry, len(cleaned))
		for _, service := range cleaned {
			base.Output.Info("  - %s", service)
		}
	}

	return nil
}

// getRegistry creates a registry manager
func (h *CleanupHandler) getRegistry() (*registry.Manager, error) {
	detector, err := clicontext.NewDetector()
	if err != nil {
		return nil, err
	}

	execCtx, err := detector.Detect()
	if err != nil {
		return nil, err
	}

	return registry.NewManager(execCtx.Shared.Root), nil
}

// confirmOrphanCleanup prompts user to confirm orphan cleanup
func (h *CleanupHandler) confirmOrphanCleanup(base *base.BaseCommand, orphans []registry.OrphanInfo) bool {
	base.Output.Warning(messages.OrphanWillRemove)
	for _, orphan := range orphans {
		base.Output.Info("  - %s", orphan.Service)
	}
	return h.confirmCleanup(base)
}

// ValidateArgs validates the command arguments
func (h *CleanupHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *CleanupHandler) GetRequiredFlags() []string {
	return []string{}
}
