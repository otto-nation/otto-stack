package services

// Services-specific error constants
const (
	// Components
	ComponentServices       = "services"
	ComponentServiceManager = "service-manager"
	ComponentFormatter      = "formatter"

	// Actions
	ActionLoadServices          = "load services"
	ActionCreateManager         = "create service manager"
	ActionResolveServices       = "resolve services"
	ActionLoadCatalog           = "load service catalog"
	ActionReadServicesDirectory = "read services directory"
	ActionLoadCategory          = "load category"
	ActionReadCategoryDirectory = "read category directory"
	ActionLoadService           = "load service"
	ActionParseServiceYAML      = "parse service YAML"
)
