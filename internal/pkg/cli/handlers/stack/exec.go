package stack

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
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
	if err := h.ValidateArgs(args); err != nil {
		return err
	}

	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	serviceName := args[0]
	command := args[1:]

	flags, err := core.ParseExecFlags(cmd)
	if err != nil {
		return err
	}

	dockerArgs := h.buildDockerArgs(setup.Config.Project.Name, serviceName, command, flags)
	return setup.DockerClient.RunCommand(ctx, dockerArgs...)
}

// buildDockerArgs constructs the docker compose exec command arguments
func (h *ExecHandler) buildDockerArgs(projectName, serviceName string, command []string, flags *core.ExecFlags) []string {
	args := []string{"compose", "-f", docker.DockerComposeFilePath, "-p", projectName, "exec"}

	if flags.User != "" {
		args = append(args, "--user", flags.User)
	}
	if flags.Workdir != "" {
		args = append(args, "--workdir", flags.Workdir)
	}

	args = append(args, serviceName)
	return append(args, command...)
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
