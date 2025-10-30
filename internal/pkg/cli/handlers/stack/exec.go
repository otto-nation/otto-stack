package stack

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/types"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	pkgTypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/spf13/cobra"
)

// ExecHandler handles the exec command
type ExecHandler struct{}

// NewExecHandler creates a new exec handler
func NewExecHandler() *ExecHandler {
	return &ExecHandler{}
}

// Handle executes the exec command
func (h *ExecHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	// Get CI-friendly flags
	ciFlags := utils.GetCIFlags(cmd)

	if len(args) < 2 {
		utils.HandleError(ciFlags, fmt.Errorf("%s", constants.MsgRequiresServiceAndCommand.Content))
		return nil
	}

	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		utils.HandleError(ciFlags, err)
		return nil
	}
	defer cleanup()

	serviceName := args[0]
	command := args[1:]

	// Get flags
	interactive, _ := cmd.Flags().GetBool(constants.FlagInteractive)
	tty, _ := cmd.Flags().GetBool(constants.FlagTTY)
	user, _ := cmd.Flags().GetString(constants.FlagUser)
	workdir, _ := cmd.Flags().GetString(constants.FlagWorkdir)

	// Create exec options
	options := pkgTypes.ExecOptions{
		Interactive: interactive,
		TTY:         tty,
		User:        user,
		WorkingDir:  workdir,
	}

	// Execute command using Docker client
	return setup.DockerClient.Containers().Exec(ctx, setup.Config.Project.Name, serviceName, command, options)
}

// ValidateArgs validates the command arguments
func (h *ExecHandler) ValidateArgs(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("%s", constants.MsgRequiresServiceAndCommand.Content)
	}
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *ExecHandler) GetRequiredFlags() []string {
	return []string{}
}
