package constants

// Default configuration values
const (
	DefaultEnvironment = "local"
	DefaultProjectName = AppName
)

// Service names
const (
	ServiceKafkaTopics    = "kafka-topics"
	ServiceLocalstackInit = "localstack-init"
)

// User action responses
const (
	ActionProceed = "proceed"
	ActionBack    = "back"
	ActionCancel  = "cancel"
)

// Default configuration values
const (
	DefaultSkipWarnings      = false
	DefaultAllowMultipleDBs  = true
	DefaultAutoStart         = true
	DefaultPullLatestImages  = true
	DefaultCleanupOnRecreate = false
)
