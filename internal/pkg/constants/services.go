package constants

// Service type constants
const (
	ServiceTypeContainer     = "container"
	ServiceTypeConfiguration = "configuration"
	ServiceTypeComposite     = "composite"
)

// Service category constants
const (
	CategoryDatabase      = "database"
	CategoryCache         = "cache"
	CategoryMessaging     = "messaging"
	CategoryObservability = "observability"
	CategoryCloud         = "cloud"
)

// Service category display names and icons
var CategoryDisplayInfo = map[string]struct {
	Name string
	Icon string
}{
	CategoryDatabase:      {"Database", "📊"},
	CategoryCache:         {"Cache", "💾"},
	CategoryMessaging:     {"Messaging", "📨"},
	CategoryObservability: {"Observability", "🔍"},
	CategoryCloud:         {"Cloud", "☁️"},
}

// Service display format constants
const (
	ServiceCatalogTableFormat = "table"
	ServiceCatalogGroupFormat = "group"
	ServiceCatalogJSONFormat  = "json"
	ServiceCatalogYAMLFormat  = "yaml"
)

// Service catalog messages
const (
	MsgServiceCatalogHeader = "Available Services by Category"
	MsgNoServicesInCategory = "No services found in category: %s"
	MsgServiceCount         = "%s (%d service%s)"
)
