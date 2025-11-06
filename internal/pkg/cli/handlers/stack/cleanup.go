package stack

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
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
	if err := utils.CheckInitialization(); err != nil {
		return err
	}

	// Parse all flags with validation - single line!
	flags, err := constants.ParseCleanupFlags(cmd)
	if err != nil {
		return err
	}

	// Get CI-friendly flags
	ciFlags := utils.GetCIFlags(cmd)

	if !ciFlags.Quiet {
		base.Output.Header(constants.MsgCleaning)
	}

	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		utils.HandleError(ciFlags, err)
		return nil
	}
	defer cleanup()

	// If --all is specified, enable all cleanup options
	if flags.All {
		flags.Volumes = true
		flags.Images = true
		flags.Networks = true
	}

	// Show what will be cleaned up
	if flags.DryRun {
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
		utils.HandleError(ciFlags, fmt.Errorf("cleanup failed: %w", err))
		return nil
	}

	if !ciFlags.Quiet {
		base.Output.Success("Cleanup completed successfully")
	}

	return nil
}

// performCleanup executes the actual cleanup operations
func (h *CleanupHandler) performCleanup(ctx context.Context, setup *CoreSetup, cmd *cobra.Command, base *base.BaseCommand) error {
	// Parse cleanup flags
	flags, err := constants.ParseCleanupFlags(cmd)
	if err != nil {
		return fmt.Errorf("failed to parse cleanup flags: %w", err)
	}

	ciFlags := utils.GetCIFlags(cmd)

	if !ciFlags.Quiet {
		base.Output.Info("Starting cleanup operations...")
	}

	// Stop and remove containers
	if err := setup.DockerClient.ComposeDown(ctx, setup.Config.Project.Name, docker.StopOptions{
		Remove:        true,
		RemoveVolumes: flags.Volumes,
	}); err != nil {
		return fmt.Errorf("failed to stop containers: %w", err)
	}

	// Clean up additional resources if requested
	if flags.Images {
		if err := setup.DockerClient.RemoveResources(ctx, docker.ResourceImage, setup.Config.Project.Name); err != nil && !ciFlags.Quiet {
			base.Output.Warning("Failed to remove images: %v", err)
		}
	}

	if flags.Networks {
		if err := setup.DockerClient.RemoveResources(ctx, docker.ResourceNetwork, setup.Config.Project.Name); err != nil && !ciFlags.Quiet {
			base.Output.Warning("Failed to remove networks: %v", err)
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
