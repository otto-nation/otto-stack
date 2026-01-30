package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/command"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/middleware"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
)

// InitHandler handles the init command
type InitHandler struct {
	forceOverwrite          bool
	promptManager           *PromptManager
	validationManager       *ValidationManager
	projectManager          *ProjectManager
	serviceSelectionManager *ServiceSelectionManager
}

// NewInitHandler creates a new InitHandler
func NewInitHandler() *InitHandler {
	handler := &InitHandler{
		validationManager: NewValidationManager(),
		projectManager:    NewProjectManager(),
	}
	handler.promptManager = NewPromptManager(handler)
	handler.serviceSelectionManager = NewServiceSelectionManager(handler.promptManager, handler.validationManager)
	return handler
}

// ValidateProjectName implements ProjectValidator interface
func (h *InitHandler) ValidateProjectName(name string) error {
	return h.validateProjectName(name)
}

// Handle executes the init command
func (h *InitHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	logger.Info(logger.LogMsgStartingOperation, logger.LogFieldOperation, logger.OperationInit)
	defer h.handlePanic()

	if err := h.validateInitFlags(cmd); err != nil {
		return err
	}

	ciFlags := ci.GetFlags(cmd)
	initFlags, _ := core.ParseInitFlags(cmd)

	h.logDebugInfo(base, cmd, initFlags)

	if err := h.setDefaultProjectName(initFlags, base, cmd); err != nil {
		return err
	}

	h.forceOverwrite = initFlags.Force

	projectCtx, err := h.processMode(ciFlags, initFlags, base, cmd)
	if err != nil {
		return err
	}

	if err := h.executeInit(ctx, projectCtx, base); err != nil {
		return err
	}

	h.displaySuccessMessage(projectCtx.Project.Name, base)
	return nil
}

func (h *InitHandler) handlePanic() {
	if r := recover(); r != nil {
		logger.Error(logger.LogMsgOperationFailed, logger.LogFieldOperation, logger.OperationInit, logger.LogFieldError, fmt.Errorf("panic: %v", r))
		panic(r)
	}
}

func (h *InitHandler) logDebugInfo(base *base.BaseCommand, cmd *cobra.Command, initFlags *core.InitFlags) {
	if base.GetVerbose(cmd) {
		logger.Debug("Initializing project", "projectName", initFlags.ProjectName, "force", initFlags.Force)
	}
}

func (h *InitHandler) setDefaultProjectName(initFlags *core.InitFlags, base *base.BaseCommand, cmd *cobra.Command) error {
	if initFlags.ProjectName != "" {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, "failed to get current directory", err)
	}

	initFlags.ProjectName = filepath.Base(cwd)
	if base.GetVerbose(cmd) {
		logger.Debug("Using current directory as project name", "projectName", initFlags.ProjectName)
	}
	return nil
}

func (h *InitHandler) processMode(ciFlags ci.Flags, initFlags *core.InitFlags, base *base.BaseCommand, cmd *cobra.Command) (clicontext.Context, error) {
	processor := NewModeProcessor(ciFlags.NonInteractive, h)
	projectCtx, err := processor.Process(initFlags, base)
	if err != nil {
		return clicontext.Context{}, err
	}

	if base.GetVerbose(cmd) {
		logger.Debug("Project context created", "services", len(projectCtx.Services.Names))
	}

	return projectCtx, nil
}

func (h *InitHandler) executeInit(ctx context.Context, projectCtx clicontext.Context, base *base.BaseCommand) error {
	initCommand := NewInitCommand(h.projectManager)
	validationMiddleware := middleware.NewValidationMiddleware()
	loggingMiddleware := middleware.NewLoggingMiddleware()

	handler := command.NewHandler(initCommand, loggingMiddleware, validationMiddleware)
	return handler.Execute(ctx, projectCtx, base)
}

func (h *InitHandler) validateInitFlags(cmd *cobra.Command) error {
	initFlags, err := core.ParseInitFlags(cmd)
	if err != nil {
		logger.Error(logger.LogMsgOperationFailed, logger.LogFieldOperation, logger.OperationInit, logger.LogFieldError, err)
		return err
	}

	return h.validateNonInteractiveMode(cmd, initFlags)
}

func (h *InitHandler) validateNonInteractiveMode(cmd *cobra.Command, initFlags *core.InitFlags) error {
	ciFlags := ci.GetFlags(cmd)
	if ciFlags.NonInteractive && initFlags.Services == "" {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, "services are required in non-interactive mode", nil)
	}
	return nil
}

func (h *InitHandler) displaySuccessMessage(_ string, base *base.BaseCommand) {
	base.Output.Success("%s", core.MsgSuccess_init)
	base.Output.Info("%s", core.MsgInit_next_steps)
	h.displayNextSteps(base)
}

func (h *InitHandler) displayNextSteps(base *base.BaseCommand) {
	base.Output.Info(core.MsgInit_step_review_config, core.OttoStackDir, core.ConfigFileName)
	base.Output.Info(core.MsgInit_step_start_stack, core.AppName)
	base.Output.Info(core.MsgInit_step_check_status, core.AppName)
}

// ValidateArgs validates the command arguments
func (h *InitHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *InitHandler) GetRequiredFlags() []string {
	return []string{}
}

// validateProjectName validates a project name
func (h *InitHandler) validateProjectName(name string) error {
	return h.validationManager.ValidateProjectName(name)
}

// buildSharingConfig creates sharing configuration from init flags
func (h *InitHandler) buildSharingConfig(initFlags *core.InitFlags) *clicontext.SharingSpec {
	// If --no-shared-containers is set, disable all sharing
	if initFlags.NoSharedContainers {
		return &clicontext.SharingSpec{
			Enabled: false,
		}
	}

	sharingSpec := &clicontext.SharingSpec{
		Enabled:  true,
		Services: make(map[string]bool),
	}

	// If --shared-services is specified, parse it
	if initFlags.SharedServices != "" {
		//nolint:modernize // SplitSeq requires Go 1.24+
		for _, svc := range strings.Split(initFlags.SharedServices, ",") {
			svc = strings.TrimSpace(svc)
			if svc != "" {
				sharingSpec.Services[svc] = true
			}
		}
	}

	return sharingSpec
}
