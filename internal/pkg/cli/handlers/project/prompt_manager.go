package project

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ErrGoBack is returned when user wants to go back
var ErrGoBack = errors.New("go back")

// PromptManager handles all user prompts for project initialization
type PromptManager struct {
	validator ProjectValidator
}

// ProjectValidator interface for validating project inputs
type ProjectValidator interface {
	ValidateProjectName(name string) error
}

// NewPromptManager creates a new prompt manager
func NewPromptManager(validator ProjectValidator) *PromptManager {
	return &PromptManager{
		validator: validator,
	}
}

// PromptForProjectDetails prompts user for project configuration
func (pm *PromptManager) PromptForProjectDetails() (string, error) {
	var projectName string

	// Get current directory name as default project name
	currentDir, err := filepath.Abs(".")
	if err != nil {
		return "", pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, "validation failed", err)
	}
	defaultName := filepath.Base(currentDir)

	// Project name prompt
	namePrompt := &survey.Input{
		Message: core.PromptProjectName,
		Default: defaultName,
		Help:    core.HelpProjectName,
	}

	if err := survey.AskOne(namePrompt, &projectName, survey.WithValidator(func(ans any) error {
		return pm.validator.ValidateProjectName(ans.(string))
	})); err != nil {
		return "", pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, "validation failed", err)
	}

	return projectName, nil
}

// PromptForServiceConfigs prompts user to select services and returns ServiceConfigs
func (pm *PromptManager) PromptForServiceConfigs() ([]types.ServiceConfig, error) {
	categories, err := pm.loadServiceCategories()
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return nil, pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldProjectName, "validation failed", nil)
	}

	// New approach: Show all services in one list with category labels
	return pm.selectServicesFromAllCategories(categories)
}

// PromptForAdvancedOptions prompts for validation and advanced options
func (pm *PromptManager) PromptForAdvancedOptions() (map[string]bool, map[string]bool, error) {
	validation := make(map[string]bool)
	advanced := make(map[string]bool)

	enableValidation, err := pm.askToEnableValidation()
	if err != nil {
		return validation, advanced, err
	}

	if enableValidation {
		validation, err = pm.selectValidationOptions()
		if err != nil {
			return validation, advanced, err
		}
	}

	return validation, advanced, nil
}

func (pm *PromptManager) askToEnableValidation() (bool, error) {
	validationPrompt := &survey.Confirm{
		Message: "Select validation options:",
		Default: true,
		Help:    "Choose which validations to run",
	}

	var enableValidation bool
	if err := survey.AskOne(validationPrompt, &enableValidation); err != nil {
		return false, pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, "validation", "validation failed", err)
	}
	return enableValidation, nil
}

func (pm *PromptManager) selectValidationOptions() (map[string]bool, error) {
	validationOptions := pm.buildValidationOptionsList()
	selectedValidations, err := pm.promptForValidationSelection(validationOptions)
	if err != nil {
		return nil, err
	}
	return pm.mapValidationSelections(selectedValidations), nil
}

func (pm *PromptManager) buildValidationOptionsList() []string {
	validationOptions := make([]string, 0, len(core.ValidationOptions))
	for _, description := range core.ValidationOptions {
		validationOptions = append(validationOptions, description)
	}
	return validationOptions
}

func (pm *PromptManager) promptForValidationSelection(validationOptions []string) ([]string, error) {
	var selectedValidations []string
	validationSelectPrompt := &survey.MultiSelect{
		Message: "Select validation options:",
		Options: validationOptions,
		Default: validationOptions,
		Help:    "Choose which validations to run",
	}

	if err := survey.AskOne(validationSelectPrompt, &selectedValidations); err != nil {
		return nil, pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, "validation", "validation failed", err)
	}
	return selectedValidations, nil
}

func (pm *PromptManager) mapValidationSelections(selectedValidations []string) map[string]bool {
	validation := make(map[string]bool)
	for _, selectedDesc := range selectedValidations {
		for key, description := range core.ValidationOptions {
			if description == selectedDesc {
				validation[key] = true
				break
			}
		}
	}
	return validation
}

