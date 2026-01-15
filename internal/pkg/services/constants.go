package services

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
