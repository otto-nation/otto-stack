package project

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/validation"
)

// ValidateHandler handles the validate command
type ValidateHandler struct{}

// NewValidateHandler creates a new validate handler
func NewValidateHandler() *ValidateHandler {
	return &ValidateHandler{}
}

// Handle executes the validate command
func (h *ValidateHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	flags, err := core.ParseValidateFlags(cmd)
	if err != nil {
		return err
	}
	quiet := ci.GetFlags(cmd).Quiet

	if err := validation.CheckInitialization(); err != nil {
		return err
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentProject, messages.ErrorsConfigLoadFailed, err)
	}
	if !quiet {
		base.Output.Success(messages.ValidateCheckConfigSyntax)
	}

	vm := NewValidationManager()
	if err := vm.ValidateProjectName(cfg.Project.Name); err != nil {
		return err
	}
	if !quiet {
		base.Output.Success(messages.ValidateCheckProjectName)
	}

	if _, err := services.ResolveUpServices(cfg.Stack.Enabled, cfg); err != nil {
		return err
	}
	if !quiet {
		base.Output.Success(messages.ValidateCheckServices)
	}

	if flags.Strict {
		if !isCommandAvailable(docker.DockerCmd) {
			return pkgerrors.NewSystemError(pkgerrors.ErrCodeInvalid, messages.ValidateStrictDockerUnavailable, nil)
		}
		if !quiet {
			base.Output.Success(messages.ValidateCheckDocker)
		}

		if _, err := os.Stat(".git"); os.IsNotExist(err) {
			base.Output.Warning("%s", messages.WarningsNotGitRepository)
		}
	}

	base.Output.Success(messages.SuccessConfigurationValid)
	return nil
}

// ValidateArgs validates the command arguments
func (h *ValidateHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *ValidateHandler) GetRequiredFlags() []string {
	return []string{}
}
