package project

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/spf13/cobra"
)

// InitHandler handles the init command
type InitHandler struct {
	selectedServices  []string
	promptManager     *PromptManager
	validationManager *ValidationManager
	projectManager    *ProjectManager
}

// NewInitHandler creates a new InitHandler
func NewInitHandler() *InitHandler {
	handler := &InitHandler{
		validationManager: NewValidationManager(),
		projectManager:    NewProjectManager(),
	}
	handler.promptManager = NewPromptManager(handler)
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

	projectName, err := h.promptManager.PromptForProjectDetails()
	if err != nil {
		return pkgerrors.NewValidationError(pkgerrors.FieldProjectName, MsgFailedToGetProjectDetails, err)
	}

	services, validation, advanced, err := h.runServiceSelectionWorkflow(base)
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
		return pkgerrors.NewValidationError(FieldConfig, core.MsgNon_interactive_mode_requires_config, nil)
	}
	return nil
}

func (h *InitHandler) runServiceSelectionWorkflow(base *base.BaseCommand) ([]string, map[string]bool, map[string]bool, error) {
	base.Output.Header("%s", core.MsgProcess_initializing)
	logger.Info(logger.LogMsgProjectAction, logger.LogFieldAction, core.CommandInit, logger.LogFieldProject, "current_directory")

	for {
		services, err := h.promptManager.PromptForServices()
		if err != nil {
			return nil, nil, nil, pkgerrors.NewValidationError(pkgerrors.FieldServiceName, ActionSelectServices, err)
		}
		h.selectedServices = services

		if err := h.validateServices(services); err != nil {
			return nil, nil, nil, pkgerrors.NewValidationError(pkgerrors.FieldServiceName, ActionValidateServices, err)
		}

		validation, advanced, err := h.promptManager.PromptForAdvancedOptions()
		if err != nil {
			return nil, nil, nil, pkgerrors.NewValidationError(FieldOptions, ActionGetOptions, err)
		}

		if err := h.validationManager.RunValidations(validation, h, base); err != nil {
			return nil, nil, nil, pkgerrors.NewValidationError(FieldValidation, ActionValidation, err)
		}

		action, err := h.promptManager.ConfirmInitialization("", services, validation, advanced, base)
		if err != nil {
			return nil, nil, nil, pkgerrors.NewValidationError(FieldConfirmation, ActionGetConfirmation, err)
		}

		switch action {
		case core.ActionProceed:
			return services, validation, advanced, nil
		case core.ActionBack:
			base.Output.Info("%s", core.MsgInit_going_back)
			continue
		default:
			base.Output.Info("%s", core.MsgInit_cancelled)
			return nil, nil, nil, nil
		}
	}
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
