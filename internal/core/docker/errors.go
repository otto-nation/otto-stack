package docker

// Docker-specific error constants
const (
	// Components
	ComponentDocker = "docker"

	// Actions
	ActionCreateClient         = "create client"
	ActionCreateCLI            = "create CLI"
	ActionInitializeCLI        = "initialize CLI"
	ActionCreateComposeService = "create compose service"
	ActionLoadComposeProject   = "load compose project"
)
