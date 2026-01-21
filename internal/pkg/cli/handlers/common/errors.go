package common

// Stack-specific error constants
const (
	// Components
	ComponentStack    = "stack"
	ComponentService  = "service"
	ComponentRegistry = "registry"

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
	ActionLoadService      = "load service"
	ActionRegister         = "register"
	ActionUnregister       = "unregister"
	ActionSaveRegistry     = "save registry"
	ActionLoadRegistry     = "load registry"

	// Docker operations
	OpShowLogs        = "show logs"
	OpListContainers  = "list containers"
	OpRemoveResources = "remove resources"

	// Legacy messages (to be removed - use core.Msg* instead)
	MsgFailedCreateStackService = "failed to create stack service"
	MsgUnsupportedService       = "unsupported service"

	// Context types
	ContextGlobal = "global"
)
