package project

import (
	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	svc "github.com/otto-nation/otto-stack/internal/pkg/services"
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
func (ssm *ServiceSelectionManager) RunWorkflow(handler *InitHandler, base *base.BaseCommand) ([]svc.ServiceConfig, map[string]bool, map[string]bool, error) {
	base.Output.Header("%s", core.MsgProcess_initializing)
	logger.Info(logger.LogMsgProjectAction, logger.LogFieldAction, core.CommandInit, logger.LogFieldProject, "current_directory")

	for {
		serviceConfigs, err := ssm.promptManager.PromptForServiceConfigs()
		if err != nil {
			return nil, nil, nil, pkgerrors.NewValidationError(pkgerrors.FieldServiceName, ActionSelectServices, err)
		}
		if err := handler.validateServiceConfigs(serviceConfigs); err != nil {
			return nil, nil, nil, pkgerrors.NewValidationError(pkgerrors.FieldServiceName, ActionValidateServices, err)
		}

		validation, advanced, err := ssm.promptManager.PromptForAdvancedOptions()
		if err != nil {
			return nil, nil, nil, pkgerrors.NewValidationError(FieldOptions, ActionGetOptions, err)
		}

		if err := ssm.validationManager.RunValidations(validation, handler, serviceConfigs, base); err != nil {
			return nil, nil, nil, pkgerrors.NewValidationError("validation-check", ActionValidation, err)
		}

		serviceNames := svc.ExtractServiceNames(serviceConfigs)
		action, err := ssm.promptManager.ConfirmInitialization("", serviceNames, validation, advanced, base)
		if err != nil {
			return nil, nil, nil, pkgerrors.NewValidationError(FieldConfirmation, ActionGetConfirmation, err)
		}

		switch action {
		case core.ActionProceed:
			return serviceConfigs, validation, advanced, nil
		case core.ActionBack:
			base.Output.Info("%s", core.MsgInit_going_back)
			continue
		default:
			base.Output.Info("%s", core.MsgInit_cancelled)
			return nil, nil, nil, nil
		}
	}
}
