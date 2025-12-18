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

// PromptForServices prompts user to select services with category navigation
func (pm *PromptManager) PromptForServices() ([]string, error) {
	categories, err := pm.loadServiceCategories()
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return nil, pkgerrors.NewValidationError(pkgerrors.FieldServiceName, MsgNoServicesAvailable, nil)
	}

	categoryNames, categoryServicesList := pm.prepareCategoryNavigation(categories)
	return pm.navigateServiceCategories(categoryNames, categoryServicesList)
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
		validationOptions := []string{
			"Docker Compose validation",
			"Service health checks",
			"Port conflict detection",
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

		// Convert to map
		for _, option := range selectedValidations {
			validation[option] = true
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

// navigateServiceCategories handles the interactive category navigation
func (pm *PromptManager) navigateServiceCategories(categoryNames []string, categoryServicesList [][]services.ServiceConfig) ([]string, error) {
	var allSelectedServices []string

	for {
		// Category selection prompt
		categoryPrompt := &survey.Select{
			Message: "Select a service category (or finish selection):",
			Options: append(categoryNames, "Finish selection"),
		}

		var selectedCategory string
		if err := survey.AskOne(categoryPrompt, &selectedCategory); err != nil {
			return nil, pkgerrors.NewValidationError(pkgerrors.FieldServiceName, "failed to select category", err)
		}

		if selectedCategory == "Finish selection" {
			break
		}

		// Find category index
		categoryIndex := -1
		for i, name := range categoryNames {
			if name == selectedCategory {
				categoryIndex = i
				break
			}
		}

		if categoryIndex == -1 {
			continue
		}

		// Prompt for services in this category
		categoryServices := categoryServicesList[categoryIndex]
		serviceOptions := pm.buildServiceOptions(categoryServices, true)
		selectedServices, goBack, err := pm.promptCategoryServices(selectedCategory, serviceOptions)

		if err != nil {
			return nil, err
		}

		if goBack {
			continue
		}

		allSelectedServices = append(allSelectedServices, selectedServices...)
	}

	if len(allSelectedServices) == 0 {
		return nil, pkgerrors.NewValidationError(pkgerrors.FieldServiceName, MsgNoServicesSelected, nil)
	}

	return allSelectedServices, nil
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
		options = append(options, "Go Back")
	}

	return options
}

// promptCategoryServices prompts for service selection within a category
func (pm *PromptManager) promptCategoryServices(categoryName string, serviceOptions []string) ([]string, bool, error) {
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
	if slices.Contains(selected, "Go Back") {
		return nil, true, nil
	}

	// Extract service names from display names
	var serviceNames []string
	for _, selection := range selected {
		if selection != "Go Back" {
			// Extract service name (before " - ")
			parts := strings.Split(selection, " - ")
			if len(parts) > 0 {
				serviceNames = append(serviceNames, parts[0])
			}
		}
	}

	return serviceNames, false, nil
}
