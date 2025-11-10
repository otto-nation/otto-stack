package ui

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

// ServiceOption represents a service option for selection
type ServiceOption struct {
	Name         string
	Description  string
	Dependencies []string
	Category     string
}

// PromptInput prompts for text input
func PromptInput(message, defaultValue string) (string, error) {
	var result string
	prompt := &survey.Input{
		Message: message,
		Default: defaultValue,
	}
	return result, survey.AskOne(prompt, &result)
}

// PromptConfirm prompts for yes/no confirmation
func PromptConfirm(message string, defaultValue bool) (bool, error) {
	var result bool
	prompt := &survey.Confirm{
		Message: message,
		Default: defaultValue,
	}
	return result, survey.AskOne(prompt, &result)
}

// PromptMultiSelect prompts for multiple selections
func PromptMultiSelect(message string, options []string) ([]string, error) {
	var result []string
	prompt := &survey.MultiSelect{
		Message: message,
		Options: options,
	}
	return result, survey.AskOne(prompt, &result)
}

// PromptServiceSelection prompts for service selection within categories
func PromptServiceSelection(categoryName string, services []ServiceOption) ([]string, error) {
	if len(services) == 0 {
		return []string{}, nil
	}

	options := make([]string, len(services))
	for i, service := range services {
		options[i] = service.Name + " - " + service.Description
	}

	selected, err := PromptMultiSelect(fmt.Sprintf("Select %s services:", categoryName), options)
	if err != nil {
		return nil, err
	}

	// Extract service names from selected options
	result := make([]string, len(selected))
	for i, sel := range selected {
		result[i] = strings.Split(sel, " - ")[0]
	}
	return result, nil
}

// PromptFinalConfirmation shows a summary and asks for final confirmation
func PromptFinalConfirmation(projectName, environment string, services, dependencies []string) (bool, error) {
	summary := fmt.Sprintf(`Project: %s
Environment: %s
Services: %s
Dependencies: %s
Total: %d services`,
		projectName,
		environment,
		strings.Join(services, ", "),
		strings.Join(dependencies, ", "),
		len(services)+len(dependencies))

	DefaultOutput.Box("Configuration Summary", summary)
	return PromptConfirm("Continue with this configuration?", true)
}
