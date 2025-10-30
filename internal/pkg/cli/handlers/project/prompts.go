package project

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/types"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
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
		Message: constants.PromptProjectName,
		Default: defaultName,
		Help:    constants.HelpProjectName,
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
	serviceUtils := utils.NewServiceUtils()

	// Get available services by category
	categories, err := serviceUtils.GetServicesByCategory()
	if err != nil {
		return nil, fmt.Errorf("failed to load services: %w", err)
	}

	if len(categories) == 0 {
		return nil, fmt.Errorf("no services available")
	}

	// Convert map to ordered slice for navigation
	var categoryNames []string
	var categoryServicesList [][]types.ServiceInfo
	for categoryName, services := range categories {
		if len(services) > 0 {
			categoryNames = append(categoryNames, categoryName)
			categoryServicesList = append(categoryServicesList, services)
		}
	}

	var allSelectedServices []string
	categoryIndex := 0

	// Navigate through categories with back support
	for categoryIndex < len(categoryNames) {
		categoryName := categoryNames[categoryIndex]
		services := categoryServicesList[categoryIndex]

		constants.SendMessage(constants.MsgSelectServices, categoryName)

		var serviceOptions []string
		for _, service := range services {
			description := service.Description
			if description == "" {
				description = constants.MsgNoDescription.Content
			}
			serviceOptions = append(serviceOptions, fmt.Sprintf("%s - %s", service.Name, description))
		}

		if categoryIndex > 0 {
			serviceOptions = append(serviceOptions, constants.PromptGoBack)
		}

		var selected []string
		prompt := &survey.MultiSelect{
			Message: fmt.Sprintf(constants.MsgSelectServices.Content, categoryName),
			Options: serviceOptions,
			Help:    constants.HelpServiceSelection,
		}

		// Set custom template with back instructions
		survey.MultiSelectQuestionTemplate = constants.MultiSelectTemplateWithBack

		if err := survey.AskOne(prompt, &selected); err != nil {
			return nil, fmt.Errorf("failed to select %s services: %w", categoryName, err)
		}

		// Check if user selected "Go Back"
		goBack := false
		var categoryServices []string
		for _, selection := range selected {
			if selection == constants.PromptGoBack {
				goBack = true
				break
			}
			serviceName := strings.Split(selection, " - ")[0]
			categoryServices = append(categoryServices, serviceName)
		}

		if goBack && categoryIndex > 0 {
			categoryIndex--
			continue
		}

		// Add selected services to the total
		allSelectedServices = append(allSelectedServices, categoryServices...)
		categoryIndex++
	}

	if len(allSelectedServices) == 0 {
		return nil, fmt.Errorf("no services selected")
	}

	return allSelectedServices, nil
}

// promptForAdvancedOptions prompts for advanced configuration options
func (h *InitHandler) promptForAdvancedOptions() (map[string]bool, map[string]bool, error) {
	validation := make(map[string]bool)
	advanced := make(map[string]bool)

	// Ask if user wants advanced options
	var wantsAdvanced bool
	advancedPrompt := &survey.Confirm{
		Message: constants.PromptAdvancedConfig,
		Default: false,
		Help:    constants.HelpAdvancedConfig,
	}

	if err := survey.AskOne(advancedPrompt, &wantsAdvanced); err != nil {
		return validation, advanced, fmt.Errorf("failed to get advanced options preference: %w", err)
	}

	if !wantsAdvanced {
		return validation, advanced, nil
	}

	// Validation options
	var validationKeys []string
	for key := range constants.ValidationOptions {
		validationKeys = append(validationKeys, key)
	}

	var selectedValidation []string
	validationPrompt := &survey.MultiSelect{
		Message: constants.PromptValidationOptions,
		Options: validationKeys,
		Help:    constants.HelpValidationOptions,
	}

	if err := survey.AskOne(validationPrompt, &selectedValidation); err != nil {
		return validation, advanced, fmt.Errorf("failed to get validation options: %w", err)
	}

	for _, option := range selectedValidation {
		if key, exists := constants.ValidationOptions[option]; exists {
			validation[key] = true
		}
	}

	// Advanced options
	var advancedKeys []string
	for key := range constants.AdvancedOptions {
		advancedKeys = append(advancedKeys, key)
	}

	var selectedAdvanced []string
	advancedPrompt2 := &survey.MultiSelect{
		Message: constants.PromptAdvancedFeatures,
		Options: advancedKeys,
		Help:    constants.HelpAdvancedFeatures,
	}

	if err := survey.AskOne(advancedPrompt2, &selectedAdvanced); err != nil {
		return validation, advanced, fmt.Errorf("failed to get advanced features: %w", err)
	}

	for _, option := range selectedAdvanced {
		if key, exists := constants.AdvancedOptions[option]; exists {
			advanced[key] = true
		}
	}

	return validation, advanced, nil
}

// confirmInitializationWithBack shows a summary and asks for confirmation with back option
func (h *InitHandler) confirmInitializationWithBack(projectName string, services []string, validation, advanced map[string]bool) (string, error) {
	constants.SendMessage(constants.MsgInitSummary)
	constants.SendMessage(constants.MsgProject, projectName)
	constants.SendMessage(constants.MsgServices, strings.Join(services, ", "))

	if len(validation) > 0 {
		var validationFeatures []string
		for feature := range validation {
			validationFeatures = append(validationFeatures, feature)
		}
		constants.SendMessage(constants.MsgValidation, strings.Join(validationFeatures, ", "))
	}

	if len(advanced) > 0 {
		var advancedFeatures []string
		for feature := range advanced {
			advancedFeatures = append(advancedFeatures, feature)
		}
		constants.SendMessage(constants.MsgAdvanced, strings.Join(advancedFeatures, ", "))
	}

	var action string
	actionPrompt := &survey.Select{
		Message: constants.PromptActionSelect,
		Options: constants.ActionOptions,
		Default: constants.PromptProceedInit,
	}

	if err := survey.AskOne(actionPrompt, &action); err != nil {
		return "", fmt.Errorf("failed to get action: %w", err)
	}

	if mappedAction, exists := constants.ActionOptionMap[action]; exists {
		return mappedAction, nil
	}

	return constants.ActionCancel, nil
}
