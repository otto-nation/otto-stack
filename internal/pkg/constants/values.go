package constants

import (
	"os"
	"strings"
)

// File and directory permissions
const (
	FilePermReadWrite    = 0644 // Standard file read/write permissions
	DirPermReadWriteExec = 0755 // Standard directory permissions
	FilePermReadWriteAll = 0666 // File permissions for all users
)

// Validation limits
const (
	MinProjectNameLength = 2
	MaxProjectNameLength = 50
	MinArgumentCount     = 2 // For command parsing
	MinFieldCount        = 2 // For field parsing
	MaxCategoryCommands  = 10
	GitCommitHashLength  = 7
)

// Timeouts and intervals (in seconds)
const (
	DefaultStopTimeoutSeconds    = 10
	DefaultStartTimeoutSeconds   = 30
	DefaultExecTimeoutSeconds    = 30
	DefaultConnectTimeoutSeconds = 10
	DefaultBackupTimeoutSeconds  = 300
	DefaultScaleTimeoutSeconds   = 60
	HealthCheckIntervalSeconds   = 2
	SpinnerIntervalMilliseconds  = 100
)

// Display formatting
const (
	SeparatorLength       = 50
	StatusSeparatorLength = 45
	TableWidth42          = 42
	TableWidth75          = 75
	TableWidth80          = 80
	TableWidth85          = 85
	TableWidth90          = 90
	HoursPerDay           = 24
	UIPadding             = 2
)

// Validation thresholds (percentages)
const (
	PercentageMultiplier   = 100
	MinExampleCoverage     = 80
	MinTipsCoverage        = 50
	MinDescriptionCoverage = 60
	MinConfigurationScore  = 80
	MaxValidationErrors    = 5
	BaseValidationScore    = 100.0
	ErrorWeight            = 10
	WarningWeight          = 2
)

// Protocol constants
const (
	ProtocolTCP = "tcp"
	ProtocolUDP = "udp"
)

// Volume mode constants
const (
	VolumeModeReadOnly  = "ro"
	VolumeModeReadWrite = "rw"
)

// Format constants
const (
	FormatJSON  = "json"
	FormatYAML  = "yaml"
	FormatTable = "table"
)

// Version and parsing
const (
	MaxVersionNumber = 999
	KeyValueParts    = 2 // For splitting "key=value" strings
	PortSearchRange  = 100
	HexDivisor       = 2 // For hex string conversion
)

// Version defaults
const (
	DefaultVersion   = "dev"
	DefaultCommit    = "unknown"
	DefaultBuildDate = "unknown"
	DefaultBuildBy   = "unknown"
	DevelVersion     = "(devel)"
)

// Version comparison results
const (
	VersionEqual   = 0
	VersionNewer   = 1
	VersionOlder   = -1
	VersionInvalid = -999
)

// UI ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[32m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorGray   = "\033[90m"
	ColorBold   = "\033[1m"
)

// UI message prefixes
const (
	IconSuccess = "✅"
	IconError   = "❌"
	IconWarning = "⚠️ "
	IconInfo    = "ℹ️ "
	IconHeader  = "🚀"
	IconBox     = "📦"
)

// Log levels
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

// Log formats
const (
	LogFormatText = "text"
	LogFormatJSON = "json"
)

// Operation names for logging
const (
	OperationStackUp   = "stack_up"
	OperationStackDown = "stack_down"
	OperationInit      = "project_init"
	OperationDoctor    = "project_doctor"
)

// Action names for logging
const (
	ActionStart = "start"
	ActionStop  = "stop"
	ActionCheck = "check"
)

// Log message templates
const (
	LogMsgStartingOperation  = "Starting operation"
	LogMsgOperationCompleted = "Operation completed"
	LogMsgOperationFailed    = "Operation failed"
	LogMsgServiceAction      = "Service action"
	LogMsgProjectAction      = "Project action"
)

// Log field names
const (
	LogFieldError     = "error"
	LogFieldOperation = "operation"
	LogFieldAction    = "action"
	LogFieldProject   = "project"
	LogFieldService   = "service"
	LogFieldServices  = "services"
	LogFieldResult    = "result"
	LogFieldVersion   = "version"
	LogFieldFormat    = "format"
	LogFieldBuildInfo = "build_info"
)

// Environment file constants
const (
	EnvFileExtension = ".env"
	EnvLocalFile     = ".env.local"
	EnvExampleFile   = ".env.example"
)

// Command constants
const (
	CmdPgrep    = "pgrep"
	CmdTasklist = "tasklist"
	CmdLsof     = "lsof"
	CmdNetstat  = "netstat"
	CmdTaskkill = "taskkill"
)

// OS constants (using Go's runtime values)
const (
	OSWindows = "windows"
	OSLinux   = "linux"
	OSDarwin  = "darwin"
)

