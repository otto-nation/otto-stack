package stack

// Stack-specific error constants
const (
	// Components
	ComponentStack          = "stack"
	ComponentServiceManager = "service-manager"
	ComponentServices       = "services"

	// Actions
	ActionStartServices     = "start services"
	ActionStopServices      = "stop services"
	ActionRestartServices   = "restart services"
	ActionShowLogs          = "show logs"
	ActionShowStatus        = "show status"
	ActionCleanupResources  = "cleanup resources"
	ActionCreateService     = "create service"
	ActionGetManager        = "get services manager"
	ActionCreateGenerator   = "create generator"
	ActionGenerateCompose   = "generate compose"
	ActionCreateDirectory   = "create directory"
	ActionGenerateEnv       = "generate env content"
	ActionRunInitContainer  = "run init container"
	ActionCreateInitConfig  = "create init container config"
	ActionLoadServiceConfig = "load service configuration"
	ActionCreateManager     = "create service manager"
	ActionResolveServices   = "resolve services"
	ActionParseFlags        = "parse command flags"

	// Docker operations
	OpListContainers  = "list containers"
	OpShowLogs        = "show logs"
	OpRemoveResources = "remove resources"
)
