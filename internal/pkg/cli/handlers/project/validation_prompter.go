package project

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
)

// ValidationPrompter handles validation option prompts
type ValidationPrompter struct{}

// NewValidationPrompter creates a new validation prompter
func NewValidationPrompter() *ValidationPrompter {
	return &ValidationPrompter{}
}

// PromptForValidationOptions prompts user for validation options
func (vp *ValidationPrompter) PromptForValidationOptions() (map[string]bool, error) {
	enableValidation, err := vp.askToEnableValidation()
	if err != nil {
		return nil, err
	}

	if !enableValidation {
		return make(map[string]bool), nil
	}

	return vp.selectValidationOptions()
}

func (vp *ValidationPrompter) askToEnableValidation() (bool, error) {
	validationPrompt := &survey.Confirm{
		Message: messages.PromptsSelectValidationOptions,
		Default: true,
		Help:    messages.PromptsSelectValidationOptionsHelp,
	}

	var enableValidation bool
	if err := survey.AskOne(validationPrompt, &enableValidation); err != nil {
		return false, err
	}
	return enableValidation, nil
}

func (vp *ValidationPrompter) selectValidationOptions() (map[string]bool, error) {
	validationOptions := vp.buildValidationOptionsList()
	selectedValidations, err := vp.promptForValidationSelection(validationOptions)
	if err != nil {
		return nil, err
	}
	return vp.mapValidationSelections(selectedValidations), nil
}

func (vp *ValidationPrompter) buildValidationOptionsList() []string {
	validationOptions := make([]string, 0, len(core.ValidationOptions))
	for _, description := range core.ValidationOptions {
		validationOptions = append(validationOptions, description)
	}
	return validationOptions
}

func (vp *ValidationPrompter) promptForValidationSelection(validationOptions []string) ([]string, error) {
	var selectedValidations []string
	validationSelectPrompt := &survey.MultiSelect{
		Message: messages.PromptsSelectValidationOptions,
		Options: validationOptions,
		Default: validationOptions,
		Help:    messages.PromptsSelectValidationOptionsHelp,
	}

	if err := survey.AskOne(validationSelectPrompt, &selectedValidations); err != nil {
		return nil, err
	}
	return selectedValidations, nil
}

func (vp *ValidationPrompter) mapValidationSelections(selectedValidations []string) map[string]bool {
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
