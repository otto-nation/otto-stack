package services

// Service constants
const (
	ServiceLocalhost = "localhost"
)

// Init container constants
const (
	InitServiceEndpointURL = "SERVICE_ENDPOINT_URL"
	InitServiceName        = "INIT_SERVICE_NAME"
	InitConfigDir          = "CONFIG_DIR"
)

// Network naming constants
const (
	NetworkNameSuffix = "-network"
)

// Service catalog formats
const (
	ServiceCatalogJSONFormat  = "json"
	ServiceCatalogYAMLFormat  = "yaml"
	ServiceCatalogTableFormat = "table"
)

// Service catalog messages
const (
	MsgServiceCatalogHeader = "Available Services"
	MsgServiceCount         = "Total services: %d"
	MsgCategoryServiceCount = "%s: %d service%s"
	SummaryTotal            = "Total"
	SummaryRunning          = "Running"
	SummaryHealthy          = "Healthy"
)

// Directory names
const (
	ServicesDir         = "internal/config/services"
	EmbeddedServicesDir = "services"
)
