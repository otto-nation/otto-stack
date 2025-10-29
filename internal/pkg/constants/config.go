package constants

// Default configuration values
const (
	DefaultEnvironment = "local"
	DefaultProjectName = "otto-stack"
)

// Configuration sections
const (
	ProjectSection              = "project"
	StackSection                = "stack"
	ServiceConfigurationSection = "service-configuration"
	ValidationSection           = "validation"
	AdvancedSection             = "advanced"
)

// Service types
const (
	ServiceTypeContainer     = "container"
	ServiceTypeConfiguration = "configuration"
	ServiceTypeComposite     = "composite"
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
