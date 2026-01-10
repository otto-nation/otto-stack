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
	// Check initialization first

	// Parse all flags with validation - single line!
	flags, err := core.ParseCleanupFlags(cmd)
	if err != nil {
		return err
	}

	// Get CI-friendly flags
	ciFlags := ci.GetFlags(cmd)

	if !ciFlags.Quiet {
		base.Output.Header(core.MsgCleaning)
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

	// Show what will be cleaned up
	dryRun, _ := cmd.Flags().GetBool(core.FlagDryRun)
	if dryRun {
		base.Output.Info("Dry run mode - showing what would be cleaned")
		if flags.Volumes {
			base.Output.Info("Would clean unused volumes")
		}
		if flags.Images {
			base.Output.Info("Would clean unused images")
		}
		if flags.Networks {
			base.Output.Info("Would clean unused networks")
		}
		base.Output.Info("Would clean stopped containers")
		return nil
	}

	// Confirm cleanup unless forced
	if !flags.Force && !ciFlags.NonInteractive {
		base.Output.Warning("This will remove all containers, networks, and volumes")
		// Note: Need to implement proper confirmation with base.Output
		confirmed := true // For now, assume confirmed
		if !confirmed {
			// Cleanup operation
			return nil
		}
	}

	// Perform cleanup operations
	if err := h.performCleanup(ctx, setup, cmd, base); err != nil {
		return ci.FormatError(ciFlags, fmt.Errorf("cleanup failed: %w", err))
	}

	if !ciFlags.Quiet {
		base.Output.Success("Cleanup completed successfully")
	}

	return nil
}

// performCleanup executes the actual cleanup operations
func (h *CleanupHandler) performCleanup(ctx context.Context, setup *common.CoreSetup, cmd *cobra.Command, base *base.BaseCommand) error {
	flags, err := core.ParseCleanupFlags(cmd)
	if err != nil {
		return pkgerrors.NewValidationError(pkgerrors.FieldFlags, "failed to parse cleanup flags", err)
	}

	ciFlags := ci.GetFlags(cmd)

	// Get project name from flag or config
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

	// Create stack service
	stackService, err := common.NewServiceManager(false)
	if err != nil {
		return pkgerrors.NewServiceError(common.ComponentStack, common.MsgFailedCreateStackService, err)
	}

	// List containers to clean
	var containers []docker.ContainerInfo
	if projectName != "" {
		containers, err = stackService.DockerClient.ListContainers(ctx, projectName)
	} else {
		// For all Otto containers, use project name pattern
		containers, err = stackService.DockerClient.ListContainers(ctx, "")
	}

	if err != nil {
		return pkgerrors.NewDockerError(common.OpListContainers, "", err)
	}

	if len(containers) == 0 {
		if !ciFlags.Quiet {
			base.Output.Info("No containers to clean")
		}
		return nil
	}

	// Stop and remove containers
	for _, container := range containers {
		if !ciFlags.Quiet {
			base.Output.Info("Removing container: %s", container.Name)
		}
		if err := stackService.DockerClient.RemoveContainer(ctx, container.ID, flags.Force); err != nil {
			base.Output.Warning("Failed to remove container %s: %v", container.Name, err)
		}
	}

	// Clean up volumes if requested
	if flags.Volumes {
		if err := stackService.DockerClient.RemoveResources(ctx, docker.ResourceVolume, projectName); err != nil && !ciFlags.Quiet {
			base.Output.Warning("Failed to clean volumes: %v", err)
		}
	}

	// Clean up networks if requested
	if flags.Networks {
		if err := stackService.DockerClient.RemoveResources(ctx, docker.ResourceNetwork, projectName); err != nil && !ciFlags.Quiet {
			base.Output.Warning("Failed to clean networks: %v", err)
		}
	}

	// Clean up images if requested
	if flags.Images {
		if err := stackService.DockerClient.RemoveResources(ctx, docker.ResourceImage, projectName); err != nil && !ciFlags.Quiet {
			base.Output.Warning("Failed to remove images: %v", err)
		}
	}

	return nil
}

// ValidateArgs validates the command arguments
func (h *CleanupHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *CleanupHandler) GetRequiredFlags() []string {
	return []string{}
}
