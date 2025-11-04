package project

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/spf13/cobra"
)

// InitHandler handles the init command
type InitHandler struct {
	serviceUtils *utils.ServiceUtils
}

// NewInitHandler creates a new InitHandler
func NewInitHandler() *InitHandler {
	return &InitHandler{
		serviceUtils: utils.NewServiceUtils(),
	}
}

// Handle executes the init command
func (h *InitHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	logger.Info(constants.LogMsgStartingOperation, constants.LogFieldOperation, constants.OperationInit)
	defer func() {
		if r := recover(); r != nil {
			logger.Error(constants.LogMsgOperationFailed, constants.LogFieldOperation, constants.OperationInit, constants.LogFieldError, fmt.Errorf("panic: %v", r))
			panic(r)
		}
	}()

	// Parse all flags with validation - single line!
	flags, err := constants.ParseInitFlags(cmd)
	if err != nil {
		logger.Error(constants.LogMsgOperationFailed, constants.LogFieldOperation, constants.OperationInit, constants.LogFieldError, err)
		return err
	}

	base.Output.Header("%s", constants.Messages[constants.MsgProcess_initializing])
	logger.Info(constants.LogMsgProjectAction, constants.LogFieldAction, constants.CommandInit, constants.LogFieldProject, "current_directory")

	// Validate environment
	if err := h.validateInitEnvironment(base); err != nil && !flags.Force {
		return fmt.Errorf(constants.Messages[constants.MsgValidation_failed], err)
	}

	// Validate directory structure
	if err := h.validateDirectoryStructure(base); err != nil && !flags.Force {
		return fmt.Errorf(constants.Messages[constants.MsgValidation_directory_failed], err)
	}

	// Prompt for project details
	projectName, err := h.promptForProjectDetails()
	if err != nil {
		return fmt.Errorf("failed to get project details: %w", err)
	}

	// Service selection loop (allows going back)
	var services []string
	var validation, advanced map[string]bool

	for {
		// Prompt for services
		services, err = h.promptForServices(base)
		if err != nil {
			return fmt.Errorf("failed to select services: %w", err)
		}

		// Validate selected services
		if err := h.validateServices(services); err != nil {
			return fmt.Errorf("service validation failed: %w", err)
		}

		// Prompt for advanced options
		validation, advanced, err = h.promptForAdvancedOptions()
		if err != nil {
			return fmt.Errorf("failed to get advanced options: %w", err)
		}

		// Confirm initialization (with back option)
		action, err := h.confirmInitializationWithBack(projectName, services, validation, advanced, base)
		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}

		switch action {
		case constants.ActionProceed:
			goto exitLoop
		case constants.ActionBack:
			base.Output.Info("%s", constants.Messages[constants.MsgInit_going_back])
			continue
		default:
			base.Output.Info("%s", constants.Messages[constants.MsgInit_cancelled])
			return nil
		}
	}
exitLoop:

	// Create directory structure
	if err := h.createDirectoryStructure(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Create configuration file
	if err := h.createConfigFile(projectName, services, base); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	// Generate initial compose files
	if err := h.generateInitialComposeFiles(services, projectName, validation, advanced, base); err != nil {
		return fmt.Errorf("failed to generate compose files: %w", err)
	}

	// Create .gitignore entries
	if err := h.createGitignoreEntries(base); err != nil {
		base.Output.Warning(constants.Messages[constants.MsgWarnings_failed_gitignore], err)
	}

	// Create README
	if err := h.createReadme(projectName, services, base); err != nil {
		base.Output.Warning(constants.Messages[constants.MsgWarnings_failed_readme], err)
	}

	base.Output.Success("%s", constants.Messages[constants.MsgSuccess_init])
	base.Output.Info("%s", constants.Messages[constants.MsgInit_next_steps])
	base.Output.Info(constants.Messages[constants.MsgInit_step_review_config], constants.OttoStackDir, constants.ConfigFileName)
	base.Output.Info(constants.Messages[constants.MsgInit_step_start_stack], constants.AppName)
	base.Output.Info(constants.Messages[constants.MsgInit_step_check_status], constants.AppName)

	return nil
}

// ValidateArgs validates the command arguments
func (h *InitHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *InitHandler) GetRequiredFlags() []string {
	return []string{}
}
