package project

import (
	"context"
	"fmt"
	"slices"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/spf13/cobra"
)

// InitHandler handles the init command
type InitHandler struct {
	serviceUtils     *services.ServiceUtils
	configService    services.ConfigService
	selectedServices []string
}

// NewInitHandler creates a new InitHandler
func NewInitHandler() *InitHandler {
	return &InitHandler{
		serviceUtils:  services.NewServiceUtils(),
		configService: services.NewConfigService(),
	}
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

	projectName, err := h.promptForProjectDetails()
	if err != nil {
		return pkgerrors.NewValidationError(pkgerrors.FieldProjectName, MsgFailedToGetProjectDetails, err)
	}

	services, validation, advanced, err := h.runServiceSelectionWorkflow(base)
	if err != nil {
		return err
	}

	if err := h.createProjectStructure(projectName, services, validation, advanced, base); err != nil {
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
		services, err := h.promptForServices()
		if err != nil {
			return nil, nil, nil, pkgerrors.NewValidationError(pkgerrors.FieldServiceName, ActionSelectServices, err)
		}
		h.selectedServices = services

		if err := h.validateServices(services); err != nil {
			return nil, nil, nil, pkgerrors.NewValidationError(pkgerrors.FieldServiceName, ActionValidateServices, err)
		}

		validation, advanced, err := h.promptForAdvancedOptions()
		if err != nil {
			return nil, nil, nil, pkgerrors.NewValidationError(FieldOptions, ActionGetOptions, err)
		}

		if err := h.runValidations(validation, base); err != nil {
			return nil, nil, nil, pkgerrors.NewValidationError(FieldValidation, ActionValidation, err)
		}

		action, err := h.confirmInitializationWithBack("", services, validation, advanced, base)
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

func (h *InitHandler) createProjectStructure(projectName string, services []string, validation, advanced map[string]bool, base *base.BaseCommand) error {
	if err := h.createDirectoryStructure(); err != nil {
		return pkgerrors.NewServiceError(ComponentProject, ActionCreateDirectories, err)
	}

	if err := h.createConfigFile(projectName, services, validation, base); err != nil {
		return pkgerrors.NewConfigError("", ActionCreateConfigFile, err)
	}

	h.generateServiceConfigs(services, base)

	if err := h.generateInitialComposeFiles(services, projectName, validation, advanced, base); err != nil {
		return pkgerrors.NewServiceError(ComponentCompose, ActionGenerateFiles, err)
	}

	if err := h.createGitignoreEntries(base); err != nil {
		base.Output.Warning(core.MsgWarnings_failed_gitignore, err)
	}

	if err := h.createReadme(projectName, services, base); err != nil {
		base.Output.Warning(core.MsgWarnings_failed_readme, err)
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
func (h *InitHandler) runValidations(selectedValidations map[string]bool, base *base.BaseCommand) error {
	// Always run required validations
	for validationKey := range core.ValidationOptions {
		validationFunc, exists := ValidationRegistry[validationKey]
		if !exists {
			continue
		}

		// Run if it's required OR if user selected it
		isRequired := isRequiredValidation(validationKey)
		isSelected := selectedValidations[validationKey]

		if isRequired || isSelected {
			if err := validationFunc(h, base); err != nil {
				return pkgerrors.NewValidationError(FieldValidation, "validation failed", err)
			}
		}
	}
	return nil
}

// isRequiredValidation checks if a validation is required based on YAML config
func isRequiredValidation(key string) bool {
	requiredValidations := []string{
		core.ValidationDocker,
		core.ValidationDockerCompose,
		core.ValidationConfigSyntax,
		core.ValidationServiceDefinitions,
	}

	return slices.Contains(requiredValidations, key)
}
