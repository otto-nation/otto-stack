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

	// New approach: Show all services in one list with category labels
	return pm.selectServicesFromAllCategories(categories)
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

// selectServicesFromAllCategories shows all services in one list with category labels
func (pm *PromptManager) selectServicesFromAllCategories(categories map[string][]services.ServiceConfig) ([]services.ServiceConfig, error) {
	// Build flat list of all services with category labels
	var allServices []services.ServiceConfig
	var serviceOptions []string

	for categoryName, categoryServices := range categories {
		for _, service := range categoryServices {
			// Format: [Category] ServiceName - Description
			caser := cases.Title(language.English)
			displayName := fmt.Sprintf("[%s] %s - %s",
				caser.String(categoryName), service.Name, service.Description)
			serviceOptions = append(serviceOptions, displayName)
			allServices = append(allServices, service)
		}
	}

	// Sort options for better UX
	// (Could add sorting by category then name if needed)

	prompt := &survey.MultiSelect{
		Message: "Select services for your project:",
		Options: serviceOptions,
		Help:    "Use space to select, enter to confirm. Services are grouped by category.",
	}

	var selected []string
	if err := survey.AskOne(prompt, &selected); err != nil {
		return nil, pkgerrors.NewValidationError(pkgerrors.FieldServiceName, "failed to select services", err)
	}

	if len(selected) == 0 {
		return nil, pkgerrors.NewValidationError(pkgerrors.FieldServiceName, MsgNoServicesSelected, nil)
	}

	// Map selected display names back to ServiceConfigs
	var selectedConfigs []services.ServiceConfig
	selectedMap := make(map[string]bool)

	const minParts = 2
	for _, selection := range selected {
		// Extract service name from "[Category] ServiceName - Description"
		parts := strings.Split(selection, "] ")
		if len(parts) < minParts {
			continue
		}
		serviceNamePart := strings.Split(parts[1], " - ")
		if len(serviceNamePart) < 1 {
			continue
		}
		serviceName := serviceNamePart[0]

		// Prevent duplicates
		if selectedMap[serviceName] {
			continue
		}
		selectedMap[serviceName] = true

		// Find the corresponding ServiceConfig
		for _, service := range allServices {
			if service.Name == serviceName {
				selectedConfigs = append(selectedConfigs, service)
				break
			}
		}
	}

	return selectedConfigs, nil
}
