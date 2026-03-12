package lifecycle

import (
	"context"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/middleware"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
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
	if ciFlags.DryRun {
		base.Output.Info("%s", messages.DryRunShowingWhatWouldHappen)
		base.Output.Info(messages.DryRunWouldClean, flags.Project)
		return nil
	}

	if !ciFlags.Quiet {
		base.Output.Header(messages.LifecycleCleaning)
	}

	// Reconcile registry — works from any directory, does not need project config.
	if err := h.reconcileRegistry(ctx, base, &ciFlags); err != nil {
		// Log but don't fail — reconciliation is best-effort
		if !ciFlags.Quiet {
			base.Output.Warning(messages.WarningsRegistryReconcileFailed, err)
		}
	}

	if flags.Orphans {
		// Orphan cleanup works from any directory — no project config needed.
		dockerClient, err := docker.NewClient(nil)
		if err != nil {
			return ci.FormatError(ciFlags, pkgerrors.NewDockerError(pkgerrors.ErrCodeOperationFail, messages.ErrorsDockerClientCreateFailed, err))
		}
		defer func() { _ = dockerClient.Close() }()

		if !ciFlags.Quiet {
			if err := h.checkOrphans(ctx, dockerClient, base, &ciFlags); err != nil {
				logger.Error("orphan check failed", "error", err)
			}
		}
		return h.handleOrphanCleanup(ctx, dockerClient, base, &ciFlags, flags.Force)
	}

	if flags.All || flags.Volumes || flags.Images || flags.Networks {
		// Resource cleanup (volumes, images, networks) requires a project context.
		setup, cleanup, err := middleware.CoreSetupOrCreate(ctx, base)
		if err != nil {
			return ci.FormatError(ciFlags, err)
		}
		defer cleanup()
		return h.handleResourceCleanup(ctx, setup, flags, &ciFlags, base)
	}

	if !ciFlags.Quiet {
		base.Output.Info(messages.InfoNoCleanupOptions)
	}
	return nil
}

func (h *CleanupHandler) handleOrphanCleanup(ctx context.Context, dockerClient *docker.Client, base *base.BaseCommand, ciFlags *ci.Flags, force bool) error {
	if err := h.cleanOrphans(ctx, dockerClient, base, ciFlags, force); err != nil {
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
		Project:        projectName,
		Force:          flags.Force,
		RemoveVolumes:  flags.Volumes,
		RemoveImages:   flags.Images,
		RemoveNetworks: flags.Networks,
	})
	if err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, projectName, messages.ErrorsServiceCleanupFailed, err)
	}

	if !ciFlags.Quiet {
		base.Output.Success(messages.SuccessCleanupCompleted)
	}

	return nil
}

// confirmCleanup asks for user confirmation before cleaning
func (h *CleanupHandler) confirmCleanup(base *base.BaseCommand) bool {
	base.Output.Warning(messages.WarningsCleanupWarning)

	prompt := &survey.Confirm{
		Message: messages.PromptsCleanupConfirm,
		Default: false,
	}

	var confirmed bool
	if err := survey.AskOne(prompt, &confirmed); err != nil {
		return false
	}

	return confirmed
}

// checkOrphans checks for and reports orphaned shared containers
func (h *CleanupHandler) checkOrphans(ctx context.Context, dockerClient *docker.Client, base *base.BaseCommand, ciFlags *ci.Flags) error {
	if ciFlags.Quiet {
		return nil
	}

	reg, err := h.getRegistry()
	if err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsRegistryGetFailed, err)
	}

	orphans, err := reg.FindOrphansWithChecks(ctx, dockerClient)
	if err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsRegistryFindOrphansFailed, err)
	}

	if len(orphans) == 0 {
		return nil
	}

	h.displayOrphans(base, orphans)
	return nil
}

func (h *CleanupHandler) displayOrphans(base *base.BaseCommand, orphans []registry.OrphanInfo) {
	display := registry.NewOrphanDisplay(base.Output)
	display.Display(orphans)
}

// cleanOrphans stops and removes orphaned shared containers, then clears them from the registry.
func (h *CleanupHandler) cleanOrphans(ctx context.Context, dockerClient *docker.Client, base *base.BaseCommand, ciFlags *ci.Flags, force bool) error {
	reg, err := h.getRegistry()
	if err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsRegistryGetFailed, err)
	}

	orphans, err := reg.FindOrphansWithChecks(ctx, dockerClient)
	if err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsRegistryFindOrphansFailed, err)
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

	// Stop and remove each orphaned container before clearing the registry entry.
	for _, orphan := range orphans {
		if orphan.ContainerState != registry.ContainerStateNotFound {
			// Force-remove handles both stop and removal in one call.
			if err := dockerClient.RemoveContainer(ctx, orphan.Container, true); err != nil && !ciFlags.Quiet {
				base.Output.Warning(messages.WarningsOrphanRemoveContainerFailed, orphan.Container, err)
			}
		}
	}

	cleaned, err := reg.CleanOrphans()
	if err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsRegistryCleanOrphansFailed, err)
	}

	if !ciFlags.Quiet {
		base.Output.Success(messages.OrphanRemovedFromRegistry, len(cleaned))
		for _, service := range cleaned {
			base.Output.Info(messages.InfoListItem, service)
		}
	}

	return nil
}

// getRegistry creates a registry manager pointed at the global shared root.
// This works from any directory — no project config required.
func (h *CleanupHandler) getRegistry() (*registry.Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, pkgerrors.NewSystemError(pkgerrors.ErrCodeInternal, messages.ErrorsContextDetectFailed, err)
	}
	sharedRoot := filepath.Join(homeDir, core.OttoStackDir, core.SharedDir)
	return registry.NewManager(sharedRoot), nil
}

// reconcileRegistry syncs registry with actual Docker container state
func (h *CleanupHandler) reconcileRegistry(ctx context.Context, base *base.BaseCommand, ciFlags *ci.Flags) error {
	reg, err := h.getRegistry()
	if err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsRegistryGetFailed, err)
	}

	dockerClient, err := docker.NewClient(nil)
	if err != nil {
		return pkgerrors.NewDockerError(pkgerrors.ErrCodeOperationFail, messages.ErrorsDockerClientCreateFailed, err)
	}
	defer func() { _ = dockerClient.Close() }()

	result, err := reg.Reconcile(ctx, dockerClient)
	if err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsRegistryReconcileFailed, err)
	}

	if len(result.Removed) > 0 && !ciFlags.Quiet {
		base.Output.Info(messages.InfoReconciledRegistry, len(result.Removed))
	}

	if warnings := reg.ValidateAgainstDocker(ctx, dockerClient); len(warnings) > 0 && !ciFlags.Quiet {
		for _, w := range warnings {
			base.Output.Warning(w)
		}
	}

	return nil
}

// confirmOrphanCleanup prompts user to confirm orphan cleanup
func (h *CleanupHandler) confirmOrphanCleanup(base *base.BaseCommand, orphans []registry.OrphanInfo) bool {
	base.Output.Warning(messages.OrphanWillRemove)
	for _, orphan := range orphans {
		base.Output.Info(messages.InfoListItem, orphan.Service)
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
