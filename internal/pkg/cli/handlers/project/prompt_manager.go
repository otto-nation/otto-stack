package project

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
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
		return "", pkgerrors.NewValidationError(pkgerrors.FieldProjectPath, ActionGetCurrentDir, err)
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
		return "", pkgerrors.NewValidationError(pkgerrors.FieldProjectName, ActionGetProjectName, err)
	}

	return projectName, nil
}

// PromptForServiceConfigs prompts user to select services and returns ServiceConfigs
func (pm *PromptManager) PromptForServiceConfigs() ([]services.ServiceConfig, error) {
	categories, err := pm.loadServiceCategories()
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return nil, pkgerrors.NewValidationError(pkgerrors.FieldServiceName, MsgNoServicesAvailable, nil)
	}

	categoryNames, categoryServicesList := pm.prepareCategoryNavigation(categories)
	return pm.navigateServiceCategoriesForConfigs(categoryNames, categoryServicesList)
}

// PromptForAdvancedOptions prompts for validation and advanced options
func (pm *PromptManager) PromptForAdvancedOptions() (map[string]bool, map[string]bool, error) {
	validation := make(map[string]bool)
	advanced := make(map[string]bool)

	// Validation options prompt
	validationPrompt := &survey.Confirm{
		Message: "Select validation options:",
		Default: true,
		Help:    "Choose which validations to run",
	}

	var enableValidation bool
	if err := survey.AskOne(validationPrompt, &enableValidation); err != nil {
		return validation, advanced, pkgerrors.NewValidationError(FieldValidation, MsgFailedToGetValidationPreference, err)
	}

	if enableValidation {
		// Individual validation options
		// Use ValidationOptions from generated constants
		validationOptions := make([]string, 0, len(core.ValidationOptions))
		for _, description := range core.ValidationOptions {
			validationOptions = append(validationOptions, description)
		}

		var selectedValidations []string
		validationSelectPrompt := &survey.MultiSelect{
			Message: "Select validation options:",
			Options: validationOptions,
			Default: validationOptions, // All enabled by default
			Help:    "Choose which validations to run",
		}

		if err := survey.AskOne(validationSelectPrompt, &selectedValidations); err != nil {
			return validation, advanced, pkgerrors.NewValidationError(FieldValidation, MsgFailedToGetValidationOptions, err)
		}

		// Convert descriptions back to keys
		for _, selectedDesc := range selectedValidations {
			for key, description := range core.ValidationOptions {
				if description == selectedDesc {
					validation[key] = true
					break
				}
			}
		}
	}

	return validation, advanced, nil
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
		return "", pkgerrors.NewValidationError(FieldAction, MsgFailedToGetAction, err)
	}

	return action, nil
}

// loadServiceCategories loads available service categories
func (pm *PromptManager) loadServiceCategories() (map[string][]services.ServiceConfig, error) {
	utils := services.NewServiceUtils()
	return utils.GetServicesByCategory()
}

// prepareCategoryNavigation prepares data structures for category navigation
func (pm *PromptManager) prepareCategoryNavigation(categories map[string][]services.ServiceConfig) ([]string, [][]services.ServiceConfig) {
	categoryNames := make([]string, 0, len(categories))
	categoryServicesList := make([][]services.ServiceConfig, 0, len(categories))

	for categoryName, categoryServices := range categories {
		categoryNames = append(categoryNames, categoryName)
		categoryServicesList = append(categoryServicesList, categoryServices)
	}

	return categoryNames, categoryServicesList
}

