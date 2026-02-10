package operations

import (
	"context"
	"fmt"

	"github.com/docker/compose/v5/pkg/api"
	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/registry"
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
	execCtx, err := h.detectContext()
	if err != nil {
		return err
	}

	switch mode := execCtx.(type) {
	case *clicontext.ProjectMode:
		return h.handleProjectContext(ctx, cmd, args, base)
	case *clicontext.SharedMode:
		return h.handleSharedContext(ctx, cmd, args, mode)
	default:
		return fmt.Errorf("unknown execution mode: %T", execCtx)
	}
}

func (h *ExecHandler) detectContext() (clicontext.ExecutionMode, error) {
	detector, err := clicontext.NewDetector()
	if err != nil {
		return nil, err
	}
	return detector.DetectContext()
}

func (h *ExecHandler) handleProjectContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
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

func (h *ExecHandler) handleSharedContext(ctx context.Context, cmd *cobra.Command, args []string, mode *clicontext.SharedMode) error {
	serviceName := args[0]
	command := args[1:]

	if err := h.verifyServiceInRegistry(serviceName, mode); err != nil {
		return err
	}

	user, _ := cmd.Flags().GetString(docker.FlagUser)
	workdir, _ := cmd.Flags().GetString(docker.FlagWorkdir)
	interactive, _ := cmd.Flags().GetBool(docker.FlagInteractive)
	tty, _ := cmd.Flags().GetBool(docker.FlagTTY)

	composeManager, err := docker.NewManager()
	if err != nil {
		return err
	}

	options := api.RunOptions{
		Service:     serviceName,
		Command:     command,
		User:        user,
		WorkingDir:  workdir,
		Interactive: interactive,
		Tty:         tty,
		Index:       1,
	}

	_, err = composeManager.Exec(ctx, "shared", options)
	return err
}

func (h *ExecHandler) verifyServiceInRegistry(serviceName string, mode *clicontext.SharedMode) error {
	reg := registry.NewManager(mode.Shared.Root)
	registryData, err := reg.Load()
	if err != nil {
		return err
	}

	if _, exists := registryData.Containers[serviceName]; !exists {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, fmt.Sprintf(messages.SharedServiceNotInRegistry, serviceName), nil)
	}
	return nil
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
