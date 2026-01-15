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
