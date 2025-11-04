package constants

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

// Version format templates
const (
	AppNameTemplate   = AppName + " %s"
	UserAgentTemplate = AppName + "/%s (%s/%s)"
)

// File names
const (
	ConfigFileName         = AppName + "-config.yml"
	LocalConfigFileName    = AppName + "-config.local.yml"
	DockerComposeFileName  = "docker-compose.yml"
	EnvGeneratedFileName   = ".env.generated"
	GitignoreFileName      = ".gitignore"
	ReadmeFileName         = "README.md"
	ServiceConfigExtension = ".yaml"
	ServiceConfigExtAlt    = ".yml"
	KafkaTopicsInitScript  = "kafka-topics-init.sh"
	LocalstackInitScript   = "localstack-init.sh"
	StateFileName          = "state.json"
)

// Directory names
const (
	OttoStackDir        = AppName
	DataDir             = "data"
	LogsDir             = "logs"
	TmpDir              = "tmp"
	ScriptsDir          = "scripts"
	ServicesDir         = "internal/config/services"
	EmbeddedServicesDir = "services" // Directory name in embedded FS
)

// Configuration URLs
const (
	ConfigDocsURL    = "https://github.com/otto-nation/otto-stack/tree/main/docs-site/content/configuration.md"
	ServiceConfigURL = "https://github.com/otto-nation/otto-stack/tree/main/internal/config/services"
	DockerInstallURL = "https://docs.docker.com/get-docker/"
)

// Git entries
var GitignoreEntries = []string{
	"",
	"# Otto Stack",
	OttoStackDir + "/" + EnvGeneratedFileName,
	OttoStackDir + "/" + DataDir + "/",
	OttoStackDir + "/" + LogsDir + "/",
	OttoStackDir + "/" + TmpDir + "/",
	"",
	"# Local config overrides",
	OttoStackDir + "/" + LocalConfigFileName,
}

// Common application messages
const (
	MsgInitializing   = "Initializing " + AppNameTitle
	MsgStarting       = "Starting " + AppNameTitle
	MsgStopping       = "Stopping " + AppNameTitle
	MsgRestarting     = "Restarting " + AppNameTitle
	MsgStatus         = AppNameTitle + " Status"
	MsgLogs           = AppNameTitle + " Logs"
	MsgConnecting     = "Connecting to %s"
	MsgCleaning       = "Cleaning up " + AppNameTitle
	MsgInitSuccess    = AppNameLower + " initialized successfully!"
	MsgStartSuccess   = AppNameLower + " started successfully"
	MsgStopSuccess    = AppNameLower + " stopped successfully"
	MsgRestartSuccess = AppNameLower + " restarted successfully"
)

// Exit codes
const (
	ExitSuccess = 0
	ExitError   = 1
)
