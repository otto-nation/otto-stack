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
	ConfigFileName       = "otto-stack-config.yml"
	LocalConfigFileName  = "otto-stack-config.local.yml"
	ReadmeFileName       = "README.md"
	GitIgnoreFileName    = ".gitignore"
	StateFileName        = "state.json"
	OttoStackDir         = AppName
	EnvGeneratedFileName = ".env.generated"
)

// File and directory permissions
const (
	PermReadWrite     = 0644
	PermReadWriteExec = 0755
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
	MsgStopping       = "Stopping " + AppNameTitle
	MsgStopSuccess    = AppNameLower + " stopped successfully"
	MsgLogs           = "Logs"
	MsgRestarting     = "Restarting"
	MsgRestartSuccess = "Restart successful"
	MsgStatus         = "Status"
	MsgStarting       = "Starting " + AppNameTitle
	MsgStartSuccess   = AppNameTitle + " started successfully"
)

// Prompt constants
const (
	PromptProjectName       = "Enter project name:"
	HelpProjectName         = "The project name will be used to identify your stack"
	PromptGoBack            = "← Go back"
	HelpServiceSelection    = "Select services to include in your stack"
	PromptValidationOptions = "Select validation options:"
	HelpValidationOptions   = "Choose which validations to run"
	PromptActionSelect      = "What would you like to do?"
	PromptProceedInit       = "Proceed with initialization?"
)

// Template constants
const (
	MultiSelectTemplateWithBack = "multiselect_with_back"
)

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

// Action constants
const (
	ActionProceed = "proceed"
	ActionBack    = "back"
	ActionCancel  = "cancel"
)

// GitignoreEntries contains default .gitignore entries for otto-stack projects
var GitignoreEntries = []string{
	"# " + AppNameTitle + " generated files",
	EnvGeneratedFileName,
	AppName + "/",
	"*.log",
	".DS_Store",
}
