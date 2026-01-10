package core

import (
	"fmt"
	"os"
	"strings"
)

// Application identity
const (
	AppName      = "otto-stack" // CLI command name
	AppNameTitle = "Otto Stack" // Title case for headers
	AppNameLower = "otto stack" // Sentence case for messages
)

// Host constants
const (
	ServiceLocalhost = "localhost"
)

// AWS constants
const (
	DefaultAWSRegion = "us-east-1"
)

// GitHub repository information
const (
	GitHubOrg  = "otto-nation"
	GitHubRepo = AppName
)

// GitHub URLs
const (
	GitHubRepoURL     = "https://github.com/" + GitHubOrg + "/" + GitHubRepo
	GitHubReleaseURL  = GitHubRepoURL + "/releases/tag/v%s"
	GitHubDownloadURL = GitHubRepoURL + "/releases/download/v%s/" + AppName
)

// Documentation URLs
const (
	DocsURL = "https://" + GitHubOrg + ".github.io/" + GitHubRepo + "/"
)

// File names
const (
	ConfigFileName       = "otto-stack-config.yml"
	LocalConfigFileName  = "otto-stack-config" + LocalFileExtension + ".yml"
	ServiceConfigsDir    = "service-configs"
	ScriptsDir           = "scripts"
	GeneratedDir         = "generated"
	ReadmeFileName       = "README.md"
	GitIgnoreFileName    = ".gitignore"
	StateFileName        = "state.json"
	OttoStackDir         = "." + AppName
	EnvGeneratedFileName = ".env.generated"
	LocalFileExtension   = ".local"
	YAMLFileExtension    = ".yaml"
	YMLFileExtension     = ".yml"
)

// File paths
const (
	StateFilePath        = OttoStackDir + "/" + GeneratedDir + "/" + StateFileName
	EnvGeneratedFilePath = OttoStackDir + "/" + GeneratedDir + "/" + EnvGeneratedFileName
)

var YAMLExtensions = []string{YMLFileExtension, YAMLFileExtension}

// File and directory permissions
const (
	PermReadWrite     = 0644
	PermReadWriteExec = 0755
)

// Environment variables
const (
	EnvOttoNonInteractive = "OTTO_NON_INTERACTIVE"
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
	DefaultHTTPTimeoutSeconds  = 5
)

// HTTP constants
const (
	HTTPOKStatusThreshold = 400
)

// Log constants
const (
	DefaultLogTailLines = "100"
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
	PromptFinishSelection   = "Finish selection"
	PromptSelectCategory    = "Select a service category (or finish selection):"
	PromptGoBackOption      = "Go Back"
)

// Template constants
const (
	MultiSelectTemplateWithBack = "multiselect_with_back"
)

// IsYAMLFile checks if a filename has a YAML extension
func IsYAMLFile(filename string) bool {
	return strings.HasSuffix(filename, YAMLFileExtension) || strings.HasSuffix(filename, YMLFileExtension)
}

// TrimYAMLExt removes YAML extension from filename
func TrimYAMLExt(filename string) string {
	if name, found := strings.CutSuffix(filename, YAMLFileExtension); found {
		return name
	}
	if name, found := strings.CutSuffix(filename, YMLFileExtension); found {
		return name
	}
	return filename
}

// FindYAMLFile finds a YAML file with either .yml or .yaml extension
func FindYAMLFile(dir, baseName string) (string, error) {
	for _, ext := range YAMLExtensions {
		path := fmt.Sprintf("%s/%s%s", dir, baseName, ext)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("YAML file not found: %s", baseName)
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
	AppName + "/" + GeneratedDir + "/",
	"*.log",
	".DS_Store",
}
