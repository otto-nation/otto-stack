package project

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
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
		return "", err
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
		return "", err
	}

	return projectName, nil
}

// PromptForServiceConfigs prompts user to select services and returns ServiceConfigs
func (pm *PromptManager) PromptForServiceConfigs() ([]types.ServiceConfig, error) {
	selector := NewServiceSelector()
	return selector.SelectServices()
}

// PromptForAdvancedOptions prompts for validation and advanced options
func (pm *PromptManager) PromptForAdvancedOptions() (map[string]bool, map[string]bool, error) {
	advanced := make(map[string]bool)

	prompter := NewValidationPrompter()
	validation, err := prompter.PromptForValidationOptions()
	if err != nil {
		return validation, advanced, err
	}

	return validation, advanced, nil
}

// InitConfirmation encapsulates initialization confirmation parameters
type InitConfirmation struct {
	ProjectName string
	Services    []string
	Validation  map[string]bool
	Advanced    map[string]bool
	Base        *base.BaseCommand
}

// ConfirmInitialization shows final confirmation with option to go back
func (pm *PromptManager) ConfirmInitialization(projectName string, services []string, validation, advanced map[string]bool, base *base.BaseCommand) (string, error) {
	conf := InitConfirmation{
		ProjectName: projectName,
		Services:    services,
		Validation:  validation,
		Advanced:    advanced,
		Base:        base,
	}
	return pm.confirmInitializationWithConfig(conf)
}

func (pm *PromptManager) confirmInitializationWithConfig(conf InitConfirmation) (string, error) {
	// Display summary
	conf.Base.Output.Info(messages.InfoProjectConfigSummary)
	conf.Base.Output.Info("  Project Name: %s", conf.ProjectName)
	conf.Base.Output.Info("  Services: %s", strings.Join(conf.Services, ", "))

	if len(conf.Validation) > 0 {
		conf.Base.Output.Info(messages.InfoValidationOptions)
		// Sort keys for consistent display order
		keys := make([]string, 0, len(conf.Validation))
		for k := range conf.Validation {
			keys = append(keys, k)
		}
		for _, option := range keys {
			conf.Base.Output.Info("    - %s", option)
		}
	}

	// Confirmation prompt with back option
	confirmPrompt := &survey.Select{
		Message: messages.PromptsProceedInitialization,
		Options: []string{
			core.ActionProceed,
			core.ActionBack,
		},
		Default: core.ActionProceed,
	}

	var action string
	if err := survey.AskOne(confirmPrompt, &action); err != nil {
		return "", err
	}

	return action, nil
}
