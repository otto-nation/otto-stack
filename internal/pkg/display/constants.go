package display

const (
	// Duration formatting
	HoursPerDay = 24

	// Port display limits
	MaxPortsDisplay = 12

	// Display values
	NotApplicable = "n/a"

	// Scope values
	ScopeLocal  = "local"
	ScopeShared = "shared"

	// State values
	StateNotFound = "not found"
	StateUnknown  = "unknown"

	// Table headers - Status
	HeaderService    = "SERVICE"
	HeaderScope      = "SCOPE"
	HeaderContainer  = "CONTAINER"
	HeaderProvidedBy = "PROVIDED BY"
	HeaderState      = "STATE"
	HeaderHealth     = "HEALTH"
	HeaderUptime     = "UPTIME"
	HeaderPorts      = "PORTS"
	HeaderUpdated    = "UPDATED"

	// Table headers - Catalog
	HeaderCategory    = "CATEGORY"
	HeaderDescription = "DESCRIPTION"

	// Table headers - Dependencies
	HeaderDependencies = "DEPENDENCIES"

	// Table headers - Web Interfaces
	HeaderInterface = "INTERFACE"
	HeaderURL       = "URL"
	HeaderStatus    = "STATUS"
)
