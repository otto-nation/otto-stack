package lifecycle

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
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

	// If --all is specified, enable all cleanup options
	if flags.All {
		flags.Volumes = true
		flags.Images = true
		flags.Networks = true
	}

	// Confirm cleanup unless forced or in non-interactive mode
	if !flags.Force && !ciFlags.NonInteractive {
		if !h.confirmCleanup(base) {
			return nil
		}
	}

	// Perform cleanup operations
	if err := h.performCleanup(ctx, setup, flags, &ciFlags, base); err != nil {
		return ci.FormatError(ciFlags, fmt.Errorf("cleanup failed: %w", err))
	}

	if !ciFlags.Quiet {
		base.Output.Success("Cleanup completed successfully")
	}

	return nil
}

// performCleanup executes the actual cleanup operations
func (h *CleanupHandler) performCleanup(ctx context.Context, setup *common.CoreSetup, flags *core.CleanupFlags, ciFlags *ci.Flags, base *base.BaseCommand) error {
	projectName := flags.Project
	if projectName == "" {
		projectName = setup.Config.Project.Name
	}

	if !ciFlags.Quiet {
		if projectName != "" {
			base.Output.Info("Cleaning up project: %s", projectName)
		} else {
			base.Output.Info("Cleaning up all Otto Stack containers")
		}
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

	// Remove containers
	h.removeContainers(ctx, stackService, containers, flags.Force, ciFlags, base)

	// Clean up additional resources
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

// ValidateArgs validates the command arguments
func (h *CleanupHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *CleanupHandler) GetRequiredFlags() []string {
	return []string{}
}
