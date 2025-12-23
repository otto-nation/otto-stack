package stack

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
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

	// Create stack service
	stackService, err := NewStackService(false)
	if err != nil {
		return fmt.Errorf("failed to create stack service: %w", err)
	}

	// Create exec request
	execRequest := services.ExecRequest{
		Project:    setup.Config.Project.Name,
		Service:    serviceName,
		Command:    command,
		User:       flags.User,
		WorkingDir: flags.Workdir,
	}

	return stackService.Exec(ctx, execRequest)
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
