package project

import (
	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
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
	ProjectName    string
}

// NewServiceSelectionManager creates a new service selection manager
func NewServiceSelectionManager(promptManager *PromptManager, validationManager *ValidationManager) *ServiceSelectionManager {
	return &ServiceSelectionManager{
		promptManager:     promptManager,
		validationManager: validationManager,
	}
}

// RunWorkflow executes the complete service selection workflow
func (ssm *ServiceSelectionManager) RunWorkflow(handler *InitHandler, projectName string, base *base.BaseCommand) (*SelectionResult, error) {
	ssm.logWorkflowStart(base)

	for {
		result, err := ssm.runSelectionCycle(handler, projectName, base)
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

func (ssm *ServiceSelectionManager) runSelectionCycle(handler *InitHandler, projectName string, base *base.BaseCommand) (*SelectionResult, error) {
	selectedConfigs, err := ssm.promptManager.PromptForServiceConfigs()
	if err != nil {
		return nil, err
	}

	// Resolve dependencies for selected services
	serviceNames := services.ExtractServiceNames(selectedConfigs)
	serviceConfigs, err := services.ResolveUpServices(serviceNames, nil)
	if err != nil {
		return nil, err
	}

	validator := services.NewValidator()
	if err := validator.ValidateServiceConfigs(serviceConfigs); err != nil {
		return nil, err
	}

	validation, advanced, err := ssm.promptManager.PromptForAdvancedOptions()
	if err != nil {
		return nil, err
	}

	if err := ssm.validationManager.RunValidations(validation, handler, serviceConfigs, base); err != nil {
		return nil, err
	}

	action, err := ssm.promptManager.ConfirmInitialization(projectName, serviceNames, validation, advanced, base)
	if err != nil {
		return nil, err
	}

	return &SelectionResult{
		ServiceConfigs: serviceConfigs,
		Validation:     validation,
		Advanced:       advanced,
		Action:         action,
		ProjectName:    projectName,
	}, nil
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
