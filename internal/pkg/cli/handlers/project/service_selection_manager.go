package project

import (
	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
)

// ServiceSelectionManager handles service selection workflow
type ServiceSelectionManager struct {
	promptManager     *PromptManager
	validationManager *ValidationManager
}

// NewServiceSelectionManager creates a new service selection manager
func NewServiceSelectionManager(promptManager *PromptManager, validationManager *ValidationManager) *ServiceSelectionManager {
	return &ServiceSelectionManager{
		promptManager:     promptManager,
		validationManager: validationManager,
	}
}

// RunWorkflow executes the complete service selection workflow
func (ssm *ServiceSelectionManager) RunWorkflow(handler *InitHandler, base *base.BaseCommand) ([]string, map[string]bool, map[string]bool, error) {
	base.Output.Header("%s", core.MsgProcess_initializing)
	logger.Info(logger.LogMsgProjectAction, logger.LogFieldAction, core.CommandInit, logger.LogFieldProject, "current_directory")

	for {
		services, err := ssm.promptManager.PromptForServices()
		if err != nil {
			return nil, nil, nil, pkgerrors.NewValidationError(pkgerrors.FieldServiceName, ActionSelectServices, err)
		}
		handler.selectedServices = services

		if err := handler.validateServices(services); err != nil {
			return nil, nil, nil, pkgerrors.NewValidationError(pkgerrors.FieldServiceName, ActionValidateServices, err)
		}

		validation, advanced, err := ssm.promptManager.PromptForAdvancedOptions()
		if err != nil {
			return nil, nil, nil, pkgerrors.NewValidationError(FieldOptions, ActionGetOptions, err)
		}

		if err := ssm.validationManager.RunValidations(validation, handler, base); err != nil {
			return nil, nil, nil, pkgerrors.NewValidationError(FieldValidation, ActionValidation, err)
		}

		action, err := ssm.promptManager.ConfirmInitialization("", services, validation, advanced, base)
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
