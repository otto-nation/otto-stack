package lifecycle

import (
	"context"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/registry"
	"github.com/otto-nation/otto-stack/internal/pkg/services"

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
		base.Output.Header(core.MsgLifecycle_cleaning)
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
		base.Output.Info("No cleanup options specified. Use --help for available options")
	}
	return nil
}

func (h *CleanupHandler) handleOrphanCleanup(setup *common.CoreSetup, base *base.BaseCommand, ciFlags *ci.Flags, force bool) error {
	if err := h.cleanOrphans(setup, base, ciFlags, force); err != nil {
		return ci.FormatError(*ciFlags, err)
	}
	if !ciFlags.Quiet {
		base.Output.Success(core.MsgOrphan_cleanup_success)
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
		return pkgerrors.NewServiceError(common.ComponentStack, common.MsgFailedCreateStackService, err)
	}

	containers, err := stackService.DockerClient.ListContainers(ctx, projectName)
	if err != nil {
		return pkgerrors.NewDockerError(common.OpListContainers, "", err)
	}

	if len(containers) == 0 {
		if !ciFlags.Quiet {
			base.Output.Info("No containers to clean")
		}
		return nil
	}

	h.removeContainers(ctx, stackService, containers, flags.Force, ciFlags, base)
	h.cleanupResources(ctx, stackService, flags, projectName, ciFlags, base)
	return nil
}

// confirmCleanup asks for user confirmation before cleaning
func (h *CleanupHandler) confirmCleanup(base *base.BaseCommand) bool {
	base.Output.Warning("This will remove all containers, networks, and volumes")
	// TODO: Implement proper confirmation with base.Output
	return true
}

// removeContainers removes all containers in the list
func (h *CleanupHandler) removeContainers(ctx context.Context, stackService *services.Service, containers []docker.ContainerInfo, force bool, ciFlags *ci.Flags, base *base.BaseCommand) {
	for _, container := range containers {
		if !ciFlags.Quiet {
			base.Output.Info("Removing container: %s", container.Name)
		}
		if err := stackService.DockerClient.RemoveContainer(ctx, container.ID, force); err != nil {
			base.Output.Warning("Failed to remove container %s: %v", container.Name, err)
		}
	}
}

// cleanupResources cleans up volumes, networks, and images if requested
func (h *CleanupHandler) cleanupResources(ctx context.Context, stackService *services.Service, flags *core.CleanupFlags, projectName string, ciFlags *ci.Flags, base *base.BaseCommand) {
	if flags.Volumes {
		if err := stackService.DockerClient.RemoveResources(ctx, docker.ResourceVolume, projectName); err != nil && !ciFlags.Quiet {
			base.Output.Warning("Failed to clean volumes: %v", err)
		}
	}

	if flags.Networks {
		if err := stackService.DockerClient.RemoveResources(ctx, docker.ResourceNetwork, projectName); err != nil && !ciFlags.Quiet {
			base.Output.Warning("Failed to clean networks: %v", err)
		}
	}

	if flags.Images {
		if err := stackService.DockerClient.RemoveResources(ctx, docker.ResourceImage, projectName); err != nil && !ciFlags.Quiet {
			base.Output.Warning("Failed to remove images: %v", err)
		}
	}
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

	base.Output.Warning(core.MsgOrphan_found, len(orphans))
	for _, orphan := range orphans {
		base.Output.Info("  - %s (%s)", orphan.Service, orphan.Reason)
	}
	base.Output.Info(core.MsgOrphan_run_cleanup_hint)
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
			base.Output.Info(core.MsgOrphan_none_found)
		}
		return nil
	}

	if !force && !ciFlags.NonInteractive && !h.confirmOrphanCleanup(base, orphans) {
		base.Output.Info(core.MsgOrphan_cleanup_cancelled)
		return nil
	}

	cleaned, err := reg.CleanOrphans()
	if err != nil {
		return err
	}

	if !ciFlags.Quiet {
		base.Output.Success(core.MsgOrphan_removed_from_registry, len(cleaned))
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
	base.Output.Warning(core.MsgOrphan_will_remove)
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
