package stack

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
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

	if len(args) < constants.MinArgumentCount {
		utils.HandleError(ciFlags, fmt.Errorf("%s", constants.Messages[constants.MsgErrors_requires_service_and_command]))
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

	// Parse all flags with validation - single line!
	flags, err := constants.ParseExecFlags(cmd)
	if err != nil {
		return err
	}

	// Create exec options - clean usage with no repetitive error handling
	options := types.ExecOptions{
		Interactive: flags.Interactive,
		TTY:         flags.TTY,
		User:        flags.User,
		WorkingDir:  flags.Workdir,
	}

	// Execute command using Docker client
	return setup.DockerClient.Containers().Exec(ctx, setup.Config.Project.Name, serviceName, command, options)
}

// ValidateArgs validates the command arguments
func (h *ExecHandler) ValidateArgs(args []string) error {
	if len(args) < constants.MinArgumentCount {
		return fmt.Errorf("%s", constants.Messages[constants.MsgErrors_requires_service_and_command])
	}
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *ExecHandler) GetRequiredFlags() []string {
	return []string{}
}
