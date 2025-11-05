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
	if len(args) < constants.MinArgumentCount {
		return fmt.Errorf("%s", constants.MsgErrors_requires_service_and_command)
	}

	// Check initialization first
	if err := utils.CheckInitialization(); err != nil {
		return err
	}

	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return err
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
	// Execute command in service container
	dockerArgs := []string{"compose", "-f", constants.DockerComposeFile, "-p", setup.Config.Project.Name, "exec"}
	if options.User != "" {
		dockerArgs = append(dockerArgs, "--user", options.User)
	}
	if options.WorkingDir != "" {
		dockerArgs = append(dockerArgs, "--workdir", options.WorkingDir)
	}
	dockerArgs = append(dockerArgs, serviceName)
	dockerArgs = append(dockerArgs, command...)

	return setup.DockerClient.RunCommand(ctx, dockerArgs...)
}

// ValidateArgs validates the command arguments
func (h *ExecHandler) ValidateArgs(args []string) error {
	if len(args) < constants.MinArgumentCount {
		return fmt.Errorf("%s", constants.MsgErrors_requires_service_and_command)
	}
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *ExecHandler) GetRequiredFlags() []string {
	return []string{}
}
