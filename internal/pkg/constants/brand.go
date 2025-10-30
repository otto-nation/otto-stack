package constants

// Brand constants for consistent naming
const (
	// Application name variations
	AppName      = "otto-stack" // CLI command name
	AppNameTitle = "Otto Stack" // Title case for headers
	AppNameLower = "otto stack" // Sentence case for messages

	// Common messages
	MsgInitializing = "Initializing " + AppNameTitle
	MsgStarting     = "Starting " + AppNameTitle
	MsgStopping     = "Stopping " + AppNameTitle
	MsgRestarting   = "Restarting " + AppNameTitle
	MsgStatus       = AppNameTitle + " Status"
	MsgLogs         = AppNameTitle + " Logs"
	MsgConnecting   = "Connecting to %s"
	MsgCleaning     = "Cleaning up " + AppNameTitle

	// Success messages
	MsgInitSuccess    = AppNameLower + " initialized successfully!"
	MsgStartSuccess   = AppNameLower + " started successfully"
	MsgStopSuccess    = AppNameLower + " stopped successfully"
	MsgRestartSuccess = AppNameLower + " restarted successfully"

	// GitHub repository information
	GitHubOrg  = "otto-nation"
	GitHubRepo = AppName

	// GitHub URL templates
	GitHubRepoURL     = "https://github.com/" + GitHubOrg + "/" + GitHubRepo
	GitHubReleaseURL  = GitHubRepoURL + "/releases/tag/v%s"
	GitHubDownloadURL = GitHubRepoURL + "/releases/download/v%s/" + AppName
)

// Version file patterns
var VersionFilePatterns = []string{
	".otto-stack-version",
}

// Version search paths
var VersionSearchPaths = []string{
	".",
	".otto-stack",
	".config",
	"config",
}

// Version prefixes for cleaning
var VersionPrefixes = []string{
	"otto-stack-",
	"v",
}

// Version format templates
const (
	AppNameTemplate   = "otto-stack %s"
	UserAgentTemplate = "otto-stack/%s (%s/%s)"
)
