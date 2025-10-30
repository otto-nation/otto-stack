package stack

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/types"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	pkgTypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"github.com/spf13/cobra"
)

// CleanupHandler handles the cleanup command
type CleanupHandler struct{}

// NewCleanupHandler creates a new cleanup handler
func NewCleanupHandler() *CleanupHandler {
	return &CleanupHandler{}
}

// Handle executes the cleanup command
func (h *CleanupHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	// Get CI-friendly flags
	ciFlags := utils.GetCIFlags(cmd)

	if !ciFlags.Quiet {
		ui.Header(constants.MsgCleaning)
	}

	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		utils.HandleError(ciFlags, err)
		return nil
	}
	defer cleanup()

	// Get flags
	all, _ := cmd.Flags().GetBool(constants.FlagAll)
	volumes, _ := cmd.Flags().GetBool(constants.FlagVolumes)
	images, _ := cmd.Flags().GetBool(constants.FlagImages)
	networks, _ := cmd.Flags().GetBool(constants.FlagNetworks)
	force, _ := cmd.Flags().GetBool(constants.FlagForce)
	dryRun, _ := cmd.Flags().GetBool(constants.FlagDryRun)

	// If --all is specified, enable all cleanup options
	if all {
		volumes = true
		images = true
		networks = true
	}

	// Show what will be cleaned up
	if dryRun {
		constants.SendMessage(constants.MsgCleanupDryRun)
		if volumes {
			constants.SendMessage(constants.MsgCleanupUnusedVolumes)
		}
		if images {
			constants.SendMessage(constants.MsgCleanupUnusedImages)
		}
		if networks {
			constants.SendMessage(constants.MsgCleanupUnusedNetworks)
		}
		constants.SendMessage(constants.MsgCleanupStoppedContainers)
		return nil
	}

	// Confirm cleanup unless forced
	if !force && !ciFlags.NonInteractive {
		constants.SendMessage(constants.MsgCleanupWarning)
		confirmed, err := ui.PromptConfirm(constants.MsgCleanupConfirm.Content, false)
		if err != nil {
			utils.HandleError(ciFlags, fmt.Errorf(constants.MsgFailedGetConfirmation.Content, err))
			return nil
		}
		if !confirmed {
			constants.SendMessage(constants.MsgCleanupCancelled)
			return nil
		}
	}

	// Perform cleanup operations
	if err := h.performCleanup(ctx, setup); err != nil {
		utils.HandleError(ciFlags, fmt.Errorf(constants.MsgCleanupFailed.Content, err))
		return nil
	}

	if !ciFlags.Quiet {
		constants.SendMessage(constants.MsgCleanupSuccess)
	}

	return nil
}

// performCleanup executes the actual cleanup operations
func (h *CleanupHandler) performCleanup(ctx context.Context, setup *CoreSetup) error {
	// Clean up stopped containers
	constants.SendMessage(constants.MsgRemovingContainers)
	if err := setup.DockerClient.Containers().Stop(ctx, setup.Config.Project.Name, []string{}, pkgTypes.StopOptions{
		Remove: true,
	}); err != nil {
		constants.SendMessage(constants.MsgFailedRemoveContainers, err)
	}

	constants.SendMessage(constants.MsgCleanupOperationsCompleted)
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
