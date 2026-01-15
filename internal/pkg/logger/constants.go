package logger

// Log levels
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

// Log formats
const (
	LogFormatJSON = "json"
	LogFormatText = "text"
)

// Log field names
const (
	LogFieldResult    = "result"
	LogFieldAction    = "action"
	LogFieldProject   = "project"
	LogFieldVersion   = "version"
	LogFieldFormat    = "format"
	LogFieldBuildInfo = "build_info"
	LogFieldOperation = "operation"
	LogFieldError     = "error"
)

// Log message templates
const (
	LogMsgProjectAction     = "project_action"
	LogMsgStartingOperation = "starting_operation"
	LogMsgOperationFailed   = "operation_failed"
)

// Operation constants
const (
	OperationInit      = "init"
	OperationStackDown = "stack_down"
	OperationStackUp   = "stack_up"
)

// Log field constants
const (
	LogFieldServices = "services"
	LogFieldService  = "service"
)

// Log message constants
const (
	LogMsgServiceAction      = "service_action"
	LogMsgOperationCompleted = "operation_completed"
)

// Action constants
const (
	ActionStop  = "stop"
	ActionStart = "start"
)