// Error message constants
const (
	ErrUnsupportedOS    = "unsupported operating system: %s"
	ErrProcessNotFound  = "process not found: %s"
	ErrNoFreePort       = "no free port found in range %d-%d"
	ErrOperationTimeout = "operation timed out after %v"
	ErrFailedAfterRetry = "failed after %d attempts: %w"
	ErrInvalidChoice    = "Invalid choice. Please select 1-%d."
)

// Format constants
const (
	ByteUnit     = 1024
	ByteUnits    = "KMGTPE"
	TimeFormat1s = "%.1fs"
	TimeFormat1m = "%.1fm"
	TimeFormat1h = "%.1fh"
	TimeFormat1d = "%.1fd"
	ByteFormatB  = "%d B"
	ByteFormatKB = "%.1f %cB"
)

// Type constants
const (
	ProjectTypeGo        = "go"
	ProjectTypeNode      = "node"
	ProjectTypePython    = "python"
	ProjectTypeDocker    = "docker"
	ProjectTypeFullStack = "fullstack"
)

// Error codes
const (
	ErrCodeInvalidConfig    = "INVALID_CONFIG"
	ErrCodeServiceNotFound  = "SERVICE_NOT_FOUND"
	ErrCodeConnectionFailed = "CONNECTION_FAILED"
	ErrCodeOperationFailed  = "OPERATION_FAILED"
)

// Default configuration values
const (
	DefaultEnvironment       = "local"
	DefaultProjectName       = AppName
	DefaultProjectType       = "docker"
	DefaultLogLevel          = "info"
	DefaultColorOutput       = true
	DefaultCheckUpdates      = true
	DefaultSkipWarnings      = false
	DefaultAllowMultipleDBs  = true
	DefaultAutoStart         = true
	DefaultPullLatestImages  = true
	DefaultCleanupOnRecreate = false
)

// User action responses
const (
	ActionProceed = "proceed"
	ActionBack    = "back"
	ActionCancel  = "cancel"
)

// Shell types for completion
const (
	ShellBash       = "bash"
	ShellZsh        = "zsh"
	ShellFish       = "fish"
	ShellPowerShell = "powershell"
)

// Docker commands
const (
	DockerCmd        = "docker"
	DockerInfoCmd    = "info"
	DockerComposeCmd = "compose"
	DockerVersionCmd = "version"
)

// Docker container states
const (
	StateRunning  = "running"
	StateStopped  = "exited"
	StateCreated  = "created"
	StateStarting = "starting"
	StatePaused   = "paused"
)

// Health statuses
const (
	HealthHealthy   = "healthy"
	HealthUnhealthy = "unhealthy"
	HealthStarting  = "starting"
	HealthNone      = "none"
)

// Summary keys
const (
	SummaryTotal   = "total"
	SummaryRunning = "running"
	SummaryHealthy = "healthy"
)

// Docker Compose labels
const (
	ComposeProjectLabel = "com.docker.compose.project"
	ComposeServiceLabel = "com.docker.compose.service"
)

// Docker file paths
const (
	DockerComposeFile = OttoStackDir + "/" + DockerComposeFileName
)

// Docker Compose field names
const (
	ComposeFieldServices    = "services"
	ComposeFieldNetworks    = "networks"
	ComposeFieldVolumes     = "volumes"
	ComposeFieldImage       = "image"
	ComposeFieldPorts       = "ports"
	ComposeFieldEnvironment = "environment"
	ComposeFieldRestart     = "restart"
	ComposeFieldCommand     = "command"
	ComposeFieldDependsOn   = "depends_on"
	ComposeFieldName        = "name"
)

// Generic service references
const (
	ServiceLocalhost = "localhost"
)

// Status formatting
const (
	StatusHeaderService = "SERVICE"
	StatusHeaderState   = "STATE"
	StatusHeaderHealth  = "HEALTH"
	StatusSeparator     = "-"
)

// Service category display names and icons
var CategoryDisplayInfo = map[string]struct {
	Name string
	Icon string
}{
	CategoryDatabase:      {"Database", "📊"},
	CategoryCache:         {"Cache", "💾"},
	CategoryMessaging:     {"Messaging", "📨"},
	CategoryObservability: {"Observability", "🔍"},
	CategoryCloud:         {"Cloud", "☁️"},
}

// Service display format constants
const (
	ServiceCatalogTableFormat = "table"
	ServiceCatalogGroupFormat = "group"
	ServiceCatalogJSONFormat  = "json"
	ServiceCatalogYAMLFormat  = "yaml"
)

// Service catalog messages
const (
	MsgServiceCatalogHeader = "Available Services by Category"
	MsgNoServicesInCategory = "No services found in category: %s"
	MsgServiceCount         = "%s (%d service%s)"
)

// Service restart policies
const (
	RestartPolicyNo        = "no"
	RestartPolicyAlways    = "always"
	RestartPolicyOnFailure = "on-failure"
	RestartPolicyUnless    = "unless-stopped"
)

