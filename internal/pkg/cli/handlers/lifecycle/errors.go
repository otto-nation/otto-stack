package lifecycle

// Stack-specific error constants
const (
	// Components
	ComponentStack = "stack"

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

	// File operations
	ComponentFile       = "file"
	ActionMarshalData   = "marshal data"
	ActionWriteFile     = "write file"
	ActionReadFile      = "read file"
	ActionUnmarshalData = "unmarshal data"

	// Template operations
	ComponentTemplate     = "template"
	ActionProcessTemplate = "process template"
	ActionExecuteTemplate = "execute template"

	// Docker operations
	OpListContainers  = "list containers"
	OpShowLogs        = "show logs"
	OpRemoveResources = "remove resources"

	// Error messages
	MsgUnsupportedService       = "unsupported service"
	MsgFailedCreateStackService = "failed to create stack service"
)
