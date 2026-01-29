package operations

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

// ExecHandler handles the exec command
type ExecHandler struct{}

// NewExecHandler creates a new exec handler
func NewExecHandler() *ExecHandler {
	return &ExecHandler{}
}

// Handle executes the exec command
func (h *ExecHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	setup, cleanup, err := common.SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	serviceName := args[0]
	command := args[1:]

	user, _ := cmd.Flags().GetString(docker.FlagUser)
	workdir, _ := cmd.Flags().GetString(docker.FlagWorkdir)
	interactive, _ := cmd.Flags().GetBool(docker.FlagInteractive)
	tty, _ := cmd.Flags().GetBool(docker.FlagTTY)

	stackService, err := common.NewServiceManager(false)
	if err != nil {
		return err
	}

	req := services.ExecRequest{
		Project:     setup.Config.Project.Name,
		Service:     serviceName,
		Command:     command,
		User:        user,
		WorkingDir:  workdir,
		Interactive: interactive,
		TTY:         tty,
	}

	return stackService.Exec(ctx, req)
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
