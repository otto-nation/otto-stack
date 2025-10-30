package project

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/types"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
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
	force, _ := cmd.Flags().GetBool(constants.FlagForce)

	ui.Header(constants.MsgInitializing)

	// Validate environment
	if err := h.validateInitEnvironment(); err != nil && !force {
		return fmt.Errorf("%s", constants.MsgValidationFailed.Content)
	}

	// Validate directory structure
	if err := h.validateDirectoryStructure(); err != nil && !force {
		return fmt.Errorf("%s", constants.MsgDirectoryValidationFailed.Content)
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
		services, err = h.promptForServices()
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
		action, err := h.confirmInitializationWithBack(projectName, services, validation, advanced)
		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}

		switch action {
		case constants.ActionProceed:
			goto exitLoop
		case constants.ActionBack:
			constants.SendMessage(constants.MsgGoingBack)
			continue
		default:
			constants.SendMessage(constants.MsgInitCancelled)
			return nil
		}
	}
exitLoop:

	// Create directory structure
	if err := h.createDirectoryStructure(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Create configuration file
	if err := h.createConfigFile(projectName, services, validation, advanced); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	// Generate initial compose files
	if err := h.generateInitialComposeFiles(services, projectName, validation, advanced); err != nil {
		return fmt.Errorf("failed to generate compose files: %w", err)
	}

	// Create .gitignore entries
	if err := h.createGitignoreEntries(); err != nil {
		constants.SendMessage(constants.MsgFailedGitignore, err)
	}

	// Create README
	if err := h.createReadme(projectName, services); err != nil {
		constants.SendMessage(constants.MsgFailedReadme, err)
	}

	ui.Success(constants.MsgInitSuccess)
	constants.SendMessage(constants.MsgNextSteps)
	constants.SendMessage(constants.MsgStep1, constants.DevStackDir, constants.ConfigFileName)
	constants.SendMessage(constants.MsgStep2, constants.AppName+" up")
	constants.SendMessage(constants.MsgStep3, constants.AppName+" status")

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
