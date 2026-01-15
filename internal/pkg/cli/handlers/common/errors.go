package common

// Stack-specific error constants
const (
	// Components
	ComponentStack = "stack"

	// Actions
	ActionStartServices    = "start services"
	ActionStopServices     = "stop services"
	ActionRestartServices  = "restart services"
	ActionShowLogs         = "show logs"
	ActionShowStatus       = "show status"
	ActionCleanupResources = "cleanup resources"
	ActionCreateService    = "create service"
	ActionGetManager       = "get services manager"
	ActionCreateGenerator  = "create generator"
	ActionGenerateCompose  = "generate compose"
	ActionCreateDirectory  = "create directory"
	ActionGenerateEnv      = "generate env content"
	ActionCreateManager    = "create manager"
	ActionLoadProject      = "load project"
	ActionCreateClient     = "create client"
	ActionGetServiceStatus = "get service status"
	ActionConnectToService = "connect to service"
	ActionExecuteCommand   = "execute command"
	ActionGetLogs          = "get logs"
	ActionValidateArgs     = "validate arguments"
	ActionBuildContext     = "build context"
	ActionResolveServices  = "resolve services"
	ActionFilterServices   = "filter services"

	// Docker operations
	OpShowLogs        = "show logs"
	OpListContainers  = "list containers"
	OpRemoveResources = "remove resources"

	// Messages
	MsgFailedCreateStackService = "failed to create stack service"
	MsgUnsupportedService       = "unsupported service"
)
