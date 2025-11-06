package core

import "strings"

// Application identity
const (
	AppName      = "otto-stack" // CLI command name
	AppNameTitle = "Otto Stack" // Title case for headers
	AppNameLower = "otto stack" // Sentence case for messages
)

// GitHub repository information
const (
	GitHubOrg  = "otto-nation"
	GitHubRepo = AppName
)

// GitHub URL templates
const (
	GitHubRepoURL     = "https://github.com/" + GitHubOrg + "/" + GitHubRepo
	GitHubReleaseURL  = GitHubRepoURL + "/releases/tag/v%s"
	GitHubDownloadURL = GitHubRepoURL + "/releases/download/v%s/" + AppName
)

// File names
const (
	ConfigFileName      = "otto-stack-config.yml"
	LocalConfigFileName = ".otto-stack.yaml"
	EnvFileName         = ".env.generated"
	ReadmeFileName      = "README.md"
	GitIgnoreFileName   = ".gitignore"
	StateFileName       = "state.json"
	OttoStackDir        = "otto-stack"
)

// File permissions
const (
	FilePermReadWrite    = 0644
	DirPermReadWriteExec = 0755
)

// Exit codes
const (
	ExitSuccess = 0
	ExitError   = 1
)

// Messages
const (
	MsgCleaning = "Cleaning up " + AppNameTitle
)

// Timeouts
const (
	DefaultStartTimeoutSeconds = 30
	DefaultStopTimeoutSeconds  = 10
)

// Time constants
const (
	HoursPerDay = 24
)

// Display constants
const (
	MaxCategoryCommands = 3
)

// System constants
const (
	MinArgumentCount     = 2
	MinProjectNameLength = 3
	MaxProjectNameLength = 50
)

// Messages
const (
	MsgStopping           = "Stopping Otto Stack"
	MsgStopSuccess        = "otto stack stopped successfully"
	MsgLogs               = "Logs"
	MsgRestarting         = "Restarting"
	MsgRestartSuccess     = "Restart successful"
	MsgStatus             = "Status"
	MsgStarting           = "Starting Otto Stack"
	MsgStartSuccess       = "Otto Stack started successfully"
	EnvGeneratedFileName  = ".env.generated"
	DockerInstallURL      = "https://docs.docker.com/get-docker/"
	DockerComposeFileName = "docker-compose.yml"
)

// Status display constants
const (
	StatusHeaderService   = "Service"
	StatusHeaderState     = "State"
	StatusHeaderHealth    = "Health"
	StatusSeparator       = "-"
	StatusSeparatorLength = 50
)

// Action constants
const (
	ActionProceed = "proceed"
	ActionBack    = "back"
	ActionCancel  = "cancel"
)

// Prompt constants
const (
	PromptProjectName       = "Enter project name:"
	HelpProjectName         = "The project name will be used to identify your stack"
	PromptGoBack            = "← Go back"
	HelpServiceSelection    = "Select services to include in your stack"
	PromptAdvancedConfig    = "Configure advanced settings?"
	HelpAdvancedConfig      = "Advanced configuration allows you to customize service settings"
	PromptValidationOptions = "Select validation options:"
	HelpValidationOptions   = "Choose which validations to run"
	PromptAdvancedFeatures  = "Select advanced features:"
	HelpAdvancedFeatures    = "Choose additional features to enable"
	PromptActionSelect      = "What would you like to do?"
	PromptProceedInit       = "Proceed with initialization?"
)

// Template constants
const (
	MultiSelectTemplateWithBack = "multiselect_with_back"
)

// Options constants
const (
	ValidationOptionsKey = "validation_options"
	AdvancedOptionsKey   = "advanced_options"
	ActionOptionsKey     = "action_options"
)

// CLI Flags
const (
// Flags are now generated in cli_generated.go
)

// CLI Messages
const (
// Messages are now generated in cli_generated.go
)

// Functions are now generated in cli_generated.go

// IsYAMLFile checks if a filename has a YAML extension
func IsYAMLFile(filename string) bool {
	return strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml")
}

// TrimYAMLExt removes YAML extension from filename
func TrimYAMLExt(filename string) string {
	if name, found := strings.CutSuffix(filename, ".yaml"); found {
		return name
	}
	if name, found := strings.CutSuffix(filename, ".yml"); found {
		return name
	}
	return filename
}

// ActionOptionMap maps action keys to display names
var ActionOptionMap = map[string]string{
	ActionProceed: "Proceed",
	ActionBack:    "Go Back",
	ActionCancel:  "Cancel",
}

// ValidationOptions maps validation keys to descriptions
var ValidationOptions = map[string]string{
	"docker":   "Docker availability",
	"ports":    "Port conflicts",
	"services": "Service configuration",
}

// AdvancedOptions maps advanced option keys to descriptions
var AdvancedOptions = map[string]string{
	"monitoring": "Enable monitoring",
	"logging":    "Enable logging",
	"backup":     "Enable backup",
}

// ActionOptions contains available action options
var ActionOptions = []string{
	ActionProceed,
	ActionBack,
	ActionCancel,
}

// GitignoreEntries contains default .gitignore entries for otto-stack projects
var GitignoreEntries = []string{
	"# Otto Stack generated files",
	".env.generated",
	"otto-stack/",
	"*.log",
	".DS_Store",
}
