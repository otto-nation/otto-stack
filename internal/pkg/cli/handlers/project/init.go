package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/spf13/cobra"
)

// InitHandler handles the init command
type InitHandler struct {
	selectedServices        []string
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
	defer func() {
		if r := recover(); r != nil {
			logger.Error(logger.LogMsgOperationFailed, logger.LogFieldOperation, logger.OperationInit, logger.LogFieldError, fmt.Errorf("panic: %v", r))
			panic(r)
		}
	}()

	if err := h.validateInitFlags(cmd); err != nil {
		return err
	}

	ciFlags := ci.GetFlags(cmd)
	initFlags, _ := core.ParseInitFlags(cmd)

	// Default project name to current directory if not provided
	if initFlags.ProjectName == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return pkgerrors.NewValidationError(core.FlagProjectName, "failed to get current directory", err)
		}
		initFlags.ProjectName = filepath.Base(cwd)
	}

	// Set force flag in handler for validation functions
	h.forceOverwrite = initFlags.Force

	// Check if project is already initialized
	if _, err := os.Stat(core.OttoStackDir); err == nil && !initFlags.Force {
		return fmt.Errorf("project already initialized. Use --%s to overwrite", core.FlagForce)
	}

	var projectName string
	var services []string
	var validation map[string]bool
	var advanced map[string]bool
	var err error

	processor := NewModeProcessor(ciFlags.NonInteractive, h)
	projectName, services, validation, advanced, err = processor.Process(initFlags, base)
	if err != nil {
		return err
	}

	if err := h.projectManager.CreateProjectStructure(projectName, services, validation, advanced, base); err != nil {
		return err
	}

	h.displaySuccessMessage(projectName, base)
	return nil
}

func (h *InitHandler) validateInitFlags(cmd *cobra.Command) error {
	_, err := core.ParseInitFlags(cmd)
	if err != nil {
		logger.Error(logger.LogMsgOperationFailed, logger.LogFieldOperation, logger.OperationInit, logger.LogFieldError, err)
		return err
	}

	ciFlags := ci.GetFlags(cmd)
	if ciFlags.NonInteractive {
		// In non-interactive mode, we need services flag
		initFlags, _ := core.ParseInitFlags(cmd)
		if initFlags.Services == "" {
			return pkgerrors.NewValidationError(core.FlagServices, "services are required in non-interactive mode", nil)
		}
	}
	return nil
}

func (h *InitHandler) displaySuccessMessage(_ string, base *base.BaseCommand) {
	base.Output.Success("%s", core.MsgSuccess_init)
	base.Output.Info("%s", core.MsgInit_next_steps)
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

// runValidations executes selected validation functions