// navigateServiceCategoriesForConfigs handles interactive category navigation and returns ServiceConfigs
func (pm *PromptManager) navigateServiceCategoriesForConfigs(categoryNames []string, categoryServicesList [][]services.ServiceConfig) ([]services.ServiceConfig, error) {
	var allSelectedServices []services.ServiceConfig

	for {
		selectedCategory, err := pm.promptForCategory(categoryNames)
		if err != nil {
			return nil, err
		}

		if selectedCategory == core.PromptFinishSelection {
			break
		}

		categoryIndex := pm.findCategoryIndex(categoryNames, selectedCategory)
		if categoryIndex == -1 {
			continue
		}

		selectedServices, shouldGoBack, err := pm.selectServicesFromCategory(categoryIndex, categoryServicesList)
		if err != nil {
			return nil, err
		}

		if shouldGoBack {
			continue
		}

		allSelectedServices = append(allSelectedServices, selectedServices...)
	}

	if len(allSelectedServices) == 0 {
		return nil, pkgerrors.NewValidationError(pkgerrors.FieldServiceName, MsgNoServicesSelected, nil)
	}

	return allSelectedServices, nil
}

// promptForCategory prompts user to select a category
func (pm *PromptManager) promptForCategory(categoryNames []string) (string, error) {
	categoryPrompt := &survey.Select{
		Message: core.PromptSelectCategory,
		Options: append(categoryNames, core.PromptFinishSelection),
	}

	var selectedCategory string
	if err := survey.AskOne(categoryPrompt, &selectedCategory); err != nil {
		return "", pkgerrors.NewValidationError(pkgerrors.FieldServiceName, "failed to select category", err)
	}

	return selectedCategory, nil
}

// findCategoryIndex finds the index of a category by name
func (pm *PromptManager) findCategoryIndex(categoryNames []string, selectedCategory string) int {
	for i, name := range categoryNames {
		if name == selectedCategory {
			return i
		}
	}
	return -1
}

// selectServicesFromCategory handles service selection for a specific category
func (pm *PromptManager) selectServicesFromCategory(categoryIndex int, categoryServicesList [][]services.ServiceConfig) ([]services.ServiceConfig, bool, error) {
	categoryServices := categoryServicesList[categoryIndex]
	serviceOptions := pm.buildServiceOptions(categoryServices, true)
	return pm.promptCategoryServicesForConfigs(categoryServices[0].Category, serviceOptions, categoryServices)
}

// buildServiceOptions creates service options for prompts
func (pm *PromptManager) buildServiceOptions(categoryServices []services.ServiceConfig, allowGoBack bool) []string {
	const extraOptions = 2 // Space for "Go Back" and potential future options
	options := make([]string, 0, len(categoryServices)+extraOptions)

	for _, service := range categoryServices {
		displayName := fmt.Sprintf("%s - %s", service.Name, service.Description)
		options = append(options, displayName)
	}

	if allowGoBack {
		options = append(options, core.PromptGoBackOption)
	}

	return options
}

// promptCategoryServicesForConfigs prompts for service selection and returns ServiceConfigs
func (pm *PromptManager) promptCategoryServicesForConfigs(categoryName string, serviceOptions []string, categoryServices []services.ServiceConfig) ([]services.ServiceConfig, bool, error) {
	prompt := &survey.MultiSelect{
		Message: fmt.Sprintf("Select services from %s category:", categoryName),
		Options: serviceOptions,
		Help:    "Use space to select, enter to confirm",
	}

	var selected []string
	if err := survey.AskOne(prompt, &selected); err != nil {
		return nil, false, pkgerrors.NewValidationError(pkgerrors.FieldServiceName, "failed to select services", err)
	}

	// Check if user selected "Go Back"
	if slices.Contains(selected, core.PromptGoBackOption) {
		return nil, true, nil
	}

	// Map selected display names back to ServiceConfigs
	var selectedConfigs []services.ServiceConfig
	for _, selection := range selected {
		if selection != core.PromptGoBackOption {
			// Extract service name (before " - ")
			parts := strings.Split(selection, " - ")
			if len(parts) > 0 {
				serviceName := parts[0]
				// Find the corresponding ServiceConfig
				for _, serviceConfig := range categoryServices {
					if serviceConfig.Name == serviceName {
						selectedConfigs = append(selectedConfigs, serviceConfig)
						break
					}
				}
			}
		}
	}

	return selectedConfigs, false, nil
}
