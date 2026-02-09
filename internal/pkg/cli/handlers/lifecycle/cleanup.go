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

	// Reconcile registry with Docker state before checking orphans
	if err := h.reconcileRegistry(ctx, base, &ciFlags); err != nil {
		// Log but don't fail - reconciliation is best-effort
		if !ciFlags.Quiet {
			base.Output.Warning("Failed to reconcile registry: %v", err)
		}
	}

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
func (h *CleanupHandler) checkOrphans(setup *common.CoreSetup, base *base.BaseCommand, ciFlags *ci.Flags) error {
	if ciFlags.Quiet {
		return nil
	}

	reg, err := h.getRegistry()
	if err != nil {
		return err
	}

	orphans, err := reg.FindOrphansWithChecks(context.Background(), setup.DockerClient)
	if err != nil {
		return err
	}

	if len(orphans) == 0 {
		return nil
	}

	h.displayOrphans(base, orphans)
	return nil
}

func (h *CleanupHandler) displayOrphans(base *base.BaseCommand, orphans []registry.OrphanInfo) {
	safe, warning, critical := h.groupBySeverity(orphans)

	base.Output.Warning(messages.OrphanFound, len(orphans))
	h.displayCritical(base, critical)
	h.displayWarning(base, warning)
	h.displaySafe(base, safe)
	base.Output.Info(messages.OrphanRunCleanupHint)
}

func (h *CleanupHandler) displayCritical(base *base.BaseCommand, orphans []registry.OrphanInfo) {
	if len(orphans) == 0 {
		return
	}
	base.Output.Error(messages.OrphanSeverityCritical, len(orphans))
	for _, o := range orphans {
		base.Output.Info("    - %s: %s", o.Service, o.Reason)
	}
}

func (h *CleanupHandler) displayWarning(base *base.BaseCommand, orphans []registry.OrphanInfo) {
	if len(orphans) == 0 {
		return
	}
	base.Output.Warning(messages.OrphanSeverityWarning, len(orphans))
	for _, o := range orphans {
		base.Output.Info("    - %s: %s", o.Service, o.Reason)
		if len(o.ProjectsFound) > 0 {
			base.Output.Info("      "+messages.OrphanRemainingProjects, o.ProjectsFound)
		}
	}
}

func (h *CleanupHandler) displaySafe(base *base.BaseCommand, orphans []registry.OrphanInfo) {
	if len(orphans) == 0 {
		return
	}
	base.Output.Info(messages.OrphanSeveritySafe, len(orphans))
	for _, o := range orphans {
		base.Output.Info("    - %s: %s", o.Service, o.Reason)
	}
}

func (h *CleanupHandler) groupBySeverity(orphans []registry.OrphanInfo) (safe, warning, critical []registry.OrphanInfo) {
	for _, o := range orphans {
		switch o.Severity {
		case registry.OrphanSeveritySafe:
			safe = append(safe, o)
		case registry.OrphanSeverityWarning:
			warning = append(warning, o)
		case registry.OrphanSeverityCritical:
			critical = append(critical, o)
		}
	}
	return
}

// cleanOrphans removes orphaned shared containers
func (h *CleanupHandler) cleanOrphans(setup *common.CoreSetup, base *base.BaseCommand, ciFlags *ci.Flags, force bool) error {
	reg, err := h.getRegistry()
	if err != nil {
		return err
	}

	orphans, err := reg.FindOrphansWithChecks(context.Background(), setup.DockerClient)
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

	return registry.NewManager(execCtx.SharedContainers.Root), nil
}

// reconcileRegistry syncs registry with actual Docker container state
func (h *CleanupHandler) reconcileRegistry(ctx context.Context, base *base.BaseCommand, ciFlags *ci.Flags) error {
	reg, err := h.getRegistry()
	if err != nil {
		return err
	}

	// Get service manager which has Docker client
	svc, err := common.NewServiceManager(false)
	if err != nil {
		return err
	}

	result, err := reg.Reconcile(ctx, svc.DockerClient)
	if err != nil {
		return err
	}

	if len(result.Removed) > 0 && !ciFlags.Quiet {
		base.Output.Info("Reconciled registry: removed %d stale entries", len(result.Removed))
	}

	return nil
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
