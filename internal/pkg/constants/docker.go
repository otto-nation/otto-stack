package constants

// Docker container states
const (
	StateRunning  = "running"
	StateStopped  = "exited"
	StateCreated  = "created"
	StateStarting = "starting"
	StatePaused   = "paused"
)

// Health statuses
const (
	HealthHealthy   = "healthy"
	HealthUnhealthy = "unhealthy"
	HealthStarting  = "starting"
	HealthNone      = "none"
)

// Summary keys
const (
	SummaryTotal   = "total"
	SummaryRunning = "running"
	SummaryHealthy = "healthy"
)

// Docker Compose labels
const (
	ComposeProjectLabel = "com.docker.compose.project"
	ComposeServiceLabel = "com.docker.compose.service"
)

// Docker file paths
const (
	DockerComposeFile = OttoStackDir + "/" + DockerComposeFileName
)

// Docker Compose field names
const (
	ComposeFieldServices    = "services"
	ComposeFieldNetworks    = "networks"
	ComposeFieldVolumes     = "volumes"
	ComposeFieldImage       = "image"
	ComposeFieldPorts       = "ports"
	ComposeFieldEnvironment = "environment"
	ComposeFieldRestart     = "restart"
	ComposeFieldCommand     = "command"
	ComposeFieldDependsOn   = "depends_on"
	ComposeFieldName        = "name"
)
