package constants

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

// URLs
const (
	DockerInstallURL = "https://docs.docker.com/get-docker/"
)

// Generic service references
const (
	ServiceLocalhost = "localhost"
)

// Error messages
var (
	ErrNotInitialized = AppName + " not initialized. Run '" + AppName + " init' first"
)

// Status formatting
const (
	StatusHeaderService = "SERVICE"
	StatusHeaderState   = "STATE"
	StatusHeaderHealth  = "HEALTH"
	StatusSeparator     = "-"
)
