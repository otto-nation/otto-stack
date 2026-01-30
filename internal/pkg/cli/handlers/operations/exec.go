package operations

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
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
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, "args", messages.ErrorsRequiresServiceAndCommand, nil)
	}
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *ExecHandler) GetRequiredFlags() []string {
	return []string{}
}
