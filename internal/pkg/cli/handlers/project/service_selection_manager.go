package project

import (
	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	svc "github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// ServiceSelectionManager handles service selection workflow
type ServiceSelectionManager struct {
	promptManager     *PromptManager
	validationManager *ValidationManager
}

// SelectionResult holds the result of service selection
type SelectionResult struct {
	ServiceConfigs []types.ServiceConfig
	Validation     map[string]bool
	Advanced       map[string]bool
	Action         string
}

// NewServiceSelectionManager creates a new service selection manager
func NewServiceSelectionManager(promptManager *PromptManager, validationManager *ValidationManager) *ServiceSelectionManager {
	return &ServiceSelectionManager{
		promptManager:     promptManager,
		validationManager: validationManager,
	}
}

// RunWorkflow executes the complete service selection workflow
func (ssm *ServiceSelectionManager) RunWorkflow(handler *InitHandler, base *base.BaseCommand) (*SelectionResult, error) {
	ssm.logWorkflowStart(base)

	for {
		result, err := ssm.runSelectionCycle(handler, base)
		if err != nil {
			return nil, err
		}

		if shouldProceed := ssm.handleAction(result.Action, base); shouldProceed {
			return result, nil
		}

		if result.Action != core.ActionBack {
			return nil, nil
		}
	}
}

func (ssm *ServiceSelectionManager) logWorkflowStart(base *base.BaseCommand) {
	base.Output.Header("%s", messages.ProcessInitializing)
	logger.Info(logger.LogMsgProjectAction, logger.LogFieldAction, core.CommandInit, logger.LogFieldProject, "current_directory")
}

func (ssm *ServiceSelectionManager) runSelectionCycle(handler *InitHandler, base *base.BaseCommand) (*SelectionResult, error) {
	serviceConfigs, err := ssm.selectAndValidateServices(handler)
	if err != nil {
		return nil, err
	}

	validation, advanced, err := ssm.getAdvancedOptions()
	if err != nil {
		return nil, err
	}

	if err := ssm.runValidationChecks(validation, handler, serviceConfigs, base); err != nil {
		return nil, err
	}

	action, err := ssm.confirmSelection(serviceConfigs, validation, advanced, base)
	if err != nil {
		return nil, err
	}

	return &SelectionResult{
		ServiceConfigs: serviceConfigs,
		Validation:     validation,
		Advanced:       advanced,
		Action:         action,
	}, nil
}

func (ssm *ServiceSelectionManager) selectAndValidateServices(handler *InitHandler) ([]types.ServiceConfig, error) {
	serviceConfigs, err := ssm.promptManager.PromptForServiceConfigs()
	if err != nil {
		return nil, pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, "validation failed", err)
	}

	if err := handler.validateServiceConfigs(serviceConfigs); err != nil {
		return nil, pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, "validation failed", err)
	}

	return serviceConfigs, nil
}

func (ssm *ServiceSelectionManager) getAdvancedOptions() (map[string]bool, map[string]bool, error) {
	validation, advanced, err := ssm.promptManager.PromptForAdvancedOptions()
	if err != nil {
		return nil, nil, pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, "validation", "validation failed", err)
	}
	return validation, advanced, nil
}

func (ssm *ServiceSelectionManager) runValidationChecks(validation map[string]bool, handler *InitHandler, serviceConfigs []types.ServiceConfig, base *base.BaseCommand) error {
	if err := ssm.validationManager.RunValidations(validation, handler, serviceConfigs, base); err != nil {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, "validation-check", "validation failed", err)
	}
	return nil
}

func (ssm *ServiceSelectionManager) confirmSelection(serviceConfigs []types.ServiceConfig, validation, advanced map[string]bool, base *base.BaseCommand) (string, error) {
	serviceNames := svc.ExtractServiceNames(serviceConfigs)
	action, err := ssm.promptManager.ConfirmInitialization("", serviceNames, validation, advanced, base)
	if err != nil {
		return "", pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, "validation", "validation failed", err)
	}
	return action, nil
}

func (ssm *ServiceSelectionManager) handleAction(action string, base *base.BaseCommand) bool {
	switch action {
	case core.ActionProceed:
		return true
	case core.ActionBack:
		base.Output.Info("%s", messages.InitGoingBack)
		return false
	default:
		base.Output.Info("%s", messages.InitCancelled)
		return false
	}
}
