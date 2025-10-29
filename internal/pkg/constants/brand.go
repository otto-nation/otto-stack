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
)
