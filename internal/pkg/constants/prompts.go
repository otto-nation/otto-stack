package constants

// Prompt messages
const (
	// Project setup prompts
	PromptProjectName    = "Project name:"
	PromptEnvironment    = "Environment:"
	PromptAdvancedConfig = "Configure advanced options?"
	PromptConfirmInit    = "Proceed with initialization?"

	// Action prompts
	PromptActionSelect = "What would you like to do?"
	PromptProceedInit  = "Proceed with initialization"
	PromptGoBack       = "Go back to change services"
	PromptCancel       = "Cancel"

	// Service selection prompts
	PromptValidationOptions = "Validation options:"
	PromptAdvancedFeatures  = "Advanced features:"
)

// Prompt help text
const (
	HelpProjectName       = "Enter a name for your project (letters, numbers, hyphens, underscores only)"
	HelpEnvironment       = "Select the environment for this project"
	HelpServiceSelection  = "Use space to select/deselect, enter to confirm. You can go back at the confirmation step."
	HelpAdvancedConfig    = "Enable additional configuration options for validation, monitoring, etc."
	HelpValidationOptions = "Select validation features to enable"
	HelpAdvancedFeatures  = "Select additional features to enable"
)

// Default values
const (
	DefaultEnvironmentValue = "local"
)

// Option lists
var (
	EnvironmentOptions = []string{DefaultEnvironmentValue}

	ActionOptions = []string{
		PromptProceedInit,
		PromptGoBack,
		PromptCancel,
	}
)

// Option mappings for easier processing
var (
	ValidationOptions = map[string]string{
		"Skip dependency warnings": "skip_warnings",
		"Allow multiple databases": "allow_multiple_databases",
	}

	AdvancedOptions = map[string]string{
		"Auto-start services on boot":        "auto_start",
		"Pull latest images before starting": "pull_latest_images",
		"Clean up containers on recreate":    "cleanup_on_recreate",
	}

	ActionOptionMap = map[string]string{
		PromptProceedInit: ActionProceed,
		PromptGoBack:      ActionBack,
		PromptCancel:      ActionCancel,
	}
)
