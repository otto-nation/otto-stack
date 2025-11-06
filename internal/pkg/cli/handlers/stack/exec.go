package stack

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/validation"
	"github.com/spf13/cobra"
)

// ExecHandler handles the exec command
type ExecHandler struct{}

// NewExecHandler creates a new exec handler
func NewExecHandler() *ExecHandler {
	return &ExecHandler{}
}

// Handle executes the exec command
func (h *ExecHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	if len(args) < core.MinArgumentCount {
		return fmt.Errorf("%s", core.MsgErrors_requires_service_and_command)
	}

	// Check initialization first
	if err := validation.CheckInitialization(); err != nil {
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
	flags, err := core.ParseExecFlags(cmd)
	if err != nil {
		return err
	}

	// Create exec options - clean usage with no repetitive error handling
	options := docker.ExecOptions{
		Interactive: flags.Interactive,
		TTY:         flags.TTY,
		User:        flags.User,
		WorkingDir:  flags.Workdir,
	}

	// Execute command using Docker client
	// Execute command in service container
	dockerArgs := []string{"compose", "-f", docker.DockerComposeFile, "-p", setup.Config.Project.Name, "exec"}
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
	if len(args) < core.MinArgumentCount {
		return fmt.Errorf("%s", core.MsgErrors_requires_service_and_command)
	}
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *ExecHandler) GetRequiredFlags() []string {
	return []string{}
}
