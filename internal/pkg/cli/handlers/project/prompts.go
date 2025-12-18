package project

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

// ErrGoBack is returned when user wants to go back
var ErrGoBack = errors.New("go back")

// promptForProjectDetails prompts user for project configuration
func (h *InitHandler) promptForProjectDetails() (string, error) {
	var projectName string

	// Get current directory name as default project name
	currentDir, err := filepath.Abs(".")
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	defaultName := filepath.Base(currentDir)

	// Project name prompt
	namePrompt := &survey.Input{
		Message: core.PromptProjectName,
		Default: defaultName,
		Help:    core.HelpProjectName,
	}

	if err := survey.AskOne(namePrompt, &projectName, survey.WithValidator(func(ans any) error {
		return h.validateProjectName(ans.(string))
	})); err != nil {
		return "", fmt.Errorf("failed to get project name: %w", err)
	}

	return projectName, nil
}

// promptForServices prompts user to select services with category navigation
func (h *InitHandler) promptForServices() ([]string, error) {
	categories, err := h.loadServiceCategories()
	if err != nil {
		return nil, err
	}

	categoryNames, categoryServicesList := h.prepareCategoryNavigation(categories)
	if len(categoryNames) == 0 {
		return nil, fmt.Errorf("no services available")
	}

	return h.navigateServiceCategories(categoryNames, categoryServicesList)
}

func (h *InitHandler) loadServiceCategories() (map[string][]services.ServiceConfig, error) {
	serviceUtils := services.NewServiceUtils()
	categories, err := serviceUtils.GetServicesByCategory()
	if err != nil {
		return nil, fmt.Errorf("failed to load services: %w", err)
	}
	return categories, nil
}

func (h *InitHandler) prepareCategoryNavigation(categories map[string][]services.ServiceConfig) ([]string, [][]services.ServiceConfig) {
	var categoryNames []string
	var categoryServicesList [][]services.ServiceConfig

	for categoryName, categoryServices := range categories {
		if len(categoryServices) > 0 {
			categoryNames = append(categoryNames, categoryName)
			categoryServicesList = append(categoryServicesList, categoryServices)
		}
	}
	return categoryNames, categoryServicesList
}

func (h *InitHandler) navigateServiceCategories(categoryNames []string, categoryServicesList [][]services.ServiceConfig) ([]string, error) {
	var allSelectedServices []string
	categoryIndex := 0

	for categoryIndex < len(categoryNames) {
		categoryName := categoryNames[categoryIndex]
		categoryServices := categoryServicesList[categoryIndex]

		serviceOptions := h.buildServiceOptions(categoryServices, categoryIndex > 0)

		selectedServiceNames, goBack, err := h.promptCategoryServices(categoryName, serviceOptions)
		if err != nil {
			return nil, err
		}

		if goBack && categoryIndex > 0 {
			categoryIndex--
			continue
		}

		allSelectedServices = append(allSelectedServices, selectedServiceNames...)
		categoryIndex++
	}

	if len(allSelectedServices) == 0 {
		return nil, fmt.Errorf("no services selected")
	}
	return allSelectedServices, nil
}

func (h *InitHandler) buildServiceOptions(categoryServices []services.ServiceConfig, allowGoBack bool) []string {
	var serviceOptions []string
	for _, service := range categoryServices {
		description := service.Description
		if description == "" {
			description = core.MsgServices_no_description
		}
		serviceOptions = append(serviceOptions, fmt.Sprintf("%s - %s", service.Name, description))
	}

	if allowGoBack {
		serviceOptions = append(serviceOptions, core.PromptGoBack)
	}
	return serviceOptions
}

func (h *InitHandler) promptCategoryServices(categoryName string, serviceOptions []string) ([]string, bool, error) {
	var selected []string
	prompt := &survey.MultiSelect{
		Message: fmt.Sprintf(core.MsgServices_select, categoryName),
		Options: serviceOptions,
		Help:    core.HelpServiceSelection,
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		return nil, false, fmt.Errorf("failed to select %s services: %w", categoryName, err)
	}

	var selectedServiceNames []string
	for _, selection := range selected {
		if selection == core.PromptGoBack {
			return nil, true, nil
		}
		serviceName := strings.Split(selection, " - ")[0]
		selectedServiceNames = append(selectedServiceNames, serviceName)
	}

	return selectedServiceNames, false, nil
}

// promptForAdvancedOptions prompts for advanced configuration options
func (h *InitHandler) promptForAdvancedOptions() (map[string]bool, map[string]bool, error) {
	validation := make(map[string]bool)
	advanced := make(map[string]bool)

	// Check if there are any validation options available
	if len(core.ValidationOptions) == 0 {
		return validation, advanced, nil
	}

	// Ask if user wants validation options
	var wantsValidation bool
	validationPrompt := &survey.Confirm{
		Message: "Enable validation checks?",
		Default: true,
		Help:    "Run validation checks to catch potential issues early",
	}

	if err := survey.AskOne(validationPrompt, &wantsValidation); err != nil {
		return validation, advanced, fmt.Errorf("failed to get validation preference: %w", err)
	}

	if !wantsValidation {
		return validation, advanced, nil
	}

	// Build validation options dynamically from available options (only optional ones)
	var validationOptions []string
	var descriptionToKey = make(map[string]string)

	optionalValidations := []string{
		core.ValidationPorts,
		core.ValidationResourceLimits,
		core.ValidationEnvironmentVariables,
		core.ValidationFilePermissions,
		core.ValidationNetworkConnectivity,
		core.ValidationStorageRequirements,
	}

	for _, key := range optionalValidations {
		if description, exists := core.ValidationOptions[key]; exists {
			validationOptions = append(validationOptions, description)
			descriptionToKey[description] = key
		}
	}

	var selectedValidation []string
	validationSelectPrompt := &survey.MultiSelect{
		Message: "Select validation checks:",
		Options: validationOptions,
		Help:    "Choose which validation checks to run during initialization",
	}

	if err := survey.AskOne(validationSelectPrompt, &selectedValidation); err != nil {
		return validation, advanced, fmt.Errorf("failed to get validation options: %w", err)
	}

	// Map selected descriptions back to keys
	for _, description := range selectedValidation {
		if key, exists := descriptionToKey[description]; exists {
			validation[key] = true
		}
	}

	return validation, advanced, nil
}

// confirmInitializationWithBack shows a summary and asks for confirmation with back option
func (h *InitHandler) confirmInitializationWithBack(projectName string, services []string, validation, advanced map[string]bool, base *base.BaseCommand) (string, error) {
	base.Output.Header(core.MsgInit_summary)
	base.Output.Info(core.MsgInit_project_summary, projectName)
	base.Output.Info(core.MsgInit_services_summary, strings.Join(services, ", "))

	if len(validation) > 0 {
		var validationFeatures []string
		for feature := range validation {
			validationFeatures = append(validationFeatures, feature)
		}
		base.Output.Info(core.MsgInit_validation_summary, strings.Join(validationFeatures, ", "))
	}

	if len(advanced) > 0 {
		var advancedFeatures []string
		for feature := range advanced {
			advancedFeatures = append(advancedFeatures, feature)
		}
		base.Output.Info(core.MsgInit_advanced_summary, strings.Join(advancedFeatures, ", "))
	}

	var action string
	actionPrompt := &survey.Select{
		Message: core.PromptActionSelect,
		Options: []string{core.ActionProceed, core.ActionBack, core.ActionCancel},
		Default: core.ActionProceed,
	}

	if err := survey.AskOne(actionPrompt, &action); err != nil {
		return "", fmt.Errorf("failed to get action: %w", err)
	}

	return action, nil
}
