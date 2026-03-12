package project

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
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

// PromptForAdvancedOptions prompts for validation, sharing, and advanced options.
// It receives the resolved service configs so it can show which services are shareable.
func (pm *PromptManager) PromptForAdvancedOptions(serviceConfigs []types.ServiceConfig) (map[string]bool, *clicontext.SharingSpec, *clicontext.AdvancedSpec, error) {
	prompter := NewValidationPrompter()
	validation, err := prompter.PromptForValidationOptions()
	if err != nil {
		return validation, nil, nil, err
	}

	sharing, err := pm.promptForSharing(serviceConfigs)
	if err != nil {
		return validation, nil, nil, err
	}

	autoStart, err := pm.promptForAutoStart()
	if err != nil {
		return validation, nil, nil, err
	}

	return validation, sharing, &clicontext.AdvancedSpec{AutoStart: autoStart}, nil
}

// promptForSharing prompts the user whether to enable shared containers for the
// shareable services in their selection. Returns a SharingSpec with Enabled=false
// if no shareable services are present (no prompt shown in that case).
func (pm *PromptManager) promptForSharing(serviceConfigs []types.ServiceConfig) (*clicontext.SharingSpec, error) {
	var shareableNames []string
	for _, cfg := range serviceConfigs {
		if cfg.Shareable {
			shareableNames = append(shareableNames, cfg.Name)
		}
	}

	if len(shareableNames) == 0 {
		return &clicontext.SharingSpec{Enabled: false}, nil
	}

	prompt := &survey.Confirm{
		Message: fmt.Sprintf(messages.PromptsEnableSharing, strings.Join(shareableNames, ", ")),
		Help:    messages.PromptsEnableSharingHelp,
		Default: true,
	}
	var enabled bool
	if err := survey.AskOne(prompt, &enabled); err != nil {
		return nil, err
	}
	return &clicontext.SharingSpec{Enabled: enabled}, nil
}

func (pm *PromptManager) promptForAutoStart() (bool, error) {
	prompt := &survey.Confirm{
		Message: messages.PromptsAutoStart,
		Help:    messages.PromptsAutoStartHelp,
		Default: false,
	}
	var autoStart bool
	if err := survey.AskOne(prompt, &autoStart); err != nil {
		return false, err
	}
	return autoStart, nil
}

// InitConfirmation encapsulates initialization confirmation parameters
type InitConfirmation struct {
	ProjectName string
	Services    []string
	Validation  map[string]bool
	Base        *base.BaseCommand
}

// ConfirmInitialization shows final confirmation with option to go back
func (pm *PromptManager) ConfirmInitialization(projectName string, services []string, validation map[string]bool, base *base.BaseCommand) (string, error) {
	conf := InitConfirmation{
		ProjectName: projectName,
		Services:    services,
		Validation:  validation,
		Base:        base,
	}
	return pm.confirmInitializationWithConfig(conf)
}

func (pm *PromptManager) confirmInitializationWithConfig(conf InitConfirmation) (string, error) {
	conf.Base.Output.Info(messages.InfoProjectConfigSummary)
	conf.Base.Output.Info(messages.InfoProjectNameLabel, conf.ProjectName)
	conf.Base.Output.Info(messages.InfoServicesLabel, strings.Join(conf.Services, ", "))

	if len(conf.Validation) > 0 {
		conf.Base.Output.Info(messages.InfoValidationOptions)
		// Sort keys for consistent display order
		keys := make([]string, 0, len(conf.Validation))
		for k := range conf.Validation {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, option := range keys {
			conf.Base.Output.Info(messages.InfoListItemIndented, option)
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
