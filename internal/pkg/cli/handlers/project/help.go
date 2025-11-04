package project

import (
	"context"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/spf13/cobra"
)

// HelpHandler handles help command
type HelpHandler struct{}

// NewHelpHandler creates a new help handler
func NewHelpHandler() *HelpHandler {
	return &HelpHandler{}
}

// ValidateArgs validates help command arguments
func (h *HelpHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for help command
func (h *HelpHandler) GetRequiredFlags() []string {
	return []string{}
}

// Handle handles the help command
func (h *HelpHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	if len(args) > 0 {
		// Show help for specific command
		return cmd.Root().Help()
	}

	// Show general help
	base.Output.Info("%s", constants.Messages[constants.MsgHelp_app_description])
	base.Output.Info("")
	base.Output.Info("%s", constants.Messages[constants.MsgHelp_usage])
	base.Output.Info("%s", constants.Messages[constants.MsgHelp_usage_example])
	base.Output.Info("")
	base.Output.Info("%s", constants.Messages[constants.MsgHelp_available_commands])
	base.Output.Info("%s", constants.Messages[constants.MsgHelp_command_up])
	base.Output.Info("%s", constants.Messages[constants.MsgHelp_command_down])
	base.Output.Info("%s", constants.Messages[constants.MsgHelp_command_status])
	base.Output.Info("%s", constants.Messages[constants.MsgHelp_command_init])
	base.Output.Info("%s", constants.Messages[constants.MsgHelp_command_services])
	base.Output.Info("")
	base.Output.Info("%s", constants.Messages[constants.MsgHelp_more_info])

	return nil
}