// Prompt messages
const (
	PromptProjectName       = "Project name:"
	PromptEnvironment       = "Environment:"
	PromptAdvancedConfig    = "Configure advanced options?"
	PromptConfirmInit       = "Proceed with initialization?"
	PromptActionSelect      = "What would you like to do?"
	PromptProceedInit       = "Proceed with initialization"
	PromptGoBack            = "Go back to change services"
	PromptCancel            = "Cancel"
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

// Type constants
const (
	TypeBool        = "bool"
	TypeString      = "string"
	TypeStringArray = "[]string"
	TypeInt         = "int"
)

// Severity constants
const (
	SeverityLow = "low"
)

// Boolean string constants
const (
	BoolTrue = "true"
)

// Survey templates
const (
	MultiSelectTemplateWithBack = `
{{- define "option"}}
    {{- if eq .SelectedIndex .CurrentIndex }}{{color .Config.Icons.SelectFocus.Format }}{{ .Config.Icons.SelectFocus.Text }}{{color "reset"}}{{else}} {{end}}
    {{- if index .Checked .CurrentOpt.Index }}{{color .Config.Icons.MarkedOption.Format }} {{ .Config.Icons.MarkedOption.Text }} {{else}}{{color .Config.Icons.UnmarkedOption.Format }} {{ .Config.Icons.UnmarkedOption.Text }} {{end}}
    {{- color "reset"}}
    {{- " "}}{{- .CurrentOpt.Value}}{{ if ne ($.GetDescription .CurrentOpt) "" }} - {{color "cyan"}}{{ $.GetDescription .CurrentOpt }}{{color "reset"}}{{end}}
{{end}}
{{- if .ShowHelp }}{{- color .Config.Icons.Help.Format }}{{ .Config.Icons.Help.Text }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color .Config.Icons.Question.Format }}{{ .Config.Icons.Question.Text }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }}{{ .FilterMessage }}{{color "reset"}}
{{- if .ShowAnswer}}{{color "cyan"}} {{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
	{{- "  "}}{{- color "cyan"}}[Use arrows to move, space to select,{{- if not .Config.RemoveSelectAll }} <right> to all,{{end}}{{- if not .Config.RemoveSelectNone }} <left> to none,{{end}} type to filter, go back at confirmation{{- if and .Help (not .ShowHelp)}}, {{ .Config.HelpInput }} for more help{{end}}]{{color "reset"}}
  {{- "\n"}}
  {{- range $ix, $option := .PageEntries}}
    {{- template "option" $.IterateOption $ix $option}}
  {{- end}}
{{- end}}`
)

// Schema constants
const (
	YAMLTypeString  = "string"
	YAMLTypeBoolean = "boolean"
	YAMLTypeArray   = "array"
	YAMLTypeObject  = "object"

	SectionStack         = "stack"
	SectionProject       = "project"
	SectionValidation    = "validation"
	SectionAdvanced      = "advanced"
	SectionVersionConfig = "version_config"
	SectionServiceConfig = "service_configuration"

	PropertyEnabled = "enabled"
	PropertyName    = "name"

	TemplateProjectName             = "{{project_name}}"
	TemplateOttoVersion             = "{{otto_version}}"
	TemplateConfigDocsURL           = "{{config_docs_url}}"
	TemplateServiceConfigURL        = "{{service_config_url}}"
	TemplateDefaultSkipWarnings     = "{{default_skip_warnings}}"
	TemplateDefaultAllowMultipleDBs = "{{default_allow_multiple_dbs}}"
	TemplateDefaultAutoStart        = "{{default_auto_start}}"
	TemplateDefaultPullLatest       = "{{default_pull_latest_images}}"
	TemplateDefaultCleanupRecreate  = "{{default_cleanup_on_recreate}}"

	YAMLIndent   = "  "
	YAMLComment  = "# %s\n"
	YAMLProperty = "%s%s: "
	YAMLSection  = "%s:\n"
	YAMLListItem = "    - %s\n"
	YAMLEnabled  = "  enabled:\n"
)

// IsYAMLFile checks if a filename has a YAML extension
func IsYAMLFile(filename string) bool {
	return strings.HasSuffix(filename, ServiceConfigExtension) || strings.HasSuffix(filename, ServiceConfigExtAlt)
}

// TrimYAMLExt removes YAML extension from filename
func TrimYAMLExt(filename string) string {
	if name, found := strings.CutSuffix(filename, ServiceConfigExtension); found {
		return name
	}
	if name, found := strings.CutSuffix(filename, ServiceConfigExtAlt); found {
		return name
	}
	return filename
}

// FindYAMLFile returns the path to a YAML file, preferring .yaml over .yml
func FindYAMLFile(basePath, name string) string {
	yamlPath := basePath + "/" + name + ServiceConfigExtension
	ymlPath := basePath + "/" + name + ServiceConfigExtAlt

	// Check if .yaml exists first, fallback to .yml
	if fileExists(yamlPath) {
		return yamlPath
	}
	return ymlPath
}

// fileExists checks if a file exists (helper for FindYAMLFile)
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