// ConfirmInitialization shows final confirmation with option to go back
func (pm *PromptManager) ConfirmInitialization(projectName string, services []string, validation, advanced map[string]bool, base *base.BaseCommand) (string, error) {
	// Display summary
	base.Output.Info("Project Configuration Summary:")
	base.Output.Info("  Project Name: %s", projectName)
	base.Output.Info("  Services: %s", strings.Join(services, ", "))

	if len(validation) > 0 {
		base.Output.Info("  Validation Options:")
		for option := range validation {
			base.Output.Info("    - %s", option)
		}
	}

	// Confirmation prompt with back option
	confirmPrompt := &survey.Select{
		Message: "Proceed with initialization?",
		Options: []string{
			core.ActionProceed,
			core.ActionBack,
		},
		Default: core.ActionProceed,
	}

	var action string
	if err := survey.AskOne(confirmPrompt, &action); err != nil {
		return "", pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, "validation", "validation failed", err)
	}

	return action, nil
}

// loadServiceCategories loads available service categories
func (pm *PromptManager) loadServiceCategories() (map[string][]types.ServiceConfig, error) {
	utils := services.NewServiceUtils()
	return utils.GetServicesByCategory()
}

// selectServicesFromAllCategories shows all services in one list with category labels
func (pm *PromptManager) selectServicesFromAllCategories(categories map[string][]types.ServiceConfig) ([]types.ServiceConfig, error) {
	allServices, serviceOptions := pm.buildServiceList(categories)

	selected, err := pm.promptForServiceSelection(serviceOptions)
	if err != nil {
		return nil, err
	}

	if len(selected) == 0 {
		return nil, pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldProjectName, "validation failed", nil)
	}

	return pm.mapSelectedServices(selected, allServices), nil
}

func (pm *PromptManager) buildServiceList(categories map[string][]types.ServiceConfig) ([]types.ServiceConfig, []string) {
	var allServices []types.ServiceConfig
	var serviceOptions []string
	caser := cases.Title(language.English)

	for categoryName, categoryServices := range categories {
		for _, service := range categoryServices {
			displayName := fmt.Sprintf("[%s] %s - %s",
				caser.String(categoryName), service.Name, service.Description)
			serviceOptions = append(serviceOptions, displayName)
			allServices = append(allServices, service)
		}
	}

	return allServices, serviceOptions
}

func (pm *PromptManager) promptForServiceSelection(serviceOptions []string) ([]string, error) {
	prompt := &survey.MultiSelect{
		Message: "Select services for your project:",
		Options: serviceOptions,
		Help:    "Use space to select, enter to confirm. Services are grouped by category.",
	}

	var selected []string
	if err := survey.AskOne(prompt, &selected); err != nil {
		return nil, pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, "failed to select services", err)
	}

	return selected, nil
}

func (pm *PromptManager) mapSelectedServices(selected []string, allServices []types.ServiceConfig) []types.ServiceConfig {
	var selectedConfigs []types.ServiceConfig
	selectedMap := make(map[string]bool)

	for _, selection := range selected {
		serviceName := pm.extractServiceName(selection)
		if serviceName == "" || selectedMap[serviceName] {
			continue
		}
		selectedMap[serviceName] = true

		if config := pm.findServiceConfig(serviceName, allServices); config != nil {
			selectedConfigs = append(selectedConfigs, *config)
		}
	}

	return selectedConfigs
}

func (pm *PromptManager) extractServiceName(selection string) string {
	const minParts = 2
	parts := strings.Split(selection, "] ")
	if len(parts) < minParts {
		return ""
	}

	serviceNamePart := strings.Split(parts[1], " - ")
	if len(serviceNamePart) < 1 {
		return ""
	}

	return serviceNamePart[0]
}

func (pm *PromptManager) findServiceConfig(serviceName string, allServices []types.ServiceConfig) *types.ServiceConfig {
	for _, service := range allServices {
		if service.Name == serviceName {
			return &service
		}
	}
	return nil
}
